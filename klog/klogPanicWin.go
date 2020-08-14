// +build windows

package klog

import (
	"os"
	"syscall"
)

const (
	kernel32dll  = "kernel32.dll"
	stderrHandle = uint32(-12 & 0xFFFFFFFF)
)

func InitPanicFile(logFile string) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	v, _, err := setStdHandle.Call(uintptr(stderrHandle), file.Fd())
	if v == 0 {
		return err
	}
	return nil
}
