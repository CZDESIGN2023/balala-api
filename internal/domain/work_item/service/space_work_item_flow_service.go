package service

import (
	"context"
	"errors"
	"go-cs/internal/consts"
	"go-cs/internal/domain/pkg/flow_simulator"
	domain_message "go-cs/internal/domain/pkg/message"
	tplt_conf "go-cs/internal/domain/work_flow/flow_tplt_config"
	domain "go-cs/internal/domain/work_item"
	"go-cs/internal/domain/work_item/facade"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"

	"github.com/spf13/cast"
)

type UpgradeTaskWorkFlowRequest_Directors struct {
	RoleId    string
	RoleKey   string
	Directors domain.Directors
}
type UpgradeTaskFlowRequest struct {
	UpgradeTemplateId             int64                                   //升级后的模板id
	Directors                     []*UpgradeTaskWorkFlowRequest_Directors //升级后需要补充的负责人信息
	WorkItemStatusFacade          *facade.WorkItemStatusFacade
	WorkFlowTemplateServiceFacade *facade.WorkFlowTemplateServiceFacade
}

type UpgradeTaskFlowResult struct {
	NewWorkItemFlowRoles    domain.WorkItemFlowRoles
	DeleteWorkItemFlowRoles domain.WorkItemFlowRoles

	NewWorkItemFlowNodes    domain.WorkItemFlowNodes
	DeleteWorkItemFlowNodes domain.WorkItemFlowNodes
}

