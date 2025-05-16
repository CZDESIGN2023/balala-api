package event

import (
	"go-cs/api/notify"
	"go-cs/internal/domain/space"
	"go-cs/internal/domain/work_item"
)

type FieldUpdate struct {
	Field    string
	OldValue any
	NewValue any
}

type ChangeWorkItemField struct {
	Event    notify.Event
	Space    *space.Space
	WorkItem *work_item.WorkItem
	Operator int64
	Updates  []FieldUpdate
}
type ChangeWorkItemDirector struct {
	Event             notify.Event
	Space             *space.Space
	WorkItem          *work_item.WorkItem
	Operator          int64
	OldValues         []int64 //当前节点负责人
	NewValues         []int64
	Nodes             []NodeDirectorOp
	ViaCreateWorkItem bool
}

type NodeDirectorOp struct {
	NodeName  string
	OldValues []int64
	NewValues []int64
}

type ChangeWorkItemTag struct {
	Event     notify.Event
	Space     *space.Space
	WorkItem  *work_item.WorkItem
	Operator  int64
	OldValues []int64
	NewValues []int64
}

type SetWorkItemFiles_FileInfo struct {
	FileName string
}

type SetWorkItemFiles struct {
	Event    notify.Event
	Space    *space.Space
	WorkItem *work_item.WorkItem
	Operator int64
	Adds     []*SetWorkItemFiles_FileInfo
	Deletes  []*SetWorkItemFiles_FileInfo
}

type CreateChildWorkItem struct {
	Event         notify.Event
	Operator      int64
	Space         *space.Space
	WorkItem      *work_item.WorkItem //父级任务
	ChildWorkItem *work_item.WorkItem //子级任务
}

type WorkItemExpired struct {
	Event      notify.Event
	WorkItemId int64
	Operator   int64
}

type CreateWorkItem struct {
	Event    notify.Event
	Operator int64
	Space    *space.Space
	WorkItem *work_item.WorkItem
	Nodes    []NodeDirectorOp
}

type DeleteWorkItem struct {
	Event        notify.Event
	Operator     int64
	Space        *space.Space
	WorkItem     *work_item.WorkItem
	SubWorkItems []*work_item.WorkItem
}

type DeleteChildWorkItem struct {
	Event         notify.Event
	Operator      int64
	Space         *space.Space
	WorkItem      *work_item.WorkItem //父级任务
	ChildWorkItem *work_item.WorkItem //子级任务
}

type WorkItemFlowNodeExpired struct {
	Event         notify.Event
	WorkItemId    int64
	Operator      int64
	NodeName      string
	NodeDirectors []int64
}
