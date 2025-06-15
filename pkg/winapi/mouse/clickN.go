package mouse

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Windows API 常量
const (
	INPUT_MOUSE          = 0
	MOUSEEVENTF_MOVE     = 0x0001
	MOUSEEVENTF_ABSOLUTE = 0x8000
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
	SM_CXSCREEN          = 0
	SM_CYSCREEN          = 1
)

// Windows API 函数声明
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procSendInput        = user32.NewProc("SendInput")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procGetCursorPos     = user32.NewProc("GetCursorPos")
)

// POINT 结构体表示屏幕坐标
type POINT struct {
	X int32
	Y int32
}

// MOUSEINPUT 结构体
type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// INPUT 结构体
type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

// ClickAt 在指定坐标执行鼠标点击 (使用 SendInput)
func ClickAt(x, y int32) {
	// 获取屏幕分辨率
	cxScreen, _, _ := procGetSystemMetrics.Call(SM_CXSCREEN)
	cyScreen, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)

	// 转换为绝对坐标 (0-65535)
	absX := x * 65535 / int32(cxScreen)
	absY := y * 65535 / int32(cyScreen)
	fmt.Println("绝对坐标:", absX, absY)

	// 创建三个输入事件: 移动、左键按下、左键释放
	inputs := []INPUT{
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				Dx:          absX,
				Dy:          absY,
				DwFlags:     MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags:     MOUSEEVENTF_LEFTDOWN,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags:     MOUSEEVENTF_LEFTUP,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	}

	// 发送输入事件
	size := unsafe.Sizeof(INPUT{})
	for _, input := range inputs {
		r, _, err := procSendInput.Call(
			1, // cInputs
			uintptr(unsafe.Pointer(&input)),
			uintptr(size),
		)

		if r == 0 {
			fmt.Printf("SendInput 失败: %v\n", err)
		}
	}
}

// GetCursorPos 获取当前鼠标位置
func GetCursorPos() (int32, int32) {
	var pt POINT
	_, _, _ = procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return pt.X, pt.Y
}

// 示例：点击屏幕坐标 (500, 300)
// x, y := int32(500), int32(300)

// // 点击前位置
// startX, startY := GetCursorPos()
// fmt.Printf("点击前位置: (%d, %d)\n", startX, startY)

// // 执行点击
// ClickAt(x, y)

// // 短暂延迟
// time.Sleep(100 * time.Millisecond)

// // 点击后位置
// endX, endY := GetCursorPos()
// fmt.Printf("点击后位置: (%d, %d)\n", endX, endY)
