package service

import (
	"context"
	"errors"
	"go-cs/internal/consts"
	"go-cs/internal/domain/pkg/flow_simulator"
	domain_message "go-cs/internal/domain/pkg/message"
	flow_config "go-cs/internal/domain/work_flow/flow_tplt_config"
	domain "go-cs/internal/domain/work_item"
	"go-cs/internal/domain/work_item/facade"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/stream"
	"slices"
	"time"

	"github.com/spf13/cast"
)

type WorkItemService struct {
	repo      repo.WorkItemRepo
	idService *biz_id.BusinessIdService
}

func NewWorkItemService(
	repo repo.WorkItemRepo,
	idService *biz_id.BusinessIdService,
) *WorkItemService {
	return &WorkItemService{
		repo:      repo,
		idService: idService,
	}
}

type CreateWorkItemRequest_Directors struct {
	RoleId    string
	RoleKey   string
	Directors domain.Directors
}

type CreateWorkItemRequest struct {
	SpaceId                int64
	UserId                 int64
	WorkObjectId           int64
	VersionId              int64
	WorkItemTypeId         int64
	WorkItemTypeKey        string
	Name                   string
	PlanTime               domain.PlanTime
	ProcessRate            int32
	Remark                 string
	Describe               string
	Priority               string
	IconFlag               domain.IconFlag
	Tags                   domain.Tags
	Files                  []int64
	Followers              []int64
	Directors              []CreateWorkItemRequest_Directors
	WorkFlowFacade         *facade.WorkFlowFacade
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
	FileInfoFacade         *facade.FileInfoFacade
}

func (s *WorkItemService) CreateSpaceTask(ctx context.Context, req *CreateWorkItemRequest, oper shared.Oper) (*domain.WorkItem, error) {

	flowInfo := req.WorkFlowFacade.Flow()
	if flowInfo == nil {
		return nil, errs.Business(ctx, "任务单流程模板不存在")
	}

	flowTplt := req.WorkFlowTemplateFacade.Template()
	if flowTplt == nil {
		return nil, errs.Business(ctx, "任务单流程模板不存在")
	}

	//校验配置角色是否存在

	directorRoleMap := make(map[string]CreateWorkItemRequest_Directors, 0)
	for _, director := range req.Directors {
		directorRoleMap[director.RoleId] = director
	}

	//获取配置的流程节点
	//因为负责人是必填项,所以有角色的负责人全部都要校验
	var hasRoleChecked bool

	switch req.WorkFlowFacade.Flow().FlowMode {
	case consts.FlowMode_WorkFlow:
		for _, nodeConf := range flowTplt.WorkFlowConf().Nodes {
			if directorRole, isOk := directorRoleMap[nodeConf.GetOwnerRoleId()]; isOk {
				if !nodeConf.CheckOwnerRule(directorRole.Directors.ToStrings()) {
					return nil, errs.Business(ctx, "指定的负责人不在范围内")
				}
				//检查角色相关配置
				hasRoleChecked = true
			}
		}
	case consts.FlowMode_StateFlow:
		for _, nodeConf := range flowTplt.StateFlowConf().StateFlowNodes {
			if directorRole, isOk := directorRoleMap[nodeConf.GetOwnerRoleId()]; isOk {
				if !nodeConf.CheckOwnerRule(directorRole.Directors.ToStrings()) {
					return nil, errs.Business(ctx, "指定的负责人不在范围内")
				}
				//检查角色相关配置
				hasRoleChecked = true
			}
		}
	}

	if !hasRoleChecked {
		return nil, errs.Business(ctx, "指定的负责人不符合要求")
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkItem)
	if bizId == nil {
		return nil, errs.Business(ctx, "分配任务单ID失败")
	}

	//创建角色负责人
	flowRolesTag := make(map[string]bool)
	flowRoles := make([]*domain.WorkItemFlowRole, 0)
	for _, v := range req.Directors {

		_, isOk := flowRolesTag[v.RoleKey]
		if isOk {
			//重复的不添加
			continue
		}

		flowRolesTag[v.RoleKey] = true
		flowRole := domain.NewWorkItemFlowRole(
			req.SpaceId,
			bizId.Id,
			flowTplt.WorkFlowId,
			flowTplt.Id,
			cast.ToInt64(v.RoleId),
			v.RoleKey,
			v.Directors,
			time.Now().Unix(),
		)

		flowRoles = append(flowRoles, flowRole)
	}

	//先通过配置，创建出对应的工作项流程信息
	flowNodes := make(domain.WorkItemFlowNodes, 0)
	switch req.WorkFlowFacade.Flow().FlowMode {
	case consts.FlowMode_WorkFlow:
		for _, flowNodeConf := range flowTplt.WorkFlowConf().Nodes {

			var directors domain.Directors
			needRoleId := flowNodeConf.GetOwnerRoleId()
			needRoleKey := ""

			directorRole, isOk := directorRoleMap[needRoleId]
			if isOk {
				directors = directorRole.Directors
				needRoleKey = directorRole.RoleKey
			}

			flowNode := domain.NewWorkItemFlowNode(
				flowTplt.WorkFlowId,
				flowTplt.Id,
				req.SpaceId,
				bizId.Id,
				flowNodeConf.Key,
				cast.ToInt64(needRoleId),
				needRoleKey,
				directors,
				time.Now().Unix(),
			)

			flowNodes = append(flowNodes, flowNode)
		}
	case consts.FlowMode_StateFlow:
		for _, flowNodeConf := range flowTplt.StateFlowConf().StateFlowNodes {

			var directors domain.Directors
			needRoleId := flowNodeConf.GetOwnerRoleId()
			needRoleKey := ""

			directorRole, isOk := directorRoleMap[needRoleId]
			if isOk {
				directors = directorRole.Directors
				needRoleKey = directorRole.RoleKey
			}

			flowNode := domain.NewWorkItemFlowNode(
				flowTplt.WorkFlowId,
				flowTplt.Id,
				req.SpaceId,
				bizId.Id,
				flowNodeConf.Key,
				cast.ToInt64(needRoleId),
				needRoleKey,
				directors,
				time.Now().Unix(),
			)

			flowNodes = append(flowNodes, flowNode)
		}
	}

	//文件信息

	workItemFiles := make([]*domain.WorkItemFile, 0)
	fileInfos, _ := req.FileInfoFacade.GetFileInfos(ctx, req.Files)
	for _, v := range fileInfos {
		workItemFile := domain.NewWorkItemFile(
			req.SpaceId,
			bizId.Id,
			domain.FileInfo{
				FileInfoId: v.FileInfoId,
				FileName:   v.FileName,
				FileUri:    v.FileUri,
				FileSize:   v.FileSize,
			},
		)
		workItemFiles = append(workItemFiles, workItemFile)
	}

	workItem := domain.NewWorkItem(bizId.Id,
		req.SpaceId,
		req.UserId,
		req.WorkObjectId,
		req.VersionId,
		req.Name,
		domain.WorkItemStatus{
			Val: "",
			Key: "",
			Id:  0,
		},
		req.PlanTime,
		req.ProcessRate,
		req.Remark,
		req.Describe,
		req.Priority,
		req.IconFlag,
		req.Tags,
		domain.Directors{},
		req.Followers,
		flowNodes,
		flowRoles,
		workItemFiles,
		oper,
	)

	workItem.WorkFlowId = flowInfo.Id
	workItem.WorkFlowKey = flowInfo.Key
	workItem.WorkFlowTemplateId = flowTplt.Id
	workItem.WorkFlowTemplateVersion = int64(flowInfo.Version)
	workItem.WorkItemTypeId = req.WorkItemTypeId
	workItem.WorkItemTypeKey = consts.WorkItemTypeKey(req.WorkItemTypeKey)
	workItem.WorkItemFlowId = flowInfo.Id
	workItem.WorkItemFlowKey = flowInfo.Key
	workItem.FlowMode = flowInfo.FlowMode

	switch flowInfo.FlowMode {
	case consts.FlowMode_WorkFlow:
		graph, err := flow_simulator.NewWorkFlowGraph(workItem, flowTplt)
		if err != nil {
			return nil, err
		}
		graph.ConfirmNode("started")
		graph.ReCalculateWorkItemStatus()
		if err != nil {
			return nil, err
		}
		workItem.ReCalcDirectors()
		workItem.UpdateParticipators()
	case consts.FlowMode_StateFlow:
		flowConf := req.WorkFlowTemplateFacade.Template().StateFlowConfig
		initStateNode := flowConf.GetInitStateNode()
		// 任务状态
		workItem.UpdateStatus(domain.WorkItemStatus{
			Val: initStateNode.SubStateVal,
			Key: initStateNode.Key,
			Id:  cast.ToInt64(initStateNode.SubStateId),
		})

		// 当前节点设置为进行中
		workItem.WorkItemFlowNodes.GetNodeByCode(initStateNode.Key).ResetProgressStatus()

		workItem.ReCalcDirectors()
		workItem.UpdateParticipators()
	}

	return workItem, nil
}

