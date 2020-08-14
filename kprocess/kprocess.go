package kprocess

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetCurrentModulePath() string {
	procFile, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}

	return procFile
}

func GetCurrentModuleName() string {
	procFile := GetCurrentModulePath()
	if procFile == "" {
		return ""
	}

	return filepath.Base(procFile)
}

func GetCurrentModuleNameNoExt() string {
	procName := GetCurrentModuleName()
	for i := len(procName) - 1; i >= 0; i-- {
		if procName[i] == '.' {
			return procName[:i]
		}
	}

	return ""
}
