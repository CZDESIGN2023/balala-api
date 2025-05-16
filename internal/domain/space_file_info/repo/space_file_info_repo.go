package repo

import (
	"context"
	domain "go-cs/internal/domain/space_file_info"
)

type SpaceFileInfoRepo interface {
	SpaceFileInfoQueryRepo

	CreateSpaceFileInfo(ctx context.Context, in *domain.SpaceFileInfo) error

	GetSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) (*domain.SpaceFileInfo, error)
	GetSpaceWorkItemFileInfoById(ctx context.Context, id int64) (*domain.SpaceFileInfo, error)
	GetSpaceWorkItemFileInfoByWorkItemIds(ctx context.Context, spaceId int64, workItemIds []int64) (domain.SpaceFileInfos, error)

	SoftDelSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) error
	HardDelSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) error
	SoftDelSpaceWorkItemsAllFile(ctx context.Context, workItemIds []int64) error
	SoftDelFileBySpaceId(ctx context.Context, spaceId int64) error

	CountWorkItemFileNum(ctx context.Context, workItemId int64) (int64, error)

	SaveFileDownToken(ctx context.Context, token string, fileInfoId int64, expSecond int64) error
	GetFileDownToken(ctx context.Context, token string) (string, error)
}

type SpaceFileInfoQueryRepo interface {
	QFileInfoList(ctx context.Context, spaceId int64, workItemIds int64) ([]*domain.SpaceFileInfo, error)
}
