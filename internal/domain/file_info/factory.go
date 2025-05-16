package file_info

import (
	"go-cs/internal/utils"
	"time"
)

func NewFileInfo(
	id int64,
	hash string,
	name string,
	typ int32,
	size int64,
	uri string,
	cover string,
	owner int64,
	meta map[string]string,
	uploadDomain string,
	uploadMd5 string,
	uploadPath string,
) *FileInfo {
	return &FileInfo{
		Id:           id,
		Hash:         hash,
		Name:         name,
		Typ:          typ,
		Size:         size,
		Uri:          uri,
		Cover:        cover,
		Owner:        owner,
		Meta:         utils.ToJSON(meta),
		UploadDomain: uploadDomain,
		UploadMd5:    uploadMd5,
		UploadPath:   uploadPath,
		CreatedAt:    time.Now().Unix(),
		Status:       3,
	}
}
