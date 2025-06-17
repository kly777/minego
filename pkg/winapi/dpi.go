package winapi

import (

	"syscall"
)
func SetDPIAware() {
	user32 := syscall.NewLazyDLL("user32.dll")
	procSetProcessDPIAware := user32.NewProc("SetProcessDPIAware")
	procSetProcessDPIAware.Call()
}
func GetDPI(hwnd uintptr) (int, error) {
	user32 := syscall.NewLazyDLL("user32.dll")
	getDpiForWindow := user32.NewProc("GetDpiForWindow")
	if getDpiForWindow.Find() != nil {
		hdc, _, _ := syscall.SyscallN(user32.NewProc("GetDC").Addr(), 1, hwnd, 0, 0)
		if hdc == 0 {
			return 96, nil
		}
		defer syscall.SyscallN(user32.NewProc("ReleaseDC").Addr(), 2, hwnd, hdc, 0)
		logPixelsX, _, _ := syscall.SyscallN(syscall.NewLazyDLL("gdi32.dll").NewProc("GetDeviceCaps").Addr(), 2, hdc, 88, 0)
		return int(logPixelsX), nil
	}
	dpi, _, _ := getDpiForWindow.Call(hwnd)
	return int(dpi), nil
}

func LogicalToPhysical(hwnd uintptr, x, y int) (int, int) {
	dpi, _ := GetDPI(hwnd)
	if dpi == 96 {
		return x, y
	}
	scale := float64(dpi) / 96.0
	return int(float64(x) * scale), int(float64(y) * scale)
}

