package repo

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/domain/search/search2"
	esV8 "go-cs/internal/utils/es/v8"
)

type SearchRepo interface {
	QueryWorkItem(ctx context.Context, spaceIds []int64, group *search2.ConditionGroup, selectFields string) ([]*search2.Model, error)
	PendingWorkItem(ctx context.Context, userId int64, spaceIds []int64) ([]int64, error)
	QueryWorkItemEs(ctx context.Context, query *esV8.SearchSource) (*esV8.SearchResult, error)
	QueryWorkItemEsByPid(ctx context.Context, pid []int64) ([]*search2.Model, error)
	SearchByName(ctx context.Context, spaceId int64, keyword string) ([]*search2.Model, error)
	QueryWorkFlowNode(ctx context.Context, spaceIds []int64, group *search2.ConditionGroup) ([]*db.SpaceWorkItemFlowV2, error)
	GetWorkItemIdsByQueryWorkFlowNode(ctx context.Context, spaceIds []int64, group *search2.ConditionGroup) ([]int64, error)
}
