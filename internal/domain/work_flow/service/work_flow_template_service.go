package service

import (
	"context"
	"errors"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	domain "go-cs/internal/domain/work_flow"
	"go-cs/internal/domain/work_flow/facade"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	wf_vaildate "go-cs/internal/domain/work_flow/vaildate"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type GenerateWorkFlowTemplateReq struct {
	SpaceId            int64
	WorkItemTypeId     int64
	WorkFlowId         int64
	WorkItemStatusInfo *facade.WorkItemStatusInfo
	WorkItemRoleInfo   *facade.WorkItemRoleInfo
	UserId             int64
}

// 组织系统默认的模版
func (s *WorkFlowService) newXuQiuWorkFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo

	processState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))
	checkingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowCheckingDefaultKey))
	testingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowTestingDefaultKey))

	roleInfo := req.WorkItemRoleInfo
	roleDevelper := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Developer)
	roleQa := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Qa)
	roleAccepter := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Acceptor)

	flowTpltConf := config.NewWorkFlow("xuqiu")

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.FormFields = append(flowTpltConf.FormFields, &config.WorkFlowFormField{
		Name: "describe", Value: `<p>- [任务描述]： </p><p>- [思维导图 | 如有]：</p><p>- [需求文档 | 如有]：</p><p>- [设计稿件 | 如有]：</p>`,
	})

	//开始节点 自动完成
	startNode := config.NewStartWorkFlowNode()

	//开发节点 进入切换状态
	node_0 := config.NewWorkFlowNode("开发", "state_0")
	node_0.DoneOperationDisplayName = "提交测试"
	node_0.FillDefaultReasonOptions()
	node_0.BelongStatus = config.WorkflowNodeCode_Started
	node_0.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_0.PassMode = config.WorkflowNodePassMode_Single
	node_0.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_0.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleDevelper.Id), Key: roleDevelper.Key, Uuid: roleDevelper.Uuid},
	}
	node_0.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_0.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_0.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(processState.Id),
			Key:  processState.Key,
			Uuid: processState.Uuid,
			Val:  processState.Val,
		},
	})

	//QA节点 进入切换状态
	node_1 := config.NewWorkFlowNode("测试", "state_1")
	node_1.Uuid = uuid.NewString()
	node_1.Name = "测试"
	node_1.Key = "state_1"
	node_1.DoneOperationDisplayName = "提交验收"
	node_1.EnableRollback = true
	node_1.FillDefaultReasonOptions()
	node_1.BelongStatus = config.WorkflowNodeCode_Started
	node_1.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_1.PassMode = config.WorkflowNodePassMode_Single
	node_1.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_1.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleQa.Id), Key: roleQa.Key, Uuid: roleQa.Uuid},
	}

	node_1.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_1.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_1.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(testingState.Id),
			Key:  testingState.Key,
			Uuid: testingState.Uuid,
			Val:  testingState.Val,
		},
	})

	//验收节点 进入切换状态
	node_2 := config.NewWorkFlowNode("验收", "state_2")
	node_2.DoneOperationDisplayName = "完成验收"
	node_2.FillDefaultReasonOptions()
	node_2.BelongStatus = config.WorkflowNodeCode_Started
	node_2.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_2.PassMode = config.WorkflowNodePassMode_Single
	node_2.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_2.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleAccepter.Id), Key: roleAccepter.Key, Uuid: roleAccepter.Uuid},
	}
	node_2.Owner.Value = &config.OwnerConf_UsageMode_None{
		FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
			{Type: config.FillOwnerType_Role, Value: "_creator"},
		},
	}

	node_2.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_2.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_2.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(checkingState.Id),
			Key:  checkingState.Key,
			Uuid: checkingState.Uuid,
			Val:  checkingState.Val,
		},
	})

	//完成节点 自动完成
	endNode := config.NewEndWorkFlowNode()
	endNode.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	endNode.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(completeState.Id),
			Key:  completeState.Key,
			Uuid: completeState.Uuid,
			Val:  completeState.Val,
		},
	})

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.Nodes = append(flowTpltConf.Nodes, startNode, node_0, node_1, node_2, endNode)

	//--节点关系
	flowTpltConf.Connections = make([]*config.WorkFlowConnection, 0)
	flowTpltConf.Connections = append(flowTpltConf.Connections,
		&config.WorkFlowConnection{StartNode: startNode.Key, EndNode: node_0.Key},
		&config.WorkFlowConnection{StartNode: node_0.Key, EndNode: node_1.Key},
		&config.WorkFlowConnection{StartNode: node_1.Key, EndNode: node_2.Key},
		&config.WorkFlowConnection{StartNode: node_2.Key, EndNode: endNode.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_WorkFlow, flowTpltConf, nil, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newBugWorkFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo
	processState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))
	checkingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowCheckingDefaultKey))
	testingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowTestingDefaultKey))
	confirmState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowWaitConfirmDefaultKey))

	roleInfo := req.WorkItemRoleInfo
	roleDevelper := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Developer)
	roleQa := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Qa)
	roleAccepter := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Acceptor)

	flowTpltConf := config.NewWorkFlow("bug")
	flowTpltConf.FormFields = append(flowTpltConf.FormFields, &config.WorkFlowFormField{
		Name: "describe", Value: `<p>- [当前版本/环境]： </p><p>- [问题描述]：</p><p>- [复现步骤]：</p><p>- [期望效果]：</p>`,
	})

	//开始节点 自动完成
	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)

	//开始节点 自动完成
	startNode := config.NewStartWorkFlowNode()
	//审查
	node_0 := config.NewWorkFlowNode("审查", "state_0")
	node_0.DoneOperationDisplayName = "确认问题"
	node_0.EnableClose = true
	node_0.FillDefaultReasonOptions()
	node_0.BelongStatus = config.WorkflowNodeCode_Started
	node_0.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_0.PassMode = config.WorkflowNodePassMode_Single
	node_0.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_0.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleQa.Id), Uuid: roleQa.Uuid, Key: roleQa.Key},
	}

	node_0.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_0.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_0.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(confirmState.Id),
			Key:  confirmState.Key,
			Uuid: confirmState.Uuid,
			Val:  confirmState.Val,
		},
	})

	//开发节点 进入切换状态
	node_1 := config.NewWorkFlowNode("开发", "state_1")
	node_1.DoneOperationDisplayName = "提交测试"
	node_1.FillDefaultReasonOptions()
	node_1.BelongStatus = config.WorkflowNodeCode_Started
	node_1.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_1.PassMode = config.WorkflowNodePassMode_Single
	node_1.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_1.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleDevelper.Id), Uuid: roleDevelper.Uuid, Key: roleDevelper.Key},
	}

	node_1.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_1.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_1.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(processState.Id),
			Key:  processState.Key,
			Uuid: processState.Uuid,
			Val:  processState.Val,
		},
	})

	//QA节点 进入切换状态
	node_2 := config.NewWorkFlowNode("测试", "state_2")
	node_2.DoneOperationDisplayName = "提交验收"
	node_2.EnableRollback = true
	node_2.FillDefaultReasonOptions()
	node_2.BelongStatus = config.WorkflowNodeCode_Started
	node_2.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_2.PassMode = config.WorkflowNodePassMode_Single
	node_2.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_2.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleQa.Id), Uuid: roleQa.Uuid, Key: roleQa.Key},
	}

	node_2.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_2.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_2.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(testingState.Id),
			Key:  testingState.Key,
			Uuid: testingState.Uuid,
			Val:  testingState.Val,
		},
	})

	//验收节点 进入切换状态
	node_3 := config.NewWorkFlowNode("验收", "state_3")
	node_3.DoneOperationDisplayName = "完成验收"
	node_3.FillDefaultReasonOptions()
	node_3.BelongStatus = config.WorkflowNodeCode_Started
	node_3.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_3.PassMode = config.WorkflowNodePassMode_Single
	node_3.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_3.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleAccepter.Id), Uuid: roleAccepter.Uuid, Key: roleAccepter.Key},
	}
	node_3.Owner.Value = &config.OwnerConf_UsageMode_None{
		FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
			{Type: config.FillOwnerType_Role, Value: "_creator"},
		},
	}

	node_3.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_3.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_3.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(checkingState.Id),
			Key:  checkingState.Key,
			Uuid: checkingState.Uuid,
			Val:  checkingState.Val,
		},
	})

	//完成节点 自动完成
	endNode := config.NewEndWorkFlowNode()
	endNode.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	endNode.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(completeState.Id),
			Key:  completeState.Key,
			Uuid: completeState.Uuid,
			Val:  completeState.Val,
		},
	})

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.Nodes = append(flowTpltConf.Nodes, startNode, node_0, node_1, node_2, node_3, endNode)

	//--节点关系
	flowTpltConf.Connections = make([]*config.WorkFlowConnection, 0)
	flowTpltConf.Connections = append(flowTpltConf.Connections,
		&config.WorkFlowConnection{StartNode: startNode.Key, EndNode: node_0.Key},
		&config.WorkFlowConnection{StartNode: node_0.Key, EndNode: node_1.Key},
		&config.WorkFlowConnection{StartNode: node_1.Key, EndNode: node_2.Key},
		&config.WorkFlowConnection{StartNode: node_2.Key, EndNode: node_3.Key},
		&config.WorkFlowConnection{StartNode: node_3.Key, EndNode: endNode.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_WorkFlow, flowTpltConf, nil, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newZouChaWorkFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo
	processState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))
	checkingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowCheckingDefaultKey))
	testingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowTestingDefaultKey))
	confirmState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowWaitConfirmDefaultKey))

	roleInfo := req.WorkItemRoleInfo
	roleDevelper := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Developer)
	roleQa := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Qa)
	roleAccepter := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Acceptor)

	flowTpltConf := config.NewWorkFlow("zoucha")
	flowTpltConf.FormFields = append(flowTpltConf.FormFields, &config.WorkFlowFormField{
		Name: "describe", Value: `<p>- [当前版本/环境]： </p><p>- [问题描述]：</p><p>- [复现步骤]：</p><p>- [期望效果]：</p>`,
	})

	//开始节点 自动完成
	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)

	//开始节点 自动完成
	startNode := config.NewStartWorkFlowNode()

	//审查
	node_0 := config.NewWorkFlowNode("审查", "state_0")
	node_0.DoneOperationDisplayName = "确认问题"
	node_0.EnableClose = true
	node_0.FillDefaultReasonOptions()
	node_0.BelongStatus = config.WorkflowNodeCode_Started
	node_0.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_0.PassMode = config.WorkflowNodePassMode_Single
	node_0.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_0.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleQa.Id), Uuid: roleQa.Uuid, Key: roleQa.Key},
	}

	node_0.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_0.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_0.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(confirmState.Id),
			Key:  confirmState.Key,
			Uuid: confirmState.Uuid,
			Val:  confirmState.Val,
		},
	})

	//开发节点 进入切换状态
	node_1 := config.NewWorkFlowNode("开发", "state_1")
	node_1.DoneOperationDisplayName = "提交测试"
	node_1.FillDefaultReasonOptions()
	node_1.BelongStatus = config.WorkflowNodeCode_Started
	node_1.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_1.PassMode = config.WorkflowNodePassMode_Single
	node_1.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_1.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleDevelper.Id), Uuid: roleDevelper.Uuid, Key: roleDevelper.Key},
	}

	node_1.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_1.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_1.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(processState.Id),
			Key:  processState.Key,
			Uuid: processState.Uuid,
			Val:  processState.Val,
		},
	})

	//QA节点 进入切换状态
	node_2 := config.NewWorkFlowNode("测试", "state_2")
	node_2.DoneOperationDisplayName = "提交验收"
	node_2.EnableRollback = true
	node_2.FillDefaultReasonOptions()
	node_2.BelongStatus = config.WorkflowNodeCode_Started
	node_2.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_2.PassMode = config.WorkflowNodePassMode_Single
	node_2.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_2.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleQa.Id), Uuid: roleQa.Uuid, Key: roleQa.Key},
	}

	node_2.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_2.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_2.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(testingState.Id),
			Key:  testingState.Key,
			Uuid: testingState.Uuid,
			Val:  testingState.Val,
		},
	})

	//验收节点 进入切换状态
	node_3 := config.NewWorkFlowNode("验收", "state_3")
	node_3.DoneOperationDisplayName = "完成验收"
	node_3.FillDefaultReasonOptions()
	node_3.BelongStatus = config.WorkflowNodeCode_Started
	node_3.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_3.PassMode = config.WorkflowNodePassMode_Single
	node_3.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_3.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleAccepter.Id), Uuid: roleAccepter.Uuid, Key: roleAccepter.Key},
	}
	node_3.Owner.Value = &config.OwnerConf_UsageMode_None{
		FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
			{Type: config.FillOwnerType_Role, Value: "_creator"},
		},
	}

	node_3.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_3.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_3.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(checkingState.Id),
			Key:  checkingState.Key,
			Uuid: checkingState.Uuid,
			Val:  checkingState.Val,
		},
	})

	//完成节点 自动完成
	endNode := config.NewEndWorkFlowNode()
	endNode.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	endNode.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(completeState.Id),
			Key:  completeState.Key,
			Uuid: completeState.Uuid,
			Val:  completeState.Val,
		},
	})

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.Nodes = append(flowTpltConf.Nodes, startNode, node_0, node_1, node_2, node_3, endNode)

	//--节点关系
	flowTpltConf.Connections = make([]*config.WorkFlowConnection, 0)
	flowTpltConf.Connections = append(flowTpltConf.Connections,
		&config.WorkFlowConnection{StartNode: startNode.Key, EndNode: node_0.Key},
		&config.WorkFlowConnection{StartNode: node_0.Key, EndNode: node_1.Key},
		&config.WorkFlowConnection{StartNode: node_1.Key, EndNode: node_2.Key},
		&config.WorkFlowConnection{StartNode: node_2.Key, EndNode: node_3.Key},
		&config.WorkFlowConnection{StartNode: node_3.Key, EndNode: endNode.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_WorkFlow, flowTpltConf, nil, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newDesignWorkFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo

	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))
	checkingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowCheckingDefaultKey))
	designingState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowDesigningDefaultKey))
	planningState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkflowPlanningDefaultKey))

	roleInfo := req.WorkItemRoleInfo
	roleAcceptor := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Acceptor)
	roleProducer := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_Producer)
	roleUiDesigner := roleInfo.GetRoleByKey(consts.WorkflowOwnerRole_UIDesigner)

	/*
		设计类型任务流程

		>“开始”，（新建/创建中...）
		>“产品”，（策划中，提交设计）
		>“美术”，（设计中，回滚/提交审核）
		>“产品”，（审核中，回滚/完成审核）
		>“完成”，（已完成）（已编辑）
	*/

	flowTpltConf := config.NewWorkFlow("design")
	flowTpltConf.FormFields = append(flowTpltConf.FormFields, &config.WorkFlowFormField{
		Name: "describe", Value: `<p>- [设计描述]： </p><p>- [思维导图 | 如有]：</p>`,
	})

	//开始节点 自动完成
	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)

	//开始节点 自动完成
	startNode := config.NewStartWorkFlowNode()

	//审查
	node_0 := config.NewWorkFlowNode("策划", "state_0")
	node_0.DoneOperationDisplayName = "提交设计"
	node_0.FillDefaultReasonOptions()
	node_0.BelongStatus = config.WorkflowNodeCode_Started
	node_0.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_0.PassMode = config.WorkflowNodePassMode_Single
	node_0.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_0.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleProducer.Id), Uuid: roleProducer.Uuid, Key: roleProducer.Key},
	}

	node_0.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_0.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_0.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(planningState.Id),
			Key:  planningState.Key,
			Uuid: planningState.Uuid,
			Val:  planningState.Val,
		},
	})

	//开发节点 进入切换状态
	node_1 := config.NewWorkFlowNode("设计", "state_1")
	node_1.DoneOperationDisplayName = "提交验收"
	node_1.EnableRollback = true
	node_1.FillDefaultReasonOptions()
	node_1.RollbackReasonOptions = []string{
		"策划 基础功能 的完成度，不足以完成美术设计",
		"策划 交互操作逻辑 的完成度，不足以完成美术设计",
	}
	node_1.BelongStatus = config.WorkflowNodeCode_Started
	node_1.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_1.PassMode = config.WorkflowNodePassMode_Single
	node_1.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_1.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleUiDesigner.Id), Uuid: roleUiDesigner.Uuid, Key: roleUiDesigner.Key},
	}

	node_1.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_1.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_1.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(designingState.Id),
			Key:  designingState.Key,
			Uuid: designingState.Uuid,
			Val:  designingState.Val,
		},
	})

	//QA节点 进入切换状态
	node_2 := config.NewWorkFlowNode("验收", "state_2")
	node_2.DoneOperationDisplayName = "完成验收"
	node_2.EnableRollback = true
	node_2.FillDefaultReasonOptions()
	node_2.RollbackReasonOptions = []string{
		"设计稿完成度未达可交付标准",
	}
	node_2.BelongStatus = config.WorkflowNodeCode_Started
	node_2.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_2.PassMode = config.WorkflowNodePassMode_Single
	node_2.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_2.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Id: cast.ToString(roleAcceptor.Id), Uuid: roleAcceptor.Uuid, Key: roleAcceptor.Key},
	}
	node_2.Owner.Value = &config.OwnerConf_UsageMode_None{
		FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
			{Type: config.FillOwnerType_Role, Value: "_creator"},
		},
	}

	node_2.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_2.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_2.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(checkingState.Id),
			Key:  checkingState.Key,
			Uuid: checkingState.Uuid,
			Val:  checkingState.Val,
		},
	})

	//完成节点 自动完成
	endNode := config.NewEndWorkFlowNode()
	endNode.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	endNode.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(completeState.Id),
			Key:  completeState.Key,
			Uuid: completeState.Uuid,
			Val:  completeState.Val,
		},
	})

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.Nodes = append(flowTpltConf.Nodes, startNode, node_0, node_1, node_2, endNode)

	//--节点关系
	flowTpltConf.Connections = make([]*config.WorkFlowConnection, 0)
	flowTpltConf.Connections = append(flowTpltConf.Connections,
		&config.WorkFlowConnection{StartNode: startNode.Key, EndNode: node_0.Key},
		&config.WorkFlowConnection{StartNode: node_0.Key, EndNode: node_1.Key},
		&config.WorkFlowConnection{StartNode: node_1.Key, EndNode: node_2.Key},
		&config.WorkFlowConnection{StartNode: node_2.Key, EndNode: endNode.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_WorkFlow, flowTpltConf, nil, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newSubTaskStateFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo
	processState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))

	flowTpltConf := config.NewStateFlow("sub_task")
	flowTpltConf.ResumeReasonOptions = nil
	flowTpltConf.TerminatedReasonOptions = nil

	flowTpltConf.StateFlowNodes = make([]*config.StateFlowNode, 0)
	flowTpltConf.StateFlowNodes = append(flowTpltConf.StateFlowNodes,
		&config.StateFlowNode{
			Key: "PROCESSING", Name: "进行中", IsInitState: true,
			SubStateId:   cast.ToString(processState.Id),
			SubStateVal:  processState.Val,
			SubStateKey:  processState.Key,
			SubStateUuid: processState.Uuid,
		},
		&config.StateFlowNode{
			Key: "COMPLETED", Name: "已完成", IsArchivedState: true,
			SubStateId:   cast.ToString(completeState.Id),
			SubStateVal:  completeState.Val,
			SubStateKey:  completeState.Key,
			SubStateUuid: completeState.Uuid,
		},
	)

	flowTpltConf.StateFlowTransitionRule = make([]*config.StateFlowTransitionRule, 0)
	flowTpltConf.StateFlowTransitionRule = append(flowTpltConf.StateFlowTransitionRule,
		&config.StateFlowTransitionRule{SourceStateKey: "PROCESSING", TargetStateKey: "COMPLETED"},
		&config.StateFlowTransitionRule{SourceStateKey: "COMPLETED", TargetStateKey: "PROCESSING"},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_StateFlow, nil, flowTpltConf, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newIssueStateFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	roleInfo := req.WorkItemRoleInfo
	issueOperator := roleInfo.GetRoleByKey(consts.StateflowOwnerRole_Operator)
	issueReporter := roleInfo.GetRoleByKey(consts.StateflowOwnerRole_Reporter)
	issueReviewer := roleInfo.GetRoleByKey(consts.StateflowOwnerRole_Reviewer)

	statusInfo := req.WorkItemStatusInfo
	stPendingStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowPendingDefaultKey))
	stFixingStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowFixingDefaultKey))
	stPendingVerificationStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowPendingVerificationDefaultKey))
	stClosedStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowClosedDefaultKey))
	stRestartStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowRestartDefaultKey))
	stConvertToStoryStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowConvertToStoryDefaultKey))
	stDoNotProcessStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowDoNotProcessDefaultKey))

	flowTpltConf := config.NewStateFlow("issue")

	issueReporterOwnerConf := &config.OwnerConf{
		UsageMode:  config.UsageMode_None,
		ForceOwner: true,
		OwnerRole: []*config.OwnerConf_Role{
			{Id: cast.ToString(issueReporter.Id), Uuid: issueReporter.Uuid, Key: issueReporter.Key},
		},
		Value: &config.OwnerConf_UsageMode_None{
			FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
				{Type: config.FillOwnerType_Role, Value: "_creator"},
			},
		},
	}

	issueOperatorOwnerConf := &config.OwnerConf{
		UsageMode:  config.UsageMode_None,
		ForceOwner: true,
		OwnerRole: []*config.OwnerConf_Role{
			{Id: cast.ToString(issueOperator.Id), Uuid: issueOperator.Uuid, Key: issueOperator.Key},
		},
	}

	issueReviewerOwnerConf := &config.OwnerConf{
		UsageMode:  config.UsageMode_None,
		ForceOwner: true,
		OwnerRole: []*config.OwnerConf_Role{
			{Id: cast.ToString(issueReviewer.Id), Uuid: issueReviewer.Uuid, Key: issueReviewer.Key},
		},
	}

	flowTpltConf.StateFlowNodes = make([]*config.StateFlowNode, 0)
	flowTpltConf.StateFlowNodes = append(flowTpltConf.StateFlowNodes,
		&config.StateFlowNode{
			IsInitState: true,
			Key:         stPendingStatus.Key, Name: stPendingStatus.Name,
			SubStateId:      cast.ToString(stPendingStatus.Id),
			SubStateVal:     stPendingStatus.Val,
			SubStateKey:     stPendingStatus.Key,
			SubStateUuid:    stPendingStatus.Uuid,
			IsArchivedState: stPendingStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReviewerOwnerConf,
		},
		&config.StateFlowNode{
			Key: stFixingStatus.Key, Name: stFixingStatus.Name,
			SubStateId:      cast.ToString(stFixingStatus.Id),
			SubStateVal:     stFixingStatus.Val,
			SubStateKey:     stFixingStatus.Key,
			SubStateUuid:    stFixingStatus.Uuid,
			IsArchivedState: stFixingStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueOperatorOwnerConf,
		},
		&config.StateFlowNode{
			Key: stPendingVerificationStatus.Key, Name: stPendingVerificationStatus.Name,
			SubStateId:      cast.ToString(stPendingVerificationStatus.Id),
			SubStateVal:     stPendingVerificationStatus.Val,
			SubStateKey:     stPendingVerificationStatus.Key,
			SubStateUuid:    stPendingVerificationStatus.Uuid,
			IsArchivedState: stPendingVerificationStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReporterOwnerConf,
		},
		&config.StateFlowNode{
			Key: stClosedStatus.Key, Name: stClosedStatus.Name,
			SubStateId:      cast.ToString(stClosedStatus.Id),
			SubStateVal:     stClosedStatus.Val,
			SubStateKey:     stClosedStatus.Key,
			SubStateUuid:    stClosedStatus.Uuid,
			IsArchivedState: stClosedStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReporterOwnerConf,
		},
		&config.StateFlowNode{
			Key: stRestartStatus.Key, Name: stRestartStatus.Name,
			SubStateId:      cast.ToString(stRestartStatus.Id),
			SubStateVal:     stRestartStatus.Val,
			SubStateKey:     stRestartStatus.Key,
			SubStateUuid:    stRestartStatus.Uuid,
			IsArchivedState: stRestartStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReviewerOwnerConf,
		},
		&config.StateFlowNode{
			Key: stConvertToStoryStatus.Key, Name: stConvertToStoryStatus.Name,
			SubStateId:      cast.ToString(stConvertToStoryStatus.Id),
			SubStateVal:     stConvertToStoryStatus.Val,
			SubStateKey:     stConvertToStoryStatus.Key,
			SubStateUuid:    stConvertToStoryStatus.Uuid,
			IsArchivedState: stConvertToStoryStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReporterOwnerConf,
		},
		&config.StateFlowNode{
			Key: stDoNotProcessStatus.Key, Name: stDoNotProcessStatus.Name,
			SubStateId:      cast.ToString(stDoNotProcessStatus.Id),
			SubStateVal:     stDoNotProcessStatus.Val,
			SubStateKey:     stDoNotProcessStatus.Key,
			SubStateUuid:    stDoNotProcessStatus.Uuid,
			IsArchivedState: stDoNotProcessStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReporterOwnerConf,
		},
	)

	flowTpltConf.StateFlowTransitionRule = make([]*config.StateFlowTransitionRule, 0)
	flowTpltConf.StateFlowTransitionRule = append(flowTpltConf.StateFlowTransitionRule,
		//待确认 Pending -> 修复中,不予处理,转需求
		&config.StateFlowTransitionRule{SourceStateKey: stPendingStatus.Key, TargetStateKey: stFixingStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stPendingStatus.Key, TargetStateKey: stDoNotProcessStatus.Key, ConfirmForm: []string{
			"重复 BUG",
			"符合需求",
			"暂时搁置",
			"无需解决",
		}},
		&config.StateFlowTransitionRule{SourceStateKey: stPendingStatus.Key, TargetStateKey: stConvertToStoryStatus.Key},
		//修复中 Fixing -> 待验证,转需求
		&config.StateFlowTransitionRule{SourceStateKey: stFixingStatus.Key, TargetStateKey: stPendingVerificationStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stFixingStatus.Key, TargetStateKey: stConvertToStoryStatus.Key},
		//待验证 Pending_Verification -> 关闭，修复中
		&config.StateFlowTransitionRule{SourceStateKey: stPendingVerificationStatus.Key, TargetStateKey: stClosedStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stPendingVerificationStatus.Key, TargetStateKey: stFixingStatus.Key},
		//关闭 Closed -》 重启
		&config.StateFlowTransitionRule{SourceStateKey: stClosedStatus.Key, TargetStateKey: stRestartStatus.Key},
		//重启 Restart -》 修复中， 不予处理
		&config.StateFlowTransitionRule{SourceStateKey: stRestartStatus.Key, TargetStateKey: stFixingStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stRestartStatus.Key, TargetStateKey: stDoNotProcessStatus.Key},
		//转需求 Convert_To_Story -》 关闭，不予处理
		&config.StateFlowTransitionRule{SourceStateKey: stConvertToStoryStatus.Key, TargetStateKey: stClosedStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stConvertToStoryStatus.Key, TargetStateKey: stDoNotProcessStatus.Key},
		//不予处理 Do_Not_Process-> 关闭，重启
		&config.StateFlowTransitionRule{SourceStateKey: stDoNotProcessStatus.Key, TargetStateKey: stClosedStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stDoNotProcessStatus.Key, TargetStateKey: stRestartStatus.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_StateFlow, nil, flowTpltConf, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newDefaultWorkFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	statusInfo := req.WorkItemStatusInfo

	processState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowProgressingDefaultKey))
	completeState := statusInfo.GetItemByKey(string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey))

	roleInfo := req.WorkItemRoleInfo.GetRoles()
	flowTpltConf := config.NewWorkFlow("tplt_" + rand.Letters(5))

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)

	//开始节点 自动完成
	startNode := config.NewStartWorkFlowNode()

	//开发节点 进入切换状态
	node_0 := config.NewWorkFlowNode("未命名节点", "state_0")
	node_0.DoneOperationDisplayName = "确认完成"
	node_0.FillDefaultReasonOptions()
	node_0.BelongStatus = config.WorkflowNodeCode_Started
	node_0.StartMode = config.WorkflowNodeStartMode_PreAllDone
	node_0.PassMode = config.WorkflowNodePassMode_Single
	node_0.Owner = &config.OwnerConf{UsageMode: config.UsageMode_None, ForceOwner: true}
	node_0.Owner.OwnerRole = []*config.OwnerConf_Role{
		{Key: roleInfo[0].Key, Id: cast.ToString(roleInfo[0].Id), Uuid: roleInfo[0].Uuid},
	}
	node_0.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}
	node_0.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	node_0.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(processState.Id),
			Key:  processState.Key,
			Uuid: processState.Uuid,
			Val:  processState.Val,
		},
	})

	//完成节点 自动完成
	endNode := config.NewEndWorkFlowNode()
	endNode.OnReach = make([]*config.WorkFlowNodeEvent, 0)
	endNode.OnReach = append(startNode.OnReach, &config.WorkFlowNodeEvent{
		EventType: "changeStoryStage",
		TargetSubState: &config.WorkFlowSubState{
			Id:   cast.ToString(completeState.Id),
			Key:  completeState.Key,
			Uuid: completeState.Uuid,
			Val:  completeState.Val,
		},
	})

	flowTpltConf.Nodes = make([]*config.WorkFlowNode, 0)
	flowTpltConf.Nodes = append(flowTpltConf.Nodes, startNode, node_0, endNode)

	//--节点关系
	flowTpltConf.Connections = make([]*config.WorkFlowConnection, 0)
	flowTpltConf.Connections = append(flowTpltConf.Connections,
		&config.WorkFlowConnection{StartNode: startNode.Key, EndNode: node_0.Key},
		&config.WorkFlowConnection{StartNode: node_0.Key, EndNode: endNode.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_WorkFlow, flowTpltConf, nil, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) newDefaultStateFlowTemplate(ctx context.Context, req *GenerateWorkFlowTemplateReq) *domain.WorkFlowTemplate {

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil
	}

	roleInfo := req.WorkItemRoleInfo
	issueOperator := roleInfo.GetRoleByKey(consts.StateflowOwnerRole_Operator)
	issueReporter := roleInfo.GetRoleByKey(consts.StateflowOwnerRole_Reporter)

	statusInfo := req.WorkItemStatusInfo
	stProgressingStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowProgressingDefaultKey))
	stClosedStatus := statusInfo.GetItemByKey(string(consts.WorkItemStatus_StateflowClosedDefaultKey))

	flowTpltConf := config.NewStateFlow("tplt_st_" + rand.Letters(5))

	issueReporterOwnerConf := &config.OwnerConf{
		UsageMode:  config.UsageMode_None,
		ForceOwner: true,
		OwnerRole: []*config.OwnerConf_Role{
			{Id: cast.ToString(issueReporter.Id), Uuid: issueReporter.Uuid, Key: issueReporter.Key},
		},
		Value: &config.OwnerConf_UsageMode_None{
			FillOwner: []*config.OwnerConf_UsageMode_FillOwner{
				{Type: config.FillOwnerType_Role, Value: "_creator"},
			},
		},
	}

	issueOperatorOwnerConf := &config.OwnerConf{
		UsageMode:  config.UsageMode_None,
		ForceOwner: true,
		OwnerRole: []*config.OwnerConf_Role{
			{Id: cast.ToString(issueOperator.Id), Uuid: issueOperator.Uuid, Key: issueOperator.Key},
		},
	}

	flowTpltConf.StateFlowNodes = make([]*config.StateFlowNode, 0)
	flowTpltConf.StateFlowNodes = append(flowTpltConf.StateFlowNodes,
		&config.StateFlowNode{
			IsInitState: true,
			Key:         stProgressingStatus.Key, Name: stProgressingStatus.Name,
			SubStateId:      cast.ToString(stProgressingStatus.Id),
			SubStateVal:     stProgressingStatus.Val,
			SubStateKey:     stProgressingStatus.Key,
			SubStateUuid:    stProgressingStatus.Uuid,
			IsArchivedState: stProgressingStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueOperatorOwnerConf,
		},
		&config.StateFlowNode{
			Key: stClosedStatus.Key, Name: stClosedStatus.Name,
			SubStateId:      cast.ToString(stClosedStatus.Id),
			SubStateVal:     stClosedStatus.Val,
			SubStateKey:     stClosedStatus.Key,
			SubStateUuid:    stClosedStatus.Uuid,
			IsArchivedState: stClosedStatus.IsArchivedTypeState(),
			OperationRole:   []string{},
			Owner:           issueReporterOwnerConf,
		},
	)

	flowTpltConf.StateFlowTransitionRule = make([]*config.StateFlowTransitionRule, 0)
	flowTpltConf.StateFlowTransitionRule = append(flowTpltConf.StateFlowTransitionRule,
		&config.StateFlowTransitionRule{SourceStateKey: stProgressingStatus.Key, TargetStateKey: stClosedStatus.Key},
		&config.StateFlowTransitionRule{SourceStateKey: stClosedStatus.Key, TargetStateKey: stProgressingStatus.Key},
	)

	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, req.SpaceId, req.WorkItemTypeId, req.WorkFlowId, 1, consts.FlowMode_StateFlow, nil, flowTpltConf, domain.WorkFlowTemplateStatus_Enable, req.UserId, nil)
	return flowTplt
}

