package facade

import (
	"context"
	file_service "go-cs/internal/domain/file_info/service"
	domain "go-cs/internal/domain/work_item"
)

type FileInfoFacade struct {
	fileService *file_service.FileInfoService
}

func (f *FileInfoFacade) GetFileInfos(ctx context.Context, fileInfoIds []int64) ([]domain.FileInfo, error) {
	fileInfos, err := f.fileService.GetFileInfos(ctx, fileInfoIds)
	if err != nil {
		return nil, err
	}

	var workItemFiles []domain.FileInfo
	for _, v := range fileInfos {
		workItemFiles = append(workItemFiles, domain.FileInfo{
			FileInfoId: v.Id,
			FileName:   v.Name,
			FileUri:    v.Uri,
			FileSize:   v.Size,
		})
	}

	return workItemFiles, nil
}

func BuildFileInfoFacade(fileService *file_service.FileInfoService) *FileInfoFacade {

	return &FileInfoFacade{
		fileService: fileService,
	}
}
