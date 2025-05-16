package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/domain/work_item"
)

type CooperateWorkItemEndedEvData struct {
	SpaceId    int64
	OperatorId int64
	WorkItem   *work_item.WorkItem
}

func (s *Notify) CooperateWorkItemStatusChange(evType notify.Event, evData *CooperateWorkItemEndedEvData) {
	ctx := context.Background()

	spaceId := evData.SpaceId
	operatorId := evData.OperatorId
	workItem := evData.WorkItem

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": evType,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"workItem": map[string]any{
				"id":   workItem.Id,
				"pid":  workItem.Pid,
				"name": workItem.WorkItemName,
			},
		},
	}, spaceId)

}

func (s *Notify) CooperateDeleteWorkItem(e *event.DeleteWorkItem) {
	ctx := context.Background()

	spaceId := e.Space.Id
	operatorId := e.Operator
	workItem := e.WorkItem

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": notify.Event_DeleteWorkItem,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"workItem": map[string]any{
				"id":   workItem.Id,
				"pid":  workItem.Pid,
				"name": workItem.WorkItemName,
			},
		},
	}, spaceId)

	if len(e.SubWorkItems) > 0 {
		for _, subWorkItem := range e.SubWorkItems {
			s.SendAllMember(map[string]any{
				"type": notify.Event_DeleteWorkItem,
				"data": map[string]any{
					"operator": map[string]any{
						"id":       operator.Id,
						"username": operator.UserName,
						"nickname": operator.UserNickname,
					},
					"workItem": map[string]any{
						"id":   subWorkItem.Id,
						"pid":  subWorkItem.Pid,
						"name": subWorkItem.WorkItemName,
					},
				},
			}, spaceId)
		}
	}
}

func (s *Notify) CooperateAddMember(e *event.AddMember) {
	ctx := context.Background()

	space := e.Space
	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendToMember(map[string]any{
		"type": notify.Event_AddMember,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   space.Id,
				"name": space.SpaceName,
			},
		},
	}, e.TargetId)
}

func (s *Notify) CooperateRemoveMember(e *event.RemoveMember) {
	ctx := context.Background()

	space := e.Space
	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendToMember(map[string]any{
		"type": notify.Event_RemoveMember,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   space.Id,
				"name": space.SpaceName,
			},
		},
	}, e.TargetId)
}

func (s *Notify) CooperateComment(e *event.CooperateComment) {
	ctx := context.Background()

	spaceId := e.Space.Id
	operatorId := e.Operator
	workItem := e.WorkItem

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": e.Event,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"workItem": map[string]any{
				"id":   workItem.Id,
				"pid":  workItem.Pid,
				"name": workItem.WorkItemName,
			},
		},
	}, spaceId)
}

func (s *Notify) CooperateTag(e *event.CooperateTag) {
	ctx := context.Background()

	spaceId := e.Space.Id
	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": e.Event,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   e.Space.Id,
				"name": e.Space.SpaceName,
			},
			"tag": map[string]any{
				"id":   e.Tag.Id,
				"name": e.Tag.TagName,
			},
		},
	}, spaceId)
}

func (s *Notify) CooperateDisableWorkFlow(e *event.CooperateDisableWorkFlow) {
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": e.Event,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   space.Id,
				"name": space.SpaceName,
			},
			"workFlow": map[string]any{
				"id":   e.WorkFlowId,
				"name": e.WorkFlowName,
			},
		},
	}, e.SpaceId)
}

func (s *Notify) CooperateDeleteWorkFlow(e *event.CooperateDeleteWorkFlow) {
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": e.Event,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   space.Id,
				"name": space.SpaceName,
			},
			"workFlow": map[string]any{
				"id":   e.WorkFlowId,
				"name": e.WorkFlowName,
			},
		},
	}, e.SpaceId)
}

func (s *Notify) CooperateUpgradeWorkFlow(e *event.CooperateUpgradeWorkFlow) {
	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	operatorId := e.Operator

	operator, err := s.userRepo.GetUserByUserId(ctx, operatorId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.SendAllMember(map[string]any{
		"type": e.Event,
		"data": map[string]any{
			"operator": map[string]any{
				"id":       operator.Id,
				"username": operator.UserName,
				"nickname": operator.UserNickname,
			},
			"space": map[string]any{
				"id":   space.Id,
				"name": space.SpaceName,
			},
			"workFlow": map[string]any{
				"id":   e.WorkFlowId,
				"name": e.WorkFlowName,
			},
			"workItem": map[string]any{
				"id":   e.WorkItemId,
				"name": "",
			},
		},
	}, e.SpaceId)
}

func (s *Notify) CooperateAdminCancelUser(e *event.AdminCancelUser) {
	s.SendToMember(map[string]any{
		"type": notify.Event_AdminCancelUser,
	}, e.UserId)
}
