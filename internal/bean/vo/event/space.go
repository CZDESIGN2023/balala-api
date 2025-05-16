package event

import (
	"go-cs/api/notify"
	"go-cs/internal/domain/space"
	"go-cs/internal/domain/space_member"
	"go-cs/internal/domain/space_tag"
)

type AddMember struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
	RoleId   int64
}

type RemoveMember struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
}

type ChangeRole struct {
	Event     notify.Event
	Space     *space.Space
	Operator  int64
	TargetId  int64
	OldRoleId int32
	NewRoleId int32
}

type ChangeSpaceName struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	OldValue string
	NewValue string
}

type ChangeSpaceDescribe struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	OldValue string
	NewValue string
}

type CreateSpace struct {
	Event    notify.Event
	Operator int64
	Space    *space.Space
	Members  []*space_member.SpaceMember
}

type DeleteSpace struct {
	Event     notify.Event
	Space     *space.Space
	Operator  int64
	MemberIds []int64
}

type QuitSpace struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
}

type TransferWorkItem struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
	Num      int //被转移的任务数
}

type TransferSpace struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
}

type SetSpaceNotify struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	Notify   int64
}

type SetSpaceWorkingDay struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	WeekDays []int64
}

type SetCommentDeletable struct {
	Event     notify.Event
	Space     *space.Space
	Operator  int64
	Deletable int64
}

type SetCommentDeletableWhenArchived struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	Value    int64
}

type SetCommentShowPos struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	Value    int64
}

type SpaceAbnormal struct {
	Event      notify.Event
	Space      *space.Space
	ExpiredNum int64
}

type AddSpaceManager struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
}

type RemoveSpaceManager struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	TargetId int64
}

type DeleteTag struct {
	Event    notify.Event
	Space    *space.Space
	Operator int64
	Tag      *space_tag.SpaceTag
}
