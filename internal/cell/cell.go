package cell

import (
	"image"
	"image/color"
)

type CellState int

const (
	Mine CellState = iota - 4
	Flagged
	Unknown
	Locked
	Empty
	Number1
	Number2
	Number3
	Number4
	Number5
	Number6
	Number7
	Number8
)

type MineField struct {
	Bounds image.Rectangle
	Grid   [][]GridCell
}

type GridCell struct {
	Offset       image.Point
	State        CellState
	X, Y         int // 坐标位置
	Width, Hight int
	Position     image.Point
	Color        color.Color
}

func NewMineField(bounds image.Rectangle, cells [][]GridCell) *MineField {
	return &MineField{
		Grid:   cells,
		Bounds: bounds,
	}
}

func (gc *GridCell) ScreenPos() image.Point {
	screenX := gc.Offset.X + gc.X // 补偿窗口边框
	screenY := gc.Offset.Y + gc.Y
	return image.Point{
		X: screenX,
		Y: screenY,
	}
}
