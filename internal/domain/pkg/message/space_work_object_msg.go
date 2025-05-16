package message

import shared "go-cs/internal/pkg/domain"

const (
	Message_Type_WorkObject_Create shared.MessageType = "Domain_Message.WorkObject.Create"
	Message_Type_WorkObject_Modify shared.MessageType = "Domain_Message.WorkObject.Modify"
	Message_Type_WorkObject_Delete shared.MessageType = "Domain_Message.WorkObject.Delete"
)

// --- 创建模块
type CreateWorkObject struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkObjectName string
	WorkObjectId   int64
}

func (ops *CreateWorkObject) MessageType() shared.MessageType {
	return Message_Type_WorkObject_Create
}

// --- 创建模块
type DeleteWorkObject struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkObjectName string
	WorkObjectId   int64
}

func (ops *DeleteWorkObject) MessageType() shared.MessageType {
	return Message_Type_WorkObject_Delete
}

// --- 编辑模块
type ModifyWorkObject struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkObjectName string
	WorkObjectId   int64

	Updates []FieldUpdate
}

func (ops *ModifyWorkObject) MessageType() shared.MessageType {
	return Message_Type_WorkObject_Modify
}
