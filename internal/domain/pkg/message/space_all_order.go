package message

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

const (
	Message_Type_SpaceAllOrder_WorkObject       shared.MessageType = "Domain_Message.SpaceAllOrder.WorkObject"
	Message_Type_SpaceAllOrder_Version          shared.MessageType = "Domain_Message.SpaceAllOrder.Version"
	Message_Type_SpaceAllOrder_WorkFlow         shared.MessageType = "Domain_Message.SpaceAllOrder.WorkFlow"
	Message_Type_SpaceAllOrder_Role             shared.MessageType = "Domain_Message.SpaceAllOrder.Role"
	Message_Type_SpaceAllOrder_Status           shared.MessageType = "Domain_Message.SpaceAllOrder.Status"
	Message_Type_SpaceAllOrder_OverviewDataItem shared.MessageType = "Domain_Message.SpaceAllOrder.OverviewDataItem"
	Message_Type_SpaceAllOrder_OverviewBlock    shared.MessageType = "Domain_Message.SpaceAllOrder.OverviewBlock"
	Message_Type_SpaceAllOrder_View             shared.MessageType = "Domain_Message.SpaceAllOrder.View"
)

type ChangeWorkObjectOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeWorkObjectOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_WorkObject
}

type ChangeVersionOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeVersionOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_Version
}

type ChangeWorkFlowOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeWorkFlowOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_WorkFlow
}

type ChangeRoleOrder struct {
	shared.DomainMessageBase

	SpaceId int64
	consts.FlowScope
}

func (ops *ChangeRoleOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_Role
}

type ChangeStatusOrder struct {
	shared.DomainMessageBase

	SpaceId   int64
	FlowScope consts.FlowScope
}

func (ops *ChangeStatusOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_Status
}

type ChangeOverviewDataItemOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeOverviewDataItemOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_OverviewDataItem
}

type ChangeOverviewBlockOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeOverviewBlockOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_OverviewBlock
}

type ChangeSpaceViewOrder struct {
	shared.DomainMessageBase

	SpaceId int64
}

func (ops *ChangeSpaceViewOrder) MessageType() shared.MessageType {
	return Message_Type_SpaceAllOrder_View
}
