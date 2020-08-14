// +build linux

package klog

func InitPanicFile(logFile string) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return err
	}
	return nil
}
