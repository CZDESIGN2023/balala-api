package repo

import (
	"context"

	domain "go-cs/internal/domain/file_info"
)

type FileInfoRepo interface {
	CreateFileInfo(ctx context.Context, info *domain.FileInfo) error
	GetFileInfo(context.Context, int64) (*domain.FileInfo, error)
	GetFileInfoByOwner(ctx context.Context, id int64, ownerId int64) (*domain.FileInfo, error)
	GetFileInfoByIds(ctx context.Context, ids []int64) ([]*domain.FileInfo, error)
}
