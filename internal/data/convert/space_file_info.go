package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_file_info"
)

func SpaceFileInfoEntityToPo(fileInfo *domain.SpaceFileInfo) *db.SpaceFileInfo {
	// 转换逻辑
	return &db.SpaceFileInfo{
		Id:         fileInfo.Id,
		SpaceId:    fileInfo.SpaceId,
		FileInfoId: fileInfo.FileInfo.FileInfoId,
		FileName:   fileInfo.FileInfo.FileName,
		FileSize:   fileInfo.FileInfo.FileSize,
		FileUri:    fileInfo.FileInfo.FileUri,
		Status:     int32(fileInfo.Status),
		SourceId:   fileInfo.FileSource.SourceId,
		SourceType: int32(fileInfo.FileSource.SourceType),
		CreatedAt:  fileInfo.CreatedAt,
		UpdatedAt:  fileInfo.UpdatedAt,
		DeletedAt:  fileInfo.DeletedAt,
	}
}

func SpaceFileInfoPoSliceToEntitys(fileInfoList []*db.SpaceFileInfo) []*domain.SpaceFileInfo {
	// 转换逻辑
	var entityList []*domain.SpaceFileInfo
	for _, fileInfo := range fileInfoList {
		entityList = append(entityList, SpaceFileInfoPoToEntity(fileInfo))
	}
	return entityList
}

func SpaceFileInfoPoToEntity(fileInfo *db.SpaceFileInfo) *domain.SpaceFileInfo {
	// 转换逻辑
	return &domain.SpaceFileInfo{
		Id:      fileInfo.Id,
		SpaceId: fileInfo.SpaceId,
		FileInfo: domain.FileInfo{
			FileInfoId: fileInfo.FileInfoId,
			FileName:   fileInfo.FileName,
			FileSize:   fileInfo.FileSize,
			FileUri:    fileInfo.FileUri,
		},
		Status: domain.FileStatus(fileInfo.Status),
		FileSource: domain.FileSource{
			SourceId:   fileInfo.SourceId,
			SourceType: domain.FileSourceType(fileInfo.SourceType),
		},
		CreatedAt: fileInfo.CreatedAt,
		UpdatedAt: fileInfo.UpdatedAt,
		DeletedAt: fileInfo.DeletedAt,
	}
}

func SpaceFileInfoEntityToPos(fileInfoList []*domain.SpaceFileInfo) []*db.SpaceFileInfo {
	// 转换逻辑
	var poList []*db.SpaceFileInfo
	for _, fileInfo := range fileInfoList {
		poList = append(poList, SpaceFileInfoEntityToPo(fileInfo))
	}
	return poList
}
