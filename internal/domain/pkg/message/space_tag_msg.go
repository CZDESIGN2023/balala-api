package message

import (
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_SpaceTag_Create shared.MessageType = "Domain_Message.SpaceTag.Create"
	Message_Type_SpaceTag_Modify shared.MessageType = "Domain_Message.SpaceTag.Modify"
	Message_Type_SpaceTag_Delete shared.MessageType = "Domain_Message.SpaceTag.Delete"
)

// --- 创建任务状态
type CreateSpaceTag struct {
	shared.DomainMessageBase

	SpaceId int64

	SpaceTagName string
	SpaceTagId   int64
}

func (ops *CreateSpaceTag) MessageType() shared.MessageType {
	return Message_Type_SpaceTag_Create
}

// --- 创建任务状态
type DeleteSpaceTag struct {
	shared.DomainMessageBase

	SpaceId int64

	SpaceTagName string
	SpaceTagId   int64
}

func (ops *DeleteSpaceTag) MessageType() shared.MessageType {
	return Message_Type_SpaceTag_Delete
}

// --- 编辑工作流
type ModifySpaceTag struct {
	shared.DomainMessageBase

	SpaceId int64

	SpaceTagName string
	SpaceTagId   int64
	Updates      []FieldUpdate
}

func (ops *ModifySpaceTag) MessageType() shared.MessageType {
	return Message_Type_SpaceTag_Modify
}
