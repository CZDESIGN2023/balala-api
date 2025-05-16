package consts

type WorkFlowMode string

const (
	//工作流程负责人角色

	WorkflowOwnerRole_Developer  = "_developer"    //开发者
	WorkflowOwnerRole_Qa         = "_qa"           //QA
	WorkflowOwnerRole_Acceptor   = "_accepter"     //验收人
	WorkflowOwnerRole_Producer   = "_productor"    //产品
	WorkflowOwnerRole_UIDesigner = "_ui_designner" //美术
	WorkflowOwnerRole_Reviewer   = "_reviewer"     //审核人
	WorkflowOwnerRole_Evaluator  = "_evaluator"    //评审组

	//工作流模式

	StateflowOwnerRole_Operator = "_st_operator" //经办人
	StateflowOwnerRole_Reporter = "_st_reporter" //报告人
	StateflowOwnerRole_Reviewer = "_st_reviewer" //审核人

	//工作流程节点状态

	FlowNodeStatus_Unknown     = 0
	FlowNodeStatus_Waiting     = 1 //等待中,未到达
	FlowNodeStatus_Progressing = 2 //进行中
	FlowNodeStatus_Finished    = 3 //完成

	FlowMode_StateFlow WorkFlowMode = "state_flow" //状态模式
	FlowMode_WorkFlow  WorkFlowMode = "work_flow"  //工作流模式
)

func IsValidRole(role string) bool {
	switch role {
	case WorkflowOwnerRole_Developer:
	case WorkflowOwnerRole_Qa:
	case WorkflowOwnerRole_Acceptor:
	case WorkflowOwnerRole_Producer:
	case WorkflowOwnerRole_UIDesigner:
	case WorkflowOwnerRole_Reviewer:
	default:
		return false
	}
	return true
}
