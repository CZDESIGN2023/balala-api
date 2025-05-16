package convert

import (
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_role"
)

func WorkItemRoleEntityToPo(role *domain.WorkItemRole) *db.WorkItemRole {
	return &db.WorkItemRole{
		Id:             role.Id,
		Uuid:           role.Uuid,
		UserId:         role.UserId,
		SpaceId:        role.SpaceId,
		WorkItemTypeId: role.WorkItemTypeId,
		Key:            role.Key,
		Name:           role.Name,
		Status:         role.Status,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
		DeletedAt:      role.DeletedAt,
		Ranking:        role.Ranking,
		IsSys:          role.IsSys,
		FlowScope:      string(role.FlowScope),
	}
}

func WorkItemRolePoToEntity(role *db.WorkItemRole) *domain.WorkItemRole {
	return &domain.WorkItemRole{
		Id:             role.Id,
		Uuid:           role.Uuid,
		UserId:         role.UserId,
		SpaceId:        role.SpaceId,
		WorkItemTypeId: role.WorkItemTypeId,
		Key:            role.Key,
		Name:           role.Name,
		Status:         role.Status,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
		DeletedAt:      role.DeletedAt,
		Ranking:        role.Ranking,
		IsSys:          role.IsSys,
		FlowScope:      consts.FlowScope(role.FlowScope),
	}
}
