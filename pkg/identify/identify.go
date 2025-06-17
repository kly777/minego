package identify

import (
	"fmt"
	"image"

	"image/color"
	"minego/pkg/colorutil"
	"os"
	"strconv"
)

type CellState int

const (
	Mine CellState = iota - 3
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

var (
	UnknownColor = color.RGBA{88, 118, 220, 255}
	Number1Color = color.RGBA{64, 80, 189, 255}
	Number2Color = color.RGBA{29, 103, 5, 255}
	Number3Color = color.RGBA{175, 31, 34, 255}
	EmptyColor   = color.RGBA{198, 209, 231, 255}
)

type MineSize struct {
	Cols, Rows int
}

type GridCell struct {
	State CellState
	X, Y  int // 坐标位置
}

func IdentifyMinesweeper(img image.Image, horizontalLines, verticalLines []int) [][]GridCell {

	rows := len(horizontalLines) - 1
	cols := len(verticalLines) - 1

	// 初始化二维切片
	result := make([][]GridCell, rows)
	for i := range result {
		result[i] = make([]GridCell, cols) // 初始化每行的列切片
		for j := range result[i] {
			x := (horizontalLines[i] + horizontalLines[i+1]) / 2
			y := (verticalLines[j] + verticalLines[j+1]) / 2
			state := recognizeColor(img.At(x, y))
			result[i][j] = GridCell{
				State: state,
				X:     x,
				Y:     y,
			}
		}
	}

	SaveResultToFile(result, "GridcellRec.txt")
	return result
}

func recognizeColor(c color.Color) CellState {
	// 实现具体颜色匹配逻辑
	// 此处需要根据实际截图的颜色值进行调整
	fmt.Println(c)
	if colorutil.ColorsClose(c, Number1Color, 10*256) {
		return Number1
	} else if colorutil.ColorsClose(c, Number2Color, 10*256) {
		return Number2
	} else if colorutil.ColorsClose(c, Number3Color, 10*256) {
		return Number3
	} else if colorutil.ColorsClose(c, EmptyColor, 10*256) {
		return Empty
	}

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
