package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	domain_message "go-cs/internal/domain/pkg/message"
	"go-cs/internal/domain/work_item/repo"
)

func (s *Notify) reachWorkItemFlowNodeByDomainMessage(e *domain_message.ReachWorkItemFlowNode) {
	s.log.Infof("reachWorkItemFlowNodeByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.changeWorkItemFlowNode(&event.ChangeWorkItemFlowNode{
		Event: notify.Event_ChangeWorkItemFlowNode,
		Data: &event.ChangeWorkItemFlowNodeData{
			Operator: &event.ChangeWorkItemFlowNodeData_User{
				Id: e.Oper.GetId(),
			},
			Space: &event.ChangeWorkItemFlowNodeData_Space{
				Id:   space.Id,
				Name: space.SpaceName,
			},
			WorkItem: &event.ChangeWorkItemFlowNodeData_WorkItem{
				Id:   e.WorkItemId,
				Name: e.WorkItemName,
			},
			ToNode: &event.ChangeWorkItemFlowNodeData_Node{
				Code: e.FlowNodeCode,
				Id:   e.FlowNodeId,
				Name: e.FlowNodeName,
			},
		},
	})
}

func (s *Notify) confirmWorkItemFlowNodeByDomainMessage(e *domain_message.ConfirmWorkItemFlowNode) {
	s.log.Infof("ConfirmWorkItemFlowNode: %+v", e)

	// ctx := context.Background()

	// space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	// if err != nil {
	// 	s.log.Error(err)
	// 	return
	// }

	// s.changeWorkItemFlowNode(&event.ChangeWorkItemFlowNode{
	// 	Event: notify.Event_ChangeWorkItemFlowNode,
	// 	Data: &event.ChangeWorkItemFlowNodeData{
	// 		Operator: &event.ChangeWorkItemFlowNodeData_User{
	// 			Id: e.Oper.GetId(),
	// 		},
	// 		Space: &event.ChangeWorkItemFlowNodeData_Space{
	// 			Id:   space.Id,
	// 			Name: space.SpaceName,
	// 		},
	// 		WorkItem: &event.ChangeWorkItemFlowNodeData_WorkItem{
	// 			Id:   e.WorkItemId,
	// 			Name: e.WorkItemName,
	// 		},
	// 		ToNode: &event.ChangeWorkItemFlowNodeData_Node{
	// 			Code: e.FlowNodeCode,
	// 			Id:   e.FlowNodeId,
	// 		},
	// 	},
	// })
}

func (s *Notify) rollbackWorkItemFlowNodeByDomainMessage(e *domain_message.RollbackWorkItemFlowNode) {

	s.log.Infof("RollbackWorkItemFlowNode: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.rollbackWorkItemFlowNode(&event.RollbackWorkItemFlowNode{
		Event: notify.Event_RollbackWorkItem,
		Data: &event.RollbackWorkItemFlowNodeData{
			Operator: &event.RollbackWorkItemFlowNodeData_User{
				Id: e.Oper.GetId(),
			},
			Space: &event.RollbackWorkItemFlowNodeData_Space{
				Id:   space.Id,
				Name: space.SpaceName,
			},
			WorkItem: &event.RollbackWorkItemFlowNodeData_WorkItem{
				Id:   e.WorkItemId,
				Name: e.WorkItemName,
			},
			ToNode: &event.RollbackWorkItemFlowNodeData_Node{
				Code: e.FlowNodeCode,
				Id:   e.FlowNodeId,
			},
			Reason: e.Reason,
		},
	})

	// s.changeWorkItemFlowNode(&event.ChangeWorkItemFlowNode{
	// 	Event: notify.Event_ChangeWorkItemFlowNode,
	// 	Data: &event.ChangeWorkItemFlowNodeData{
	// 		Operator: &event.ChangeWorkItemFlowNodeData_User{
	// 			Id: e.Oper.GetId(),
	// 		},
	// 		Space: &event.ChangeWorkItemFlowNodeData_Space{
	// 			Id:   space.Id,
	// 			Name: space.SpaceName,
	// 		},
	// 		WorkItem: &event.ChangeWorkItemFlowNodeData_WorkItem{
	// 			Id:   e.WorkItemId,
	// 			Name: e.WorkItemName,
	// 		},
	// 		ToNode: &event.ChangeWorkItemFlowNodeData_Node{
	// 			Code: e.FlowNodeCode,
	// 			Id:   e.FlowNodeId,
	// 		},
	// 	},
	// })
}

func (s *Notify) modifymWorkItemFlowNodeByDomainMessage(e *domain_message.ModifyWorkItemFlowNode) {
	s.log.Infof("modifymWorkItemFlowNodeByDomainMessage: %+v", e)

	ctx := context.Background()

	space, err := s.spaceRepo.GetSpace(ctx, e.SpaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	workItem, err := s.workItemRepo.GetWorkItem(ctx, e.WorkItemId, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
		Tags:      true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	for _, v := range e.Updates {
		if v.Field == "planTime" {
			oldValues := v.OldValue.([]any)
			NewValues := v.NewValue.([]any)

			s.ChangeWorkFlowNodePlanTime(&event.ChangeWorkFlowNodePlanTime{
				Event:     notify.Event_ChangeWorkFlowNodePlanTime,
				Space:     space,
				WorkItem:  workItem,
				Operator:  e.Oper.GetId(),
				OldValues: []any{oldValues[0], oldValues[1]},
				NewValues: []any{NewValues[0], NewValues[1]},
			})
		}
	}

}
