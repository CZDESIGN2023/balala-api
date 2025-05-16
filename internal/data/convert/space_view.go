package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_view"
)

func SpaceUserViewEntityToPo(entity *domain.SpaceUserView) *db.SpaceUserView {
	// 转换为数据库对象
	return &db.SpaceUserView{
		Id:          entity.Id,
		SpaceId:     entity.SpaceId,
		UserId:      entity.UserId,
		Key:         entity.Key,
		Name:        entity.Name,
		Type:        entity.Type,
		QueryConfig: entity.QueryConfig,
		TableConfig: entity.TableConfig,
		OuterId:     entity.OuterId,
		Status:      entity.Status,
		Ranking:     entity.Ranking,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func SpaceUserViewEntityToPos(entities []*domain.SpaceUserView) []*db.SpaceUserView {
	var po []*db.SpaceUserView
	for _, entity := range entities {
		po = append(po, SpaceUserViewEntityToPo(entity))
	}
	return po
}

func SpaceUserViewPoToEntity(entity *db.SpaceUserView) *domain.SpaceUserView {
	return &domain.SpaceUserView{
		Id:          entity.Id,
		SpaceId:     entity.SpaceId,
		UserId:      entity.UserId,
		Key:         entity.Key,
		Name:        entity.Name,
		Type:        entity.Type,
		QueryConfig: entity.QueryConfig,
		TableConfig: entity.TableConfig,
		OuterId:     entity.OuterId,
		Status:      entity.Status,
		Ranking:     entity.Ranking,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func SpaceUserViewPoToEntities(versions []*db.SpaceUserView) []*domain.SpaceUserView {
	var po []*domain.SpaceUserView
	for _, entity := range versions {
		po = append(po, SpaceUserViewPoToEntity(entity))
	}
	return po
}

func SpaceGlobalViewEntityToPo(entity *domain.SpaceGlobalView) *db.SpaceGlobalView {
	// 转换为数据库对象
	return &db.SpaceGlobalView{
		Id:          entity.Id,
		SpaceId:     entity.SpaceId,
		Key:         entity.Key,
		Name:        entity.Name,
		Type:        entity.Type,
		QueryConfig: entity.QueryConfig,
		TableConfig: entity.TableConfig,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func SpaceGlobalViewEntityToPos(entities []*domain.SpaceGlobalView) []*db.SpaceGlobalView {
	var po []*db.SpaceGlobalView
	for _, entity := range entities {
		po = append(po, SpaceGlobalViewEntityToPo(entity))
	}
	return po
}

func SpaceGlobalViewPoToEntity(entity *db.SpaceGlobalView) *domain.SpaceGlobalView {
	return &domain.SpaceGlobalView{
		Id:          entity.Id,
		SpaceId:     entity.SpaceId,
		Key:         entity.Key,
		Name:        entity.Name,
		Type:        entity.Type,
		QueryConfig: entity.QueryConfig,
		TableConfig: entity.TableConfig,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func SpaceGlobalViewPoToEntities(versions []*db.SpaceGlobalView) []*domain.SpaceGlobalView {
	var po []*domain.SpaceGlobalView
	for _, entity := range versions {
		po = append(po, SpaceGlobalViewPoToEntity(entity))
	}
	return po
}
