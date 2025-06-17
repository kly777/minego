package screen

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Windows 常量
const (
	INPUT_MOUSE              = 0
	MOUSEEVENTF_MOVE         = 0x0001
	MOUSEEVENTF_ABSOLUTE     = 0x8000
	MOUSEEVENTF_LEFTDOWN     = 0x0002
	MOUSEEVENTF_LEFTUP       = 0x0004
	MONITOR_DEFAULTTONEAREST = 0x00000002
)

const (
	SM_XVIRTUALSCREEN  = 76
	SM_YVIRTUALSCREEN  = 77
	SM_CXVIRTUALSCREEN = 78
	SM_CYVIRTUALSCREEN = 79
)

// Windows 结构体
type (
	POINT struct {
		X, Y int32
	}

	RECT struct {
		Left, Top, Right, Bottom int32
	}

	MOUSEINPUT struct {
		Dx, Dy      int32
		MouseData   uint32
		DwFlags     uint32
		Time        uint32
		DwExtraInfo uintptr
	}

	INPUT struct {
		Type uint32
		Mi   MOUSEINPUT
	}
)

// Windows API 函数声明
var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	procSendInput           = user32.NewProc("SendInput")
	procMonitorFromPoint    = user32.NewProc("MonitorFromPoint")
	procGetDpiForMonitor    = user32.NewProc("GetDpiForMonitor")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
)

// 回调函数类型
type MonitorEnumProc uintptr

// getVirtualDesktopSize 获取虚拟桌面尺寸
func getVirtualDesktopRect() (int32, int32, int32, int32) {
	left, _, _ := procGetSystemMetrics.Call(SM_XVIRTUALSCREEN)
	top, _, _ := procGetSystemMetrics.Call(SM_YVIRTUALSCREEN)
	width, _, _ := procGetSystemMetrics.Call(SM_CXVIRTUALSCREEN)
	height, _, _ := procGetSystemMetrics.Call(SM_CYVIRTUALSCREEN)
	fmt.Printf("left:%d, top:%d, width:%d, height:%d\n", left, top, width, height)
	return int32(left), int32(top), int32(width), int32(height)
}

// logicalToPhysicalPoint 将逻辑坐标转换为物理坐标
func logicalToPhysicalPoint(hwnd uintptr, pt *POINT) bool {
	procLogicalToPhysicalPoint := user32.NewProc("LogicalToPhysicalPoint")
	ret, _, _ := procLogicalToPhysicalPoint.Call(
		hwnd,
		uintptr(unsafe.Pointer(pt)),
	)
	return ret != 0
}

// getMonitorDPI 获取显示器DPI
func getMonitorDPI(x, y int32) (dpiX, dpiY uint32) {
	pt := POINT{X: x, Y: y}
	hmonitor, _, _ := procMonitorFromPoint.Call(
		uintptr(*(*int64)(unsafe.Pointer(&pt))),
		MONITOR_DEFAULTTONEAREST,
	)

	if hmonitor != 0 {
		procGetDpiForMonitor.Call(
			hmonitor,
			0, // DPI type: MDT_EFFECTIVE_DPI
			uintptr(unsafe.Pointer(&dpiX)),
			uintptr(unsafe.Pointer(&dpiY)),
		)
	}
	return
}

// sendMouseClick 发送鼠标点击事件
func sendMouseClick(x, y int32) {
	// 获取虚拟桌面范围（包含负坐标区域）
	virtLeft, virtTop, virtWidth, virtHeight := getVirtualDesktopRect()
	fmt.Printf("Virtual desktop: %d, %d, %d, %d\n", virtLeft, virtTop, virtWidth, virtHeight)
	// 验证坐标是否在有效范围内
	if x < virtLeft || x >= virtLeft+virtWidth ||
		y < virtTop || y >= virtTop+virtHeight {
		fmt.Printf("坐标超出范围: (%d, %d) 不在 [%d, %d]x[%d, %d]内\n",
			x, y, virtLeft, virtLeft+virtWidth-1, virtTop, virtTop+virtHeight-1)
		return
	}

	// 计算归一化坐标（支持负坐标转换）
	absX := int32((float64(x-virtLeft) * 65535.0) / float64(virtWidth-1))
	absY := int32((float64(y-virtTop) * 65535.0) / float64(virtHeight-1))

	fmt.Printf("点击: (%d, %d) -> 绝对坐标: (%d, %d)\n", x, y, absX, absY)

	// 创建鼠标事件序列
	inputs := []INPUT{
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				Dx:          -65535,
				Dy:          0,
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

	// 发送事件
	size := unsafe.Sizeof(INPUT{})
	for _, input := range inputs {
		r, _, err := procSendInput.Call(
			1,
			uintptr(unsafe.Pointer(&input)),
			uintptr(size),
		)
		if r == 0 {
			fmt.Printf("SendInput失败: %v\n", err)
		}
	}
}

// ClickAtPhysical 在物理坐标点击
func ClickAtPhysical(x, y int) {
	sendMouseClick(int32(x), int32(y))
}
