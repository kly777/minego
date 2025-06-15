package main

import (
	"fmt"

	"image/color"

	"minego/clip"
	"minego/identify"

	"minego/imgP"
	"minego/kit"
	"minego/screenshot"
)

var (
	BorderColor = color.RGBA{3, 0, 6, 255}
)

func main() {
	img, err := screenshot.ShotMineSweeper()
	if err != nil {
		panic(err)
	}
	fmt.Println(img.Bounds())

	rect := kit.FindSurroundingRect(img, BorderColor)
	rect.Min.X -= 1
	rect.Min.Y -= 1
	rect.Max.X += 1
	rect.Max.Y += 1
	fmt.Println(rect)
	img2, err := clip.ClipImage(img, rect)

	if err != nil {
		panic(err)
	}
	kit.SaveImg(&img2, "clip.png")
	rows, cols := imgP.DetectMineGrid(img2)
	fmt.Println(rows, cols)
	size := identify.MineSize{Cols: cols, Rows: rows}
	identify.RecognizeMinesweeper(img2, size)
}