type CreateSubTaskRequest struct {
	Name                   string
	PlanTime               domain.PlanTime
	ProcessRate            int32
	Directors              domain.Directors
	WorkFlowFacade         *facade.WorkFlowFacade
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

func (s *WorkItemService) CreateSpaceSubTask(ctx context.Context, parentTask *domain.WorkItem, req *CreateSubTaskRequest, oper shared.Oper) (*domain.WorkItem, error) {

	if len(req.Directors) == 0 {
		return nil, errs.Business(ctx, "负责人不能为空")
	}

	if req.Name == "" {
		return nil, errs.Business(ctx, "任务名称不能为空")
	}

	progressStatus := req.WorkItemStatusFacade.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))

	flowInfo := req.WorkFlowFacade.Flow()
	if flowInfo == nil {
		return nil, errs.Business(ctx, "任务单流程模板不存在")
	}

	flowTplt := req.WorkFlowTemplateFacade.Template()
	if flowTplt == nil {
		return nil, errs.Business(ctx, "子任务单流程模板不存在")
	}

	stateFlowConf := flowTplt.StateFlowConfig
	if stateFlowConf == nil {
		return nil, errs.Business(ctx, "子任务单流程模板不存在")
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkItem)
	if bizId == nil {
		return nil, errs.Business(ctx, "分配任务单ID失败")
	}

	//先通过配置，创建出对应的工作项流程信息
	flowNode := domain.NewWorkItemFlowNode(
		flowTplt.WorkFlowId,
		flowTplt.Id,
		parentTask.SpaceId,
		bizId.Id,
		"sub_task_"+rand.Letters(5),
		0,
		"",
		req.Directors,
		time.Now().Unix(),
	)

	subTask := domain.NewWorkItemSubTask(
		bizId.Id,
		parentTask.Id,
		parentTask.SpaceId,
		oper.GetId(),
		parentTask.WorkObjectId,
		parentTask.VersionId,
		req.Name,
		domain.WorkItemStatus{
			Val: progressStatus.Val,
			Key: progressStatus.Key,
			Id:  progressStatus.Id,
		},
		req.PlanTime,
		req.ProcessRate,
		parentTask.Doc.Priority,
		req.Directors,
		domain.WorkItemFlowNodes{flowNode},
		oper,
	)

	subTask.WorkFlowId = flowInfo.Id
	subTask.WorkFlowKey = flowInfo.Key
	subTask.WorkFlowId = flowTplt.WorkFlowId
	subTask.WorkFlowTemplateId = flowTplt.Id
	subTask.WorkFlowTemplateVersion = int64(flowInfo.Version)
	subTask.WorkItemTypeId = parentTask.WorkItemTypeId
	subTask.WorkItemTypeKey = parentTask.WorkItemTypeKey
	subTask.WorkItemFlowId = parentTask.WorkItemFlowId
	subTask.WorkItemFlowKey = parentTask.WorkItemFlowKey
	subTask.FlowMode = parentTask.FlowMode

	return subTask, nil
}

