package message

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_Admin_Add_User             shared.MessageType = "Domain_Message.Admin.Add.User"
	Message_Type_Admin_Cancel_User          shared.MessageType = "Domain_Message.Admin.Cancel.User"
	Message_Type_Admin_Change_User_Nickname shared.MessageType = "Domain_Message.Admin.Change.User.Nickname"
	Message_Type_Admin_Change_User_Role     shared.MessageType = "Domain_Message.Admin.Change.User.Role"
	Message_Type_Admin_Reset_User_Password  shared.MessageType = "Domain_Message.Admin.Reset.User.Password"
)

type AdminAddUser struct {
	shared.DomainMessageBase

	UserId int64
}

func (ops *AdminAddUser) MessageType() shared.MessageType {
	return Message_Type_Admin_Add_User
}

type AdminCancelUser struct {
	shared.DomainMessageBase

	UserId   int64
	Username string
	Nickname string
}

func (ops *AdminCancelUser) MessageType() shared.MessageType {
	return Message_Type_Admin_Cancel_User
}

type AdminChangeUserNickname struct {
	shared.DomainMessageBase

	UserId int64

	OldValue string
	NewValue string
}

func (ops *AdminChangeUserNickname) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_User_Nickname
}

type AdminChangeUserRole struct {
	shared.DomainMessageBase

	UserId   int64
	OldValue consts.SystemRole
	NewValue consts.SystemRole
}

func (ops *AdminChangeUserRole) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_User_Role
}

type AdminResetUserPassword struct {
	shared.DomainMessageBase

	UserId int64
}

func (ops *AdminResetUserPassword) MessageType() shared.MessageType {
	return Message_Type_Admin_Reset_User_Password
}
