package screenshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"minego/pkg/winapi"

	"github.com/kbinani/screenshot"
)

// 保存截图到指定路径
func saveScreenshot(img *image.RGBA, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}
	return nil
}

func ShotMineSweeper() (*image.RGBA, error) {
	hwnd, err := winapi.FindMineWindow()
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	fmt.Println("窗口句柄:", hwnd)

	// 激活窗口
	if err := winapi.ActivateWindow(hwnd); err != nil {
		fmt.Println("警告:", err)
	}
	time.Sleep(time.Second / 10)

	// 获取窗口区域
	bounds, err := winapi.GetWindowBounds(hwnd)
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	fmt.Printf("截图区域: %v\n", bounds)

	// 截图
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("截图失败:", err)
		return nil, err
	}

	// 保存图片
	outputPath := "minesweeper.png"
	if err := saveScreenshot(img, outputPath); err != nil {
		fmt.Println("保存失败:", err)
	} else {
		fmt.Printf("截图已保存至: %s\n", outputPath)
	}
	return img, nil
}
