package kfile

import (
	"fmt"
	"io"
	"os"
)

// GetFileSize 获取文件大小
func GetFileSize(path string) int64 {
	st, err := os.Stat(path)
	if err != nil || st.IsDir() == true {
		return -1
	}

	return st.Size()
}

// IsPathExist 判断文件或者目录是否存在
func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsFileCanExec 判断文件是否拥有可执行权限
func IsFileCanExec(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return uint32(73) == uint32(fileInfo.Mode().Perm()&os.FileMode(0111))
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
