package service

import (
	"context"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/file_info"
	repo "go-cs/internal/domain/file_info/repo"
	"go-cs/internal/pkg/biz_id"
	"go-cs/internal/utils/errs"
)

type FileInfoService struct {
	repo      repo.FileInfoRepo
	idService *biz_id.BusinessIdService
}

func NewFileInfoService(
	repo repo.FileInfoRepo,
	idService *biz_id.BusinessIdService,

) *FileInfoService {
	return &FileInfoService{
		repo:      repo,
		idService: idService,
	}
}

type CreateFileInfoRequest struct {
	Hash         string
	Name         string
	Typ          int32
	Size         int64
	Uri          string
	Cover        string
	Owner        int64
	Meta         map[string]string // 自定义元数据
	UploadDomain string
	UploadMd5    string
	UploadPath   string
}

func (s *FileInfoService) CreateFileInfo(ctx context.Context, req *CreateFileInfoRequest) (*domain.FileInfo, error) {

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_FileInfo)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成文件ID失败")
	}

	return domain.NewFileInfo(
		bizId.Id,
		req.Hash,
		req.Name,
		req.Typ,
		req.Size,
		req.Uri,
		req.Cover,
		req.Owner,
		req.Meta,
		req.UploadDomain,
		req.UploadMd5,
		req.UploadPath,
	), nil
}

func (s *FileInfoService) GetFileInfos(ctx context.Context, ids []int64) (domain.FileInfos, error) {
	return s.repo.GetFileInfoByIds(ctx, ids)
}
