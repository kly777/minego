package winapi

import (
	"fmt"
	"image"
	"syscall"
	"unsafe"
)

const (
	SW_RESTORE       = 9
	SW_SHOW          = 5
	SW_SHOWMAXIMIZED = 3
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	findWindow          = user32.NewProc("FindWindowW")
	setForegroundWindow = user32.NewProc("SetForegroundWindow")
	showWindow          = user32.NewProc("ShowWindow")
	getWindowRect       = user32.NewProc("GetWindowRect")
)

func FindWindow(className, windowName string) (HWND, error) {
	classPtr, _ := syscall.UTF16PtrFromString(className)
	windowPtr, _ := syscall.UTF16PtrFromString(windowName)

	hwnd, _, _ := findWindow.Call(
		uintptr(unsafe.Pointer(classPtr)),
		uintptr(unsafe.Pointer(windowPtr)),
	)

	if hwnd == 0 {
		return 0, fmt.Errorf("mine window not found")
	}
	return hwnd, nil
}

func FindMineWindow() (HWND, error) {
	className := "Minesweeper"
	windowName := "扫雷"
	return FindWindow(className, windowName)
}

func activateWindow(hwnd uintptr) error {
	showWindow.Call(hwnd, SW_RESTORE)
	_, _, err := setForegroundWindow.Call(hwnd)
	if err != syscall.Errno(0) {
		return fmt.Errorf("activate window failed: %v", err)
	}
	return nil
}

func getWindowBounds(hwnd uintptr) (image.Rectangle, error) {
	var rect [4]int32
	_, _, _ = getWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect[0])))

	left, top := LogicalToPhysical(hwnd, int(rect[0]), int(rect[1]))
	right, bottom := LogicalToPhysical(hwnd, int(rect[2]), int(rect[3]))
	return image.Rect(left, top, right, bottom), nil
}
