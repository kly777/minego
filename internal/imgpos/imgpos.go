package imgpos

import "image"

// ImageWithOffset 表示带有偏移量的图像
type ImageWithOffset struct {
	Image  image.Image   // 图像对象
	Offset image.Point   // 图像在父坐标系中的偏移量
}

// RectWithOffset 表示带有偏移量的矩形区域
type RectWithOffset struct {
	Rect   image.Rectangle // 矩形区域
	Offset image.Point     // 矩形在父坐标系中的偏移量
}

// NewImageWithOffset 创建带偏移量的图像对象
func NewImageWithOffset(img image.Image, offset image.Point) *ImageWithOffset {
	return &ImageWithOffset{img, offset}
}

// NewRectWithOffset 创建带偏移量的矩形对象
func NewRectWithOffset(rect image.Rectangle, offset image.Point) *RectWithOffset {
	return &RectWithOffset{rect, offset}
}

// PositionCalculator 定义位置计算接口
type PositionCalculator interface {
	// RelativePosition 返回相对位置（偏移量）
	RelativePosition() image.Point
	// AbsolutePosition 返回绝对位置（屏幕坐标）
	AbsolutePosition() image.Point
}

// RelativePosition 实现PositionCalculator接口
func (img *ImageWithOffset) RelativePosition() image.Point {
	return img.Offset
}

// AbsolutePosition 计算图像在屏幕中的绝对位置
func (img *ImageWithOffset) AbsolutePosition() image.Point {
	bounds := img.Image.Bounds()
	return image.Point{
		X: img.Offset.X + bounds.Min.X,
		Y: img.Offset.Y + bounds.Min.Y,
	}
}

// RelativePosition 实现PositionCalculator接口
func (rect *RectWithOffset) RelativePosition() image.Point {
	return rect.Offset
}

// AbsolutePosition 计算矩形在屏幕中的绝对位置
func (rect *RectWithOffset) AbsolutePosition() image.Point {
	return image.Point{
		X: rect.Offset.X + rect.Rect.Min.X,
		Y: rect.Offset.Y + rect.Rect.Min.Y,
	}
}