func (s *WorkFlowService) UpdateWorkFlowTemplateConf(ctx context.Context, workFlow *domain.WorkFlow, newConf *config.WorkFlow, oper shared.Oper) (*domain.WorkFlowTemplate, error) {

	//获取模版相关资料
	wfTplt, err := s.repo.GetFlowTemplate(ctx, workFlow.LastTemplateId)
	if err != nil {
		return nil, errs.Business(ctx, err)
	}

	//检查配置是否合法
	validator := wf_vaildate.NewWorkFlowConfigValidate(newConf).WithCtx(ctx)
	err = validator.Valid()
	if err != nil {
		return nil, err
	}

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil, errors.New("生成模板ID失败")
	}

	version := wfTplt.Version + 1
	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, wfTplt.SpaceId, wfTplt.WorkItemTypeId, wfTplt.WorkFlowId, version, wfTplt.FlowMode, newConf, nil, domain.WorkFlowTemplateStatus_Enable, oper.GetId(), nil)
	workFlow.UpdateLastTemplate(flowTplt)

	workFlow.AddMessage(oper, &domain_message.SaveWorkFlowTemplate{
		SpaceId:         workFlow.SpaceId,
		WorkFlowId:      workFlow.Id,
		WorkFlowName:    workFlow.Name,
		TemplateId:      tpltId.Id,
		TemplateVersion: int64(version),
	})

	return flowTplt, nil
}

