package imageproc

import (
	"bufio"
	"encoding/binary"
	"log"
	"image"

	_ "image/jpeg"
	_ "image/png"

	"os"
	"sort"
)

func DetectMineSweeperGridNum(img image.Image) (gridRows int, gridCols int) {
	// 计算格子数（行数和列数）
	horizontalLines, verticalLines := DetectMineSweeperGrid(img)
	gridRows = len(horizontalLines) - 1
	gridCols = len(verticalLines) - 1
	return gridRows, gridCols
}

// 检测扫雷网格,由始图像中识别出扫雷格子数
func DetectMineSweeperGrid(img image.Image) ([]int, []int) {
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Dx(), bounds.Dy()
	log.Println("图像尺寸:", imgWidth, "x", imgHeight)
	// 转换为灰度图
	grayImg := toGrayScale(img)
	log.Println("图像灰度处理完成")

	// 二值化处理
	binaryImg := binarize(grayImg, 60)
	if err := saveDebugImage(binaryImg, "debug_output.bmp"); err != nil {
		log.Printf("保存调试图像失败: %v\n", err)
	}

	// 检测水平和垂直线
	horizontalLines := detectHorizontalLines(binaryImg, imgWidth, imgHeight)
	verticalLines := detectVerticalLines(binaryImg, imgWidth, imgHeight)
	// fmt.Println("检测到水平线:", horizontal, "列线:", vertical)
	if len(horizontalLines) == 0 || len(verticalLines) == 0 {
		return []int{}, []int{}
	}

	// 对坐标排序
	sort.Ints(horizontalLines)
	sort.Ints(verticalLines)

	// 聚类和去重（合并相近的线）
	horizontalLines = clusterPoints(horizontalLines, 5)
	verticalLines = clusterPoints(verticalLines, 5)

	return horizontalLines, verticalLines
}

// 转换为灰度图
func toGrayScale(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	gray := make([][]uint8, height)

	for y := range height {
		gray[y] = make([]uint8, width)
		for x := range width {
			r, g, b, _ := img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y).RGBA()
			// 转换为灰度值 (0.299*R + 0.587*G + 0.114*B)
			grayValue := uint8((0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)))
			gray[y][x] = grayValue
		}
	}

	return gray
}

// 二值化处理
func binarize(gray [][]uint8, threshold uint8) [][]uint8 {
	height := len(gray)
	width := len(gray[0])
	binary := make([][]uint8, height)

	for y := range height {
		binary[y] = make([]uint8, width)
		for x := range width {
			if gray[y][x] > threshold {
				binary[y][x] = 255 // 白色
			} else {
				binary[y][x] = 0 // 黑色
			}
		}
	}

	return binary
}

// 检测水平线
func detectHorizontalLines(binary [][]uint8, width, height int) []int {
	lines := make([]int, 0, height) // 预分配容量
	if width <= 0 || height <= 0 {
		return lines // 处理无效输入
	}
	requiredLength := width / 2

	for y := range height {
		lineLength := 0
		maxLineLength := 0

		for x := range width {
			if binary[y][x] == 0 {
				lineLength++
			} else {
				maxLineLength = max(maxLineLength, lineLength)
				lineLength = 0
			}
		}
		maxLineLength = max(maxLineLength, lineLength) // 合并末尾处理

		if maxLineLength >= requiredLength {
			lines = append(lines, y)
		}
	}

	return lines
}

// 检测垂直线
func detectVerticalLines(binary [][]uint8, width, height int) []int {
	lines := make([]int, 0, width) // 预分配容量
	if width <= 0 || height <= 0 {
		return lines // 处理无效输入
	}
	requiredLength := height / 2

	for x := 0; x < width; x++ {
		lineLength := 0
		maxLineLength := 0

		for y := range height {
			if binary[y][x] == 0 { // 黑色像素表示线
				lineLength++
			} else {
				maxLineLength = max(maxLineLength, lineLength)
				lineLength = 0
			}
		}
		maxLineLength = max(maxLineLength, lineLength) // 合并末尾处理

		if maxLineLength >= requiredLength {
			lines = append(lines, x)
		}
	}

	return lines
}

// 聚类点 - 合并接近的点
func clusterPoints(points []int, tolerance int) []int {
	if len(points) == 0 {
		return points
	}

	// 排序
	sort.Ints(points)

	clustered := []int{points[0]}
	current := points[0]

	for _, p := range points[1:] {
		if p-current > tolerance {
			clustered = append(clustered, p)
			current = p
		}
	}

	return clustered
}

// 保存调试图像（可选）
func saveDebugImage(bin [][]uint8, path string) error {
	height := len(bin)
	if height == 0 {
		return nil
	}
	width := len(bin[0])

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入BMP文件头
	headerSize := 14 + 40 + 256*4 // 文件头+信息头+调色板
	fileSize := headerSize + width*height
	padding := (4 - (width % 4)) % 4

	// BMP文件头
	file.Write([]byte{'B', 'M'})                                  // 签名
	binary.Write(writer, binary.LittleEndian, uint32(fileSize))   // 文件大小
	binary.Write(writer, binary.LittleEndian, uint32(0))          // 保留
	binary.Write(writer, binary.LittleEndian, uint32(headerSize)) // 像素数据偏移

	// BMP信息头
	binary.Write(writer, binary.LittleEndian, uint32(40))    // 信息头大小
	binary.Write(writer, binary.LittleEndian, int32(width))  // 宽度
	binary.Write(writer, binary.LittleEndian, int32(height)) // 高度
	binary.Write(writer, binary.LittleEndian, uint16(1))     // 颜色平面数
	binary.Write(writer, binary.LittleEndian, uint16(8))     // 每像素位数
	binary.Write(writer, binary.LittleEndian, uint32(0))     // 压缩方式
	binary.Write(writer, binary.LittleEndian, uint32(0))     // 图像大小
	binary.Write(writer, binary.LittleEndian, int32(0))      // 水平分辨率
	binary.Write(writer, binary.LittleEndian, int32(0))      // 垂直分辨率
	binary.Write(writer, binary.LittleEndian, uint32(0))     // 调色板颜色数
	binary.Write(writer, binary.LittleEndian, uint32(0))     // 重要颜色数

	// 调色板（灰度）
	for i := range 256 {
		writer.Write([]byte{byte(i), byte(i), byte(i), 0})
	}

	// 像素数据（从底部向上）
	for y := height - 1; y >= 0; y-- {
		for x := range width {
			writer.WriteByte(bin[y][x])
		}
		// 行填充
		for range padding {
			writer.WriteByte(0)
		}
	}

	return nil
}
