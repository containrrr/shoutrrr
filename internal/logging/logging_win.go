//+build windows

package logging

import (
	"syscall"

	"golang.org/x/sys/windows"
)

// ENABLE_VIRTUAL_TERMINAL_PROCESSING
// as per https://docs.microsoft.com/en-us/windows/console/setconsolemode#parameters
const enableVirtualTerminalProcessing uint32 = 0x0004

var (
	procSetConsoleMode = windows.NewLazySystemDLL("kernel32.dll").NewProc("SetConsoleMode")
)

func initPlatform() error {
	return setEnableVirtualTerminalProcessing(syscall.Stderr, true)
}

func setConsoleMode(consoleHandle syscall.Handle, mode uint32) error {
	procAddr := procSetConsoleMode.Addr()
	if ret, _, err := syscall.Syscall(procAddr, 2, uintptr(consoleHandle), uintptr(mode), 0); ret == 0 {
		return err
	}
	return nil
}

func setEnableVirtualTerminalProcessing(screenBufferHandle syscall.Handle, enable bool) error {

	// Get current mode
	var mode uint32
	if err := syscall.GetConsoleMode(screenBufferHandle, &mode); err != nil {
		return err
	}

	// Set/unset ENABLE_VIRTUAL_TERMINAL_PROCESSING bit
	if enable {
		mode |= enableVirtualTerminalProcessing
	} else {
		mode &^= enableVirtualTerminalProcessing
	}

	// Set modified mode
	return setConsoleMode(screenBufferHandle, mode)
}
