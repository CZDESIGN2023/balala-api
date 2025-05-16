package message

import (
	"fmt"
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
	"math"
	"strings"
)

const (
	Message_Type_WorkItem_Create            shared.MessageType = "Domain_Message.WorkItem.Create"
	Message_Type_WorkItem_Delete            shared.MessageType = "Domain_Message.WorkItem.Delete"
	Message_Type_WorkItem_SubTask_Create    shared.MessageType = "Domain_Message.WorkItem.SubTask_Create"
	Message_Type_WorkItem_Modify            shared.MessageType = "Domain_Message.WorkItem.Modify"
	Message_Type_WorkItem_Status_Change     shared.MessageType = "Domain_Message.WorkItem.Status_Change"
	Message_Type_WorkItem_Version_Change    shared.MessageType = "Domain_Message.WorkItem.Version_Change"
	Message_Type_WorkItem_Director_Change   shared.MessageType = "Domain_Message.WorkItem.Director_Change"
	Message_Type_WorkItem_Tag_Change        shared.MessageType = "Domain_Message.WorkItem.Tag_Change"
	Message_Type_WorkItem_FlowNode_Modify   shared.MessageType = "Domain_Message.WorkItem.FlowNode.Modify"
	Message_Type_WorkItem_File_Change       shared.MessageType = "Domain_Message.WorkItem.File_Change"
	Message_Type_WorkItem_FlowNode_Confirm  shared.MessageType = "Domain_Message.WorkItem.FlowNode.Confirm"
	Message_Type_WorkItem_FlowNode_Rollback shared.MessageType = "Domain_Message.WorkItem.FlowNode.Rollback"
	Message_Type_WorkItem_FlowNode_Reach    shared.MessageType = "Domain_Message.WorkItem.FlowNode.Reach"
)

// --- 创建任务
type CreateWorkItem struct {
	shared.DomainMessageBase

	SpaceId      int64
	WorkItemId   int64
	WorkItemName string
}

func (ops *CreateWorkItem) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Create
}

// --- 创建子任务
type CreateWorkItemSubTask struct {
	shared.DomainMessageBase

	SpaceId int64

	ParentWorkItemId int64

	WorkItemId   int64
	WorkItemName string
}

func (ops *CreateWorkItemSubTask) MessageType() shared.MessageType {
	return Message_Type_WorkItem_SubTask_Create
}

// -- 修改任务信息
type ModifyWorkItem struct {
	shared.DomainMessageBase

	SpaceId      int64
	WorkItemId   int64
	WorkItemName string

	Updates []FieldUpdate
}

func (ops *ModifyWorkItem) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Modify
}

// -- 回滚流程节点

type RollbackWorkItemFlowNode struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemId   int64
	WorkItemName string

	FlowNodeId   int64
	FlowNodeCode string

	Reason string
}

func (ops *RollbackWorkItemFlowNode) MessageType() shared.MessageType {
	return Message_Type_WorkItem_FlowNode_Rollback
}

// -- 完成了流程节点

type ConfirmWorkItemFlowNode struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemId   int64
	WorkItemName string

	FlowNodeId   int64
	FlowNodeCode string

	Reason string
}

func (ops *ConfirmWorkItemFlowNode) MessageType() shared.MessageType {
	return Message_Type_WorkItem_FlowNode_Confirm
}

// -- 流程到达

type ReachWorkItemFlowNode struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemId   int64
	WorkItemName string

	FlowNodeId   int64
	FlowNodeName string
	FlowNodeCode string

	Reason string
}

func (ops *ReachWorkItemFlowNode) MessageType() shared.MessageType {
	return Message_Type_WorkItem_FlowNode_Reach
}

// -- 修改任务状态

type ChangeWorkItemStatus struct {
	shared.DomainMessageBase

	AffectByParent bool

	SpaceId int64
	Pid     int64

	WorkItemId      int64
	WorkItemName    string
	WorkItemTypeKey consts.WorkItemTypeKey

	OldWorkItemStatusKey string
	OldWorkItemStatusId  int64
	OldWorkItemStatusVal string

	NewWorkItemStatusKey string
	NewWorkItemStatusId  int64
	NewWorkItemStatusVal string

	FlowNodeId   int64
	FlowNodeCode string

	Reason string
	Remark string
}

func (ops *ChangeWorkItemStatus) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Status_Change
}

// --- 更改版本

type ChangeWorkItemDirector_Node struct {
	FlowNodeCode string
	OldDirectors []string
	NewDirectors []string
}

type ChangeWorkItemDirector struct {
	shared.DomainMessageBase

	SpaceId int64

	WorkItemPid  int64
	WorkItemId   int64
	WorkItemName string

	FlowTemplateId int64

	OldDirectors []string //子任务使用
	NewDirectors []string //子任务使用

	WorkItemRoleKey string
	WorkItemRoleId  int64

	Nodes []*ChangeWorkItemDirector_Node
}

func (ops *ChangeWorkItemDirector) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Director_Change
}

// --- 删除任务
type DeleteWorkItem struct {
	shared.DomainMessageBase

	SpaceId          int64
	ParentWorkItemId int64

	WorkItemId     int64
	WorkItemName   string
	PlanStartAt    int64
	PlanCompleteAt int64
	ProcessRate    int32
	Directors      []int64
}

func (ops *DeleteWorkItem) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Delete
}

// --- Tag变更
type ChangeWorkItemTag struct {
	shared.DomainMessageBase

	SpaceId      int64
	WorkItemId   int64
	WorkItemName string

	OldTags []string
	NewTags []string

	AddTags    []string
	RemoveTags []string
}

func (ops *ChangeWorkItemTag) MessageType() shared.MessageType {
	return Message_Type_WorkItem_Tag_Change
}

// -- 修改任务信息
type ModifyWorkItemFlowNode struct {
	shared.DomainMessageBase

	SpaceId    int64
	WorkItemId int64

	FlowNodeId   int64
	FlowNodeCode string

	Updates []FieldUpdate
}

func (ops *ModifyWorkItemFlowNode) MessageType() shared.MessageType {
	return Message_Type_WorkItem_FlowNode_Modify
}

//-- 设置附件

type FileInfo struct {
	Name string
	Size int64
}

func (f FileInfo) String() string {
	size := float64(f.Size)
	unit := "B"
	units := []string{"B", "KB", "MB", "GB", "TB"}
	for i := 1; size >= 1024 && i < len(units); i++ {
		size /= 1024
		unit = units[i]
	}

	size = math.Floor(size*10) / 10

	s := fmt.Sprintf(" %.1f", size)
	if strings.HasSuffix(s, ".0") {
		s = s[:len(s)-2]
	}

	return f.Name + s + unit
}

type ChangeWorkItemFile struct {
	shared.DomainMessageBase

	SpaceId      int64
	WorkItemId   int64
	WorkItemName string

	AddFiles    []FileInfo
	RemoveFiles []FileInfo
}

func (ops *ChangeWorkItemFile) MessageType() shared.MessageType {
	return Message_Type_WorkItem_File_Change
}
