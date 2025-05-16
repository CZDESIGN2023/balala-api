package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_tag"
)

func SpaceTagEntityToPo(spaceTag *domain.SpaceTag) *db.SpaceTag {
	return &db.SpaceTag{
		Id:        spaceTag.Id,
		SpaceId:   spaceTag.SpaceId,
		TagGuid:   spaceTag.TagGuid,
		TagName:   spaceTag.TagName,
		TagStatus: int32(spaceTag.TagStatus),
		CreatedAt: spaceTag.CreatedAt,
		UpdatedAt: spaceTag.UpdatedAt,
		DeletedAt: spaceTag.DeletedAt,
	}
}

func SpaceTagEntityToPos(spaceTags []*domain.SpaceTag) []*db.SpaceTag {
	// 转换为 db.SpaceTag 类型
	var spaceTagPos []*db.SpaceTag
	for _, spaceTag := range spaceTags {
		spaceTagPos = append(spaceTagPos, SpaceTagEntityToPo(spaceTag))
	}
	return spaceTagPos
}

func SpaceTagPoToEntity(spaceTag *db.SpaceTag) *domain.SpaceTag {
	return &domain.SpaceTag{
		Id:        spaceTag.Id,
		SpaceId:   spaceTag.SpaceId,
		TagGuid:   spaceTag.TagGuid,
		TagName:   spaceTag.TagName,
		TagStatus: domain.TagStatus(spaceTag.TagStatus),
		CreatedAt: spaceTag.CreatedAt,
		UpdatedAt: spaceTag.UpdatedAt,
		DeletedAt: spaceTag.DeletedAt,
	}
}

func SpaceTagPoToEntitys(spaceTags []*db.SpaceTag) []*domain.SpaceTag {
	// 转换为 domain.SpaceTag 类型
	var spaceTagEntities []*domain.SpaceTag
	for _, spaceTag := range spaceTags {
		spaceTagEntities = append(spaceTagEntities, SpaceTagPoToEntity(spaceTag))
	}
	return spaceTagEntities
}