type ConfirmSpaceTaskNodeState struct {
	NodeCode               string
	Reason                 string
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

func (s *WorkItemService) ConfirmWorkFlowMain(ctx context.Context, task *domain.WorkItem, req *ConfirmSpaceTaskNodeState, oper shared.Oper) (domain.WorkItems, error) {
	if !task.IsWorkFlowMainTask() {
		return nil, errs.Business(ctx, "当前任务不是节点主流程任务")
	}

	statusInfo := req.WorkItemStatusFacade
	flowTpltInfo := req.WorkFlowTemplateFacade

	flowConf := flowTpltInfo.Template().WorkFLowConfig
	if flowConf == nil {
		return nil, errs.Business(ctx, "任务单流程配置不存在")
	}

	flownNodeConf := flowConf.GetNode(req.NodeCode)
	if flownNodeConf == nil {
		return nil, errs.Business(ctx, "任务单流程节点不存在")
	}

	flowNode := task.WorkItemFlowNodes.GetNodeByCode(req.NodeCode)
	if flowNode == nil || !flowNode.IsInProcess() {
		return nil, errs.Business(ctx, "当前任务节点状态不允许操作")
	}

	if flownNodeConf.ForcePlanTime && !flowNode.PlanTimeHasSet() {
		return nil, errs.Business(ctx, "当前任务节点未设置排期")
	}

	curStatus := statusInfo.GetItemByKey(task.WorkItemStatus.Key)
	if curStatus == nil || curStatus.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "当前任务状态不允许操作")
	}

	graph, err := flow_simulator.NewWorkFlowGraph(task, flowTpltInfo.Template())
	if err != nil {
		return nil, err
	}

	// 完成节点
	graph.ConfirmNode(req.NodeCode)
	// 重新计算状态
	graph.ReCalculateWorkItemStatus()

	// 节点流转后，让表格中任务状态字段后的时长重置
	task.ForceUpdateLastStatusTime()

	nextStatus := req.WorkItemStatusFacade.GetItemByKey(task.WorkItemStatus.Key)
	// 如果是完成状态
	if nextStatus.IsCompleted() {
		task.UpdateProcessRate(100)
	}

	if nextStatus.IsProcessingTypeState() {
		// 更新当前负责人
		task.ReCalcDirectors()
	}

	// 如果是归档状态，子任务也归档
	var subTasks domain.WorkItems
	if nextStatus.IsArchivedTypeState() {
		subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
		for _, subTask := range subTasks {
			subCurStatus := req.WorkItemStatusFacade.GetItemByKey(subTask.WorkItemStatus.Key)
			if subCurStatus != nil && !subCurStatus.IsArchivedTypeState() {
				if nextStatus.IsCompleted() {
					subTask.UpdateProcessRate(100)
				}
				subTask.ChangeStatus(domain.WorkItemStatus{
					Key: nextStatus.Key,
					Val: nextStatus.Val,
					Id:  nextStatus.Id,
				}, req.Reason, true, oper)
			}
		}
	}

	// 任务状态有变化，加上日志
	if task.WorkItemStatus.Id != curStatus.Id {
		task.AddMessage(oper, &domain_message.ChangeWorkItemStatus{
			SpaceId:              task.SpaceId,
			WorkItemId:           task.Id,
			WorkItemName:         task.WorkItemName,
			WorkItemTypeKey:      task.WorkItemTypeKey,
			OldWorkItemStatusKey: curStatus.Key,
			OldWorkItemStatusId:  curStatus.Id,
			OldWorkItemStatusVal: curStatus.Val,
			NewWorkItemStatusKey: task.WorkItemStatus.Key,
			NewWorkItemStatusId:  task.WorkItemStatus.Id,
			NewWorkItemStatusVal: task.WorkItemStatus.Val,
			Reason:               req.Reason,
		})
	}

	task.AddMessage(oper, &domain_message.ConfirmWorkItemFlowNode{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		FlowNodeCode: flowNode.FlowNodeCode,
		FlowNodeId:   flowNode.Id,
		WorkItemName: task.WorkItemName,
		Reason:       req.Reason,
	})

	// 记录计算出当前状态的节点
	rightFirstInProcessNodeCode := graph.FindRightFirstInProcessNode()
	if rightFirstInProcessNodeCode != "" {
		rightFirstInProcessNodeInfo := graph.GetNodeInfo(rightFirstInProcessNodeCode)
		rightFirstInProcessNodeConf := graph.GetNodeConfig(rightFirstInProcessNodeCode)

		task.AddMessage(oper, &domain_message.ReachWorkItemFlowNode{
			SpaceId:      task.SpaceId,
			WorkItemId:   task.Id,
			WorkItemName: task.WorkItemName,
			FlowNodeId:   rightFirstInProcessNodeInfo.Id,
			FlowNodeCode: rightFirstInProcessNodeInfo.FlowNodeCode,
			FlowNodeName: rightFirstInProcessNodeConf.Name,
			Reason:       req.Reason,
		})
	}

	return append(subTasks, task), nil
}

type SetWorkItemFileInfoRequest struct {
	AddFileInfoIds              []int64
	RemoveFileInfoIds           []int64
	FileInfoFacade              *facade.FileInfoFacade
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) SetWorkItemFileInfo(ctx context.Context, task *domain.WorkItem, req *SetWorkItemFileInfoRequest, oper shared.Oper) (domain.WorkItemFiles, domain.WorkItemFiles, error) {

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return nil, nil, errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return nil, nil, errs.Business(ctx, "任务已归档，不允许修改")
	}

	var allFileInfoIds []int64
	for _, v := range task.WorkItemFiles {
		allFileInfoIds = append(allFileInfoIds, v.FileInfo.FileInfoId)
	}

	adds := req.AddFileInfoIds
	removes := req.RemoveFileInfoIds
	removes = stream.Filter(removes, func(id int64) bool {
		return stream.Contains(allFileInfoIds, id)
	})

	adds = stream.Filter(adds, func(v int64) bool {
		return !stream.Contains(removes, v)
	})

	adds = stream.Filter(adds, func(v int64) bool {
		return !stream.Contains(allFileInfoIds, v)
	})

	// 删除
	removeWorkFiles := task.RemoveWorkFile(removes)

	// 添加
	workItemFiles := make(domain.WorkItemFiles, 0)

	fileInfos, _ := req.FileInfoFacade.GetFileInfos(ctx, stream.Unique(stream.Concat(adds, removes)))
	fileInfoMap := stream.ToMap(fileInfos, func(_ int, v domain.FileInfo) (int64, domain.FileInfo) {
		return v.FileInfoId, v
	})

	for _, id := range adds {
		if task.HasWorkFile(id) {
			continue
		}

		fileInfo := fileInfoMap[id]

		workItemFile := domain.NewWorkItemFile(
			task.SpaceId,
			task.Id,
			fileInfo,
		)
		workItemFiles = append(workItemFiles, workItemFile)
	}

	task.AddWorkFile(workItemFiles)

	addFileInfos := stream.Map(adds, func(v int64) domain.FileInfo {
		return fileInfoMap[v]
	})
	removeFileInfos := stream.Map(removes, func(v int64) domain.FileInfo {
		return fileInfoMap[v]
	})

	//-- 日志
	opsLog := &domain_message.ChangeWorkItemFile{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
	}

	for _, v := range addFileInfos {
		opsLog.AddFiles = append(opsLog.AddFiles, domain_message.FileInfo{
			Name: v.FileName,
			Size: v.FileSize,
		})
	}

	for _, v := range removeFileInfos {
		opsLog.RemoveFiles = append(opsLog.RemoveFiles, domain_message.FileInfo{
			Name: v.FileName,
			Size: v.FileSize,
		})
	}

	task.AddMessage(oper, opsLog)

	return workItemFiles, removeWorkFiles, nil
}

