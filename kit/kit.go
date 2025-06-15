package kit

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	about = 8*256
)

func SaveImg(img *image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, *img); err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}
	return nil
}

func FindSurroundingRect(img image.Image, targetColor color.Color) image.Rectangle {
	left := FindLeftmostColor(img, targetColor)
	right := FindRightmostColor(img, targetColor)
	top := FindTopmostColor(img, targetColor)
	bottom := FindBottommostColor(img, targetColor)
	p1 := image.Point{left.X, top.Y}
	p2 := image.Point{right.X, bottom.Y}
	fmt.Println("p1", p1, "p2", p2)
	rect := image.Rectangle{p1, p2}
	return rect
}

// FindLeftmostColor 查找最左侧的指定颜色位置
// 参数：
//
//	img: 源图像
//	targetColor: 要查找的目标颜色
//
// 返回：
//
//	最左侧坐标点（未找到时返回nil）
func FindLeftmostColor(img image.Image, targetColor color.Color) *image.Point {
	bounds := img.Bounds()

	// 按列优先遍历（从左到右）
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// 同一行内从上到下遍历
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if ColorsClose(img.At(x, y), targetColor, about) {
				return &image.Point{x, y}
			}
		}
	}

	return nil
}

// FindRightmostColor 查找最右侧的指定颜色位置
func FindRightmostColor(img image.Image, targetColor color.Color) *image.Point {
	bounds := img.Bounds()

	// 从右向左遍历
	for x := bounds.Max.X - 1; x >= bounds.Min.X; x-- {
		// 同一行内从上到下遍历
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			if ColorsClose(img.At(x, y), targetColor, about) {
				return &image.Point{x, y}
			}
		}
	}

	return nil
}

// FindTopmostColor 查找最顶部的指定颜色位置
func FindTopmostColor(img image.Image, targetColor color.Color) *image.Point {
	bounds := img.Bounds()

	// 从上到下遍历
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		// 同一行内从左到右遍历
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if ColorsClose(img.At(x, y), targetColor, about) {
				return &image.Point{x, y}
			}
		}
	}

	return nil
}

// FindBottommostColor 查找最底部的指定颜色位置
func FindBottommostColor(img image.Image, targetColor color.Color) *image.Point {
	bounds := img.Bounds()

	// 从下到上遍历
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		// 同一行内从左到右遍历
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if ColorsClose(img.At(x, y), targetColor, about) {
				return &image.Point{x, y}
			}
		}
	}

	return nil
}

// 颜色比较函数（考虑不同颜色模型的转换）
func ColorsEqual(c1, c2 color.Color) bool {
	// 转换为RGBA进行比较
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func ColorsClose(c1, c2 color.Color, length int) bool {
	return ColorsDist(c1, c2) < length
}

func ColorsDistance(c1, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return int(math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2))))
}

func ColorsDist(c1, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	dist := math.Abs(float64(r1-r2)) + math.Abs(float64(g1-g2)) + math.Abs(float64(b1-b2))

	return int(dist)
}
