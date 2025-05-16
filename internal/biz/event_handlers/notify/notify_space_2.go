package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/domain/space_tag"

	"github.com/spf13/cast"
)

func (s *Notify) createSpaceByDomainMessage(e *domain_message.CreateSpace) {
	s.log.Infof("createSpaceByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	spaceMembers, err := s.memberRepo.GetSpaceMemberBySpaceId(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.CreateSpace(&event.CreateSpace{
		Event:    notify.Event_CreateSpace,
		Operator: e.Oper.GetId(),
		Space:    space,
		Members:  spaceMembers,
	})

}

func (s *Notify) modifySpaceByDomainMessage(e *domain_message.ModifySpace) {
	s.log.Infof("modifySpaceByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	for _, v := range e.Updates {
		if v.Field == "describe" {
			s.changeSpaceDescribe(&event.ChangeSpaceDescribe{
				Event:    notify.Event_ChangeSpaceDescribe,
				Space:    space,
				Operator: e.Oper.GetId(),
				OldValue: cast.ToString(v.OldValue),
				NewValue: cast.ToString(v.NewValue),
			})
		}

		if v.Field == "spaceName" {
			s.changeSpaceName(&event.ChangeSpaceName{
				Event:    notify.Event_ChangeSpaceName,
				Space:    space,
				Operator: e.Oper.GetId(),
				OldValue: cast.ToString(v.OldValue),
				NewValue: cast.ToString(v.NewValue),
			})
		}
	}

}

func (s *Notify) setSpaceNotifyByDomainMessage(e *domain_message.SetSpaceNotify) {
	s.log.Infof("setSpaceNotifyByDomainMessage: %+v", e)
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SetSpaceNotify(&event.SetSpaceNotify{
		Event:    notify.Event_SetSpaceNotify,
		Space:    space,
		Operator: e.Oper.GetId(),
		Notify:   int64(e.Notify),
	})

}

func (s *Notify) transferSpaceByDomainMessage(e *domain_message.TransferSpace) {
	s.log.Infof("transferSpaceByDomainMessage: %+v", e)
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.TransferSpace(&event.TransferSpace{
		Event:    notify.Event_TransferSpace,
		Space:    space,
		Operator: e.Oper.GetId(),
		TargetId: e.TargetUserId,
	})
}

func (s *Notify) setSpaceWorkingDayByDomainMessage(e *domain_message.SetWorkingDay) {
	s.log.Infof("SetSpaceWorkingDay: %+v", e)
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SetSpaceWorkingDay(&event.SetSpaceWorkingDay{
		Event:    notify.Event_SetSpaceWorkingDay,
		Space:    space,
		Operator: e.Oper.GetId(),
		WeekDays: e.WeekDays,
	})
}

func (s *Notify) QuitSpaceByDomainMessage(e *domain_message.QuitSpace) {
	s.log.Infof("QuitSpace: %+v", e)
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.QuitSpace(&event.QuitSpace{
		Event:    notify.Event_QuitSpace,
		Space:    space,
		Operator: e.Oper.GetId(),
	})

}

func (s *Notify) addMemberByDomainMessage(e *domain_message.AddSpaceMember) {
	s.log.Infof("addMemberByDomainMessage: %+v", e)
	ctx := context.Background()
	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.addMember(&event.AddMember{
		Event:    notify.Event_AddMember,
		Space:    space,
		Operator: e.Oper.GetId(),
		TargetId: e.MemberUid,
		RoleId:   e.RoleId,
	})

}

func (s *Notify) removeMemberByDomainMessage(e *domain_message.RemoveSpaceMember) {

	ctx := context.Background()

	s.log.Infof("removeMemberByDomainMessage: %+v", e)

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.removeMember(&event.RemoveMember{
		Event:    notify.Event_RemoveMember,
		Space:    space,
		Operator: e.Oper.GetId(),
		TargetId: e.MemberUid,
	})

}

func (s *Notify) changeRoleByDomainMessage(e *domain_message.SetSpaceMemberRole) {
	s.log.Infof("changeRoleByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.changeRole(&event.ChangeRole{
		Event:     notify.Event_ChangeRole,
		Space:     space,
		Operator:  e.Oper.GetId(),
		TargetId:  e.MemeberUid,
		OldRoleId: e.OldRole,
		NewRoleId: e.NewRole,
	})

}

func (s *Notify) setCommentDeletableByDomainMessage(e *domain_message.SetCommentDeletable) {
	s.log.Infof("setCommentDeletableByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SetCommentDeletable(&event.SetCommentDeletable{
		Event:     notify.Event_SetCommentDeletable,
		Space:     space,
		Operator:  e.Oper.GetId(),
		Deletable: e.Deletable,
	})
}
func (s *Notify) setCommentDeletableWhenArchivedByDomainMessage(e *domain_message.SetCommentDeletableWhenArchived) {
	s.log.Infof("setCommentDeletableWhenArchivedByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SetCommentDeletableWhenArchived(&event.SetCommentDeletableWhenArchived{
		Event:    notify.Event_SetCommentDeletableWhenArchived,
		Space:    space,
		Operator: e.Oper.GetId(),
		Value:    e.Value,
	})
}

func (s *Notify) setCommentShowPosByDomainMessage(e *domain_message.SetCommentShowPos) {
	s.log.Infof("setCommentShowPosByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SetCommentShowPos(&event.SetCommentShowPos{
		Event:    notify.Event_SetCommentShowPos,
		Space:    space,
		Operator: e.Oper.GetId(),
		Value:    e.Value,
	})
}
func (s *Notify) deleteSpaceTagByDomainMessage(e *domain_message.DeleteSpaceTag) {
	s.log.Infof("deleteSpaceTagByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.CooperateTag(&event.CooperateTag{
		Event:    notify.Event_DeleteTag,
		Operator: e.Oper.GetId(),
		Space:    space,
		Tag: &space_tag.SpaceTag{
			Id:      e.SpaceTagId,
			TagName: e.SpaceTagName,
		},
	})
}

func (s *Notify) setSpaceMemberRoleByDomainMessage(e *domain_message.SetSpaceMemberRole) {
	s.log.Infof("setSpaceMemberRoleByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	if e.OldRole == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
		s.removeSpaceManager(&event.RemoveSpaceManager{
			Event:    notify.Event_RemoveSpaceManager,
			Space:    space,
			Operator: e.Oper.GetId(),
			TargetId: e.MemeberUid,
		})
		return
	}

	if e.NewRole == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
		s.addSpaceManager(&event.AddSpaceManager{
			Event:    notify.Event_AddSpaceManager,
			Space:    space,
			Operator: e.Oper.GetId(),
			TargetId: e.MemeberUid,
		})
		return
	}

	s.changeRole(&event.ChangeRole{
		Event:     notify.Event_ChangeRole,
		Space:     space,
		Operator:  e.Oper.GetId(),
		TargetId:  e.MemeberUid,
		OldRoleId: e.OldRole,
		NewRoleId: e.NewRole,
	})
}
