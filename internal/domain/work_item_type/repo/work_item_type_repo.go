package repo

import (
	"context"
	"go-cs/internal/bean/vo/query"
	domain "go-cs/internal/domain/work_item_type"
)

type WorkItemTypeRepo interface {
	WorkItemTypeQueryRepo

	CreateWorkItemType(ctx context.Context, itemType *domain.WorkItemType) error
	GetWorkItemType(ctx context.Context, id int64) (*domain.WorkItemType, error)
	GetWorkItemTypeBySpaceId(ctx context.Context, spaceId int64) ([]*domain.WorkItemType, error)
	WorkItemTypeMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkItemType, error)
	DelBySpaceId(ctx context.Context, spaceId int64) error

	IsExistByName(ctx context.Context, spaceId int64, name string) (bool, error)
	IsExistByKey(ctx context.Context, spaceId int64, key string) (bool, error)
}

type WorkItemTypeQueryRepo interface {
	QWorkItemTypeInfo(ctx context.Context, req query.WorkItemTypeInfoQuery) (*query.WorkItemTypeInfoQueryResult, error)
}
