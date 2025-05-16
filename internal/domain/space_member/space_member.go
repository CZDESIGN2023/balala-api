package space_member

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

type SpaceMembers []*SpaceMember

type SpaceMember struct {
	shared.AggregateRoot

	Id            int64  `json:"id"`
	UserId        int64  `json:"user_id"`
	SpaceId       int64  `json:"space_id"`
	RoleId        int64  `json:"role_id"`
	Remark        string `json:"remark"`
	Ranking       int64  `json:"ranking"`
	Notify        int32  `json:"notify"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	DeletedAt     int64  `json:"deleted_at"`
	HistoryRoleId int64  `json:"history_role_id"`
}

func (s *SpaceMembers) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)
	for _, v := range *s {
		logs = append(logs, v.GetMessages()...)
	}
	return logs
}

func (s *SpaceMember) ChangeRoleId(roleId int64, oper shared.Oper) {
	if s.RoleId == roleId {
		return
	}

	oldValue := s.RoleId
	s.HistoryRoleId = s.RoleId
	s.RoleId = roleId

	s.AddDiff(Diff_RoleId)

	s.AddMessage(oper, &domain_message.SetSpaceMemberRole{
		SpaceId:    s.SpaceId,
		MemeberUid: s.UserId,
		OldRole:    int32(oldValue),
		NewRole:    int32(roleId),
	})

	// 历史角色不能记录创建者，否则在移除项目管理员时，会还原为创建者
	if s.HistoryRoleId == consts.MEMBER_ROLE_SPACE_CREATOR {
		s.HistoryRoleId = consts.MEMBER_ROLE_MANAGER
	}
}

func (s *SpaceMember) SetSpaceManagerRole(oper shared.Oper) {

	roleId := int64(consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER)
	if s.RoleId == roleId {
		return
	}

	s.ChangeRoleId(roleId, oper)
}

func (s *SpaceMember) RemoveSupperManagerRole(oper shared.Oper) {

	if !s.IsSpaceSupperManager() {
		return
	}

	s.ChangeRoleId(s.HistoryRoleId, oper)
}

func (s *SpaceMember) SetNotify(notify int32, oper shared.Oper) {
	s.Notify = notify

	s.AddDiff(Diff_Notify)

	s.AddMessage(oper, &domain_message.PersonalSetSpaceNotify{
		SpaceId: s.SpaceId,
		Notify:  int(notify),
	})
}

func (s *SpaceMember) IsCurrentMemberUser(userId int64) bool {
	return s.UserId == userId
}

func (s *SpaceMember) IsSpaceManager() bool {
	return s.RoleId == consts.MEMBER_ROLE_MANAGER
}

func (s *SpaceMember) IsSpaceSupperManager() bool {
	return s.RoleId == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER
}

func (s *SpaceMember) IsSpaceCreator() bool {
	return s.RoleId == consts.MEMBER_ROLE_SPACE_CREATOR
}

func (s *SpaceMember) IsSpaceEditor() bool {
	return s.RoleId == consts.MEMBER_ROLE_EDITOR
}

func (s *SpaceMember) IsSpaceWatcher() bool {
	return s.RoleId == consts.MEMBER_ROLE_WATCHER
}

func (s *SpaceMember) GetRole() int64 {
	if s.RoleId == 0 {
		return consts.MEMBER_ROLE_WATCHER
	}
	return s.RoleId
}

func (s *SpaceMember) RoleGreaterThan(otherSpaceMember *SpaceMember) bool {
	return s.GetRole() > otherSpaceMember.GetRole()
}

func (s *SpaceMember) OnQuit(workItemNum int, transferTargetId int64, oper shared.Oper) {
	s.AddMessage(oper, &domain_message.QuitSpace{
		SpaceId:          s.SpaceId,
		MemberUid:        s.UserId,
		WorkItemNum:      workItemNum,
		TransferTargetId: transferTargetId,
	})
}

func (s *SpaceMember) OnKickOut(workItemNum int, transferTargetId int64, oper shared.Oper) {

	s.AddMessage(oper, &domain_message.RemoveSpaceMember{
		SpaceId:          s.SpaceId,
		MemberUid:        s.UserId,
		WorkItemNum:      workItemNum,
		TransferTargetId: transferTargetId,
	})
}
