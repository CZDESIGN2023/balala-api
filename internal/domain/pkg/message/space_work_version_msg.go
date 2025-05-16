package message

import shared "go-cs/internal/pkg/domain"

const (
	Message_Type_WorkVersion_Create shared.MessageType = "Domain_Message.WorkVersion.Create"
	Message_Type_WorkVersion_Modify shared.MessageType = "Domain_Message.WorkVersion.Modify"
	Message_Type_WorkVersion_Delete shared.MessageType = "Domain_Message.WorkVersion.Delete"
)

// --- 创建模块
type CreateWorkVersion struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkVersionName string
	WorkVersionId   int64
}

func (ops *CreateWorkVersion) MessageType() shared.MessageType {
	return Message_Type_WorkVersion_Create
}

// --- 创建模块
type DeleteWorkVersion struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkVersionName string
	WorkVersionId   int64
}

func (ops *DeleteWorkVersion) MessageType() shared.MessageType {
	return Message_Type_WorkVersion_Delete
}

// --- 编辑模块
type ModifyWorkVersion struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkVersionName string
	WorkVersionId   int64

	Updates []FieldUpdate
}

func (ops *ModifyWorkVersion) MessageType() shared.MessageType {
	return Message_Type_WorkVersion_Modify
}