func (s *WorkItemService) DelSpaceTask(ctx context.Context, task *domain.WorkItem, oper shared.Oper) (domain.WorkItems, error) {

	task.OnDelete(oper)
	var subTasks domain.WorkItems
	if !task.IsSubTask() {
		subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, &repo.WithDocOption{Directors: true, Followers: true, Participators: true}, nil)
		for _, v := range subTasks {
			v.OnDelete(nil)
		}
	}

	return append(subTasks, task), nil
}

type SetDirectorsForSubTaskRequest struct {
	AddDirectors                domain.Directors
	RemoveDirectors             domain.Directors
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) SetDirectorsForSubTask(ctx context.Context, task *domain.WorkItem, req *SetDirectorsForSubTaskRequest, oper shared.Oper) error {

	if !task.IsSubTask() {
		return errs.Business(ctx, "不是子任务，不能设置负责人")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	oldDirectors := task.Doc.Directors
	newDirectors := stream.Diff(append(oldDirectors, req.AddDirectors...), req.RemoveDirectors)
	if len(newDirectors) == 0 {
		return errs.Business(ctx, "至少需要一个负责人")
	}

	// 子任务需要先手动移除节点和角色上的负责人
	for _, v := range task.WorkItemFlowNodes {
		v.UpdateDirectors(newDirectors)
	}
	for _, v := range task.WorkItemFlowRoles {
		v.UpdateDirectors(newDirectors)
	}

	//更新负责人
	task.UpdateDirectors(newDirectors)
	//更新参与人
	task.UpdateParticipators()

	task.AddMessage(oper, &domain_message.ChangeWorkItemDirector{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
		WorkItemPid:  task.Pid,

		OldDirectors: oldDirectors,
		NewDirectors: newDirectors,
	})

	return nil
}

type SetDirectorsForTaskByFlowNodeRequest struct {
	NodeKey                     string
	AddDirectors                domain.Directors
	RemoveDirectors             domain.Directors
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) SetDirectorsForWorkFlowMainTaskByNodeKey(ctx context.Context, task *domain.WorkItem, req *SetDirectorsForTaskByFlowNodeRequest, oper shared.Oper) error {

	if task.IsSubTask() {
		return errs.Business(ctx, "不是主任务，不能设置负责人")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	wItemFlowNode := task.WorkItemFlowNodes.GetNodeByCode(req.NodeKey)
	if wItemFlowNode == nil {
		return errs.Business(ctx, "流程节点不存在")
	}

	if !wItemFlowNode.IsInProcess() {
		return errs.Business(ctx, "流程节点不在进行中状态，不能设置负责人")
	}

	directors := append(wItemFlowNode.Directors, req.AddDirectors...)
	directors = stream.Unique(directors)
	directors = stream.Filter(directors, func(id string) bool {
		return !stream.Contains(req.RemoveDirectors, id)
	})

	if len(directors) == 0 {
		return errs.Business(ctx, "至少需要一个负责人")
	}

	adds := req.AddDirectors
	removes := req.RemoveDirectors
	flowNodeId := wItemFlowNode.Id
	oldValue := task.Doc.Directors

	//计算出可被移除的负责人
	removes = stream.Filter(removes, func(id string) bool {
		return stream.Contains(wItemFlowNode.Directors, id)
	})

	//计算出被添加的负责人
	adds = stream.Filter(adds, func(id string) bool {
		return !stream.Contains(removes, id)
	})

	adds = stream.Filter(adds, func(id string) bool {
		return !stream.Contains(wItemFlowNode.Directors, id)
	})

	//最终被处理的负责人
	directors = task.Doc.Directors
	directors = append(directors, adds...)
	directors = stream.Filter(directors, func(id string) bool {
		return !stream.Contains(removes, id)
	})
	directors = stream.Unique(directors)

	//更新各节点对应的负责人
	nodeEvts := make([]*domain_message.ChangeWorkItemDirector_Node, 0)
	var roleKey string
	var roleId int64
	for _, v := range task.WorkItemFlowNodes {
		if v.Id == flowNodeId {

			//当前节点被处理中，才需要调整当前负责人
			if v.IsInProcess() {
				task.UpdateDirectors(directors)
			}

			oldNodeValue := v.Directors.Clone()

			v.AddDirectors(adds)
			v.RemoveDirectors(removes)
			roleKey = v.WorkItemRoleKey
			roleId = v.WorkItemRoleId

			nodeEvts = append(nodeEvts, &domain_message.ChangeWorkItemDirector_Node{
				FlowNodeCode: v.FlowNodeCode,
				OldDirectors: oldNodeValue,
				NewDirectors: v.Directors,
			})

			break
		}
	}

	for _, v := range task.WorkItemFlowRoles {
		if v.WorkItemRoleKey == roleKey {
			v.AddDirectors(adds)
			v.RemoveDirectors(removes)
			break
		}
	}

	task.UpdateParticipators()

	//日志
	task.AddMessage(oper, &domain_message.ChangeWorkItemDirector{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,

		FlowTemplateId: task.WorkFlowTemplateId,

		OldDirectors: oldValue,
		NewDirectors: directors,

		WorkItemRoleKey: roleKey,
		WorkItemRoleId:  roleId,

		Nodes: nodeEvts,
	})

	return nil
}

type SetDirectorsForWorkFlowMainRequest struct {
	RoleKey                     string
	AddDirectors                domain.Directors
	RemoveDirectors             domain.Directors
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) SetDirectorsForWorkFlowMainTaskByRoleKey(ctx context.Context, task *domain.WorkItem, req *SetDirectorsForWorkFlowMainRequest, oper shared.Oper) error {

	if task.IsSubTask() {
		return errs.Business(ctx, "不是主任务，不能设置负责人")
	}

	wItemFlowRole := task.WorkItemFlowRoles.GetByRoleKey(req.RoleKey)
	if wItemFlowRole == nil {
		return errs.Business(ctx, "流程节点不存在")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	oldRoleDirectors := wItemFlowRole.Directors
	newRoleDirectors := stream.Diff(stream.Unique(append(oldRoleDirectors, req.AddDirectors...)), req.RemoveDirectors)

	if len(newRoleDirectors) == 0 {
		return errs.Business(ctx, "至少需要一个负责人")
	}

	// 更新角色负责人
	wItemFlowRole.UpdateDirectors(newRoleDirectors)

	// 更新角色关联的节点
	var needUpdateWorkItemDirector bool
	var nodeEvts []*domain_message.ChangeWorkItemDirector_Node
	for _, v := range task.WorkItemFlowNodes {
		if v.WorkItemRoleKey == wItemFlowRole.WorkItemRoleKey {
			oldNodeValue := v.Directors.Clone()
			v.UpdateDirectors(newRoleDirectors)

			if v.IsInProcess() { // 进行中的节点被处理中，才需要调整当前负责人
				needUpdateWorkItemDirector = true
			}

			nodeEvts = append(nodeEvts, &domain_message.ChangeWorkItemDirector_Node{
				FlowNodeCode: v.FlowNodeCode,
				OldDirectors: oldNodeValue,
				NewDirectors: v.Directors,
			})
		}
	}

	// 更新当前负责人
	if needUpdateWorkItemDirector {
		task.ReCalcDirectors()
	}

	//更新参与人
	task.UpdateParticipators()

	//日志
	task.AddMessage(oper, &domain_message.ChangeWorkItemDirector{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,

		FlowTemplateId:  task.WorkFlowTemplateId,
		WorkItemRoleKey: wItemFlowRole.WorkItemRoleKey,
		WorkItemRoleId:  wItemFlowRole.WorkItemRoleId,
		Nodes:           nodeEvts,
	})

	return nil
}

type ConfirmSpaceSubTaskState struct {
	NextStatusKey          string
	Reason                 string
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

// 节点流程模式下的子任务 使用状态模式 但是结合了流程模式主任务的状态情况
func (s *WorkItemService) ConfirmSpaceSubTaskState(ctx context.Context, subTask *domain.WorkItem, req *ConfirmSpaceSubTaskState, oper shared.Oper) error {

	if !subTask.IsSubTask() {
		return errs.Business(ctx, "不是子任务，不能设置状态")
	}

	task, err := s.repo.GetWorkItem(ctx, subTask.Pid, nil, nil)
	if err != nil {
		return errs.Business(ctx, "获取父级任务失败")
	}

	if req.WorkItemStatusFacade.HasArchivedItem(task.WorkItemStatus.Key) {
		return errs.Business(ctx, errors.New("主任务已完成，子任务不能为修改"))
	}

	tplt := req.WorkFlowTemplateFacade.Template()
	if tplt == nil || tplt.StateFlowConf() == nil {
		return errs.Business(ctx, errors.New("未找到流程模板"))
	}

	stateConf := tplt.StateFlowConf()
	curStateNode := stateConf.GetNode(subTask.WorkItemStatus.Key)
	nextStateNode := stateConf.GetNode(req.NextStatusKey)
	if nextStateNode == nil || curStateNode == nil {
		return errs.NoPerm(ctx)
	}

	if !stateConf.CanPass(curStateNode.Key, nextStateNode.Key) {
		return errs.Business(ctx, "不能切换到此状态")
	}

	curWorkItemStatus := req.WorkItemStatusFacade.GetItemByKey(curStateNode.SubStateKey)
	if curWorkItemStatus == nil {
		return errs.Business(ctx, "未找到状态")
	}

	nextWorkItemStatus := req.WorkItemStatusFacade.GetItemByKey(nextStateNode.SubStateKey)
	if nextWorkItemStatus == nil {
		return errs.Business(ctx, "未找到状态")
	}

	// 设置完成进度
	if nextWorkItemStatus.IsArchivedTypeState() {
		subTask.UpdateProcessRate(100)

	}

	if nextWorkItemStatus.IsProcessingTypeState() {
		subTask.UpdateRestart(oper.GetId(), 1) //标记为重启
		subTask.UpdateProcessRate(0)
	}

	// 设置状态节点信息
	if nextWorkItemStatus.IsArchivedTypeState() {
		subTask.UpdateStateFlowStateToFinished()
	} else {
		subTask.UpdateStateFlowStateToProgressing()
	}

	// 设置重置
	if curWorkItemStatus.IsArchivedTypeState() && !nextWorkItemStatus.IsArchivedTypeState() {
		progressStatus := req.WorkItemStatusFacade.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
		// 设置状态
		subTask.ChangeStatus(domain.WorkItemStatus{
			Val: progressStatus.Val,
			Key: progressStatus.Key,
			Id:  progressStatus.Id,
		}, req.Reason, false, oper)

	} else {
		// 设置状态
		subTask.ChangeStatus(domain.WorkItemStatus{
			Val: nextWorkItemStatus.Val,
			Key: nextWorkItemStatus.Key,
			Id:  nextWorkItemStatus.Id,
		}, req.Reason, false, oper)
	}

	return nil
}

type RestartTaskRequest struct {
	FlowNodeCode           string
	Reason                 string
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

func (s *WorkItemService) RestartTask(ctx context.Context, task *domain.WorkItem, req *RestartTaskRequest, oper shared.Oper) (domain.WorkItems, error) {

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务不适用此接口")
	}

	curStatus := req.WorkItemStatusFacade.GetItemByKey(task.WorkItemStatus.Key)
	if curStatus == nil {
		return nil, errs.Business(ctx, "任务状态信息错误")
	}

	if !curStatus.IsCompleted() && !curStatus.IsClose() {
		return nil, errs.Business(ctx, "当前任务状态不支持重启操作")
	}

	//flowConf := req.WorkFlowTemplateFacade.Template().WorkFlowConf()

	flowNode := task.WorkItemFlowNodes.GetNodeByCode(req.FlowNodeCode)

	subTasks, _ := s.repo.GetWorkItemByPid(ctx, task.Id, &repo.WithDocOption{
		Directors:     true,
		Participators: true,
		ProcessRate:   true,
	}, &repo.WithOption{
		FlowNodes: true,
	})

	graph, err := flow_simulator.NewWorkFlowGraph(task, req.WorkFlowTemplateFacade.Template())
	if err != nil {
		return nil, errs.Business(ctx, "流程图错误1 "+err.Error())
	}

	// 关闭-重启
	if curStatus.IsClose() {
		// 重新计算状态
		graph.ReCalculateWorkItemStatus()
		// 重新计算当前负责人
		task.ReCalcDirectors()

		task.UpdateRestart(oper.GetId(), 1)
		task.SetRestartReason(req.Reason)

		progressing := req.WorkItemStatusFacade.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
		for _, subTask := range subTasks {
			subTaskStatus := req.WorkItemStatusFacade.GetItemByKey(subTask.WorkItemStatus.Key)
			if subTaskStatus != nil && !subTaskStatus.IsCompleted() {
				subTask.UpdateProcessRate(0)
				subTask.UpdateRestart(oper.GetId(), 1)
				subTask.ChangeStatus(domain.WorkItemStatus{
					Key: progressing.Key,
					Val: progressing.Val,
					Id:  progressing.Id,
				}, req.Reason, true, oper)
				task.SetRestartReason(req.Reason)
			}
		}
	}

	// 完成-重启
	if curStatus.IsCompleted() {
		if req.FlowNodeCode == "" {
			graph.RebootToNode(flow_config.WorkflowNodeCode_Started)
		} else {
			graph.RebootToNode(req.FlowNodeCode)
		}

		// 重新计算状态
		graph.ReCalculateWorkItemStatus()
		// 重新计算当前负责人
		task.ReCalcDirectors()

		task.UpdateProcessRate(0)
		task.UpdateRestart(oper.GetId(), 1)

		task.SetRestartReason(req.Reason)

		for _, subTask := range subTasks {
			subTaskStatus := req.WorkItemStatusFacade.GetItemByKey(subTask.WorkItemStatus.Key)
			if subTaskStatus != nil && !subTaskStatus.IsCompleted() {
				progressingStatus := req.WorkItemStatusFacade.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
				if progressingStatus != nil {
					subTask.UpdateProcessRate(0)
					subTask.UpdateRestart(oper.GetId(), 1)
					subTask.ChangeStatus(domain.WorkItemStatus{
						Key: progressingStatus.Key,
						Val: progressingStatus.Val,
						Id:  progressingStatus.Id,
					}, req.Reason, true, oper)
					subTask.SetRestartReason(req.Reason)

				}
			}
		}

		//var lastReachNode *domain.WorkItemFlowNode
		//for _, v := range flowSimulateResult.GetLastReachNodes() {
		//	flowNode := task.WorkItemFlowNodes.GetNodeByCode(v.NodeKey)
		//	if flowNode != nil {
		//		lastReachNode = flowNode
		//		task.AddMessage(oper, &domain_message.ReachWorkItemFlowNode{
		//			SpaceId:      task.SpaceId,
		//			WorkItemId:   task.Id,
		//			WorkItemName: task.WorkItemName,
		//			FlowNodeCode: flowNode.FlowNodeCode,
		//			FlowNodeId:   flowNode.Id,
		//			FlowNodeName: flowConf.GetNode(flowNode.FlowNodeCode).Name,
		//			Reason:       req.Reason,
		//		})
		//	}
		//}

	}

	msg := &domain_message.ChangeWorkItemStatus{
		SpaceId:              task.SpaceId,
		WorkItemId:           task.Id,
		WorkItemName:         task.WorkItemName,
		OldWorkItemStatusKey: curStatus.Key,
		OldWorkItemStatusId:  curStatus.Id,
		OldWorkItemStatusVal: curStatus.Val,
		NewWorkItemStatusKey: task.WorkItemStatus.Key,
		NewWorkItemStatusId:  task.WorkItemStatus.Id,
		NewWorkItemStatusVal: task.WorkItemStatus.Val,
		Reason:               req.Reason,
	}
	if flowNode != nil {
		msg.FlowNodeId = flowNode.Id
		msg.FlowNodeCode = flowNode.FlowNodeCode
	}

	task.AddMessage(oper, msg)

	return append(subTasks, task), nil
}

type RollbackTaskRequest struct {
	FlowNodeCode           string
	Reason                 string
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

func (s *WorkItemService) RollbackTask(ctx context.Context, task *domain.WorkItem, req *RollbackTaskRequest, oper shared.Oper) (domain.WorkItems, error) {

	if !task.IsWorkFlowMainTask() {
		return nil, errs.Business(ctx, "非流程主任务")
	}

	curStatus := req.WorkItemStatusFacade.GetItemByKey(task.WorkItemStatus.Key)
	if curStatus == nil || curStatus.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "任务状态异常")
	}

	targetFlowNode := task.WorkItemFlowNodes.GetNodeByCode(req.FlowNodeCode)
	if targetFlowNode == nil || targetFlowNode.IsStartNode() {
		return nil, errs.Business(ctx, "回滚节点不存在")
	}

	task.SetRollbackReason(req.Reason)

	graph, err := flow_simulator.NewWorkFlowGraph(task, req.WorkFlowTemplateFacade.Template())
	if err != nil {
		return nil, errs.Business(ctx, "流程图错误1 "+err.Error())
	}
	if req.FlowNodeCode == "" {
		graph.RebootToNode("started")
	} else {
		graph.RebootToNode(req.FlowNodeCode)
	}
	graph.ReCalculateWorkItemStatus()
	task.ReCalcDirectors()

	var subTasks domain.WorkItems
	if task.IsWorkFlowMainTask() {
		lastChangeStatusInfo := req.WorkItemStatusFacade.GetItemByKey(task.WorkItemStatus.Key)
		// 如果是归档, 子任务要跟着一起归档
		if lastChangeStatusInfo.IsArchivedTypeState() {
			subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
			for _, subTask := range subTasks {
				subTask.UpdateStatus(domain.WorkItemStatus{
					Key: task.WorkItemStatus.Key,
					Val: task.WorkItemStatus.Val,
					Id:  task.WorkItemStatus.Id,
				})
				subTask.SetRollbackReason(req.Reason)
			}
		}
	}

	task.AddMessage(oper, &domain_message.RollbackWorkItemFlowNode{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
		FlowNodeId:   targetFlowNode.Id,
		FlowNodeCode: targetFlowNode.FlowNodeCode,
		Reason:       req.Reason,
	})

	task.AddMessage(oper, &domain_message.ChangeWorkItemStatus{
		SpaceId:              task.SpaceId,
		WorkItemId:           task.Id,
		WorkItemName:         task.WorkItemName,
		OldWorkItemStatusKey: curStatus.Key,
		OldWorkItemStatusId:  curStatus.Id,
		OldWorkItemStatusVal: curStatus.Val,
		NewWorkItemStatusKey: task.WorkItemStatus.Key,
		NewWorkItemStatusId:  task.WorkItemStatus.Id,
		NewWorkItemStatusVal: task.WorkItemStatus.Val,
		Reason:               req.Reason,
	})

	return append(subTasks, task), nil
}

type ResumeTaskRequest struct {
	Reason               string
	WorkItemStatusFacade *facade.WorkItemStatusFacade
}

func (s *WorkItemService) ResumeTask(ctx context.Context, task *domain.WorkItem, req *ResumeTaskRequest, oper shared.Oper) (domain.WorkItems, error) {

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务不支持恢复操作")
	}

	statusInfo := req.WorkItemStatusFacade

	curStatus := req.WorkItemStatusFacade.GetItemByKey(task.WorkItemStatus.Key)
	if curStatus == nil {
		return nil, errs.Business(ctx, "任务状态信息错误")
	}

	if !curStatus.IsTerminated() {
		return nil, errs.Business(ctx, "当前任务状态不支持重启操作")
	}

	// task.UpdateRestart(oper.GetId(), 1)
	task.ChangeStatus(domain.WorkItemStatus{
		Key: task.LastWorkItemStatus.Key,
		Val: task.LastWorkItemStatus.Val,
		Id:  task.LastWorkItemStatus.Id,
	}, req.Reason, false, oper)
	task.SetResumeReason(req.Reason)

	//当时终止任务时，最后一个进行中的节点需要更新一下时间

	var subTasks domain.WorkItems
	if task.HasChild() {
		progressingStatus := req.WorkItemStatusFacade.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))

		subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
		for _, subTask := range subTasks {
			subTaskStatus := statusInfo.GetItemByKey(subTask.WorkItemStatus.Key)
			if subTaskStatus != nil && !subTaskStatus.IsCompleted() {
				subTask.ChangeStatus(domain.WorkItemStatus{
					Key: progressingStatus.Key,
					Val: progressingStatus.Val,
					Id:  progressingStatus.Id,
				}, req.Reason, true, oper)
				subTask.SetResumeReason(req.Reason)
			}
		}
	}

	return append(subTasks, task), nil
}

