package event

import (
	"go-cs/api/notify"
	"go-cs/internal/domain/space"
	"go-cs/internal/domain/work_item"
)

// 节点流转
type ChangeWorkItemFlowNode struct {
	Event notify.Event
	Data  *ChangeWorkItemFlowNodeData
}

type ChangeWorkItemFlowNodeData struct {
	Operator *ChangeWorkItemFlowNodeData_User
	Space    *ChangeWorkItemFlowNodeData_Space
	WorkItem *ChangeWorkItemFlowNodeData_WorkItem
	ToNode   *ChangeWorkItemFlowNodeData_Node
}

type ChangeWorkItemFlowNodeData_Node struct {
	Name string
	Code string
	Id   int64
}

type ChangeWorkItemFlowNodeData_WorkItem struct {
	Id   int64
	Name string
}

type ChangeWorkItemFlowNodeData_Space struct {
	Id   int64
	Name string
}

type ChangeWorkItemFlowNodeData_User struct {
	Id int64
}

// 任务回滚
type RollbackWorkItemFlowNode struct {
	Event notify.Event
	Data  *RollbackWorkItemFlowNodeData
}

type RollbackWorkItemFlowNodeData struct {
	Reason   string
	Operator *RollbackWorkItemFlowNodeData_User
	Space    *RollbackWorkItemFlowNodeData_Space
	WorkItem *RollbackWorkItemFlowNodeData_WorkItem
	ToNode   *RollbackWorkItemFlowNodeData_Node
}

type RollbackWorkItemFlowNodeData_Node struct {
	Name string
	Code string
	Id   int64
}

type RollbackWorkItemFlowNodeData_WorkItem struct {
	Id   int64
	Name string
}

type RollbackWorkItemFlowNodeData_Space struct {
	Id   int64
	Name string
}

type RollbackWorkItemFlowNodeData_User struct {
	Id int64
}

type ChangeWorkFlowNodePlanTime struct {
	Event    notify.Event
	Space    *space.Space
	WorkItem *work_item.WorkItem
	Operator int64

	OldValues []any
	NewValues []any
}
