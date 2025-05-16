package message

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_WorkItemStatus_Create shared.MessageType = "Domain_Message.WorkItemStatus.Create"
	Message_Type_WorkItemStatus_Modify shared.MessageType = "Domain_Message.WorkItemStatus.Modify"
	Message_Type_WorkItemStatus_Delete shared.MessageType = "Domain_Message.WorkItemStatus.Delete"
)

// --- 创建任务状态
type CreateWorkItemStatus struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemStatusName string
	WorkItemStatusId   int64
	FlowScope          consts.FlowScope
}

func (ops *CreateWorkItemStatus) MessageType() shared.MessageType {
	return Message_Type_WorkItemStatus_Create
}

// --- 创建任务状态
type DeleteWorkItemStatus struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemStatusName string
	WorkItemStatusId   int64
	FlowScope          consts.FlowScope
}

func (ops *DeleteWorkItemStatus) MessageType() shared.MessageType {
	return Message_Type_WorkItemStatus_Delete
}

// --- 编辑工作流
type ModifyWorkItemStatus struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemStatusName string
	WorkItemStatusId   int64
	Updates            []FieldUpdate
}

func (ops *ModifyWorkItemStatus) MessageType() shared.MessageType {
	return Message_Type_WorkItemStatus_Modify
}