func (s *WorkItemService) UpgradeTaskWorkFlow(ctx context.Context, task *domain.WorkItem, req *UpgradeTaskFlowRequest, oper shared.Oper) (*UpgradeTaskFlowResult, error) {

	type RoleDirector struct {
		RoleId    string
		RoleKey   string
		Directors domain.Directors
	}

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务支持工作流升级")
	}

	statusInfo := req.WorkItemStatusFacade
	curState := statusInfo.GetItemByKey(task.WorkItemStatus.Key)
	if curState == nil {
		return nil, errs.Business(ctx, "获取任务状态失败")
	}

	//获取流程信息
	flow, err := req.WorkFlowTemplateServiceFacade.GetWorkFlow(ctx, task.WorkFlowId)
	if err != nil {
		return nil, err
	}

	if flow.LastTemplateId == task.WorkFlowTemplateId {
		return nil, errs.Business(ctx, "当前任务模版版本高于升级后的模版版本，无需升级")
	}

	//比对新旧两个版本的模版配置信息
	oldTplt, err := req.WorkFlowTemplateServiceFacade.GetWorkFlowTemplate(ctx, task.WorkFlowTemplateId)
	if err != nil {
		return nil, err
	}

	newTplt, err := req.WorkFlowTemplateServiceFacade.GetWorkFlowTemplate(ctx, req.UpgradeTemplateId)
	if err != nil {
		return nil, err
	}

	if oldTplt.Template.Version >= newTplt.Template.Version {
		return nil, errs.Business(ctx, "当前任务模版版本高于升级后的模版版本，无需升级")
	}

	newFlowConf := newTplt.Template.WorkFLowConfig
	newFlowNodes := newFlowConf.Nodes

	//-- 通过新的模版配置，把角色信息对应的负责人编排一下
	newRoleDirectorMap := make(map[string]*RoleDirector)
	for _, v := range newFlowNodes {
		if cast.ToInt64(v.GetOwnerRoleId()) == 0 {
			continue
		}

		newRoleDirectorMap[v.GetOwnerRoleId()] = &RoleDirector{
			RoleId:    v.GetOwnerRoleId(),
			RoleKey:   v.GetOwnerRoleKey(),
			Directors: make(domain.Directors, 0),
		}
	}

	oldRoleDirectorMap := make(map[string]*RoleDirector)
	for _, v := range task.WorkItemFlowRoles {
		//配置里有的，才要设置
		if newRoleDirectorMap[cast.ToString(v.WorkItemRoleId)] == nil {
			continue
		}
		newRoleDirectorMap[cast.ToString(v.WorkItemRoleId)].Directors = v.Directors
		oldRoleDirectorMap[cast.ToString(v.WorkItemRoleId)] = &RoleDirector{
			RoleId:    cast.ToString(v.WorkItemRoleId),
			RoleKey:   v.WorkItemRoleKey,
			Directors: v.Directors,
		}
	}

	//-- 组织新的角色负责人
	convertDirector := func(director string) string {
		if director == "_creator" {
			return cast.ToString(task.UserId)
		}
		return director
	}

	for _, v := range req.Directors {
		//配置里有的，才要设置
		if newRoleDirectorMap[cast.ToString(v.RoleId)] == nil {
			continue
		}

		directors := stream.Map(v.Directors, func(director string) string {
			return convertDirector(director)
		})

		newRoleDirectorMap[cast.ToString(v.RoleId)].Directors = stream.Unique(directors)
	}

	//检查一下负责人是不是都被设置了
	for _, v := range newRoleDirectorMap {
		if len(v.Directors) == 0 {
			return nil, errors.New("负责人角色配置错误")
		}
	}

	//-- 组织新的角色负责人
	removeFlowRoles := task.RemoveAllWorkItemFlowRoles()
	flowRolesTag := make(map[string]bool)
	flowRoles := make(domain.WorkItemFlowRoles, 0)
	for k, v := range newRoleDirectorMap {

		_, isOk := flowRolesTag[k]
		if isOk {
			//重复的不添加
			continue
		}

		flowRolesTag[v.RoleId] = true
		flowRole := domain.NewWorkItemFlowRole(
			task.SpaceId,
			task.Id,
			newTplt.Template.WorkFlowId,
			newTplt.Template.Id,
			cast.ToInt64(v.RoleId),
			v.RoleKey,
			v.Directors,
			task.CreatedAt,
		)

		flowRoles = append(flowRoles, flowRole)
	}

	task.AddWorkItemFlowRole(flowRoles...)

	//--先通过配置，创建出对应的工作项流程信息, 并且从原数据的状态中填充
	removeFlowNodes := task.RemoveAllWorkItemFlowNodes()
	taskFlowNodes := removeFlowNodes.NodeMap()
	flowNodes := make(domain.WorkItemFlowNodes, 0)
	flowNodesMap := make(map[string]*domain.WorkItemFlowNode)
	for _, flowNodeConf := range newFlowConf.Nodes {

		var directors domain.Directors
		needRoleId := flowNodeConf.GetOwnerRoleId()
		needRoleKey := ""

		directorRole, isOk := newRoleDirectorMap[needRoleId]
		if isOk {
			directors = directorRole.Directors
			needRoleKey = directorRole.RoleKey
		}

		flowNode := domain.NewWorkItemFlowNode(
			newTplt.Template.WorkFlowId,
			newTplt.Template.Id,
			task.SpaceId,
			task.Id,
			flowNodeConf.Key,
			cast.ToInt64(needRoleId),
			needRoleKey,
			directors,
			task.CreatedAt,
		)

		//从历史的填充
		taskFlowNode := taskFlowNodes[flowNodeConf.Key]
		if taskFlowNode != nil {
			flowNode.UpdatePlanTime(taskFlowNode.PlanTime)
			flowNode.ResetStatusForm(taskFlowNode)
		}

		flowNodesMap[flowNodeConf.Key] = flowNode
		flowNodes = append(flowNodes, flowNode)
	}

	task.AddWorkItemFlowNode(flowNodes...)

	//更新最终版本信息
	task.UpdateWorkFlowTemplate(newTplt.Template.Id, int64(newTplt.Template.Version))

	graph, err := flow_simulator.NewWorkFlowGraph(task, newTplt.Template)
	if err != nil {
		return nil, err
	}

	// 重新计算所有节点状态
	graph.ReCalcAllNodeStatus()

	// 重新计算任务状态
	if curState.IsProcessingTypeState() {
		// 重新计算任务状态
		graph.ReCalculateWorkItemStatus()

		//强制刷一下任务最后状态的时间
		task.ForceUpdateLastStatusTime()
	}

	// 调整当前负责人
	if curState.IsCompleted() { // 升级前是完成状态
		// 调整当前负责人为完成节点的所有前置节点负责人
		prevNodesInfo := domain.WorkItemFlowNodes(graph.GetPrevNodesInfo(tplt_conf.WorkflowNodeCode_Ended))
		task.UpdateDirectors(prevNodesInfo.GetAllDirectors())
	} else {
		task.ReCalcDirectors()
	}

	// 调整参与人
	task.UpdateParticipators()

	//添加领域消息
	taskMsg := &domain_message.UpgradeTaskWorkFlow{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
		WorkFlowId:   task.WorkFlowId,
		WorkFlowName: flow.Name,
		OldVersion:   oldTplt.Template.Version,
		NewVersion:   newTplt.Template.Version,
		NewRoles: stream.Map(stream.Values(newRoleDirectorMap), func(v *RoleDirector) domain_message.RoleDirector {
			return domain_message.RoleDirector{
				RoleId:    v.RoleId,
				RoleKey:   v.RoleKey,
				Directors: v.Directors,
			}
		}),
		OldRoles: stream.Map(stream.Values(oldRoleDirectorMap), func(v *RoleDirector) domain_message.RoleDirector {
			return domain_message.RoleDirector{
				RoleId:    v.RoleId,
				RoleKey:   v.RoleKey,
				Directors: v.Directors,
			}
		}),
	}
	task.AddMessage(oper, taskMsg)

	return &UpgradeTaskFlowResult{
		NewWorkItemFlowRoles:    flowRoles,
		DeleteWorkItemFlowRoles: removeFlowRoles,

		NewWorkItemFlowNodes:    flowNodes,
		DeleteWorkItemFlowNodes: removeFlowNodes,
	}, nil
}

