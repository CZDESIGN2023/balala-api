package file

import (
	"net/http"
	"os"
)

func DirFs(dirPath string) NoDirFileSystem {
	return NoDirFileSystem{http.Dir(dirPath)}
}

// NoDirFileSystem 是一个自定义的文件系统，它禁止列出目录
type NoDirFileSystem struct {
	fs http.FileSystem
}

// Open 实现了 http.FileSystem 接口的 Open 方法
func (fs NoDirFileSystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

// neuteredReaddirFile 禁止了 Readdir 方法，这样就不能列出目录
type neuteredReaddirFile struct {
	http.File
}

// Readdir 返回一个空的切片和一个错误，这样就禁止了目录的读取
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// 返回 nil 和一个错误来禁止列出目录
	return nil, nil
}
