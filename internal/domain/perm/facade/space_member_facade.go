package facade

import "go-cs/internal/domain/space_member"

type SpaceMemberFacade struct {
	spaceMember *space_member.SpaceMember
}

func (s *SpaceMemberFacade) GetSpaceMember() *space_member.SpaceMember {
	return s.spaceMember
}

func BuildSpaceMemberFacade(spaceMember *space_member.SpaceMember) *SpaceMemberFacade {
	return &SpaceMemberFacade{
		spaceMember: spaceMember,
	}
}
