package main

import "gomini/klog"

func main() {
	klog.LogOpen("klog\\configdefault.toml")
	klog.Info("hello world.")
	klog.LogClose()
}
