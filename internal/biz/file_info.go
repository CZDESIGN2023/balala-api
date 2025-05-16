package biz

import (
	"context"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"

	domain "go-cs/internal/domain/file_info"
	"go-cs/internal/domain/file_info/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type FileInfoUsecase struct {
	repo repo.FileInfoRepo

	log *log.Helper
	tm  trans.Transaction
}

func NewFileInfoUsecase(repo repo.FileInfoRepo, tm trans.Transaction, logger log.Logger) *FileInfoUsecase {
	moduleName := "FileInfoUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &FileInfoUsecase{
		repo: repo,
		log:  hlog,
		tm:   tm,
	}
}

func (uc *FileInfoUsecase) GetFileInfo(ctx context.Context, id int64) (*domain.FileInfo, error) {
	return uc.repo.GetFileInfo(ctx, id)
}
