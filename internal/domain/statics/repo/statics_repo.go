package repo

import (
	"context"
	"go-cs/internal/bean/vo"
	esV8 "go-cs/internal/utils/es/v8"
)

type StaticsRepo interface {
	GetWorkbenchCount(ctx context.Context, userId int64, spaceIds []int64) (*vo.UserWorkbenchCountInfo, error)
	GetSpaceWorkbenchCount(ctx context.Context, uid int64, spaceId int64) (*vo.SpaceWorkbenchCountInfo, error)
	GetSpaceWorkObjectCountByIds(ctx context.Context, spaceId int64, workObjectIds []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkObjectCountInfo, error)
	GetUserFollowCount(ctx context.Context, userId int64, spaceIds []int64) (int64, error)
	GetSpaceWorkVersionCountByIds(ctx context.Context, spaceId int64, versionIds []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkVersionCountInfo, error)
	GetSpaceUserCount(ctx context.Context, spaceId int64, startTime, endTime int64) (map[int64]*vo.SpaceUserCountInfo, error)

	GetWorkItemCountBySpaceFlowId(ctx context.Context, spaceId int64, workFlowId int64) (int64, error)

	MultiCountByEs(ctx context.Context, sources []*esV8.SearchSource) ([]vo.CountInfo, error)
}