func (s *WorkFlowService) UpdateStateFlowTemplateConf(ctx context.Context, workFlow *domain.WorkFlow, newConf *config.StateFlow, oper shared.Oper) (*domain.WorkFlowTemplate, error) {

	//获取模版相关资料
	wfTplt, err := s.repo.GetFlowTemplate(ctx, workFlow.LastTemplateId)
	if err != nil {
		return nil, errs.Business(ctx, err)
	}

	//检查配置是否合法
	validator := wf_vaildate.NewStateFlowConfigValidate(newConf).WithCtx(ctx)
	err = validator.Valid()
	if err != nil {
		return nil, err
	}

	tpltId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate)
	if tpltId == nil {
		return nil, errors.New("生成模板ID失败")
	}

	version := wfTplt.Version + 1
	flowTplt := domain.NewWorkFlowTemplate(tpltId.Id, wfTplt.SpaceId, wfTplt.WorkItemTypeId, wfTplt.WorkFlowId, version, wfTplt.FlowMode, nil, newConf, domain.WorkFlowTemplateStatus_Enable, oper.GetId(), nil)
	workFlow.UpdateLastTemplate(flowTplt)

	workFlow.AddMessage(oper, &domain_message.SaveWorkFlowTemplate{
		SpaceId:         workFlow.SpaceId,
		WorkFlowId:      workFlow.Id,
		WorkFlowName:    workFlow.Name,
		TemplateId:      tpltId.Id,
		TemplateVersion: int64(version),
	})

	return flowTplt, nil
}

