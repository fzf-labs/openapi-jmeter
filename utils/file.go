package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ReadDirFilesWithSuffix 遍历指定目录下的指定后缀文件
func ReadDirFilesWithSuffix(root, suffix string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// ReadFileToString 读取文件到string
func ReadFileToString(dir string) (string, error) {
	file, err := os.OpenFile(dir, os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

// WriteContentCover 数据写入，不存在则创建
func WriteContentCover(filePath, content string) error {
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0775); err != nil {
		return err
	}
	dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0665)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}
	return err
}
