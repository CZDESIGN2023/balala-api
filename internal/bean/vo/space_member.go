package vo

type AddSpaceMembersVo struct {
	SpaceId int64

	MemberUserIds     []int64
	MemberUserRoleIds []int64
}

func (s *AddSpaceMembersVo) AddMemberUser(userId int64, roleId int64) {
	s.MemberUserIds = append(s.MemberUserIds, userId)
	s.MemberUserRoleIds = append(s.MemberUserRoleIds, roleId)
}
