package workitem_type

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

type SpaceId int64

type WorkItemType struct {
	shared.AggregateRoot

	Id        int64               `json:"id"`
	UserId    int64               `json:"user_id"`
	Uuid      string              `json:"uuid"`
	SpaceId   SpaceId             `json:"space_id"`
	Name      string              `json:"name"`
	Key       string              `json:"key"`
	FlowMode  consts.WorkFlowMode `json:"flow_mode"`
	Ranking   int64               `json:"ranking"`
	Status    int64               `json:"status"`
	CreatedAt int64               `json:"created_at"`
	UpdatedAt int64               `json:"updated_at"`
	DeletedAt int64               `json:"deleted_at"`
	IsSys     int32               `json:"is_sys"`
}

func (w *WorkItemType) IsSysType() bool {
	return w.IsSys == 1
}

func (w *WorkItemType) IsWorkFlowTaskType() bool {
	return w.Key == string(consts.WorkItemTypeKey_Task)
}

func (w *WorkItemType) IsStateFlowTaskType() bool {
	return w.Key == string(consts.WorkItemTypeKey_StateTask)
}

func (w *WorkItemType) IsSubTaskType() bool {
	return w.Key == string(consts.WorkItemTypeKey_SubTask)
}

func (w *WorkItemType) IsSameSpace(spaceId int64) bool {
	return int64(w.SpaceId) == spaceId
}
