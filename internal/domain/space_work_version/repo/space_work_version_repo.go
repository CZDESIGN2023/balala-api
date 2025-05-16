package repo

import (
	"context"
	domain "go-cs/internal/domain/space_work_version"
)

type SpaceWorkVersionRepo interface {
	SpaceWorkVersionQueryRepo

	CreateSpaceWorkVersion(ctx context.Context, workVersion *domain.SpaceWorkVersion) error
	SaveSpaceWorkVersion(ctx context.Context, workVersion *domain.SpaceWorkVersion) error

	DelWorkVersion(ctx context.Context, workVersionId int64) (int64, error)
	DelWorkVersionBySpaceId(ctx context.Context, spaceId int64) (int64, error)

	GetSpaceWorkVersion(ctx context.Context, workVersionId int64) (*domain.SpaceWorkVersion, error)
	GetSpaceWorkVersionBySpaceId(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkVersion, error)
	GetSpaceWorkVersionBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.SpaceWorkVersion, error)
	GetSpaceWorkVersionByIds(ctx context.Context, ids []int64) ([]*domain.SpaceWorkVersion, error)
	GetSpaceWorkVersionByKey(ctx context.Context, spaceId int64, workVersionKey string) (*domain.SpaceWorkVersion, error)
	SpaceWorkVersionMap(ctx context.Context, spaceId int64) (map[int64]*domain.SpaceWorkVersion, error)
	SpaceWorkVersionMapByVersionIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkVersion, error)

	CheckSpaceWorkVersionName(ctx context.Context, spaceID int64, workVersionName string) (bool, error)
	GetSpaceWorkVersionCount(ctx context.Context, spaceId int64) (int64, error)
	IsEmpty(ctx context.Context, id int64) (bool, error)
	SetOrder(ctx context.Context, spaceId int64, fromIdx, toIdx int64) error
	GetVersionRelationCount(ctx context.Context, spaceId int64, workVersionId int64) (int64, error)

	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)
}

type SpaceWorkVersionQueryRepo interface {
	QSpaceWorkVersionList(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkVersion, error)
	QSpaceWorkVersionById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.SpaceWorkVersion, error)
}
