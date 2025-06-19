package identify

import (
	"fmt"
	"image"
	"log"

	"image/color"
	"minego/pkg/colorutil"
	"os"
	"strconv"
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

type MineSize struct {
	Cols, Rows int
}

type GridCell struct {
	State        CellState
	X, Y         int // 坐标位置
	Width, Hight int
}

func IdentifyMinesweeper(img image.Image, horizontalLines, verticalLines []int) [][]GridCell {

	rows := len(horizontalLines) - 1
	cols := len(verticalLines) - 1
	log.Println("aaa", rows, cols)
	log.Println("bbb", img.Bounds())

	// 初始化二维切片
	result := make([][]GridCell, rows)
	for i := range result {
		result[i] = make([]GridCell, cols) // 初始化每行的列切片
		for j := range result[i] {
			y := (horizontalLines[i] + horizontalLines[i+1]) / 2
			x := (verticalLines[j] + verticalLines[j+1]) / 2
			width := (verticalLines[j+1] - verticalLines[j])
			hight := (horizontalLines[i+1] - horizontalLines[i])
			// fmt.Println(img.Bounds().Min.X+x, img.Bounds().Min.Y+y)

			state := recognizeColor(img, x, y, width, hight)
			result[i][j] = GridCell{
				State: state,
				X:     x,
				Y:     y,
				Width: width,
				Hight: hight,
			}
		}
	}

	SaveResultToFile(result, "GridcellRec.txt")
	return result
}

var (
	Number1FeatureColor = color.RGBA{65, 79, 188, 255}
	Number2FeatureColor = color.RGBA{30, 105, 3, 255}
	Number3FeatureColor = color.RGBA{175, 5, 8, 255}
)

func recognizeColor(img image.Image, x, y int, width, hight int) CellState {
	if hasColor(img, x, y, 7, Number1FeatureColor) {
		return Number1
	} else if hasColorWithinRange(img, x, y, 7, Number2FeatureColor, 5) {
		return Number2
	} else if hasColorWithinRange(img, x, y, 7, Number3FeatureColor, 5) {
		return Number3
	} else if diffColor(img, x, y, x, y+hight*2/6) < 30*256 {
		return Empty
	}

	// mixC := mixColor(img, x, y, 7)
	// // 实现具体颜色匹配逻辑
	// // 此处需要根据实际截图的颜色值进行调整
	// if colorutil.ColorsCloseN(mixC, Number1Color, 60) {
	// 	return Number1
	// } else if colorutil.ColorsClose(mixC, Number2Color, 10*256) {
	// 	return Number2
	// } else if colorutil.ColorsClose(mixC, Number3Color, 10*256) {
	// 	return Number3
	// } else if colorutil.ColorsCloseN(mixC, EmptyColor, 60) {
	// 	return Empty
	// }

	return Unknown
}

// SaveResultToFile 将网格状态矩阵保存为文本文件
func SaveResultToFile(result [][]GridCell, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	for _, row := range result {
		var line string
		for _, cell := range row {
			// 将CellState转换为字符串表示
			line += cellStateToString(cell.State) + " "
		}
		// 去除末尾空格并写入文件
		_, err = file.WriteString(line[:len(line)-1] + "\n")
		if err != nil {
			return fmt.Errorf("写入文件失败: %v", err)
		}
	}
	return nil
}

// CellState字符串映射
func cellStateToString(state CellState) string {
	switch state {
	case Mine:
		return "M"
	case Flagged:
		return "F"
	case Unknown:
		return "?"
	case Empty:
		return "E"
	default:
		// 处理数字状态(1-8)
		if state >= Number1 && state <= Number8 {
			return strconv.Itoa(int(state - Number1 + 1))
		}
		return "?"
	}
}

func mixColor(img image.Image, x, y int, rang int) color.Color {
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y

	// 计算实际采样区域
	width := 2*rang + 1
	totalPixels := width * width

	var rSum, gSum, bSum uint32
	for i := -rang; i <= rang; i++ {
		for j := -rang; j <= rang; j++ {
			// 获取当前像素颜色（16-bit值）
			nr, ng, nb, _ := img.At(minX+x+i, minY+y+j).RGBA()
			rSum += nr
			gSum += ng
			bSum += nb
		}
	}

	// 计算平均值并转换为8-bit
	avgR := rSum / uint32(totalPixels) >> 8 // 右移8位转换到0-255
	avgG := gSum / uint32(totalPixels) >> 8
	avgB := bSum / uint32(totalPixels) >> 8

	return color.RGBA{
		R: uint8(avgR),
		G: uint8(avgG),
		B: uint8(avgB),
		A: 255,
	}
}

func hasColor(img image.Image, x, y int, rang int, targetColor color.Color) bool {
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y
	tr, tg, tb, _ := targetColor.RGBA()
	for i := -rang; i <= rang; i++ {
		for j := -rang; j <= rang; j++ {
			// 获取当前像素颜色（16-bit值）
			nr, ng, nb, _ := img.At(minX+x+i, minY+y+j).RGBA()
			if nr == tr && ng == tg && nb == tb {
				return true
			}
		}
	}
	return false
}

func hasColorWithinRange(img image.Image, x, y int, rang int, targetColor color.Color, colorRange int) bool {
	minX := img.Bounds().Min.X
	minY := img.Bounds().Min.Y

	for i := -rang; i <= rang; i++ {
		for j := -rang; j <= rang; j++ {
			// 获取当前像素颜色（16-bit值）

			if colorutil.ColorsCloseN(img.At(minX+x+i, minY+y+j), targetColor, colorRange) {
				return true
			}
		}
	}
	return false
}

func diffColor(img image.Image, x, y int, x2, y2 int) int {
	return colorutil.ColorsDist(img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y), img.At(img.Bounds().Min.X+x2, img.Bounds().Min.Y+y2))
}
