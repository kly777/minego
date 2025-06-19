package identify

import (
	"fmt"
	"image"

	"image/color"
	"minego/internal/cell"
	"minego/internal/imgpos"
	"minego/pkg/colorutil"
	"minego/pkg/imageproc"
	"os"
	"strconv"
)

type identifier struct {
	imgpos *imgpos.ImageWithOffset
}

func IdentifyMinesweeper(imgpos *imgpos.ImageWithOffset) [][]cell.GridCell {
	horizontalLines, verticalLines := imageproc.DetectMineSweeperGrid(imgpos.Image)
	rows := len(horizontalLines) - 1
	cols := len(verticalLines) - 1

	// 初始化二维切片
	result := make([][]cell.GridCell, rows)
	for i := range result {
		result[i] = make([]cell.GridCell, cols) // 初始化每行的列切片
		for j := range result[i] {
			y := (horizontalLines[i] + horizontalLines[i+1]) / 2
			x := (verticalLines[j] + verticalLines[j+1]) / 2
			width := (verticalLines[j+1] - verticalLines[j])
			hight := (horizontalLines[i+1] - horizontalLines[i])

			state := recognizeColor(imgpos.Image, x, y, width, hight)
			result[i][j] = cell.GridCell{
				Offset: imgpos.RelativePosition(),
				State:  state,
				X:      x,
				Y:      y,
				Width:  width,
				Hight:  hight,
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

func recognizeColor(img image.Image, x, y int, width, hight int) cell.CellState {
	if hasColor(img, x, y, 7, Number1FeatureColor) {
		return cell.Number1
	} else if hasColorWithinRange(img, x, y, 7, Number2FeatureColor, 5) {
		return cell.Number2
	} else if hasColorWithinRange(img, x, y, 7, Number3FeatureColor, 5) {
		return cell.Number3
	} else if diffColor(img, x, y, x, y+hight*2/6) < 30*256 {
		return cell.Empty
	}

	return cell.Unknown
}

// SaveResultToFile 将网格状态矩阵保存为文本文件
func SaveResultToFile(result [][]cell.GridCell, filePath string) error {
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
func cellStateToString(state cell.CellState) string {
	switch state {
	case cell.Mine:
		return "M"
	case cell.Flagged:
		return "F"
	case cell.Unknown:
		return "?"
	case cell.Empty:
		return "E"
	default:
		// 处理数字状态(1-8)
		if state >= cell.Number1 && state <= cell.Number8 {
			return strconv.Itoa(int(state - cell.Number1 + 1))
		}
		return "?"
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
