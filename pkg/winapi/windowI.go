package winapi

import "image"

type HWND = uintptr

type window struct {
	hwnd HWND
}

type Window interface {
	Activate()
	GetBounds() (image.Rectangle, error)
}

func NewWindow(hwnd HWND) Window {
	return &window{hwnd: hwnd}
}


func (w window) Activate() {
	activateWindow(w.hwnd)
}

func (w window) GetBounds() (image.Rectangle, error) {
	return getWindowBounds(w.hwnd)
}
