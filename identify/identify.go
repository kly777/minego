package identify

import (
	"fmt"
	"image"
	"image/color"
	"minego/kit"
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
)

type MineSize struct {
	Width, Height int
}

type GridCell struct {
	State CellState
	X, Y  int // 坐标位置
}

// RecognizeMinesweeper 识别扫雷游戏状态
// 参数 gridImage: 扫雷游戏的截图
// 参数 mineSize: 扫雷格子行列数
// 返回值: 二维网格状态矩阵
func RecognizeMinesweeper(gridImage image.Image, mineSize MineSize) [][]GridCell {
	bounds := gridImage.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 校验 GridSize 合理性
	if mineSize.Width <= 0 || mineSize.Height <= 0 {
		panic("GridSize 的 Width 和 Height 必须大于 0")
	}

	// 正确计算行列数
	rows := height / mineSize.Height
	cols := width / mineSize.Width

	result := make([][]GridCell, rows)
	for i := range result {
		result[i] = make([]GridCell, cols)
	}

	// 使用标准循环结构
	for row := range rows {
		for col := range cols {
			// 基于 GridSize 的宽高分别计算中心点
			x := col*mineSize.Width + mineSize.Width/2
			y := row*mineSize.Height + mineSize.Height/2

			// 边界检查
			if y >= height || x >= width {
				result[row][col].State = Unknown
				continue
			}

			c := gridImage.At(x, y)
			state := recognizeColor(c)

			result[row][col] = GridCell{
				State: state,
				X:     x,
				Y:     y,
			}
		}
	}

	err := SaveResultToFile(result, "output.txt")
	if err != nil {
		fmt.Print(err)
	}

	return result
}

// recognizeColor 将颜色转换为对应状态
func recognizeColor(c color.Color) CellState {
	// 实现具体颜色匹配逻辑
	// 此处需要根据实际截图的颜色值进行调整
	if kit.ColorsClose(c, Number1Color, 10*256) {
		return Number1
	} else if kit.ColorsClose(c, Number2Color, 10*256) {
		return Number2
	} else if kit.ColorsClose(c, Number3Color, 10*256) {
		return Number3
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
