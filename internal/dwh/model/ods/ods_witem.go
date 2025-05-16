package ods

import (
	"go-cs/internal/dwh/pkg/model"
)

type OdsWitem struct {
	model.OdsModel

	Id                  int64  `json:"id,omitempty"`
	Pid                 int64  `json:"pid,omitempty"`            //上级任务id, >0 表示当前为子任务
	SpaceId             int64  `json:"space_id,omitempty"`       //空间id
	UserId              int64  `json:"user_id,omitempty"`        //创建用户id
	WorkItemType        int32  `json:"work_item_type,omitempty"` //原任务类型,已废弃
	WorkObjectId        int64  `json:"work_object_id,omitempty"` //模块id
	WorkItemGuid        string `json:"work_item_guid,omitempty"` //任务Guid
	WorkItemName        string `json:"work_item_name,omitempty"` //任务名称
	CreatedAt           int64  `json:"created_at,omitempty"`     //创建时间
	UpdatedAt           int64  `json:"updated_at,omitempty"`     //更新时间
	DeletedAt           int64  `json:"deleted_at,omitempty"`     //删除时间
	Doc                 string `json:"doc,omitempty"`            //
	IsRestart           int32  `json:"is_restart,omitempty"`     //是否为重启任务
	RestartAt           int64  `json:"restart_at,omitempty"`     //重启时间
	IconFlags           uint32 `json:"icon_flags,omitempty"`     //图标标记
	RestartUserId       int64  `json:"restart_user_id,omitempty"`
	CommentNum          int32  `json:"comment_num,omitempty"`           //评论数
	ResumeAt            int64  `json:"resume_at,omitempty"`             //恢复时间
	VersionId           int64  `json:"version_id,omitempty"`            //版本id -> 版本表
	WorkItemTypeId      int64  `json:"work_item_type_id,omitempty"`     //工作项类型Id
	WorkItemTypeKey     string `json:"work_item_type_key,omitempty"`    //工作项类型key
	FlowTemplateId      int64  `json:"flow_template_id,omitempty"`      //流程模板id
	FlowTemplateVersion int64  `json:"flow_template_version,omitempty"` //流程模板版本
	FlowId              int64  `json:"flow_id,omitempty"`               //流程Id
	FlowKey             string `json:"flow_key,omitempty"`              //流程key
	FlowMode            string `json:"flow_mode,omitempty"`             //模式类型
	FlowModeVersion     string `json:"flow_mode_version,omitempty"`     //[已废弃]模式版本
	FlowModeCode        string `json:"flow_mode_code,omitempty"`        //[已废弃]模式编码
	WorkItemStatus      string `json:"work_item_status,omitempty"`      //任务状态
	WorkItemStatusKey   string `json:"work_item_status_key,omitempty"`  //工作项状态key
	WorkItemStatusId    int64  `json:"work_item_status_id,omitempty"`   //工作项状态id
	LastStatusAt        int64  `json:"last_status_at,omitempty"`        //历史状态更新时间
	LastStatus          string `json:"last_status,omitempty"`           //历史状态
	LastStatusKey       string `json:"last_status_key,omitempty"`       //最后一次状态key
	LastStatusId        int64  `json:"last_status_id,omitempty"`        //最后一次状态id
	ChildNum            int32  `json:"child_num,omitempty"`             //子任务数量
	WorkItemFlowId      int64  `json:"work_item_flow_id,omitempty"`     //流程Id
	WorkItemFlowKey     string `json:"work_item_flow_key,omitempty"`    //流程key
}

type OdsWitemDoc struct {
	PlanStartAt    int64    `json:"plan_start_at,omitempty"`    //计划任务开始时间
	PlanCompleteAt int64    `json:"plan_complete_at,omitempty"` //计划任务完成时间
	Priority       string   `json:"priority,omitempty"`         //优先级
	Tags           []string `json:"tags,omitempty"`
	Directors      []string `json:"directors,omitempty"`
	NodeDirectors  []string `json:"node_directors,omitempty"`
	Followers      []string `json:"followers,omitempty"`
	Participators  []string `json:"participators,omitempty"`
}
