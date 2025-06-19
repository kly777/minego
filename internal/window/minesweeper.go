package window

import "minego/pkg/winapi"

func findMineWindow() (winapi.HWND, error) {
	className := "Minesweeper"
	windowName := "扫雷"
	return winapi.FindWindow(className, windowName)
}


func GetMineSweeperWindow() winapi.Window {
	hwnd, _ := findMineWindow()
	return winapi.NewWindow(hwnd)
}
