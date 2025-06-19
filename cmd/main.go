package main

import (
	"image"
	"log"

	"image/color"
	"time"

	"minego/pkg/clip"
	"minego/internal/identify"
	"minego/pkg/imageproc"
	"minego/pkg/kit"
	"minego/pkg/winapi/click"
	"minego/internal/window"

	"minego/pkg/screenshot"
)

var (
	BorderColor = color.RGBA{3, 0, 6, 255}
)

// main 函数是程序的入口点，用于执行扫雷游戏识别任务
// 主要流程包括：截图、定位扫雷区域、裁剪图像、保存中间结果、网格检测和雷区识别
const (
	windowBorderInset = 10 // 窗口边界内缩像素
	gridBorderExpand  = 3  // 雷区边界扩展像素
)

func main() {
	click.SetDPIAware()

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
	if err != nil {
		log.Fatalf("窗口截图失败: %v", err)
	}

	mineField := kit.FindSurroundingRect(windowImg, BorderColor)
	// 添加边界保护
	mineField.Min.X = max(mineField.Min.X-gridBorderExpand, 0)
	mineField.Min.Y = max(mineField.Min.Y-gridBorderExpand, 0)
	mineField.Max.X = min(mineField.Max.X+gridBorderExpand, windowImg.Bounds().Dx())
	mineField.Max.Y = min(mineField.Max.Y+gridBorderExpand, windowImg.Bounds().Dy())

	mineFieldImg, err := clip.ClipImage(windowImg, mineField)
	if err != nil {
		log.Fatalf("图像裁剪失败: %v", err)
	}

	if err := kit.SaveImg(mineFieldImg, "clip.png"); err != nil {
		log.Fatalf("保存图像失败: %v", err)
	}

	horizontalLines, verticalLines := imageproc.DetectMineSweeperGrid(mineFieldImg)
	cells := identify.IdentifyMinesweeper(mineFieldImg, horizontalLines, verticalLines)

	x, y := 4, 4
	screenX, screenY := cellToScreenPos(cells[x][y], windowBounds, mineField)
	click.PhysicalMouseClick(int32(screenX), int32(screenY))
}

// 在调用鼠标点击前转换为相对窗口坐标
func cellToScreenPos(cell identify.GridCell, bounds image.Rectangle, rect image.Rectangle) (int, int) {
	// 计算相对于窗口的坐标 = 裁剪区域偏移 + 格子中心偏移
	relX := rect.Min.X + cell.X
	relY := rect.Min.Y + cell.Y

	// 加上窗口位置和边框补偿
	screenX := bounds.Min.X + relX // 补偿窗口边框
	screenY := bounds.Min.Y + relY

	return screenX, screenY
}
