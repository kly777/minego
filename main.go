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
	fmt.Println(rect)
	img2, err := clip.ClipImage(img, rect)

	if err != nil {
		panic(err)
	}
	kit.SaveImg(&img2, "clip.png")
	x, y := imgP.DetectMineGrid(img2)
	size := identify.MineSize{Width: x, Height: y}
	identify.RecognizeMinesweeper(img2, size)
}
