package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space"
)

func SpaceEntityToPo(ent *domain.Space) *db.Space {
	return &db.Space{
		Id:        ent.Id,
		UserId:    ent.UserId,
		SpaceGuid: ent.SpaceGuid,
		SpaceName: ent.SpaceName,
		Remark:    ent.Remark,
		Describe:  ent.Describe,
		Notify:    ent.Notify,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		DeletedAt: ent.DeletedAt,
	}
}

func SpacePoToEntity(po *db.Space) *domain.Space {
	return &domain.Space{
		Id:        po.Id,
		UserId:    po.UserId,
		SpaceGuid: po.SpaceGuid,
		SpaceName: po.SpaceName,
		Remark:    po.Remark,
		Describe:  po.Describe,
		Notify:    po.Notify,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func SpaceConfigEntityToPo(ent *domain.SpaceConfig) *db.SpaceConfig {
	return &db.SpaceConfig{
		Id:                           ent.Id,
		SpaceId:                      ent.SpaceId,
		WorkingDay:                   string(ent.WorkingDay),
		CommentDeletable:             int32(ent.CommentDeletable),
		CommentDeletableWhenArchived: int32(ent.CommentDeletableWhenArchived),
		CommentShowPos:               int32(ent.CommentShowPos),
		CreatedAt:                    ent.CreatedAt,
		UpdatedAt:                    ent.UpdatedAt,
		DeletedAt:                    ent.DeletedAt,
	}
}

func SpaceConfigPoToEntity(po *db.SpaceConfig) *domain.SpaceConfig {
	return &domain.SpaceConfig{
		Id:                           po.Id,
		SpaceId:                      po.SpaceId,
		WorkingDay:                   domain.WorkingDay(po.WorkingDay),
		CommentDeletable:             int64(po.CommentDeletable),
		CommentDeletableWhenArchived: int64(po.CommentDeletableWhenArchived),
		CommentShowPos:               int64(po.CommentShowPos),
		CreatedAt:                    po.CreatedAt,
		UpdatedAt:                    po.UpdatedAt,
	}
}
