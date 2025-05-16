package convert

import (
	db "go-cs/internal/bean/biz"
	domain "go-cs/internal/domain/space_member"
)

func SpaceMemberEntityToPo(ent *domain.SpaceMember) *db.SpaceMember {
	return &db.SpaceMember{
		Id:            ent.Id,
		UserId:        ent.UserId,
		SpaceId:       ent.SpaceId,
		RoleId:        ent.RoleId,
		Remark:        ent.Remark,
		Notify:        ent.Notify,
		Ranking:       ent.Ranking,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
		DeletedAt:     ent.DeletedAt,
		HistoryRoleId: ent.HistoryRoleId,
	}
}

func SpaceMemberEntitiesToPo(ents []*domain.SpaceMember) []*db.SpaceMember {
	var list []*db.SpaceMember
	for _, v := range ents {
		list = append(list, SpaceMemberEntityToPo(v))
	}
	return list
}

func SpaceMemberPoToEntity(po *db.SpaceMember) *domain.SpaceMember {
	return &domain.SpaceMember{
		Id:            po.Id,
		UserId:        po.UserId,
		SpaceId:       po.SpaceId,
		RoleId:        po.RoleId,
		Remark:        po.Remark,
		Notify:        po.Notify,
		Ranking:       po.Ranking,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
		DeletedAt:     po.DeletedAt,
		HistoryRoleId: po.HistoryRoleId,
	}
}

func SpaceMemberPosToEntity(pos []*db.SpaceMember) []*domain.SpaceMember {
	var list []*domain.SpaceMember
	for _, v := range pos {
		list = append(list, SpaceMemberPoToEntity(v))
	}
	return list
}
