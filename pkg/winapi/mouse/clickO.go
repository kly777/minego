package mouse

import (
	"fmt"
	"time"
	"unsafe"
)

// Windows API 函数声明
var (
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procMouseEvent   = user32.NewProc("mouse_event")
)

// Point 结构体表示屏幕坐标
type Point struct {
	X int32
	Y int32
}

// ClickAt 在指定坐标执行鼠标点击
func ClickAtO(x, y int32) {
	// 获取屏幕分辨率
	cxScreen, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)
	cyScreen, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)
	fmt.Println(cxScreen, cyScreen)
	// 转换为绝对坐标 (0-65535)
	absX := x * 65535 / int32(cxScreen)
	absY := y * 65535 / int32(cyScreen)
	fmt.Println(absX, absY)

	// 移动鼠标
	_, _, _ = procMouseEvent.Call(
		uintptr(MOUSEEVENTF_ABSOLUTE|MOUSEEVENTF_MOVE),
		uintptr(absX),
		uintptr(absY),
		0, 0)

	// 左键按下
	_, _, _ = procMouseEvent.Call(
		uintptr(MOUSEEVENTF_ABSOLUTE|MOUSEEVENTF_LEFTDOWN),
		uintptr(absX),
		uintptr(absY),
		0, 0)

	// 短暂延迟模拟真实点击
	time.Sleep(50 * time.Millisecond)

	// 左键释放
	_, _, _ = procMouseEvent.Call(
		uintptr(MOUSEEVENTF_ABSOLUTE|MOUSEEVENTF_LEFTUP),
		uintptr(absX),
		uintptr(absY),
		0, 0)
}

// GetCursorPos 获取当前鼠标位置
func GetCursorPosO() (int32, int32) {
	var pt Point
	_, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

func mainO() {
	// 示例：点击屏幕坐标 (500, 300)
	x, y := int32(500), int32(300)

	// 点击前位置
	startX, startY := GetCursorPos()
	fmt.Printf("点击前位置: (%d, %d)\n", startX, startY)

	// 执行点击
	ClickAt(x, y)

	// 点击后位置
	endX, endY := GetCursorPos()
	fmt.Printf("点击后位置: (%d, %d)\n", endX, endY)
}
