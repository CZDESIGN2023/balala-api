package message

import (
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_Space_Quit                            shared.MessageType = "Domain_Message.Space.QuitSpace"
	Message_Type_Space_AddMember                       shared.MessageType = "Domain_Message.Space.AddSpaceMember"
	Message_Type_Space_RemoveMember                    shared.MessageType = "Domain_Message.Space.RemoveSpaceMember"
	Message_Type_Space_SetMemberRole                   shared.MessageType = "Domain_Message.Space.SetSpaceMemberRole"
	Message_Type_Space_Create                          shared.MessageType = "Domain_Message.Space.CreateSpace"
	Message_Type_Space_Copy                            shared.MessageType = "Domain_Message.Space.CopySpace"
	Message_Type_Space_Modify                          shared.MessageType = "Domain_Message.Space.ModifySpace"
	Message_Type_Space_Del                             shared.MessageType = "Domain_Message.Space.DelSpace"
	Message_Type_Space_Transfer                        shared.MessageType = "Domain_Message.Space.TransferSpace"
	Message_Type_Space_SetNotify                       shared.MessageType = "Domain_Message.Space.SetSpaceNotify"
	Message_Type_Space_SetWorkingDay                   shared.MessageType = "Domain_Message.Space.SetWorkingDay"
	Message_Type_Space_SetCommentDeletable             shared.MessageType = "Domain_Message.Space.SetCommentDeletable"
	Message_Type_Space_SetCommentDeletableWhenArchived shared.MessageType = "Domain_Message.Space.SetCommentDeletableWhenArchived"
	Message_Type_Space_SetCommentShowPos               shared.MessageType = "Domain_Message.Space.SetCommentShowPos"
	Message_Type_Space_SetTempConfig                   shared.MessageType = "Domain_Message.Space.SetTempConfig"
)

// --------- 退出空间 ----------
type QuitSpace struct {
	shared.DomainMessageBase

	SpaceId          int64
	RoleId           int64
	MemberUid        int64
	WorkItemNum      int
	TransferTargetId int64
}

func (ops *QuitSpace) MessageType() shared.MessageType {
	return Message_Type_Space_Quit
}

// ------ 空间添加成员

type AddSpaceMember struct {
	shared.DomainMessageBase

	SpaceId   int64
	RoleId    int64
	MemberUid int64
}

func (ops *AddSpaceMember) MessageType() shared.MessageType {
	return Message_Type_Space_AddMember
}

// ------ 空间添加成员

type RemoveSpaceMember struct {
	shared.DomainMessageBase

	SpaceId          int64
	MemberUid        int64
	WorkItemNum      int
	TransferTargetId int64
}

func (ops *RemoveSpaceMember) MessageType() shared.MessageType {
	return Message_Type_Space_RemoveMember
}

// ---- 设置空间角色权限
type SetSpaceMemberRole struct {
	shared.DomainMessageBase

	SpaceId    int64
	MemeberUid int64
	OldRole    int32
	NewRole    int32
}

func (ops *SetSpaceMemberRole) MessageType() shared.MessageType {
	return Message_Type_Space_SetMemberRole
}

// --- 创建空间
type CreateSpace struct {
	shared.DomainMessageBase

	SpaceId   int64
	SpaceName string

	SrcSpaceId   int64
	SrcSpaceName string
}

func (ops *CreateSpace) MessageType() shared.MessageType {
	return Message_Type_Space_Create
}

// -- 修改空间信息
type ModifySpace struct {
	shared.DomainMessageBase

	SpaceId   int64
	SpaceName string

	Updates []FieldUpdate
}

func (p *ModifySpace) MessageType() shared.MessageType {
	return Message_Type_Space_Modify
}

// --- 删除空间
type DelSpace struct {
	shared.DomainMessageBase

	SpaceId   int64
	SpaceName string
}

func (p *DelSpace) MessageType() shared.MessageType {
	return Message_Type_Space_Del
}

// --- 转移项目
type TransferSpace struct {
	shared.DomainMessageBase

	SpaceId   int64
	SpaceName string

	UserId       int64
	TargetUserId int64
}

func (p *TransferSpace) MessageType() shared.MessageType {
	return Message_Type_Space_Transfer
}

// --- 设置通知开关
type SetSpaceNotify struct {
	shared.DomainMessageBase

	SpaceId   int64
	SpaceName string

	Notify int
}

func (p *SetSpaceNotify) MessageType() shared.MessageType {
	return Message_Type_Space_SetNotify
}

// --- 设置工作日
type SetWorkingDay struct {
	shared.DomainMessageBase

	SpaceId     int64
	WeekDays    []int64
	OldWeekDays []int64
}

func (p *SetWorkingDay) MessageType() shared.MessageType {
	return Message_Type_Space_SetWorkingDay
}

type SetCommentDeletable struct {
	shared.DomainMessageBase

	SpaceId   int64
	Deletable int64
}

func (p *SetCommentDeletable) MessageType() shared.MessageType {
	return Message_Type_Space_SetCommentDeletable
}

type SetCommentDeletableWhenArchived struct {
	shared.DomainMessageBase

	SpaceId int64
	Value   int64
}

func (p *SetCommentDeletableWhenArchived) MessageType() shared.MessageType {
	return Message_Type_Space_SetCommentDeletableWhenArchived
}

type SetCommentShowPos struct {
	shared.DomainMessageBase

	SpaceId int64
	Value   int64
}

func (p *SetCommentShowPos) MessageType() shared.MessageType {
	return Message_Type_Space_SetCommentShowPos
}

type SetSpaceTempConfig struct {
	shared.DomainMessageBase

	SpaceId int64

	OldValues map[string]string
	NewValues map[string]string
}

func (p *SetSpaceTempConfig) MessageType() shared.MessageType {
	return Message_Type_Space_SetTempConfig
}
