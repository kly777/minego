package main

import (
	"fmt"
	"image/png"
	"os"
	"github.com/kbinani/screenshot"
)


func main() {
	// 获取显示器数量
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		panic("未找到可用显示器")
	}

	// 捕获主显示器（显示器索引从0开始）
	img, err := screenshot.CaptureDisplay(0)
	if err != nil {
		panic(fmt.Sprintf("截图失败: %v", err))
	}

	// 创建截图文件
	file, err := os.Create("screenshot.png")
	if err != nil {
		panic(fmt.Sprintf("文件创建失败: %v", err))
	}
	defer file.Close()

	// 保存为PNG格式
	if err := png.Encode(file, img); err != nil {
		panic(fmt.Sprintf("编码失败: %v", err))
	}

	fmt.Println("截图成功保存为 screenshot.png")
}
