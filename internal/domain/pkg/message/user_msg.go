package message

import (
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_Personal_ChangeNickname         shared.MessageType = "Domain_Message.Personal.ChangeNickname"
	Message_Type_Personal_ResetPwd               shared.MessageType = "Domain_Message.Personal.ResetPwd"
	Message_Type_Personal_ChangeAvatar           shared.MessageType = "Domain_Message.Personal.ChangeAvatar"
	Message_Type_Personal_SetSpaceNotify         shared.MessageType = "Domain_Message.Personal.SetSpaceNotify"
	Message_Type_Personal_BindThirdPlatform      shared.MessageType = "Domain_Message.Personal.BindThirdPlatform"
	Message_Type_Personal_UnBindThirdPlatform    shared.MessageType = "Domain_Message.Personal.UnBindThirdPlatform"
	Message_Type_Personal_SetTempConfig          shared.MessageType = "Domain_Message.Personal.SetTempConfig"
	Message_Type_Personal_SetUserConfig          shared.MessageType = "Domain_Message.Personal.SetUserConfig"
	Message_Type_Personal_SetThirdPlatformNotify shared.MessageType = "Domain_Message.Personal.SetThirdPlatformNotify"
)

// --- 修改用户昵称
type PersonalChangeNickName struct {
	shared.DomainMessageBase

	UserId int64

	OldNickName string
	NewNickName string
}

func (ops *PersonalChangeNickName) MessageType() shared.MessageType {
	return Message_Type_Personal_ChangeNickname
}

// --- 修改用户密码
type PersonalResetPwd struct {
	shared.DomainMessageBase

	UserId       int64
	UserNickname string
}

func (ops *PersonalResetPwd) MessageType() shared.MessageType {
	return Message_Type_Personal_ResetPwd
}

type PersonalChangeAvatar struct {
	shared.DomainMessageBase

	UserId int64
}

func (ops *PersonalChangeAvatar) MessageType() shared.MessageType {
	return Message_Type_Personal_ChangeAvatar
}

type PersonalSetSpaceNotify struct {
	shared.DomainMessageBase

	SpaceId int64
	Notify  int
}

func (ops *PersonalSetSpaceNotify) MessageType() shared.MessageType {
	return Message_Type_Personal_SetSpaceNotify
}

type PersonalBindThirdPlatform struct {
	shared.DomainMessageBase

	PlatformName string
}

func (ops *PersonalBindThirdPlatform) MessageType() shared.MessageType {
	return Message_Type_Personal_BindThirdPlatform
}

type PersonalUnBindThirdPlatform struct {
	shared.DomainMessageBase

	PlatformName string
}

func (ops *PersonalUnBindThirdPlatform) MessageType() shared.MessageType {
	return Message_Type_Personal_UnBindThirdPlatform
}

type PersonalSetTempConfig struct {
	shared.DomainMessageBase

	OldValues map[string]string
	NewValues map[string]string
}

func (p *PersonalSetTempConfig) MessageType() shared.MessageType {
	return Message_Type_Personal_SetTempConfig
}

type PersonalSetUserConfig struct {
	shared.DomainMessageBase

	UserId int64

	Key      string
	OldValue string
	NewValue string
}

func (p *PersonalSetUserConfig) MessageType() shared.MessageType {
	return Message_Type_Personal_SetUserConfig
}

type PersonalSetThirdPlatformNotify struct {
	shared.DomainMessageBase

	PlatformCode int32
	Notify       int32
}

func (p *PersonalSetThirdPlatformNotify) MessageType() shared.MessageType {
	return Message_Type_Personal_SetThirdPlatformNotify
}