func (s *WorkFlowService) FindWorkFlowTemplateByAppointedOwnerRule(ctx context.Context, spaceId int64, userId int64) ([]int64, error) {

	tpltIds, err := s.repo.SearchTaskWorkFlowLastTemplateByOwnerRule(ctx, spaceId, cast.ToString(userId))
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	if len(tpltIds) == 0 {
		return tpltIds, nil
	}

	var relationTpltIds []int64
	for _, v := range tpltIds {
		tplt, err := s.repo.GetWorkFlowTemplateFormMemoryCache(ctx, v)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		tpltConf := tplt.WorkFlowConf()
		if tpltConf == nil {
			continue
		}

		for _, nodeConf := range tpltConf.Nodes {
			if nodeConf.Owner == nil || !nodeConf.Owner.IsAppointedUsageMode() {
				continue
			}

			val := nodeConf.Owner.GetAppointedUsageModeVal()
			if val == nil || len(val.AppointedOwner) != 1 {
				continue
			}

			var isRelation bool
			for _, owner := range val.AppointedOwner {
				if owner.IsUserType() {
					ownerId, _ := owner.Value.(string)
					if ownerId == cast.ToString(userId) {
						isRelation = true
						break
					}
				}
			}

			if isRelation {
				relationTpltIds = append(relationTpltIds, tplt.Id)
				break
			}
		}

	}

	return relationTpltIds, nil
}
