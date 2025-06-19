package imgpos

import "image"

type Imgpos struct {
	Image image.Image
	position image.Point
}

type Rectpos struct {
	Rect image.Rectangle
	position image.Point
}

func NewImgPos(img image.Image, pos image.Point) *Imgpos {
	return &Imgpos{img, pos}
}

func NewRectPos(rect image.Rectangle, pos image.Point) *Rectpos {
	return &Rectpos{rect, pos}
}
func (rp *Rectpos) Position() image.Point {
	return rp.position
}

func (ip *Imgpos) Position() image.Point {
	return ip.position
}

func (ip *Imgpos) AsPosition() image.Point {
	x:=ip.position.X+ip.Image.Bounds().Min.X
	y:=ip.position.Y+ip.Image.Bounds().Min.Y
	return image.Point{
		X:x,
		Y:y,
	}
}

func (ip *Rectpos) AsPosition() image.Point {
	x:=ip.position.X+ip.Rect.Min.X
	y:=ip.position.Y+ip.Rect.Min.Y
	return image.Point{
		X:x,
		Y:y,
	}
}