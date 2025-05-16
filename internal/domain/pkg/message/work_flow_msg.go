package message

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_WorkFlow_Create        shared.MessageType = "Domain_Message.WorkFlow.Create"
	Message_Type_WorkFlow_Modify        shared.MessageType = "Domain_Message.WorkFlow.Modify"
	Message_Type_WorkFlow_Delete        shared.MessageType = "Domain_Message.WorkFlow.Delete"
	Message_Type_WorkFlow_Template_Save shared.MessageType = "Domain_Message.WorkFlow.Template_Save"
	Message_Type_Task_WorkFlow_Upgrade  shared.MessageType = "Domain_Message.Task.WorkFlow.Upgrade"
)

// --- 创建任务
type CreateWorkFlow struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkFlowName string
	WorkFlowId   int64
	FlowMode     consts.WorkFlowMode

	// 从其他项目复制
	SrcSpaceName    string
	SrcWorkFlowName string
}

func (ops *CreateWorkFlow) MessageType() shared.MessageType {
	return Message_Type_WorkFlow_Create
}

// --- 创建任务
type DeleteWorkFlow struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkFlowName string
	WorkFlowId   int64
}

func (ops *DeleteWorkFlow) MessageType() shared.MessageType {
	return Message_Type_WorkFlow_Delete
}

// --- 编辑工作流
type ModifyWorkFlow struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkFlowName string
	WorkFlowId   int64

	Updates []FieldUpdate
}

func (ops *ModifyWorkFlow) MessageType() shared.MessageType {
	return Message_Type_WorkFlow_Modify
}

// -- 保存模版
type SaveWorkFlowTemplate struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkFlowName string
	WorkFlowId   int64

	TemplateId      int64
	TemplateVersion int64
}

func (ops *SaveWorkFlowTemplate) MessageType() shared.MessageType {
	return Message_Type_WorkFlow_Template_Save
}

// -- 保存模版

type RoleDirector struct {
	RoleId    string
	RoleKey   string
	Directors []string
}

type UpgradeTaskWorkFlow struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkFlowName string
	WorkFlowId   int64
	WorkItemId   int64
	WorkItemName string

	OldVersion int32
	NewVersion int32

	NewRoles []RoleDirector
	OldRoles []RoleDirector
}

func (ops *UpgradeTaskWorkFlow) MessageType() shared.MessageType {
	return Message_Type_Task_WorkFlow_Upgrade
}
