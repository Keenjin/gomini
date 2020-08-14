package klog

import (
	"gomini/kfile"
	"gomini/kprocess"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// klogImpl 日志实现类
type klogImpl struct {
	cfg klogCfg // 配置文件

	debug *log.Logger // 记录所有日志
	info  *log.Logger // 重要的信息
	warn  *log.Logger // 需要注意的信息
	error *log.Logger // 非常严重的问题

	stdoutLog *log.Logger // 常规的控制台输出
	stderrLog *log.Logger

	fileFd *os.File
}

func (k *klogImpl) setLogger(file *os.File, level string) {
	prefix := k.cfg.data.CustomPrefix + "[" + level + "] "
	flag := log.Ldate | log.Ltime | log.Lshortfile
	if level == "_STDOUT" {
		k.stdoutLog = log.New(io.MultiWriter(os.Stdout), prefix, flag)
	} else if level == "_STDERR" {
		k.stderrLog = log.New(io.MultiWriter(os.Stderr), prefix, flag)
	} else if level == "ERROR" {
		k.error = log.New(io.MultiWriter(os.Stderr, file), prefix, flag)
	} else {
		logger := log.New(io.MultiWriter(os.Stdout, file), prefix, flag)
		if level == "DEBUG" {
			k.debug = logger
		} else if level == "INFO" {
			k.info = logger
		} else if level == "WARN" {
			k.warn = logger
		}
	}
}

func (k *klogImpl) createLogFile(logPath string) (*os.File, error) {
	logger, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func (k *klogImpl) Open(logCfg string) error {
	err := k.cfg.Load(logCfg)
	if err != nil {
		return err
	}

	if !kfile.IsPathExist(k.cfg.data.LogDir) {
		os.MkdirAll(k.cfg.data.LogDir, 0666)
	}
	if !kfile.IsPathExist(k.cfg.data.CrashLogDir) {
		os.MkdirAll(k.cfg.data.CrashLogDir, 0666)
	}

	logFilePath := filepath.Join(k.cfg.data.LogDir, k.cfg.GenLogName())
	k.fileFd, err = k.createLogFile(logFilePath)
	if err != nil {
		return err
	}

	k.setLogger(k.fileFd, "DEBUG")
	k.setLogger(k.fileFd, "INFO")
	k.setLogger(k.fileFd, "WARN")
	k.setLogger(k.fileFd, "ERROR")
	k.setLogger(k.fileFd, "_STDOUT")
	k.setLogger(k.fileFd, "_STDERR")

	// 初始化crash日志
	crashLogPath := filepath.Join(k.cfg.data.CrashLogDir, k.cfg.GenCrashLogName())
	_ = InitPanicFile(crashLogPath)

	// 清除过期的日志文件
	k.clearOutTimeLogs(k.cfg.data.LogDir, k.cfg.data.MaxDays)

	k.info.Printf("================== %s begin ==================", kprocess.GetCurrentModuleName())

	return nil
}

func (k *klogImpl) Close() {
	k.info.Printf("================== %s end ==================", kprocess.GetCurrentModuleName())
	if k.fileFd != nil {
		k.fileFd.Close()
	}
}

func (k *klogImpl) clearOutTimeLogs(logDir string, maxDays int) {
	// 遍历目录，把.klog后缀的排序，删除最老的
	var logs []string
	filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".klog" {
			logs = append(logs, path)
		}
		return nil
	})
	if len(logs) > maxDays {
		sort.Strings(logs)
		for i := 0; i < len(logs)-maxDays; i++ {
			os.RemoveAll(logs[i])
		}
	}
}

func (k *klogImpl) GetLogger(level string) *log.Logger {
	if k.cfg.Enable(level) {
		if level == "DEBUG" {
			return k.debug
		} else if level == "INFO" {
			return k.info
		} else if level == "WARN" {
			return k.warn
		} else if level == "ERROR" {
			return k.error
		}
	}

	if level == "ERROR" {
		return k.stderrLog
	}
	return k.stdoutLog
}
