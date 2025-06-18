package screenshot

import (
	"errors"
	"fmt"
	"image"
	"minego/pkg/kit"
	"syscall"

	"unsafe"
)

var (
	user32                     = syscall.MustLoadDLL("user32.dll")
	gdi32                      = syscall.MustLoadDLL("gdi32.dll")
	procGetDC                  = user32.MustFindProc("GetDC")
	procReleaseDC              = user32.MustFindProc("ReleaseDC")
	procCreateCompatibleDC     = gdi32.MustFindProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.MustFindProc("CreateCompatibleBitmap")
	procSelectObject           = gdi32.MustFindProc("SelectObject")
	procBitBlt                 = gdi32.MustFindProc("BitBlt")
	procDeleteObject           = gdi32.MustFindProc("DeleteObject")
	procDeleteDC               = gdi32.MustFindProc("DeleteDC")
	procGetDIBits              = gdi32.MustFindProc("GetDIBits")
)

const (
	SRCCOPY        = 0x00CC0020
	BI_RGB         = 0
	DIB_RGB_COLORS = 0
)

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors [1]uint32
}

// CaptureRect 捕获指定屏幕区域并返回 *image.RGBA
func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	hdcScreen, _, _ := procGetDC.Call(0)
	if hdcScreen == 0 {
		return nil, errors.New("GetDC failed")
	}
	defer procReleaseDC.Call(0, hdcScreen)

	hdcMem, _, _ := procCreateCompatibleDC.Call(hdcScreen)
	if hdcMem == 0 {
		return nil, errors.New("CreateCompatibleDC failed")
	}
	defer procDeleteDC.Call(hdcMem)

	width := rect.Dx()
	height := rect.Dy()

	// 创建兼容位图
	hbm, _, _ := procCreateCompatibleBitmap.Call(
		hdcScreen,
		uintptr(width),
		uintptr(height),
	)
	if hbm == 0 {
		return nil, errors.New("CreateCompatibleBitmap failed")
	}
	defer procDeleteObject.Call(hbm)

	// 将位图选入内存DC
	hOld, _, _ := procSelectObject.Call(hdcMem, hbm)
	if hOld == 0 {
		return nil, errors.New("SelectObject failed")
	}
	defer procSelectObject.Call(hdcMem, hOld)

	// 执行位块传输
	ret, _, _ := procBitBlt.Call(
		hdcMem,
		0, 0,
		uintptr(width), uintptr(height),
		hdcScreen,
		uintptr(rect.Min.X), uintptr(rect.Min.Y),
		SRCCOPY,
	)
	if ret == 0 {
		return nil, errors.New("BitBlt failed")
	}

	// 准备位图信息结构
	bi := BITMAPINFO{
		Header: BITMAPINFOHEADER{
			Size:        uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
			Width:       int32(width),
			Height:      int32(-height), // 负值表示自上而下的DIB
			Planes:      1,
			BitCount:    32,
			Compression: BI_RGB,
		},
	}

	// 创建临时缓冲区接收BGRA数据
	bgraData := make([]byte, width*height*4)
	ret, _, _ = procGetDIBits.Call(
		hdcMem,
		hbm,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&bgraData[0])), // 修改为临时缓冲区
		uintptr(unsafe.Pointer(&bi)),
		DIB_RGB_COLORS,
	)
	if ret == 0 {
		return nil, errors.New("GetDIBits failed")
	}
	// 创建RGBA图像并转换数据
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 转换BGRA到RGBA并翻转垂直方向
	for y := range height {
		srcY := y * width * 4 // 从底部开始读取

		for x := range width {
			srcIdx := srcY + x*4

			// 交换R和B通道
			img.Pix[srcIdx+0] = bgraData[srcIdx+2] // R
			img.Pix[srcIdx+1] = bgraData[srcIdx+1] // G
			img.Pix[srcIdx+2] = bgraData[srcIdx+0] // B
			// img.Pix[srcIdx+3] = bgraData[srcIdx+3] // A
			img.Pix[srcIdx+3] = 255
		}
	}

	return img, nil
}

func main() {
	fmt.Println("s")
	rect := image.Rect(0, 0, 300, 300)
	img, err := CaptureRect(rect)
	if err != nil {
		panic(err)
	}
	kit.SaveImg(img, "captured.png")
	fmt.Println(img.Bounds())
}