type TerminateWorkItemRequest struct {
	Reason               string
	WorkItemStatusFacade *facade.WorkItemStatusFacade
}

func (s *WorkItemService) TerminateWorkItem(ctx context.Context, task *domain.WorkItem, req *TerminateWorkItemRequest, oper shared.Oper) (domain.WorkItems, error) {

	statusInfo := req.WorkItemStatusFacade

	curState := statusInfo.GetItemByKey(task.WorkItemStatus.Key)
	if curState == nil ||
		curState.IsTerminated() ||
		task.IsWorkFlowMainTask() && curState.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "当前任务状态不支持终止操作")
	}

	task.ChangeStatus(domain.WorkItemStatus{
		Val: statusInfo.Keyword().Terminated.Val,
		Key: statusInfo.Keyword().Terminated.Key,
		Id:  statusInfo.Keyword().Terminated.Id,
	}, req.Reason, false, oper)
	task.SetTerminateReason(req.Reason)

	// 子任务
	var subTasks domain.WorkItems
	if task.HasChild() {
		subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
		for _, v := range subTasks {
			subTaskStatus := statusInfo.GetItemByKey(v.WorkItemStatus.Key)
			if subTaskStatus != nil && !subTaskStatus.IsArchivedTypeState() {
				v.ChangeStatus(domain.WorkItemStatus{
					Val: statusInfo.Keyword().Terminated.Val,
					Key: statusInfo.Keyword().Terminated.Key,
					Id:  statusInfo.Keyword().Terminated.Id,
				}, req.Reason, true, oper)
				v.SetTerminateReason(req.Reason)
			}
		}
	}

	return append(subTasks, task), nil
}

