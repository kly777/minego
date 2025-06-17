package tile

import "image"

type Tile struct {
	Position image.Point
	State    TileState
}

type TileState int

const (
	Mine TileState = iota - 3
	Flagged
	Unknown
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