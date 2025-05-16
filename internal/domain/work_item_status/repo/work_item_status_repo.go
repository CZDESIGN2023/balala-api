package repo

import (
	"context"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_status"
)

type WorkItemStatusRepo interface {
	WorkItemStatusQueryRepo

	GetWorkItemStatusInfo(ctx context.Context, spaceId int64) (*domain.WorkItemStatusInfo, error)
	GetWorkItemStatusItem(ctx context.Context, id int64) (*domain.WorkItemStatusItem, error)
	GetWorkItemStatusItemsBySpace(ctx context.Context, spaceId int64) (domain.WorkItemStatusItems, error)
	GetWorkItemStatusItemsBySpaceIds(ctx context.Context, spaceId []int64) (domain.WorkItemStatusItems, error)
	GetWorkItemStatusInfoBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.WorkItemStatusInfo, error)
	WorkItemStatusMap(ctx context.Context, spaceId int64) (map[int64]*domain.WorkItemStatusItem, error)
	WorkItemStatusKeyMap(ctx context.Context, spaceId int64, keys ...string) (map[string]*domain.WorkItemStatusItem, error)

	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)

	CreateWorkItemStatusItems(ctx context.Context, spaceId int64, items []*domain.WorkItemStatusItem) error
	CreateWorkItemStatusItem(ctx context.Context, item *domain.WorkItemStatusItem) error
	SaveWorkItemStatusItem(ctx context.Context, item *domain.WorkItemStatusItem) error

	DelWorkItemStatusBySpaceId(ctx context.Context, spaceId int64) error
	DelSpaceWorkItemStatusItem(ctx context.Context, spaceId int64, id int64) error
	StatusMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkItemStatusItem, error)
	StatusMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64]*domain.WorkItemStatusItem, error)

	IsExistByWorkItemStatusName(ctx context.Context, spaceId int64, name string, flowScope consts.FlowScope) (bool, error)
}

type WorkItemStatusQueryRepo interface {
	QSpaceWorkItemStatusList(ctx context.Context, spaceIds []int64, scope consts.FlowScope) ([]*domain.WorkItemStatusItem, error)
	QSpaceWorkItemStatusById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.WorkItemStatusItem, error)
}
