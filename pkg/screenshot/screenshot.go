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

// saveRGBA 将RGBA图像保存为PNG文件到指定路径
// 示例：
//
//	err := saveRGBA(img, "screenshot.png")
//	if err != nil {
//	    log.Fatal(err)
//	}
func saveRGBA(img *image.RGBA, path string) error {
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

// ShotMineSweeper 捕获扫雷游戏窗口的截图并进行处理
//
// 函数执行流程：
//  1. 查找扫雷游戏窗口句柄
//  2. 激活目标窗口确保其在前台显示
//  3. 获取窗口边界并计算有效截图区域（去除边框）
//  4. 执行屏幕截图操作并记录耗时
//  5. 将截图保存到本地文件
//
// 返回值:
//
//	*image.RGBA - 截取的图像数据
//	error      - 执行过程中遇到的错误
func ShotMineSweeper() (*image.RGBA, error) {
	// 获取扫雷窗口句柄
	hwnd, err := winapi.FindMineWindow()
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	fmt.Println("窗口句柄:", hwnd)

	// 激活目标窗口
	if err := winapi.ActivateWindow(hwnd); err != nil {
		fmt.Println("警告:", err)
	}
	time.Sleep(20 * time.Millisecond)

	// 计算有效截图区域（去除窗口边框）
	bounds, err := winapi.GetWindowBounds(hwnd)
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	// TODO: 封装为一个函数
	neoBounds := image.Rect(
		bounds.Min.X+10,
		bounds.Min.Y+10,
		bounds.Max.X-10,
		bounds.Max.Y-10,
	)
	fmt.Printf("截图区域: %v\n", neoBounds)

	// 执行截图并测量耗时

	img, err := screenshot.CaptureRect(neoBounds)
	if err != nil {
		fmt.Println("截图失败:", err)
		return nil, err
	}

	// 保存截图到本地文件
	outputPath := "minesweeper.png"
	if err := saveRGBA(img, outputPath); err != nil {
		fmt.Println("保存失败:", err)
	} else {
		fmt.Printf("截图已保存至: %s\n", outputPath)
	}
	return img, nil
}
