package khttp

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/Keenjin/gomini/kfile"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// KHttpPostJson Post一个http请求
func KHttpPostJson(url, body string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// KHttpDownloadSimple 针对小文件的下载（建议10M以内），根据链接下载文件，下载到dstPath目录。
func KHttpDownloadSimple(url, dstPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if kfile.IsPathExist(dstPath) {
		os.RemoveAll(dstPath)
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, resp.Body)

	return nil
}

// DownloadBigFile 下载大文件，根据链接下载文件，下载到dstPath。
func DownloadBigFile(user, passwd, url, dstPath string, onProgress func(err error, totalLen, downloadLen int64)) error {
	reader := strings.NewReader("hello")
	req, _ := http.NewRequest("GET", url, reader)
	if user != "" && passwd != "" {
		req.SetBasicAuth(user, passwd)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	//resp, err := client.Get(url)
	if err != nil {
		return err
	}

	// 读取服务器返回的文件大小
	fsize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return err
	}

	if fsize == 0 {
		return errors.New("文件大小为0")
	}

	// 创建文件
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var downloadLen int64

	buf := make([]byte, 8192)

	for {
		readSize, err := resp.Body.Read(buf)
		if readSize == 0 {
			break
		}

		file.Write(buf[:readSize])

		downloadLen += int64(readSize)

		if onProgress != nil {
			onProgress(nil, fsize, downloadLen)
		}

		if err != nil {
			break
		}
	}

	if downloadLen < fsize {
		err = errors.New("下载未完成")
		onProgress(err, fsize, downloadLen)
	}

	return err
}
