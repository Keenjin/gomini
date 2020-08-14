package klog

import "fmt"

var myLogger *klogImpl

// exe的main函数中，使用LogOpen
func LogOpen(logCfg string) {
	if myLogger != nil {
		return
	}

	myLogger = &klogImpl{}
	myLogger.Open(logCfg)
}

//
func LogClose() {
	if myLogger != nil {
		myLogger.Close()
	}
}

func Debug(format string, args ...interface{}) {
	if myLogger != nil {
		myLogger.GetLogger("DEBUG").Output(3, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if myLogger != nil {
		myLogger.GetLogger("INFO").Output(3, fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if myLogger != nil {
		myLogger.GetLogger("WARN").Output(3, fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if myLogger != nil {
		myLogger.GetLogger("ERROR").Output(3, fmt.Sprintf(format, args...))
	}
}
