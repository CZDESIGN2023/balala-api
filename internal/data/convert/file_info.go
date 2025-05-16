package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/file_info"
)

func FileInfoEntityToPo(fileInfo *domain.FileInfo) *db.FileInfo {
	return &db.FileInfo{
		Id:     fileInfo.Id,
		Hash:   fileInfo.Hash,
		Name:   fileInfo.Name,
		Typ:    fileInfo.Typ,
		Size:   fileInfo.Size,
		Uri:    fileInfo.Uri,
		Pwd:    fileInfo.Pwd,
		Cover:  fileInfo.Cover,
		Status: fileInfo.Status,
		Owner:  fileInfo.Owner,
		Meta:   fileInfo.Meta,

		UploadTyp:    fileInfo.UploadTyp,
		UploadDomain: fileInfo.UploadDomain,
		UploadMd5:    fileInfo.UploadMd5,
		UploadPath:   fileInfo.UploadPath,

		CreatedAt: fileInfo.CreatedAt,
		UpdatedAt: fileInfo.UpdatedAt,
		DeletedAt: fileInfo.DeletedAt,
	}
}

func FileInfoEntityToPos(fileInfos []*domain.FileInfo) []*db.FileInfo {
	// 这里需要进行类型转换
	var fileInfoPos []*db.FileInfo
	for _, fileInfo := range fileInfos {
		fileInfoPos = append(fileInfoPos, FileInfoEntityToPo(fileInfo))
	}
	return fileInfoPos
}

func FileInfoPoToEntity(fileInfo *db.FileInfo) *domain.FileInfo {
	return &domain.FileInfo{
		Id:     fileInfo.Id,
		Hash:   fileInfo.Hash,
		Name:   fileInfo.Name,
		Typ:    fileInfo.Typ,
		Size:   fileInfo.Size,
		Uri:    fileInfo.Uri,
		Cover:  fileInfo.Cover,
		Pwd:    fileInfo.Pwd,
		Status: fileInfo.Status,
		Owner:  fileInfo.Owner,
		Meta:   fileInfo.Meta,

		UploadTyp:    fileInfo.UploadTyp,
		UploadDomain: fileInfo.UploadDomain,
		UploadMd5:    fileInfo.UploadMd5,
		UploadPath:   fileInfo.UploadPath,

		CreatedAt: fileInfo.CreatedAt,
		UpdatedAt: fileInfo.UpdatedAt,
		DeletedAt: fileInfo.DeletedAt,
	}
}

func FileInfoPoToEntities(fileInfos []*db.FileInfo) []*domain.FileInfo {
	// 这里需要进行类型转换
	var fileInfoPos []*domain.FileInfo
	for _, fileInfo := range fileInfos {
		fileInfoPos = append(fileInfoPos, FileInfoPoToEntity(fileInfo))
	}
	return fileInfoPos
}
