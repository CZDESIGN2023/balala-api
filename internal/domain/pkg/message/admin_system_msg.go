package message

import shared "go-cs/internal/pkg/domain"

const (
	Message_Type_Admin_Change_System_Logo          shared.MessageType = "Domain_Message.Admin.Change.System.Logo"
	Message_Type_Admin_Change_System_Title         shared.MessageType = "Domain_Message.Admin.Change.System.Title"
	Message_Type_Admin_Change_System_AccessUrl     shared.MessageType = "Domain_Message.Admin.Change.System.AccessUrl"
	Message_Type_Admin_Change_System_LoginBg       shared.MessageType = "Domain_Message.Admin.Change.System.LoginBg"
	Message_Type_Admin_Change_System_RegisterEntry shared.MessageType = "Domain_Message.Admin.Change.System.LoginEntry"
	Message_Type_Admin_Change_System_AttachSize    shared.MessageType = "Domain_Message.Admin.Change.System.AttachSize"
)

type AdminChangeSystemLogo struct {
	shared.DomainMessageBase
}

func (ops *AdminChangeSystemLogo) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_Logo
}

type AdminChangeSystemTitle struct {
	shared.DomainMessageBase

	OldValue string
	NewValue string
}

func (ops *AdminChangeSystemTitle) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_Title
}

type AdminChangeSystemAccessUrl struct {
	shared.DomainMessageBase

	OldValue string
	NewValue string
}

func (ops *AdminChangeSystemAccessUrl) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_AccessUrl
}

type AdminChangeSystemLoginBg struct {
	shared.DomainMessageBase
}

func (ops *AdminChangeSystemLoginBg) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_LoginBg
}

type AdminChangeSystemRegisterEntry struct {
	shared.DomainMessageBase

	OldValue string
	NewValue string
}

func (ops *AdminChangeSystemRegisterEntry) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_RegisterEntry
}

type AdminChangeSystemAttachSize struct {
	shared.DomainMessageBase

	OldValue string
	NewValue string
}

func (ops *AdminChangeSystemAttachSize) MessageType() shared.MessageType {
	return Message_Type_Admin_Change_System_AttachSize
}
