package keylistener

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Windows API常量
const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
)

// KBDLLHOOKSTRUCT 结构体
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uint32
}

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procSetWindowsHookEx    = user32.MustFindProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.MustFindProc("CallNextHookEx")
	procGetMessage          = user32.MustFindProc("GetMessageW")
	procTranslateMessage    = user32.MustFindProc("TranslateMessage")
	procDispatchMessage     = user32.MustFindProc("DispatchMessageW")
	procUnhookWindowsHookEx = user32.MustFindProc("UnhookWindowsHookEx")
)

var hookHandle uintptr

//export keyboardHookProc
func keyboardHookProc(nCode int32, wParam uintptr, lParam unsafe.Pointer) uintptr {
	if nCode >= 0 {
		if wParam == WM_KEYDOWN {
			khs := (*KBDLLHOOKSTRUCT)(lParam)
			char := ConvertKeyCodeToChar(khs.VkCode)
			if char != "" {
				fmt.Printf("输入字符: %s (键码: 0x%X)\n", char, khs.VkCode)
				if char == "C" {
					panic("Cancel")
				}
			}
		}
	}
	ret, _, _ := syscall.SyscallN(procCallNextHookEx.Addr(), 4, hookHandle, uintptr(nCode), wParam, uintptr(lParam), 0, 0)
	return ret
}

type MSG struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      [2]int32
}

func Listen() error {
	cb := syscall.NewCallbackCDecl(keyboardHookProc)
	hookHandle, _, _ = syscall.SyscallN(procSetWindowsHookEx.Addr(),
		WH_KEYBOARD_LL, // idHook=WH_KEYBOARD_LL
		cb,             // lpfn=回调函数地址
		0,              // hMod=0 (全局钩子)
		0)

	if hookHandle == 0 {
		return fmt.Errorf("设置钩子失败")
	}
	defer syscall.SyscallN(procUnhookWindowsHookEx.Addr(), 1, hookHandle, 0, 0)

	// 消息循环
	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
	return nil
}

// ConvertKeyCodeToChar 将Windows虚拟键码转换为对应字符（忽略修饰键状态）
func ConvertKeyCodeToChar(vkCode uint32) string {
	// 定义键码到字符的映射表
	keyMap := map[uint32]string{
		// 字母键
		0x30: "0", 0x31: "1", 0x32: "2", 0x33: "3", 0x34: "4",
		0x35: "5", 0x36: "6", 0x37: "7", 0x38: "8", 0x39: "9",
		0x41: "A", 0x42: "B", 0x43: "C", 0x44: "D", 0x45: "E",
		0x46: "F", 0x47: "G", 0x48: "H", 0x49: "I", 0x4A: "J",
		0x4B: "K", 0x4C: "L", 0x4D: "M", 0x4E: "N", 0x4F: "O",
		0x50: "P", 0x51: "Q", 0x52: "R", 0x53: "S", 0x54: "T",
		0x55: "U", 0x56: "V", 0x57: "W", 0x58: "X", 0x59: "Y",
		0x5A: "Z",

		// 符号
		0xDB: "[", 0xDD: "]", 0xDC: "\\", 0xBA: ";", 0xDE: "'",
		0xC0: "`", 0xBB: "=", 0xBC: ",", 0xBE: ".", 0xBF: "/",

		// 特殊键
		0x0D: "[ENTER]", // 回车
		0x20: "[SPACE]", // 空格
		0x1B: "[ESC]",   // 退出键
		0x08: "[BS]",    // 退格键
		0x09: "[TAB]",   // 制表符

		// 数字键盘
		0x60: "0", 0x61: "1", 0x62: "2", 0x63: "3",
		0x64: "4", 0x65: "5", 0x66: "6", 0x67: "7",
		0x68: "8", 0x69: "9", 0x6F: "/", 0x6A: "*",
		0x6B: "+", 0x6D: "-", 0x6E: ".",

		// 功能键（示例）
		0x70: "[F1]", 0x71: "[F2]", 0x72: "[F3]", 0x73: "[F4]",
		0x74: "[F5]", 0x75: "[F6]", 0x76: "[F7]", 0x77: "[F8]",

		// 其他控制键
		0x25: "[←]", 0x26: "[↑]", 0x27: "[→]", 0x28: "[↓]",
		0x2C: "[PrtScn]", 0x2D: "[Insert]", 0x2E: "[Delete]",
		0xA0: "[Shift]", 0xA1: "[Shift]", 0xA2: "[Ctrl]", 0xA3: "[Ctrl]",
		0x5B: "[Win]",
	}

	// 查找对应字符，未找到返回空字符串
	if ch, exists := keyMap[vkCode]; exists {
		return ch
	}
	return "UNKNOWN" // 或返回 "[未知]" 标记
}