type CloseTaskRequest struct {
	Reason               string
	FlowNodeCode         string
	WorkItemStatusFacade *facade.WorkItemStatusFacade
}

func (s *WorkItemService) CloseTask(ctx context.Context, task *domain.WorkItem, req *CloseTaskRequest, oper shared.Oper) (domain.WorkItems, error) {

	statusInfo := req.WorkItemStatusFacade

	curState := statusInfo.GetItemByKey(task.WorkItemStatus.Key)
	if curState == nil || curState.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "当前任务状态不允许关闭")
	}

	taskClosedStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCloseDefaultKey))

	task.ChangeStatus(domain.WorkItemStatus{
		Val: taskClosedStatus.Val,
		Key: taskClosedStatus.Key,
		Id:  taskClosedStatus.Id,
	}, req.Reason, false, oper)
	task.SetCloseReason(req.Reason)

	// 子任务
	var subTasks domain.WorkItems
	if task.HasChild() {
		subTasks, _ = s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
		for _, v := range subTasks {
			subTaskStatus := statusInfo.GetItemByKey(v.WorkItemStatus.Key)
			if subTaskStatus != nil && !subTaskStatus.IsArchivedTypeState() {
				v.ChangeStatus(domain.WorkItemStatus{
					Val: taskClosedStatus.Val,
					Key: taskClosedStatus.Key,
					Id:  taskClosedStatus.Id,
				}, req.Reason, true, oper)
				v.SetCloseReason(req.Reason)
			}
		}
	}

	return append(subTasks, task), nil
}

