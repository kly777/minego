// pkg/colorutil/color.go
package colorutil

import (
	"image/color"
	"math"
)

func ColorsCloseN(c1, c2 color.Color, length int) bool {
	return ColorsClose(c1, c2, length*256)
}

func ColorsClose(c1, c2 color.Color, length int) bool {
	return ColorsDist(c1, c2) < length
}

func ColorsDist(c1, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return int(math.Abs(float64(r1-r2)) +
		math.Abs(float64(g1-g2)) +
		math.Abs(float64(b1-b2)))
}

// 颜色比较函数（考虑不同颜色模型的转换）
func ColorsEqual(c1, c2 color.Color) bool {
	// 转换为RGBA进行比较
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func ColorsDistance(c1, c2 color.Color) int {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return int(math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2))))
}