func (s *WorkItemService) UpgradeTaskStateFlow(ctx context.Context, task *domain.WorkItem, req *UpgradeTaskFlowRequest, oper shared.Oper) (*UpgradeTaskFlowResult, error) {

	type RoleDirector struct {
		RoleId    string
		RoleKey   string
		Directors domain.Directors
	}

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务支持工作流升级")
	}

	statusInfo := req.WorkItemStatusFacade
	curState := statusInfo.GetItemByKey(task.WorkItemStatus.Key)
	if curState == nil {
		return nil, errs.Business(ctx, "获取任务状态失败")
	}

	//获取流程信息
	flow, err := req.WorkFlowTemplateServiceFacade.GetWorkFlow(ctx, task.WorkFlowId)
	if err != nil {
		return nil, err
	}

	if flow.LastTemplateId == task.WorkFlowTemplateId {
		return nil, errs.Business(ctx, "当前任务模版版本高于升级后的模版版本，无需升级")
	}

	//比对新旧两个版本的模版配置信息
	oldTplt, err := req.WorkFlowTemplateServiceFacade.GetWorkFlowTemplate(ctx, task.WorkFlowTemplateId)
	if err != nil {
		return nil, err
	}

	newTplt, err := req.WorkFlowTemplateServiceFacade.GetWorkFlowTemplate(ctx, req.UpgradeTemplateId)
	if err != nil {
		return nil, err
	}

	oldTemplate := oldTplt.Template
	newTemplate := newTplt.Template

	if oldTemplate.Version >= newTemplate.Version {
		return nil, errs.Business(ctx, "当前任务模版版本高于升级后的模版版本，无需升级")
	}

	newFlowConf := newTemplate.StateFlowConf()
	newFlowNodes := newFlowConf.StateFlowNodes

	//-- 通过新的模版配置，把角色信息对应的负责人编排一下
	newRoleDirectorMap := make(map[string]*RoleDirector)
	for _, v := range newFlowNodes {
		if cast.ToInt64(v.GetOwnerRoleId()) == 0 {
			continue
		}

		newRoleDirectorMap[v.GetOwnerRoleId()] = &RoleDirector{
			RoleId:  v.GetOwnerRoleId(),
			RoleKey: v.GetOwnerRoleKey(),
		}
	}

	oldRoleDirectorMap := make(map[string]*RoleDirector)
	for _, v := range task.WorkItemFlowRoles {
		//配置里有的，才要设置
		if newRoleDirectorMap[cast.ToString(v.WorkItemRoleId)] == nil {
			continue
		}
		newRoleDirectorMap[cast.ToString(v.WorkItemRoleId)].Directors = v.Directors
		oldRoleDirectorMap[cast.ToString(v.WorkItemRoleId)] = &RoleDirector{
			RoleId:    cast.ToString(v.WorkItemRoleId),
			RoleKey:   v.WorkItemRoleKey,
			Directors: v.Directors,
		}
	}

	//-- 组织新的角色负责人
	convertDirector := func(director string) string {
		if director == "_creator" {
			return cast.ToString(task.UserId)
		}
		return director
	}

	for _, v := range req.Directors {
		//配置里有的，才要设置
		if newRoleDirectorMap[cast.ToString(v.RoleId)] == nil {
			continue
		}

		directors := stream.Map(v.Directors, func(director string) string {
			return convertDirector(director)
		})

		newRoleDirectorMap[cast.ToString(v.RoleId)].Directors = stream.Unique(directors)
	}

	//检查一下负责人是不是都被设置了
	for _, v := range newRoleDirectorMap {
		if len(v.Directors) == 0 {
			return nil, errors.New("负责人角色配置错误")
		}
	}

	//-- 组织新的角色负责人
	removeFlowRoles := task.RemoveAllWorkItemFlowRoles()
	newFlowRoles := stream.Map(stream.Values(newRoleDirectorMap), func(v *RoleDirector) *domain.WorkItemFlowRole {
		return domain.NewWorkItemFlowRole(
			task.SpaceId,
			task.Id,
			newTemplate.WorkFlowId,
			newTemplate.Id,
			cast.ToInt64(v.RoleId),
			v.RoleKey,
			v.Directors,
			task.CreatedAt,
		)
	})

	task.AddWorkItemFlowRole(newFlowRoles...)

	//--先通过配置，创建出对应的工作项流程信息, 并且从原数据的状态中填充
	removeFlowNodes := task.RemoveAllWorkItemFlowNodes()
	removeFlowNodeMap := removeFlowNodes.NodeMap()
	flowNodes := stream.Map(newFlowNodes, func(flowNodeConf *tplt_conf.StateFlowNode) *domain.WorkItemFlowNode {
		var directors domain.Directors
		needRoleId := flowNodeConf.GetOwnerRoleId()
		needRoleKey := ""

		directorRole, isOk := newRoleDirectorMap[needRoleId]
		if isOk {
			directors = directorRole.Directors
			needRoleKey = directorRole.RoleKey
		}

		flowNode := domain.NewWorkItemFlowNode(
			newTemplate.WorkFlowId,
			newTemplate.Id,
			task.SpaceId,
			task.Id,
			flowNodeConf.Key,
			cast.ToInt64(needRoleId),
			needRoleKey,
			directors,
			task.CreatedAt,
		)

		//从历史的填充
		taskFlowNode := removeFlowNodeMap[flowNodeConf.Key]
		if taskFlowNode != nil {
			flowNode.UpdatePlanTime(taskFlowNode.PlanTime)
			flowNode.ResetStatusForm(taskFlowNode)
		}

		return flowNode
	})

	task.AddWorkItemFlowNode(flowNodes...)

	// 更新最终版本信息
	task.UpdateWorkFlowTemplate(newTemplate.Id, int64(newTemplate.Version))

	//if !curState.IsTerminated() && newFlowConf.GetNodeByKey(curState.Key) == nil { //如果不是终止状态，并且升级后的模版里没有这个状态
	if !curState.IsTerminated() { //如果不是终止状态
		var newNode *tplt_conf.StateFlowNode

		switch curState.StatusType {
		case consts.WorkItemStatusType_Archived: // 第一个归档状态的节点
			newNode = newTemplate.StateFlowConf().GetFirstArchivedNode()
		case consts.WorkItemStatusType_Process: // 第一个处理状态的节点
			newNode = newTemplate.StateFlowConf().GetFirstProcessNode()
		}

		if newNode == nil {
			newNode = newTemplate.StateFlowConf().GetInitStateNode()
		}

		if newNode != nil {
			task.UpdateStatus(domain.WorkItemStatus{
				Id:  cast.ToInt64(newNode.SubStateId),
				Key: newNode.SubStateKey,
				Val: newNode.SubStateVal,
			})
		}
	}

	// 调整节点状态
	curFlowNode := task.WorkItemFlowNodes.GetNodeByCode(task.WorkItemStatus.Key)
	if curFlowNode != nil {
		for _, node := range task.WorkItemFlowNodes.GetProcessingNodes() {
			node.ResetStatus()
		}
		curFlowNode.ResetProgressStatus() // 重置节点状态
	}

	task.ReCalcDirectors()
	task.UpdateParticipators()

	//添加领域消息
	taskMsg := &domain_message.UpgradeTaskWorkFlow{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
		WorkFlowId:   task.WorkFlowId,
		WorkFlowName: flow.Name,
		OldVersion:   oldTemplate.Version,
		NewVersion:   newTemplate.Version,
		NewRoles: stream.Map(stream.Values(newRoleDirectorMap), func(v *RoleDirector) domain_message.RoleDirector {
			return domain_message.RoleDirector{
				RoleId:    v.RoleId,
				RoleKey:   v.RoleKey,
				Directors: v.Directors,
			}
		}),
		OldRoles: stream.Map(stream.Values(oldRoleDirectorMap), func(v *RoleDirector) domain_message.RoleDirector {
			return domain_message.RoleDirector{
				RoleId:    v.RoleId,
				RoleKey:   v.RoleKey,
				Directors: v.Directors,
			}
		}),
	}
	task.AddMessage(oper, taskMsg)

	return &UpgradeTaskFlowResult{
		NewWorkItemFlowRoles:    newFlowRoles,
		DeleteWorkItemFlowRoles: removeFlowRoles,

		NewWorkItemFlowNodes:    flowNodes,
		DeleteWorkItemFlowNodes: removeFlowNodes,
	}, nil
}
