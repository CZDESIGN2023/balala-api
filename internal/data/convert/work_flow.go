package convert

import (
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_flow"
)

func WorkFlowEntityToPo(ent *domain.WorkFlow) *db.WorkFlow {
	return &db.WorkFlow{
		Id:             ent.Id,
		Uuid:           ent.Uuid,
		UserId:         ent.UserId,
		SpaceId:        ent.SpaceId,
		WorkItemTypeId: ent.WorkItemTypeId,
		Name:           ent.Name,
		Key:            ent.Key,
		FlowMode:       string(ent.FlowMode),
		Version:        ent.Version,
		LastTemplateId: ent.LastTemplateId,
		Status:         int64(ent.Status),
		Ranking:        ent.Ranking,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
		DeletedAt:      ent.DeletedAt,
		IsSys:          ent.IsSys,
	}
}

func WorkFlowPoToEntity(po *db.WorkFlow) *domain.WorkFlow {
	return &domain.WorkFlow{
		Id:             po.Id,
		Uuid:           po.Uuid,
		UserId:         po.UserId,
		SpaceId:        po.SpaceId,
		WorkItemTypeId: po.WorkItemTypeId,
		Name:           po.Name,
		Key:            po.Key,
		FlowMode:       consts.WorkFlowMode(po.FlowMode),
		Version:        po.Version,
		LastTemplateId: po.LastTemplateId,
		Status:         domain.WorkFlowStatus(po.Status),
		Ranking:        po.Ranking,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
		DeletedAt:      po.DeletedAt,
		IsSys:          po.IsSys,
	}
}

func WorkFlowPoToEntities(po []*db.WorkFlow) []*domain.WorkFlow {
	entities := make([]*domain.WorkFlow, 0)
	for _, v := range po {
		entities = append(entities, WorkFlowPoToEntity(v))
	}
	return entities
}
