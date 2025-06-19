package main

import (
	"fmt"
	"image"

	"image/color"
	"time"

	"minego/pkg/clip"
	"minego/pkg/identify"
	"minego/pkg/imageproc"
	"minego/pkg/kit"
	"minego/pkg/winapi/click"

	"minego/pkg/winapi"

	"minego/pkg/screenshot"
)

var (
	BorderColor = color.RGBA{3, 0, 6, 255}
)

// main 函数是程序的入口点，用于执行扫雷游戏识别任务
// 主要流程包括：截图、定位扫雷区域、裁剪图像、保存中间结果、网格检测和雷区识别
func main() {
	click.SetDPIAware()

	mineSweeperWindow := winapi.GetMineSweeperWindow()
	mineSweeperWindow.Activate()
	time.Sleep(20 * time.Millisecond)
	// 获取窗口位置
	windowBounds, err := mineSweeperWindow.GetBounds()
	if err != nil {
		panic(err)
	}
	// 向内缩10像素，排除非扫雷窗口
	windowBounds.Min.X += 10
	windowBounds.Min.Y += 10
	windowBounds.Max.X -= 10
	windowBounds.Max.Y -= 10
	fmt.Printf("扫雷窗口截图区域: %v\n", windowBounds)

	// 截取扫雷窗口
	windowImg, err := screenshot.CaptureRect(windowBounds)
	if err != nil {
		panic(err)
	}
	fmt.Println("图像有效范围",windowImg.Bounds())

	// 查找扫雷雷区的边界矩形，并进行1像素扩展
	rect := kit.FindSurroundingRect(windowImg, BorderColor)
	rect.Min.X -= 3
	rect.Min.Y -= 3
	rect.Max.X += 3
	rect.Max.Y += 3
	fmt.Println("边界矩形", rect)

	// 根据边界矩形裁剪图像并保存雷区
	gridImg, err := clip.ClipImage(windowImg, rect)
	if err != nil {
		panic(err)
	}
	err = kit.SaveImg(gridImg, "clip.png")
	if err != nil {
		panic(err)
	}
	fmt.Println(gridImg.Bounds())
	// 步骤4: 检测裁剪后图像中的扫雷网格行列数
	rows, cols := imageproc.DetectMineSweeperGridNum(gridImg) // 更新函数调用

	fmt.Println(rows, cols)

	horizontalLines, verticalLines := imageproc.DetectMineSweeperGrid(gridImg)
	fmt.Println(horizontalLines, verticalLines)
	fmt.Println(gridImg.Bounds())

	cells := identify.IdentifyMinesweeper(gridImg, horizontalLines, verticalLines)

	fmt.Printf("截图区域: %v\n", windowBounds)
	fmt.Println(windowImg.Bounds())
	fmt.Println(gridImg.Bounds())
	x := 4
	y := 4

	screenX, screenY := cellToScreenPos(cells[x][y], windowBounds, rect)
	fmt.Println("state", cells[x][y].State)
	fmt.Println("点击", screenX, screenY)
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
