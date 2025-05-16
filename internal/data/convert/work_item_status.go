package convert

import (
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_status"
)

func WorkItemStatusEntityToPo(item *domain.WorkItemStatusItem) *db.WorkItemStatus {
	return &db.WorkItemStatus{
		Id:             item.Id,
		SpaceId:        item.SpaceId,
		Uuid:           item.Uuid,
		UserId:         item.UserId,
		WorkItemTypeId: item.WorkItemTypeId,
		Name:           item.Name,
		Key:            item.Key,
		Val:            item.Val,
		StatusType:     int32(item.StatusType),
		Ranking:        item.Ranking,
		Status:         item.Status,
		IsSys:          item.IsSys,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		DeletedAt:      item.DeletedAt,
		FlowScope:      string(item.FlowScope),
	}
}

func WorkItemStatusPoToEntity(item *db.WorkItemStatus) *domain.WorkItemStatusItem {
	return &domain.WorkItemStatusItem{
		Id:             item.Id,
		UserId:         item.UserId,
		SpaceId:        item.SpaceId,
		Uuid:           item.Uuid,
		WorkItemTypeId: item.WorkItemTypeId,
		Name:           item.Name,
		Key:            item.Key,
		Val:            item.Val,
		StatusType:     consts.WorkItemStatusType(item.StatusType),
		Ranking:        item.Ranking,
		IsSys:          item.IsSys,
		Status:         item.Status,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		FlowScope:      consts.FlowScope(item.FlowScope),
	}
}

func WorkItemStatusPoToEntitys(items []*db.WorkItemStatus) []*domain.WorkItemStatusItem {
	var result []*domain.WorkItemStatusItem
	for _, item := range items {
		result = append(result, WorkItemStatusPoToEntity(item))
	}
	return result
}
