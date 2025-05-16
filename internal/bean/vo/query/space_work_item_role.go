package query

import "go-cs/internal/consts"

type SpaceWorkItemRoleQuery struct {
	SpaceId   int64
	FlowScope consts.FlowScope
}

type SpaceWorkItemRoleQueryResult struct {
	Total int64
	List  []*SpaceWorkItemRoleQueryResult_ListItem // 列表
}

type SpaceWorkItemRoleQueryResult_ListItem struct {
	Id             int64
	Uuid           string
	SpaceId        int64
	WorkItemTypeId int64
	Key            string
	Name           string
	Status         int64
	Ranking        int64
	CreatedAt      int64
	UpdatedAt      int64
	IsSys          int64
	FlowScope      string
}
