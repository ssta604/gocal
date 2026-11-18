//go:build windows

package main

import "syscall"

const enableVirtualTerminalProcessing = 0x0004

// enableVT turns on ANSI escape handling for older Windows consoles.
func enableVT() {
	h, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return
	}
	var mode uint32
	if syscall.GetConsoleMode(h, &mode) != nil {
		return
	}
	setConsoleMode := syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleMode")
	setConsoleMode.Call(uintptr(h), uintptr(mode|enableVirtualTerminalProcessing))
}
