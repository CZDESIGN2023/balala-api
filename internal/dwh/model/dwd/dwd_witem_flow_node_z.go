package dwd

import (
	"encoding/json"
	"go-cs/internal/dwh/pkg/model"
	"reflect"
)

type DwdWitemFlowNode struct {
	model.DimModel
	model.ChainModel
	SpaceId        int64  `json:"space_id" gorm:"column:space_id"`
	WorkItemId     int64  `json:"work_item_id" gorm:"column:work_item_id"`
	NodeId         int64  `json:"node_id" gorm:"column:node_id"`
	NodeCode       string `json:"node_code" gorm:"column:node_code"`
	NodeStatus     int32  `json:"node_status" gorm:"column:node_status"`
	PlanStartAt    int64  `json:"plan_start_at" gorm:"column:plan_start_at"`
	PlanCompleteAt int64  `json:"plan_complete_at" gorm:"column:plan_complete_at"`
	Directors      string `json:"directors" gorm:"column:directors"`
}

func (w *DwdWitemFlowNode) DeepEqual(y *DwdWitemFlowNode) bool {
	return reflect.DeepEqual(w.KeyValue(), y.KeyValue())
}

func (w *DwdWitemFlowNode) KeyValue() map[string]interface{} {
	directors := make([]string, 0)
	json.Unmarshal([]byte(w.Directors), &directors)

	return map[string]interface{}{
		"space_id":         w.SpaceId,
		"work_item_id":     w.WorkItemId,
		"node_id":          w.NodeId,
		"node_code":        w.NodeCode,
		"node_status":      w.NodeStatus,
		"plan_start_at":    w.PlanStartAt,
		"plan_complete_at": w.PlanCompleteAt,
		"directors":        directors,
	}
}
