package screenshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
)

// Windows API常量
const (
	SW_RESTORE       = 9
	SW_SHOW          = 5
	SW_SHOWMAXIMIZED = 3
)

var (
	success             = "The operation completed successfully."
	user32              = syscall.NewLazyDLL("user32.dll")
	findWindow          = user32.NewProc("FindWindowW")
	setForegroundWindow = user32.NewProc("SetForegroundWindow")
	showWindow          = user32.NewProc("ShowWindow")
	getWindowRect       = user32.NewProc("GetWindowRect")
)

// 查找窗口句柄
func findMineWindow() (uintptr, error) {
	className := "Minesweeper" // 扫雷窗口类名为空
	windowName := "扫雷"         // 默认窗口标题

	classPtr, err := syscall.UTF16PtrFromString(className)
	if err != nil {
		return 0, fmt.Errorf("创建字符串指针失败: %v", err)
	}
	windowPtr, err := syscall.UTF16PtrFromString(windowName)
	if err != nil {
		return 0, fmt.Errorf("创建字符串指针失败: %v", err)
	}
	hwnd, _, err := findWindow.Call(
		uintptr(unsafe.Pointer(classPtr)),
		uintptr(unsafe.Pointer(windowPtr)),
	)
	if err != nil && err.Error() != success {
		return 0, fmt.Errorf("调用 FindWindow 失败: %v", err)
	}

	if hwnd == 0 {
		return 0, fmt.Errorf("未找到扫雷窗口")
	}
	return hwnd, nil
}

// 激活窗口并还原显示
func activateWindow(hwnd uintptr) error {
	// 先恢复窗口状态
	_, _, err := showWindow.Call(hwnd, SW_RESTORE)
	if err.Error() != success {
		return fmt.Errorf("恢复窗口状态失败")
	}
	// 激活窗口
	ret, _, err := setForegroundWindow.Call(hwnd)
	if ret == 0 && err.Error() != success {
		return fmt.Errorf("激活窗口失败: %v", err)
	}
	return nil
}

// 获取窗口位置和尺寸
func getWindowBounds(hwnd uintptr) (image.Rectangle, error) {
	var rect [4]int32
	_, _, err := getWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect[0])))
	if err.Error() != "The operation completed successfully." {
		return image.Rectangle{}, fmt.Errorf("获取窗口坐标失败: %v", err)
	}
	fmt.Println(rect)

	left, top := logicalToPhysical(hwnd, int(rect[0]), int(rect[1]))
	right, bottom := logicalToPhysical(hwnd, int(rect[2]), int(rect[3]))
	return image.Rect(left, top, right, bottom), nil
}

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

	hwnd, err := findMineWindow()
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	fmt.Println("窗口句柄:", hwnd)

	// 2. 激活窗口
	if err := activateWindow(hwnd); err != nil {
		fmt.Println("警告:", err)
	}
	time.Sleep(time.Second / 10)

	// 3. 获取窗口区域
	bounds, err := getWindowBounds(hwnd)
	if err != nil {
		fmt.Println("错误:", err)
		return nil, err
	}
	fmt.Printf("截图区域: %v\n", bounds)

	// 4. 截图
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("截图失败:", err)
		return nil, err
	}

	// 5. 保存图片
	outputPath := "minesweeper.png"
	if err := saveScreenshot(img, outputPath); err != nil {
		fmt.Println("保存失败:", err)
	} else {
		fmt.Printf("截图已保存至: %s\n", outputPath)
	}
	return img, nil
}

// 获取窗口 DPI 缩放比例
func getDPI(hwnd uintptr) (int, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getDpiForWindow := user32.NewProc("GetDpiForWindow")
	if getDpiForWindow.Find() != nil {
		// Windows 8.1-10 低版本回退方案
		hdc, _, _ := syscall.Syscall(user32.NewProc("GetDC").Addr(), 1, hwnd, 0, 0)
		if hdc == 0 {
			return 96, nil
		}
		defer syscall.SyscallN(user32.NewProc("ReleaseDC").Addr(), 2, hwnd, hdc, 0)
		logPixelsX, _, _ := syscall.Syscall(syscall.NewLazyDLL("gdi32.dll").NewProc("GetDeviceCaps").Addr(), 2, hdc, 88, 0)
		return int(logPixelsX), nil
	}
	dpi, _, _ := getDpiForWindow.Call(hwnd)
	if dpi == 0 {
		return 96, nil
	}
	return int(dpi), nil
}

// 将逻辑坐标转换为物理坐标
func logicalToPhysical(hwnd uintptr, x, y int) (int, int) {
	dpi, _ := getDPI(hwnd)
	if dpi == 96 {
		return x, y
	}
	scale := float64(dpi) / 96.0
	return int(float64(x) * scale), int(float64(y) * scale)
}
