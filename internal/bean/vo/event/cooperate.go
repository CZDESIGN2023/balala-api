package event

import (
	"go-cs/api/notify"
	"go-cs/internal/domain/space"
	"go-cs/internal/domain/space_tag"
	"go-cs/internal/domain/work_item"
)

type CooperateComment struct {
	Event        notify.Event
	Operator     int64
	CommentId    int64
	Space        *space.Space
	WorkItem     *work_item.WorkItem
	Content      string
	ReferUserIds []int64
}

type CooperateTag struct {
	Event    notify.Event
	Operator int64
	Space    *space.Space
	Tag      *space_tag.SpaceTag
}

type CooperateDisableWorkFlow struct {
	Event        notify.Event
	Operator     int64
	SpaceId      int64
	WorkFlowId   int64
	WorkFlowName string
}

type CooperateDeleteWorkFlow struct {
	Event        notify.Event
	Operator     int64
	SpaceId      int64
	WorkFlowId   int64
	WorkFlowName string
}

type CooperateUpgradeWorkFlow struct {
	Event        notify.Event
	Operator     int64
	SpaceId      int64
	WorkFlowId   int64
	WorkItemId   int64
	WorkFlowName string
}
