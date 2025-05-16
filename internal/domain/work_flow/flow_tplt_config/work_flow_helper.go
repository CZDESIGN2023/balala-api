package config

import "github.com/google/uuid"

func NewStartWorkFlowNode() *WorkFlowNode {
	return &WorkFlowNode{
		Uuid:                  uuid.NewString(),
		Name:                  "开始",
		Key:                   WorkflowNodeCode_Started,
		BelongStatus:          WorkflowNodeCode_Started,
		StartMode:             WorkflowNodeStartMode_PreAllDone,
		PassMode:              WorkflowNodePassMode_Auto,
		CloseReasonOptions:    make([]string, 0),
		RestartReasonOptions:  make([]string, 0),
		RollbackReasonOptions: make([]string, 0),
		OnReach:               make([]*WorkFlowNodeEvent, 0),
		OnPass:                make([]*WorkFlowNodeEvent, 0),
		Owner: &OwnerConf{
			OwnerRole: make([]*OwnerConf_Role, 0),
		},
	}
}

func NewEndWorkFlowNode() *WorkFlowNode {
	return &WorkFlowNode{
		Uuid:                  uuid.NewString(),
		Name:                  "完成",
		Key:                   WorkflowNodeCode_Ended,
		BelongStatus:          WorkflowNodeCode_Started,
		StartMode:             WorkflowNodeStartMode_PreAllDone,
		PassMode:              WorkflowNodePassMode_Auto,
		CloseReasonOptions:    make([]string, 0),
		RestartReasonOptions:  make([]string, 0),
		RollbackReasonOptions: make([]string, 0),
		OnReach:               make([]*WorkFlowNodeEvent, 0),
		OnPass:                make([]*WorkFlowNodeEvent, 0),
		Owner: &OwnerConf{
			OwnerRole: make([]*OwnerConf_Role, 0),
		},
	}
}
