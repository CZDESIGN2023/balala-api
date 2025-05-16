package repo

import (
	"context"
	domain "go-cs/internal/domain/space_work_object"
)

type SpaceWorkObjectRepo interface {
	SpaceWorkObjectQueryRepo

	CreateSpaceWorkObject(ctx context.Context, workObject *domain.SpaceWorkObject) error
	SaveSpaceWorkObject(ctx context.Context, workObject *domain.SpaceWorkObject) error

	GetSpaceWorkObject(ctx context.Context, spaceId int64, workObjectId int64) (*domain.SpaceWorkObject, error)
	GetSpaceWorkObjectCount(ctx context.Context, spaceId int64) (int64, error)
	GetSpaceWorkObjectByIds(ctx context.Context, ids []int64) ([]*domain.SpaceWorkObject, error)
	GetSpaceWorkObjectBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.SpaceWorkObject, error)

	SpaceWorkObjectMap(ctx context.Context, spaceId int64) (map[int64]*domain.SpaceWorkObject, error)
	SpaceWorkObjectMapByObjectIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkObject, error)

	DelWorkObject(ctx context.Context, workObjectId int64) (int64, error)
	DelWorkObjectBySpaceId(ctx context.Context, spaceId int64) (int64, error)

	CheckSpaceWorkObjectName(ctx context.Context, spaceID int64, workObjectName string) (bool, error)
	IsEmpty(ctx context.Context, id int64) (bool, error)
	SetOrder(ctx context.Context, spaceId int64, fromIdx, toIdx int64) error

	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)
}

type SpaceWorkObjectQueryRepo interface {
	QSpaceWorkObjectList(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkObject, error)
	QSpaceWorkObjectById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.SpaceWorkObject, error)
}
