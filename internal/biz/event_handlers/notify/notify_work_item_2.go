package notify

import (
	"context"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	domain_message "go-cs/internal/domain/pkg/message"
	tplt_conf "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
)

func (s *Notify) changeWorkItemFieldByDomainMessage(e *domain_message.ModifyWorkItem) {
	s.log.Infof("ChangeWorkItemField: %+v", e)

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
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	evt := &event.ChangeWorkItemField{
		Event:    notify.Event_ChangeWorkItemField,
		Operator: e.Oper.GetId(),
		Space:    space,
		WorkItem: workItem,
		Updates:  make([]event.FieldUpdate, 0),
	}

	for _, v := range e.Updates {
		evt.Updates = append(evt.Updates, event.FieldUpdate{Field: v.Field, OldValue: v.OldValue, NewValue: v.NewValue})
	}

	s.ChangeWorkItemField(evt)
}

func (s *Notify) createChildWorkItemByDomainMessage(e *domain_message.CreateWorkItemSubTask) {

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
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	parentWorkItem, err := s.workItemRepo.GetWorkItem(ctx, e.ParentWorkItemId, &repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
		Followers: true,
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.CreateChildWorkItem(&event.CreateChildWorkItem{
		Event:         notify.Event_CreateChildWorkItem,
		Operator:      e.Oper.GetId(),
		Space:         space,
		ChildWorkItem: workItem,
		WorkItem:      parentWorkItem,
	})
}

func (s *Notify) createWorkItemByDomainMessage(e *domain_message.CreateWorkItem) {

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
	}, &repo.WithOption{
		FlowNodes: true,
	})
	if err != nil {
		s.log.Error(err)
		return
	}

	nodeOps := make([]event.NodeDirectorOp, 0)
	if !workItem.IsSubTask() {
		flowTplt, err := s.workFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, workItem.WorkFlowTemplateId)
		if err != nil {
			s.log.Error(err)
			return
		}

		switch {
		case workItem.IsWorkFlowMainTask():
			flowTpltNodes := flowTplt.WorkFlowConf().GetNodeMap()
			for _, v := range workItem.WorkItemFlowNodes {
				nodeOps = append(nodeOps, event.NodeDirectorOp{
					NodeName:  flowTpltNodes[v.FlowNodeCode].Name,
					NewValues: v.Directors.ToInt64s(),
				})
			}
		case workItem.IsStateFlowMainTask():
			statusKeys := stream.Map(flowTplt.StateFlowConf().StateFlowNodes, func(v *tplt_conf.StateFlowNode) string {
				return v.SubStateKey
			})
			statusMap, err := s.workItemStatusRepo.WorkItemStatusKeyMap(ctx, workItem.SpaceId, statusKeys...)
			if err != nil {
				s.log.Error(err)
				return
			}

			for _, v := range workItem.WorkItemFlowNodes {
				nodeOps = append(nodeOps, event.NodeDirectorOp{
					NodeName:  statusMap[v.FlowNodeCode].Name,
					NewValues: v.Directors.ToInt64s(),
				})
			}
		}
	}

	s.CreateWorkItem(&event.CreateWorkItem{
		Event:    notify.Event_CreateWorkItem,
		Space:    space,
		WorkItem: workItem,
		Operator: e.Oper.GetId(),
		Nodes:    nodeOps,
	})
}

func (s *Notify) changeWorkItemDirectorByDomainMessage(e *domain_message.ChangeWorkItemDirector) {

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
	}, &repo.WithOption{
		FlowNodes: true,
	})
	if err != nil {
		s.log.Error(err)
		return
	}

	if workItem.IsSubTask() {
		s.ChangeChildWorkItemDirector(&event.ChangeWorkItemDirector{
			Event:     notify.Event_ChangeWorkItemDirector,
			Space:     space,
			WorkItem:  workItem,
			Operator:  e.Oper.GetId(),
			OldValues: utils.StringArrToInt64Arr(e.OldDirectors),
			NewValues: utils.StringArrToInt64Arr(e.NewDirectors),
		})
		return
	}

	nodeOps := make([]event.NodeDirectorOp, 0)
	if len(e.Nodes) > 0 {
		flowTplt, err := s.workFlowRepo.GetWorkFlowTemplateFormMemoryCache(ctx, workItem.WorkFlowTemplateId)
		if err != nil {
			s.log.Error(err)
			return
		}

		if workItem.IsWorkFlowMainTask() {
			flowTpltNodes := flowTplt.WorkFlowConf().GetNodeMap()
			for _, v := range e.Nodes {
				nodeOps = append(nodeOps, event.NodeDirectorOp{
					NodeName:  flowTpltNodes[v.FlowNodeCode].Name,
					NewValues: utils.StringArrToInt64Arr(v.NewDirectors),
					OldValues: utils.StringArrToInt64Arr(v.OldDirectors),
				})
			}
		}

		if workItem.IsStateFlowMainTask() {
			statusIds := flowTplt.StateFlowConf().GetAllStatusId()
			statusMap, err := s.workItemStatusRepo.StatusMap(ctx, statusIds)
			if err != nil {
				s.log.Error(err)
				return
			}

			statusNameMap := stream.MapKV(statusMap, func(_ int64, v *work_item_status.WorkItemStatusItem) (string, string) {
				return v.Key, v.Name
			})

			for _, v := range e.Nodes {
				nodeOps = append(nodeOps, event.NodeDirectorOp{
					NodeName:  statusNameMap[v.FlowNodeCode],
					NewValues: utils.StringArrToInt64Arr(v.NewDirectors),
					OldValues: utils.StringArrToInt64Arr(v.OldDirectors),
				})
			}
		}

	}

	s.ChangeWorkItemDirector(&event.ChangeWorkItemDirector{
		Event:     notify.Event_ChangeWorkItemDirector,
		Space:     space,
		WorkItem:  workItem,
		Operator:  e.Oper.GetId(),
		Nodes:     nodeOps,
		OldValues: utils.StringArrToInt64Arr(e.OldDirectors),
		NewValues: utils.StringArrToInt64Arr(e.NewDirectors),
	})

}

func (s *Notify) changeWorkItemTagByDomainMessage(e *domain_message.ChangeWorkItemTag) {

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

	s.ChangeWorkItemTag(&event.ChangeWorkItemTag{
		Event:     notify.Event_ChangeWorkItemTag,
		Space:     space,
		WorkItem:  workItem,
		Operator:  e.Oper.GetId(),
		OldValues: utils.StringArrToInt64Arr(e.OldTags),
		NewValues: utils.StringArrToInt64Arr(e.NewTags),
	})
}

func (s *Notify) changeWorkItemFileByDomainMessage(e *domain_message.ChangeWorkItemFile) {

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
	}, nil)
	if err != nil {
		s.log.Error(err)
		return
	}

	evt := &event.SetWorkItemFiles{
		Event:    notify.Event_SetWorkItemFile,
		Space:    space,
		WorkItem: workItem,
		Operator: e.Oper.GetId(),
		Adds:     make([]*event.SetWorkItemFiles_FileInfo, 0),
		Deletes:  make([]*event.SetWorkItemFiles_FileInfo, 0),
	}

	for _, v := range e.AddFiles {
		evt.Adds = append(evt.Adds, &event.SetWorkItemFiles_FileInfo{
			FileName: v.Name,
		})
	}

	for _, v := range e.RemoveFiles {
		evt.Deletes = append(evt.Deletes, &event.SetWorkItemFiles_FileInfo{
			FileName: v.Name,
		})
	}

	s.SetWorkItemFiles(evt)
}
