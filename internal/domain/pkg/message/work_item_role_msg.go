package message

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_WorkItemRole_Create shared.MessageType = "Domain_Message.WorkItemRole.Create"
	Message_Type_WorkItemRole_Modify shared.MessageType = "Domain_Message.WorkItemRole.Modify"
	Message_Type_WorkItemRole_Delete shared.MessageType = "Domain_Message.WorkItemRole.Delete"
)

// --- 创建流程角色
type CreateWorkItemRole struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemRoleName string
	WorkItemRoleId   int64
	FlowScope        consts.FlowScope
}

func (ops *CreateWorkItemRole) MessageType() shared.MessageType {
	return Message_Type_WorkItemRole_Create
}

// --- 创建流程角色
type DeleteWorkItemRole struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemRoleName string
	WorkItemRoleId   int64
	FlowScope        consts.FlowScope
}

func (ops *DeleteWorkItemRole) MessageType() shared.MessageType {
	return Message_Type_WorkItemRole_Delete
}

// --- 编辑工作流
type ModifyWorkItemRole struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemRoleName string
	WorkItemRoleId   int64
	Updates          []FieldUpdate
}

func (ops *ModifyWorkItemRole) MessageType() shared.MessageType {
	return Message_Type_WorkItemRole_Modify
}
