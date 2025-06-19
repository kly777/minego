package cell

import "image"

type Mouse interface {
	Click(row, col int) error
	RightClick(row, col int) error
	DoubleLeftClick(row, col int) error
}

type Identifier interface {
	GetMineField() Minefield
}

type Game struct {
	Mouse      Mouse
	Identifier Identifier
}

type Minefield struct {
	Grid      [][]Cell
	MineCount int
	Bounds    image.Rectangle
}

func NewMinefield(bounds image.Rectangle, rows, cols int) *Minefield {
	return &Minefield{
		Grid:   make([][]Cell, rows),
		Bounds: bounds,
	}
}

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

type Cell struct {
	Pos   image.Point
	State CellState
	mf    MinefieldInterface
}

type MinefieldInterface interface {
	GetBounds() image.Rectangle
	GetOffset() image.Point
	GetCell(image.Point) *Cell
	GetSurroundCell(Cell) []Cell
}

func (mf *Minefield) GetBounds() image.Rectangle {
	return mf.Bounds
}

func (mf *Minefield) GetOffset() image.Point {
	return mf.Bounds.Min
}

func (mf *Minefield) GetCell(p image.Point) *Cell {
	return &mf.Grid[p.Y][p.X]
}
