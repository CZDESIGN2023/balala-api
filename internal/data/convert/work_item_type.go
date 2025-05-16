package convert

import (
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_type"
)

func WorkItemTypeEntityToPo(item *domain.WorkItemType) *db.WorkItemType {
	return &db.WorkItemType{
		Id:        item.Id,
		Uuid:      item.Uuid,
		UserId:    item.UserId,
		SpaceId:   int64(item.SpaceId),
		Name:      item.Name,
		Key:       item.Key,
		Ranking:   item.Ranking,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		DeletedAt: item.DeletedAt,
		FlowMode:  string(item.FlowMode),
		IsSys:     item.IsSys,
	}
}

func WorkItemTypePoToEntity(item *db.WorkItemType) *domain.WorkItemType {
	return &domain.WorkItemType{
		Id:        item.Id,
		Uuid:      item.Uuid,
		UserId:    item.UserId,
		SpaceId:   domain.SpaceId(item.SpaceId),
		Name:      item.Name,
		Key:       item.Key,
		Ranking:   item.Ranking,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		FlowMode:  consts.WorkFlowMode(item.FlowMode),
		IsSys:     item.IsSys,
	}
}

func WorkItemTypePoToEntities(item []*db.WorkItemType) []*domain.WorkItemType {
	list := make([]*domain.WorkItemType, 0, len(item))
	for _, v := range item {
		list = append(list, WorkItemTypePoToEntity(v))
	}
	return list
}
