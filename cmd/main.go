package main

import (
	"fmt"
	"image"

	"log"

	"image/color"
	"time"

	"minego/internal/identify"
	"minego/internal/imgpos"
	"minego/internal/solver"
	"minego/internal/window"

	"minego/pkg/imageproc"
	"minego/pkg/kit"
	"minego/pkg/winapi/click"

	"minego/pkg/screenshot"
)

var (
	BorderColor = color.RGBA{7, 8, 9, 255}
)

// main 函数是程序的入口点，用于执行扫雷游戏识别任务
// 主要流程包括：截图、定位扫雷区域、裁剪图像、保存中间结果、网格检测和雷区识别
const (
	windowBorderInset = 10 // 窗口边界内缩像素
	gridBorderExpand  = 3  // 雷区边界扩展像素
)

func getMineFieldBounds() (image.Rectangle, error) {
	mineSweeperWindow := window.GetMineSweeperWindow()
	mineSweeperWindow.Activate()

	time.Sleep(50 * time.Millisecond)

	windowBounds, err := mineSweeperWindow.GetBounds()
	if err != nil {
		log.Fatalf("获取窗口边界失败: %v", err)
	}

	// 安全调整窗口边界
	windowBounds.Min.X += windowBorderInset
	windowBounds.Min.Y += windowBorderInset
	windowBounds.Max.X -= windowBorderInset
	windowBounds.Max.Y -= windowBorderInset
	windowImg, err := screenshot.CaptureRect(windowBounds)
	mineField := kit.FindSurroundingRect(windowImg, BorderColor)
	mineFieldBounds := image.Rect(
		windowBounds.Min.X+mineField.Min.X,
		windowBounds.Min.Y+mineField.Min.Y,
		windowBounds.Min.X+mineField.Dx()+mineField.Min.X,
		windowBounds.Min.Y+mineField.Dy()+mineField.Min.Y)

	fmt.Println("雷区边界:", mineFieldBounds)
	mineFieldBounds.Min.X = min(mineFieldBounds.Min.X-gridBorderExpand, mineFieldBounds.Min.X)
	mineFieldBounds.Min.Y = min(mineFieldBounds.Min.Y-gridBorderExpand, mineFieldBounds.Min.Y)
	mineFieldBounds.Max.X = max(mineFieldBounds.Max.X+gridBorderExpand, mineFieldBounds.Max.X)
	mineFieldBounds.Max.Y = max(mineFieldBounds.Max.Y+gridBorderExpand, mineFieldBounds.Max.Y)
	fmt.Println("最终雷区边界:", mineFieldBounds)
	return mineFieldBounds, nil
}

func main() {
	click.SetDPIAware()

	mineFieldBounds, err := getMineFieldBounds()
	if err != nil {
		log.Fatalf("获取窗口边界失败: %v", err)
	}
	mineFieldImg, err := screenshot.CaptureRect(mineFieldBounds)
	horizontalLines, verticalLines := imageproc.DetectMineSweeperGrid(mineFieldImg)
	for i := range 30 {

		log.Printf("=== 第 %d 轮迭代 ===", i+1)

		// 1. 截图阶段
		var total time.Duration
		start := time.Now()
		mineFieldImg, err := screenshot.CaptureRect(mineFieldBounds)
		if err != nil {
			panic(err)
		}
		mineFieldImgPos := imgpos.NewImageWithOffset(mineFieldImg, mineFieldBounds.Min)
		elapsed := time.Since(start)
		log.Printf("📸 截图耗时: %d ms", elapsed.Milliseconds())
		total += elapsed

		// 4. 图像保存阶段
		start = time.Now()
		go kit.SaveImg(mineFieldImg, "mineField.png")
		elapsed = time.Since(start)
		log.Printf("💾 保存耗时: %d ms", elapsed.Milliseconds())
		total += elapsed

		// 5. 雷区识别阶段
		start = time.Now()
		cells := identify.IdentifyMinesweeper(mineFieldImgPos, horizontalLines, verticalLines)
		fmt.Println(len(cells), "x", len(cells[0]))
		elapsed = time.Since(start)
		log.Printf("🧠 识别耗时: %d ms", elapsed.Milliseconds())
		total += elapsed

		// 6. 求解阶段
		start = time.Now()
		solver := solver.NewSolver(cells)
		safePoints, minePoints := solver.Solve()
		elapsed = time.Since(start)
		log.Printf("🧮 求解耗时: %d ms", elapsed.Milliseconds())
		total += elapsed

		// 7. 输出结果
		fmt.Println("✅ 安全点:", safePoints)
		fmt.Println("🚩 雷点:", minePoints)

		// 8. 点击操作阶段
		if len(safePoints) == 0 && len(minePoints) == 0 && i >= 6 {
			log.Printf("🛑 未检测到新操作，退出循环")
			break
		}

		start = time.Now()
		// 左键点击
		for _, point := range safePoints {
			p := cells[point.Y][point.X].ScreenPos()
			click.Click(p)
			time.Sleep(time.Millisecond * 20)
		}

		// 右键点击
		for _, point := range minePoints {
			p := cells[point.Y][point.X].ScreenPos()
			click.RightClick(p)
			time.Sleep(time.Millisecond * 20)
		}

		// 首次特殊点击
		p := cells[len(cells)/2][len(cells[0])/2].ScreenPos()
		click.Click(p)
		elapsed = time.Since(start)
		log.Printf("🖱️ 操作耗时: %d ms", elapsed.Milliseconds())
		total += elapsed

		log.Printf("📊 总耗时: %d ms", total.Milliseconds())
	}
}
