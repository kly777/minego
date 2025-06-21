package click

import (
	"fmt"
	"image"
	"syscall"
	"unsafe"
)

// 定义Windows API常量
const (
	INPUT_MOUSE           = 0
	MOUSEEVENTF_MOVE      = 0x0001
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010
	MOUSEEVENTF_ABSOLUTE  = 0x8000

	SM_CXSCREEN                   = 0 // 主显示器宽度
	SM_CYSCREEN                   = 1 // 主显示器高度
	SM_XVIRTUALSCREEN             = 76
	SM_YVIRTUALSCREEN             = 77
	SM_CXVIRTUALSCREEN            = 78
	SM_CYVIRTUALSCREEN            = 79
	PROCESS_PER_MONITOR_DPI_AWARE = 2
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")
	modshcore = syscall.NewLazyDLL("shcore.dll")

	procSendInput              = moduser32.NewProc("SendInput")
	procGetSystemMetrics       = moduser32.NewProc("GetSystemMetrics")
	procSetProcessDpiAwareness = modshcore.NewProc("SetProcessDpiAwareness")
)

// 定义Windows API结构体
type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type uint32
	Mi   MOUSEINPUT
}

// SetDPIAware 设置进程DPI感知
func SetDPIAware() bool {
	// 尝试使用SetProcessDpiAwareness（Windows 8.1+）
	r1, _, _ := procSetProcessDpiAwareness.Call(PROCESS_PER_MONITOR_DPI_AWARE)
	if r1 == 0 {
		return true
	}

	// 回退到SetProcessDPIAware（Windows Vista+）
	procSetProcessDPIAware := moduser32.NewProc("SetProcessDPIAware")
	r1, _, _ = procSetProcessDPIAware.Call()
	return r1 != 0
}

// GetSystemMetrics 封装系统指标获取
func GetSystemMetrics(index int) int32 {
	r1, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(r1)
}

func Click(p image.Point) {
	PhysicalMouseClick(int32(p.X), int32(p.Y))
}

func RightClick(p image.Point) {
	PhysicalRightMouseClick(int32(p.X), int32(p.Y))
}

func PhysicalRightMouseClick(x, y int32) {
	// 获取虚拟桌面范围

	primaryWidth, primaryHeight := GetPrimaryMonitorResolution()

	// 计算归一化坐标 (0-65535)
	normalizedX := int32((float64(x) / float64(primaryWidth-1) * 65535))
	normalizedY := int32((float64(y) / float64(primaryHeight-1) * 65535))

	// 创建鼠标事件序列
	inputs := []INPUT{
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				Dx:          normalizedX,
				Dy:          normalizedY,
				DwFlags:     MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags:     MOUSEEVENTF_RIGHTDOWN,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				DwFlags:     MOUSEEVENTF_RIGHTUP,
				Time:        0,
				DwExtraInfo: 0,
			},
		},
	}

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

// PhysicalMouseClick 在指定物理坐标执行鼠标点击
func PhysicalMouseClick(x, y int32) {

	primaryWidth, primaryHeight := GetPrimaryMonitorResolution()

	// 计算归一化坐标 (0-65535)
	normalizedX := int32((float64(x) / float64(primaryWidth-1) * 65535))
	normalizedY := int32((float64(y) / float64(primaryHeight-1) * 65535))

	// 创建鼠标事件序列
	inputs := []INPUT{
		{
			Type: INPUT_MOUSE,
			Mi: MOUSEINPUT{
				Dx:          normalizedX,
				Dy:          normalizedY,
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

func GetPrimaryMonitorResolution() (width, height int32) {
	width = GetSystemMetrics(SM_CXSCREEN)
	height = GetSystemMetrics(SM_CYSCREEN)
	return
}
