package klog

import (
	"fmt"
	"gomini/kprocess"
	"gomini/ktoml"
	"os"
	"time"
)

type logConfig struct {
	Level        string `toml:"Level"`
	LogDir       string `toml:"LogDir"`
	CrashLogDir  string `toml:"CrashLogDir"`
	FileBaseName string `toml:"FileBaseName"`
	MaxDays      int    `toml:"MaxDays"`
	CustomPrefix string `toml:"CustomPrefix"`
}

type klogCfg struct {
	data *logConfig
}

func (k *klogCfg) GenLogName() string {
	t := time.Now()
	logName := fmt.Sprintf(k.data.FileBaseName+"-%d-%d-%d %d-%d-%d.klog",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return logName
}

func (k *klogCfg) GenCrashLogName() string {
	t := time.Now()
	logName := fmt.Sprintf(k.data.FileBaseName+"crash-%d-%d-%d %d-%d-%d.klog",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return logName
}

func (k *klogCfg) Load(cfgfile string) error {

	k.data = &logConfig{
		"INFO",
		"logs",
		"crashes",
		kprocess.GetCurrentModuleNameNoExt(),
		7,
		"",
	}

	if cfgfile != "" {
		err := ktoml.LoadToml(cfgfile, &k.data)
		if err != nil {
			return err
		}
	}

	if k.data.FileBaseName == "" {
		k.data.FileBaseName = kprocess.GetCurrentModuleNameNoExt()
	}

	k.data.LogDir = os.ExpandEnv(k.data.LogDir)
	k.data.CrashLogDir = os.ExpandEnv(k.data.CrashLogDir)

	return nil
}

func levelInt(level string) int {
	if level == "INFO" {
		return 2
	} else if level == "WARN" {
		return 3
	} else if level == "ERROR" {
		return 4
	}

	return 1
}

func (k *klogCfg) Enable(level string) bool {
	if levelInt(level) >= levelInt(k.data.Level) {
		return true
	}
	return false
}
