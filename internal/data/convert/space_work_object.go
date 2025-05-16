package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_work_object"
)

func SpaceWorkObjectEntityToPo(workObj *domain.SpaceWorkObject) *db.SpaceWorkObject {
	// 转换逻辑
	return &db.SpaceWorkObject{
		Id:               workObj.Id,
		SpaceId:          workObj.SpaceId,
		UserId:           workObj.UserId,
		WorkObjectGuid:   workObj.WorkObjectGuid,
		WorkObjectName:   workObj.WorkObjectName,
		WorkObjectStatus: int32(workObj.WorkObjectStatus),
		Remark:           workObj.Remark,
		Describe:         workObj.Describe,
		Ranking:          workObj.Ranking,
		CreatedAt:        workObj.CreatedAt,
		UpdatedAt:        workObj.UpdatedAt,
		DeletedAt:        workObj.DeletedAt,
	}
}

func SpaceWorkObjectPoToEntitys(workObjs []*db.SpaceWorkObject) []*domain.SpaceWorkObject {
	// 转换逻辑
	var entities []*domain.SpaceWorkObject
	for _, workObj := range workObjs {
		entities = append(entities, SpaceWorkObjectPoToEntity(workObj))
	}
	return entities
}

func SpaceWorkObjectPoToEntity(workObj *db.SpaceWorkObject) *domain.SpaceWorkObject {
	// 转换逻辑
	return &domain.SpaceWorkObject{
		Id:               workObj.Id,
		SpaceId:          workObj.SpaceId,
		UserId:           workObj.UserId,
		WorkObjectGuid:   workObj.WorkObjectGuid,
		WorkObjectName:   workObj.WorkObjectName,
		WorkObjectStatus: domain.WorkObjectStatus(workObj.WorkObjectStatus),
		Remark:           workObj.Remark,
		Describe:         workObj.Describe,
		Ranking:          workObj.Ranking,
		CreatedAt:        workObj.CreatedAt,
		UpdatedAt:        workObj.UpdatedAt,
		DeletedAt:        workObj.DeletedAt,
	}
}

func SpaceWorkObjectEntityToPos(workObjs []*domain.SpaceWorkObject) []*db.SpaceWorkObject {
	// 转换逻辑
	var pos []*db.SpaceWorkObject
	for _, workObj := range workObjs {
		pos = append(pos, SpaceWorkObjectEntityToPo(workObj))
	}
	return pos
}
