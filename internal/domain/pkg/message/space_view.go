package message

import (
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_SpaceView_Create     shared.MessageType = "Domain_Message.SpaceView.Create"
	Message_Type_SpaceView_Delete     shared.MessageType = "Domain_Message.SpaceView.Delete"
	Message_Type_SpaceView_SetName    shared.MessageType = "Domain_Message.SpaceView.SetName"
	Message_Type_SpaceView_SetRanking shared.MessageType = "Domain_Message.SpaceView.SetRanking"
	Message_Type_SpaceView_SetStatus  shared.MessageType = "Domain_Message.SpaceView.SetStatus"

	Message_Type_SpaceView_Update shared.MessageType = "Domain_Message.SpaceView.Update"
)

// --- 创建任务状态
type CreateSpaceView struct {
	shared.DomainMessageBase
	SpaceId  int64
	ViewId   int64
	ViewType int64
	ViewName string
}

func (ops *CreateSpaceView) MessageType() shared.MessageType {
	return Message_Type_SpaceView_Create
}

// --- 创建任务状态
type DeleteSpaceView struct {
	shared.DomainMessageBase

	ViewId   int64
	SpaceId  int64
	ViewType int64
	ViewName string
}

func (ops *DeleteSpaceView) MessageType() shared.MessageType {
	return Message_Type_SpaceView_Delete
}

type SetSpaceViewName struct {
	shared.DomainMessageBase

	ViewId      int64
	SpaceId     int64
	ViewType    int64
	ViewOldName string
	ViewNewName string
}

func (ops *SetSpaceViewName) MessageType() shared.MessageType {
	return Message_Type_SpaceView_SetName
}

type SetSpaceViewRanking struct {
	shared.DomainMessageBase

	SpaceId  int64
	ViewId   int64
	ViewType int64
	ViewName string
	Ranking  int64
}

func (ops *SetSpaceViewRanking) MessageType() shared.MessageType {
	return Message_Type_SpaceView_SetRanking
}

type SetSpaceViewStatus struct {
	shared.DomainMessageBase
	SpaceId  int64
	ViewId   int64
	ViewType int64
	ViewName string
	Status   int64
}

func (ops *SetSpaceViewStatus) MessageType() shared.MessageType {
	return Message_Type_SpaceView_SetStatus
}

type UpdateSpaceView struct {
	shared.DomainMessageBase
	SpaceId  int64
	ViewId   int64
	ViewType int64
	ViewName string
	ViewKey  string
	Field    string
}

func (ops *UpdateSpaceView) MessageType() shared.MessageType {
	return Message_Type_SpaceView_Update
}
