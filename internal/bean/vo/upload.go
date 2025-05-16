package vo

import (
	"io"
)

type UploadFileToLocalVo struct {
	FileRootPath string
	OwnerId      int64
	FileInfo     *UploadFileInfoVo
}

type UploadFileInfoVo struct {
	File           io.Reader
	FileName       string
	FileSize       int64
	FileUri        string
	FileUploadPath string
}
