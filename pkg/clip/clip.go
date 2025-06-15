package clip

import (
	"fmt"
	"image"
	"image/color"
	"minego/pkg/colorutil"
)

// CropImage 裁剪图片（安全模式）
// 参数：
//
//	img: 源图像
//	rect: 裁剪区域（需在图像范围内）
//
// 返回：
//
//	裁剪后的图像 & 错误信息
func ClipImage(img *image.RGBA, rect image.Rectangle) (image.Image, error) {
	if !rect.In(img.Bounds()) {
		return nil, fmt.Errorf("裁剪区域 %v 超出图像范围 %v", rect, img.Bounds())
	}

	// 使用接口断言调用SubImage
	subImage := img.SubImage(rect)

	return subImage, nil
}

func ClipALine(img image.Image, y int, divColor color.Color) (int, error) {
	bounds := img.Bounds()

	// 校验 y 是否在合法范围内
	if y < bounds.Min.Y || y >= bounds.Max.Y {
		return 0, fmt.Errorf("y坐标 %d 超出图像范围 %v", y, bounds)
	}
	fmt.Println("y ", y)
	var first bool = true
	size := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		fmt.Println("x:", img.At(x, y))
		if colorutil.ColorsClose(img.At(x, y), divColor, 30*256) {
			if first {
				x += 5
				first = false
				size = x
			} else {
				size = x - size
				break
			}
		}
	}

	return size, nil
}