type ChangeVersionRequest struct {
	WorkVersionId               int64
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) ChangeWorkItemVersion(ctx context.Context, task *domain.WorkItem, req *ChangeVersionRequest, oper shared.Oper) (domain.WorkItems, error) {

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务不允许变更版本")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "任务已归档，不允许修改")
	}

	task.ChangeVersionId(req.WorkVersionId, oper)
	subTasks, _ := s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
	for _, v := range subTasks {
		v.UpdateVersionId(req.WorkVersionId)
	}

	return append(subTasks, task), nil
}

type ChangeWorkObjectRequest struct {
	WorkObjectId                int64
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) ChangeWorkItemObject(ctx context.Context, task *domain.WorkItem, req *ChangeWorkObjectRequest, oper shared.Oper) (domain.WorkItems, error) {

	if task.IsSubTask() {
		return nil, errs.Business(ctx, "子任务不允许变更模块")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return nil, errs.Business(ctx, "任务已归档，不允许修改")
	}

	task.ChangeWorkObjectId(req.WorkObjectId, oper)
	subTasks, _ := s.repo.GetWorkItemByPid(ctx, task.Id, nil, nil)
	for _, v := range subTasks {
		v.UpdateWorkObjectId(req.WorkObjectId)
	}

	return append(subTasks, task), nil
}

type ChangeWorkItemTagRequest struct {
	TagAdd                      []string
	TagRemove                   []string
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) ChangeWorkItemTag(ctx context.Context, task *domain.WorkItem, req *ChangeWorkItemTagRequest, oper shared.Oper) error {

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	task.ChangeTags(req.TagAdd, req.TagRemove, oper)

	return nil
}

type ModifyWorkItemFieldRequest struct {
	PropDiffs                   shared.PropDiffSet
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
}

