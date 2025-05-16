package vo

import (
	"encoding/json"
	"go-cs/pkg/stream"
)

type WorkItemStatusVo struct {
	Id         int64
	Name       string
	Key        string
	Val        string
	StatusType int32
}

type CreateSpaceWorkItemVoV2 struct {
	UserId         int64
	SpaceId        int64
	WorkObjectId   int64
	WorkItemName   string
	PlanStartAt    int64
	PlanCompleteAt int64
	ProcessRate    int32
	CreatedAt      int64
	Describe       string
	Remark         string
	Priority       string
	WorkFlowId     int64
	WorkVersionId  int64
	WorkItemTypeId int64
	Followers      []int64

	IconFlags []uint32
	TagAdd    []int64
	FileAdd   CreateSpaceWorkItemFileVosV2
	Owner     CreateSpaceWorkItemOwnersV2
}

type CreateSpaceWorkItemOwnersV2 []*CreateSpaceWorkItemOwnerV2

type CreateSpaceWorkItemOwnerV2 struct {
	OwnerRole string
	Directors []int64
}

func (c *CreateSpaceWorkItemOwnersV2) AllDirectors() []int64 {
	directors := make([]int64, 0)
	for _, v := range *c {
		directors = append(directors, v.Directors...)
	}
	directors = stream.Unique(directors)
	return directors
}

type CreateSpaceWorkItemFileVosV2 []*CreateSpaceWorkItemFileVoV2

type CreateSpaceWorkItemFileVoV2 struct {
	Id   int64
	Name string
	Uri  string
	Size int64
}

type CreateSpaceWorkItemTaskVoV2 struct {
	WorkItemName   string
	PlanStartAt    int64
	PlanCompleteAt int64
	ProcessRate    int32
	DirectorAdd    []int64
}

type SetSpaceWorkItemDirectorVoV2 struct {
	SpaceId    int64
	WorkItemId int64

	DirectorAdd    []int64
	DirectorRemove []int64
}

type SetSpaceWorkItemTagVoV2 struct {
	SpaceId    int64
	WorkItemId int64

	TagNew    []string
	TagAdd    []int64
	TagRemove []int64
}

type CreateSpaceWorkTaskVoV2 struct {
	UserId int64

	SpaceId        int64
	WorkItemId     int64
	WorkTaskGuid   string
	WorkTaskName   string
	WorkTaskStatus int32
	PlanStartAt    int64
	PlanCompleteAt int64
	ProcessRate    int32
	Remark         string
	Describe       string
	CreatedAt      int64

	DirectorAdd []int64
}

type SetSpaceWorkTaskDirectorVoV2 struct {
	SpaceId    int64
	WorkTaskId int64

	DirectorAdd    []int64
	DirectorRemove []int64
}

type SetSpaceWorkItemFileInfoVoV2 struct {
	WorkItemId int64

	FileInfoAdd    []int64
	FileInfoRemove []int64
}

type WorkItemDoc struct {
	PlanStartAt    int64    `json:"plan_start_at,omitempty"`
	PlanCompleteAt int64    `json:"plan_complete_at,omitempty"`
	ProcessRate    int32    `json:"process_rate,omitempty"`
	Remark         string   `json:"remark"`
	Describe       string   `json:"describe"`
	Priority       string   `json:"priority,omitempty"`
	Directors      []string `json:"directors"`     //数组不能 omitempty
	Tags           []string `json:"tags"`          //数组不能 omitempty
	Followers      []string `json:"followers"`     //数组不能 omitempty
	Participators  []string `json:"participators"` //数组不能 omitempty
}

func (d WorkItemDoc) MarshalJSON() ([]byte, error) {
	type _workItemDoc WorkItemDoc
	cpy := _workItemDoc(d)

	if cpy.Directors == nil {
		cpy.Directors = []string{}
	}
	if cpy.Tags == nil {
		cpy.Tags = []string{}
	}

	if cpy.Followers == nil {
		cpy.Followers = []string{}
	}

	if cpy.Participators == nil {
		cpy.Participators = []string{}
	}

	return json.Marshal(cpy)
}

type UpdateWorkItemTypeInfoVo struct {
	WorkItemId      int64
	FlowMode        string
	FlowModeCode    string
	FlowModeVersion string
	WorkItemType    int64
}

type SearchWorkItemTypeVo struct {
	Ids     []int64
	SpaceId []int64
}

type SearchWorkItemStatusVo struct {
	Ids []int64
}

type SpaceWorkItemEsVo struct {
	Id              int64
	Pid             int64
	SpaceId         int64
	UserId          int64
	WorkItemType    int32
	WorkObjectId    int64
	WorkItemGuid    string
	WorkItemName    string
	WorkItemStatus  string
	CreatedAt       int64
	UpdatedAt       int64
	DeletedAt       int64
	Doc             any
	FlowMode        string
	FlowModeVersion string
	FlowModeCode    string
	LastStatusAt    int64
	LastStatus      int64
	IsRestart       int32
	RestartAt       int64
	IconFlags       uint32
	RestartUserId   int64
	CommentNum      int32
	ResumeAt        int64
	VersionId       int64
}
