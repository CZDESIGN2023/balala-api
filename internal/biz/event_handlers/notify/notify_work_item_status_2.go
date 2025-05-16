package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/utils"
)

func (s *Notify) changeWorkItemStatusByDomainMessage(e *domain_message.ChangeWorkItemStatus) {
	s.log.Infof("changeWorkItemStatusByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	workItem, err := s.workItemRepo.GetWorkItem(ctx, int64(e.WorkItemId), nil, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	oldStatus, err := s.workItemStatusRepo.GetWorkItemStatusItem(ctx, e.OldWorkItemStatusId)
	if err != nil {
		s.log.Error(err)
		return
	}

	newStatus, err := s.workItemStatusRepo.GetWorkItemStatusItem(ctx, e.NewWorkItemStatusId)
	if err != nil {
		s.log.Error(err)
		return
	}

	if newStatus.IsTerminated() {

		utils.Go(func() {
			s.CooperateWorkItemStatusChange(notify.Event_TerminateWorkItem, &CooperateWorkItemEndedEvData{
				SpaceId:    space.Id,
				OperatorId: e.Oper.GetId(),
				WorkItem:   workItem,
			})
		})

		s.terminateWorkItem(&event.TerminateWorkItem{
			Event: notify.Event_TerminateWorkItem,
			Data: &event.TerminateWorkItemData{
				Reason: e.Reason,
				Operator: &event.TerminateWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.TerminateWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.TerminateWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
			},
		})
	}
	if oldStatus.IsTerminated() && !newStatus.IsArchivedTypeState() {
		s.resumeWorkItem(&event.ResumeWorkItem{
			Event: notify.Event_ResumeWorkItem,
			Data: &event.ResumeWorkItemData{
				Reason: e.Reason,
				Operator: &event.ResumeWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.ResumeWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.ResumeWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
			},
		})
	}

	if newStatus.IsCompleted() {

		utils.Go(func() {
			s.CooperateWorkItemStatusChange(notify.Event_CompleteWorkItem, &CooperateWorkItemEndedEvData{
				SpaceId:    space.Id,
				OperatorId: e.Oper.GetId(),
				WorkItem:   workItem,
			})
		})

		s.completeWorkItem(&event.CompleteWorkItem{
			Event: notify.Event_CompleteWorkItem,
			Data: &event.CompleteWorkItemData{
				Reason: e.Reason,
				Operator: &event.CompleteWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.CompleteWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.CompleteWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
			},
		})
	}

	if newStatus.IsClose() {
		utils.Go(func() {
			s.CooperateWorkItemStatusChange(notify.Event_CloseWorkItem, &CooperateWorkItemEndedEvData{
				SpaceId:    space.Id,
				OperatorId: e.Oper.GetId(),
				WorkItem:   workItem,
			})
		})

		s.closeWorkItem(&event.CloseWorkItem{
			Event: notify.Event_CloseWorkItem,
			Data: &event.CloseWorkItemData{
				Reason: e.Reason,
				Operator: &event.CloseWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.CloseWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.CloseWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
			},
		})
	}

	if oldStatus.IsClose() && !newStatus.IsArchivedTypeState() {
		s.resumeWorkItem(&event.ResumeWorkItem{
			Event: notify.Event_ResumeWorkItem,
			Data: &event.ResumeWorkItemData{
				Reason: e.Reason,
				Scene:  "formClosed",
				Operator: &event.ResumeWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.ResumeWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.ResumeWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
			},
		})
	}

	if oldStatus.IsCompleted() && !newStatus.IsArchivedTypeState() {
		s.restartWorkItem(&event.RestartWorkItem{
			Event: notify.Event_RestartWorkItem,
			Data: &event.RestartWorkItemData{
				Reason: e.Reason,
				Operator: &event.RestartWorkItemData_User{
					Id: e.Oper.GetId(),
				},
				Space: &event.RestartWorkItemData_Space{
					Id:   space.Id,
					Name: space.SpaceName,
				},
				WorkItem: &event.RestartWorkItemData_WorkItem{
					Id: e.WorkItemId,
				},
				ToNode: &event.RestartWorkItemData_Node{
					Name: "",
					Code: e.FlowNodeCode,
				},
			},
		})
	}
}
