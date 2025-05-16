package ods

import "go-cs/internal/dwh/pkg/model"

type OdsWitemStatus struct {
	model.OdsModel

	Id             int64  `json:"id,omitempty"`
	Uuid           string `json:"uuid,omitempty"`
	UserId         int64  `json:"user_id,omitempty"`
	SpaceId        int64  `json:"space_id,omitempty"`          //空间id
	WorkItemTypeId int64  `json:"work_item_type_id,omitempty"` //工作项类型id
	Name           string `json:"name,omitempty"`              //名称
	Key            string `json:"key,omitempty"`               //工作项状态key
	Val            string `json:"val,omitempty"`               //工作项状态值
	StatusType     int32  `json:"status_type,omitempty"`       //工作项状态类型 1:起始 2:过程 3:归档
	FlowScope      string `json:"flow_scope,omitempty"`        //流程范围
	Ranking        int64  `json:"ranking,omitempty"`           //排序
	Status         int64  `json:"status,omitempty"`            //状态;0:禁用,1:正常
	CreatedAt      int64  `json:"created_at,omitempty"`        //创建时间
	UpdatedAt      int64  `json:"updated_at,omitempty"`        //更新时间
	DeletedAt      int64  `json:"deleted_at,omitempty"`        //删除时间
	IsSys          int32  `json:"is_sys,omitempty"`            //是否系统预设
}
