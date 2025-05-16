package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/notify_snapshot"
)

func NotifySnapShotEntityToPo(ent *domain.NotifySnapShot) *db.Notify {
	return &db.Notify{
		Id:        ent.Id,
		UserId:    ent.UserId,
		SpaceId:   ent.SpaceId,
		Typ:       int64(ent.Typ),
		Doc:       ent.Doc,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		DeletedAt: ent.DeletedAt,
	}
}

func NotifySnapShotEntityToPos(ents []*domain.NotifySnapShot) []*db.Notify {
	var list []*db.Notify
	for _, v := range ents {
		list = append(list, NotifySnapShotEntityToPo(v))
	}
	return list
}

func NotifySnapShotPoToEntity(po *db.Notify) *domain.NotifySnapShot {
	return &domain.NotifySnapShot{
		Id:        po.Id,
		UserId:    po.UserId,
		SpaceId:   po.SpaceId,
		Typ:       domain.NotifyEventType(po.Typ),
		Doc:       po.Doc,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		DeletedAt: po.DeletedAt,
	}
}

func NotifySnapShotPoToEntitys(pos []*db.Notify) []*domain.NotifySnapShot {
	var list []*domain.NotifySnapShot
	for _, v := range pos {
		list = append(list, NotifySnapShotPoToEntity(v))
	}
	return list
}