func (s *WorkItemService) ModifyWorkItemField(ctx context.Context, task *domain.WorkItem, req *ModifyWorkItemFieldRequest, oper shared.Oper) error {

	//归档状态,不允许处理
	curStatus, err := req.WorkItemStatusServiceFacade.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curStatus.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	for k, v := range req.PropDiffs {
		switch k {
		case domain.Diff_ProcessRate:
			task.ChangeProcessRate(cast.ToInt32(v), oper)
		case domain.Diff_Name:
			task.ChangeName(cast.ToString(v), oper)
		case domain.Diff_Priority:
			task.ChangePriority(cast.ToString(v), oper)
		case domain.Diff_Remark:
			task.ChangeRemark(cast.ToString(v), oper)
		case domain.Diff_PlanTime:
			planTimes := v.([]int64)
			task.ChangePlanTime(planTimes[0], planTimes[1], oper)
		case domain.Diff_Describe:
			task.ChangeDescribe(cast.ToString(v), oper)
		case domain.Diff_IconFlags:
			task.UpdateIconFlag(v.([]*domain.IconFlagUpdate)...)
		}
	}

	return nil
}

type ConfirmStateFlowMainTaskStateByStateKey struct {
	NextStatusKey          string
	Reason                 string
	Remark                 string
	WorkFlowTemplateFacade *facade.WorkFlowTemplateFacade
	WorkItemStatusFacade   *facade.WorkItemStatusFacade
}

// 使用状态模式主任务
func (s *WorkItemService) ConfirmStateFlowMainTaskState(ctx context.Context, task *domain.WorkItem, req *ConfirmStateFlowMainTaskStateByStateKey, oper shared.Oper) error {

	if task.IsSubTask() {
		return errs.Business(ctx, "仅状态流主任务使用2")
	}

	tplt := req.WorkFlowTemplateFacade.Template()
	if tplt == nil || tplt.StateFlowConf() == nil {
		return errs.Business(ctx, errors.New("未找到流程模板"))
	}

	stateConf := tplt.StateFlowConf()
	curStateNode := stateConf.GetNode(task.WorkItemStatus.Key)
	nextStateNode := stateConf.GetNode(req.NextStatusKey)
	if nextStateNode == nil || curStateNode == nil {
		return errs.NoPerm(ctx)
	}

	if !stateConf.CanPass(curStateNode.Key, nextStateNode.Key) {
		return errs.Business(ctx, "不能切换到此状态")
	}

	curWorkItemStatus := req.WorkItemStatusFacade.GetItemByKey(curStateNode.SubStateKey)
	if curWorkItemStatus == nil {
		return errs.Business(ctx, "未找到状态")
	}

	nextWorkItemStatus := req.WorkItemStatusFacade.GetItemByKey(nextStateNode.SubStateKey)
	if nextWorkItemStatus == nil {
		return errs.Business(ctx, "未找到状态")
	}

	for _, node := range task.WorkItemFlowNodes.GetProcessingNodes() {
		node.ResetStatus() //重置节点状态
		//node.UpdatePlanTime(domain.PlanTime{}) //清空节点排期
	}

	flowNode := task.WorkItemFlowNodes.GetNodeByCode(nextWorkItemStatus.Key)
	flowNode.ResetProgressStatus()
	flowNode.UpdatePlanTime(domain.PlanTime{})

	// 设置完成进度
	if nextWorkItemStatus.IsArchivedTypeState() {
		task.UpdateProcessRate(100)
	}
	if curWorkItemStatus.IsArchivedTypeState() && nextWorkItemStatus.IsProcessingTypeState() {
		task.UpdateProcessRate(0)
	}

	// 设置状态
	task.ChangeStateFlowMainStatus(domain.WorkItemStatus{
		Val: nextWorkItemStatus.Val,
		Key: nextWorkItemStatus.Key,
		Id:  nextWorkItemStatus.Id,
	}, req.Reason, req.Remark, oper)

	//任务进行时间需要调整 countAt
	//如果状态类型相同，则不更新 countAt
	if curWorkItemStatus.StatusType != nextWorkItemStatus.StatusType {
		task.SetCountAt(time.Now())
	}

	// 调整当前负责人
	directors := stream.Unique(flowNode.Directors)
	task.UpdateDirectors(directors)

	task.AddMessage(oper, &domain_message.ReachWorkItemFlowNode{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,
		FlowNodeCode: flowNode.FlowNodeCode,
		FlowNodeId:   flowNode.Id,
		FlowNodeName: nextWorkItemStatus.Name,
		Reason:       req.Reason,
	})

	return nil
}

type SetDirectorsForStateFlowMainTaskRequest struct {
	Directors                   domain.Directors
	WorkItemStatusServiceFacade *facade.WorkItemStatusServiceFacade
	StateKeys                   []string
	RoleKeys                    []string
}

func (s *WorkItemService) SetDirectorsForStateFlowMainTask(ctx context.Context, task *domain.WorkItem, req *SetDirectorsForStateFlowMainTaskRequest, oper shared.Oper) error {

	if !task.IsStateFlowMainTask() {
		return errs.Business(ctx, "仅状态流主任务使用2")
	}

	statusInfo := req.WorkItemStatusServiceFacade
	curState, err := statusInfo.GetWorkItemStatusItem(ctx, task.WorkItemStatus.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if curState.IsArchivedTypeState() {
		return errs.Business(ctx, "任务已归档，不允许修改")
	}

	directors := stream.Unique(req.Directors)

	if len(directors) == 0 {
		return errs.Business(ctx, "至少需要一个负责人")
	}

	// 修改节点负责人
	var nodeEvts []*domain_message.ChangeWorkItemDirector_Node
	for _, v := range task.WorkItemFlowNodes {
		if slices.Contains(req.StateKeys, v.FlowNodeCode) {
			nodeEvts = append(nodeEvts, &domain_message.ChangeWorkItemDirector_Node{
				FlowNodeCode: v.FlowNodeCode,
				OldDirectors: v.Directors,
				NewDirectors: directors,
			})

			v.UpdateDirectors(directors)
		}
	}

	// 修改角色负责人
	task.ChangeRoleDirectors(directors, req.RoleKeys...)

	// 更新当前负责人
	if slices.Contains(req.StateKeys, task.WorkItemStatus.Key) {
		task.UpdateDirectors(directors)
	}

	// 更新参与人，节点负责人
	task.UpdateParticipators()

	//日志
	task.AddMessage(oper, &domain_message.ChangeWorkItemDirector{
		SpaceId:      task.SpaceId,
		WorkItemId:   task.Id,
		WorkItemName: task.WorkItemName,

		FlowTemplateId: task.WorkFlowTemplateId,
		Nodes:          nodeEvts,
	})

	return nil
}
