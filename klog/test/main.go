package main

import "../klog"

func main() {
	klog.LogOpen("klog\\configdefault.toml")
	klog.Info("hello world.")
	klog.LogClose()
}
