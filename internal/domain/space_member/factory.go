package space_member

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

func NewSpaceMember(
	spaceId int64,
	userId int64,
	roleId int64,
	oper shared.Oper,
) *SpaceMember {
	m := &SpaceMember{
		SpaceId:       spaceId,
		UserId:        userId,
		RoleId:        roleId,
		Ranking:       time.Now().UnixMilli(),
		Notify:        1,
		CreatedAt:     time.Now().Unix(),
		HistoryRoleId: roleId,
	}

	m.AddMessage(oper, &domain_message.AddSpaceMember{
		SpaceId:   spaceId,
		RoleId:    roleId,
		MemberUid: userId,
	})

	return m
}
