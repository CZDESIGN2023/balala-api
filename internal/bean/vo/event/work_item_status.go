package event

import "go-cs/api/notify"

// 恢复任务

type ResumeWorkItemScene string

const (
	ResumeWorkItemScene_formClosed    ResumeWorkItemScene = "formClosed"
	ResumeWorkItemScene_formTerminate ResumeWorkItemScene = "formTerminate"
)

type ResumeWorkItem struct {
	Event notify.Event
	Data  *ResumeWorkItemData
}

type ResumeWorkItemData struct {
	Scene    ResumeWorkItemScene
	Reason   string
	Operator *ResumeWorkItemData_User
	Space    *ResumeWorkItemData_Space
	WorkItem *ResumeWorkItemData_WorkItem
}

type ResumeWorkItemData_WorkItem struct {
	Id   int64
	Name string
}

type ResumeWorkItemData_Space struct {
	Id   int64
	Name string
}

type ResumeWorkItemData_User struct {
	Id int64
}

// 任务重启
type RestartWorkItem struct {
	Event notify.Event
	Data  *RestartWorkItemData
}

type RestartWorkItemData struct {
	Reason   string
	Operator *RestartWorkItemData_User
	Space    *RestartWorkItemData_Space
	WorkItem *RestartWorkItemData_WorkItem
	ToNode   *RestartWorkItemData_Node
}

type RestartWorkItemData_Node struct {
	Name string
	Code string
}

type RestartWorkItemData_WorkItem struct {
	Id   int64
	Name string
}

type RestartWorkItemData_Space struct {
	Id   int64
	Name string
}

type RestartWorkItemData_User struct {
	Id int64
}

// 任务关闭
type CloseWorkItem struct {
	Event notify.Event
	Data  *CloseWorkItemData
}

type CloseWorkItemData struct {
	Reason   string
	Operator *CloseWorkItemData_User
	Space    *CloseWorkItemData_Space
	WorkItem *CloseWorkItemData_WorkItem
}

type CloseWorkItemData_WorkItem struct {
	Id   int64
	Name string
}

type CloseWorkItemData_Space struct {
	Id   int64
	Name string
}

type CloseWorkItemData_User struct {
	Id int64
}

// 终止
type TerminateWorkItem struct {
	Event notify.Event
	Data  *TerminateWorkItemData
}

type TerminateWorkItemData struct {
	Reason   string
	Operator *TerminateWorkItemData_User
	Space    *TerminateWorkItemData_Space
	WorkItem *TerminateWorkItemData_WorkItem
}

type TerminateWorkItemData_WorkItem struct {
	Id   int64
	Name string
}

type TerminateWorkItemData_Space struct {
	Id   int64
	Name string
}

type TerminateWorkItemData_User struct {
	Id int64
}

// 完成
type CompleteWorkItem struct {
	Event notify.Event
	Data  *CompleteWorkItemData
}

type CompleteWorkItemData struct {
	Reason   string
	Operator *CompleteWorkItemData_User
	Space    *CompleteWorkItemData_Space
	WorkItem *CompleteWorkItemData_WorkItem
}

type CompleteWorkItemData_WorkItem struct {
	Id   int64
	Name string
}

type CompleteWorkItemData_Space struct {
	Id   int64
	Name string
}

type CompleteWorkItemData_User struct {
	Id int64
}
