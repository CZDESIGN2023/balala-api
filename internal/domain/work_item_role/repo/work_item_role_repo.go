package repo

import (
	"context"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_role"
)

type WorkItemRoleRepo interface {
	WorkItemRoleQueryRepo

	CreateWorkItemRoles(ctx context.Context, roles domain.WorkItemRoles) error
	CreateWorkItemRole(ctx context.Context, role *domain.WorkItemRole) error
	SaveWorkItemRole(ctx context.Context, role *domain.WorkItemRole) error

	GetWorkItemRole(ctx context.Context, id int64) (*domain.WorkItemRole, error)
	GetWorkItemRoles(ctx context.Context, spaceId int64) (domain.WorkItemRoles, error)
	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)
	IsExistByRoleName(ctx context.Context, spaceId int64, roleName string, scope consts.FlowScope) (bool, error)

	WorkItemRoleMap(ctx context.Context, spaceId int64) (map[int64]*domain.WorkItemRole, error)

	DelWorkItemRoleBySpaceId(ctx context.Context, spaceId int64) error
	DelWorkItemRole(ctx context.Context, id int64) error
}

type WorkItemRoleQueryRepo interface {
	QSpaceWorkItemRoleList(ctx context.Context, req *query.SpaceWorkItemRoleQuery) (*query.SpaceWorkItemRoleQueryResult, error)
}
