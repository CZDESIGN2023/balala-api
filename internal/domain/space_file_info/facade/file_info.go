package facade

import (
	"go-cs/internal/domain/file_info"
	"go-cs/internal/domain/space_file_info"
)

type FileInfoFacade struct {
	Files []space_file_info.FileInfo
}

func BuildFileInfoFacade(fileInfos []*file_info.FileInfo) *FileInfoFacade {
	sFileInfos := make([]space_file_info.FileInfo, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		// 构建空间文件信息
		sFileInfos = append(sFileInfos, space_file_info.FileInfo{
			FileInfoId: fileInfo.Id,
			FileName:   fileInfo.Name,
			FileUri:    fileInfo.Uri,
			FileSize:   fileInfo.Size,
		})
	}

	return &FileInfoFacade{
		Files: sFileInfos,
	}
}
