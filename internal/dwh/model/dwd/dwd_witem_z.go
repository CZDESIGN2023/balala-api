package dwd

import (
	"go-cs/internal/dwh/pkg/model"
)

type DwdWitem struct {
	model.DimModel
	model.ChainModel

	SpaceId         int64  `json:"space_id" gorm:"column:space_id"`
	WorkItemId      int64  `json:"work_item_id" gorm:"column:work_item_id"`
	UserId          int64  `json:"user_id" gorm:"column:user_id"`
	StatusId        int64  `json:"status_id" gorm:"column:status_id"`
	ObjectId        int64  `json:"object_id" gorm:"column:object_id"`
	VersionId       int64  `json:"version_id" gorm:"column:version_id"`
	WorkItemTypeKey string `json:"work_item_type_key" gorm:"column:work_item_type_key"`
	LastStatusAt    int64  `json:"last_status_at" gorm:"column:last_status_at"`

	PlanStartAt    int64  `json:"plan_start_at" gorm:"column:plan_start_at"`
	PlanCompleteAt int64  `json:"plan_complete_at" gorm:"column:plan_complete_at"`
	Priority       string `json:"priority" gorm:"column:priority"`
	Directors      string `json:"directors" gorm:"column:directors"`
	NodeDirectors  string `json:"node_directors" gorm:"column:node_directors"`
	Participators  string `json:"participators" gorm:"column:participators"`
}

func (w *DwdWitem) DeepEqual(y *DwdWitem) bool {
	return compareStructFields(w, y, modelFieldNames)
}
