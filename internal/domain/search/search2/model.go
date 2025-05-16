package search2

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/datatypes"
)

// Model
// query 标签：查询参数字段明
// db 标签：数据库访问路径
// gorm 标签：别名字段，当db与gorm column不同时 gorm column 将作为别名
// dt 标签：类型
type Model struct {
	Id              int64  `query:"work_item_id" db:"id" gorm:"column:id" json:"id,omitempty" es:"work_item_id"`
	Pid             int64  `query:"pid" db:"pid" gorm:"column:pid" json:"pid,omitempty" es:"pid"`
	SpaceId         int64  `query:"space_id" db:"space_id" gorm:"column:space_id" json:"space_id,omitempty" es:"space_id"`
	UserId          int64  `query:"user_id" db:"user_id" gorm:"column:user_id" json:"user_id,omitempty" es:"user_id"`
	WorkObjectId    int64  `query:"work_object_id" db:"work_object_id" gorm:"column:work_object_id" json:"work_object_id,omitempty" es:"work_object_id"`
	WorkItemGuid    string `query:"work_item_guid" db:"work_item_guid" gorm:"column:work_item_guid" json:"work_item_guid,omitempty" es:"work_item_guid"`
	WorkItemName    string `query:"work_item_name" db:"work_item_name" gorm:"column:work_item_name" json:"work_item_name,omitempty" es:"work_item_name"`
	CreatedAt       int64  `query:"created_at" db:"created_at" gorm:"column:created_at" dt:"date" json:"created_at,omitempty" es:"created_at"`
	UpdatedAt       int64  `query:"updated_at" db:"updated_at" gorm:"column:updated_at" dt:"date" json:"updated_at,omitempty" es:"updated_at"`
	DeletedAt       int64  `query:"deleted_at" db:"deleted_at" gorm:"column:deleted_at" dt:"date" json:"deleted_at,omitempty" es:"deleted_at"`
	VersionId       int64  `query:"version_id" db:"version_id"  gorm:"column:version_id" json:"version_id" es:"version_id"`
	FlowMode        string `db:"flow_mode" gorm:"column:flow_mode" json:"flow_mode" es:"flow_mode"`
	FlowModeVersion string `db:"flow_mode_version" gorm:"column:flow_mode_version"  json:"flow_mode_version" es:"flow_mode_version"`
	FlowModeCode    string `db:"flow_mode_code" gorm:"column:flow_mode_code" json:"flow_mode_code" es:"flow_mode_code"`
	LastStatusAt    int64  `query:"last_status_at" db:"last_status_at" gorm:"column:last_status_at" dt:"date" json:"last_status_at" es:"last_status_at"`
	LastStatus      int64  `db:"last_status" gorm:"column:last_status" json:"last_status" es:"last_status"`
	IsRestart       int32  `db:"is_restart"  gorm:"column:is_restart" json:"is_restart" es:"is_restart"`
	RestartAt       int64  `db:"restart_at"  gorm:"column:restart_at" json:"restart_at" es:"restart_at"`
	RestartUserId   int64  `db:"restart_user_id"  gorm:"column:restart_user_id" json:"restart_user_id" es:"restart_user_id"`
	IconFlags       uint32 `db:"icon_flags"  gorm:"column:icon_flags" json:"icon_flags" es:"icon_flags"`
	CommentNum      uint32 `db:"comment_num"  gorm:"column:comment_num" json:"comment_num" es:"comment_num"`
	ResumeAt        int64  `db:"resume_at"  gorm:"column:resume_at" json:"resume_at" es:"resume_at"`
	ChildNum        int64  `db:"child_num"  gorm:"column:child_num" json:"child_num" es:"child_num"`

	WorkItemTypeId   int64  `db:"work_item_type_id"  gorm:"column:work_item_type_id" json:"work_item_type_id" es:"work_item_type_id"`
	FlowId           int64  `query:"flow_id" db:"flow_id"  gorm:"column:flow_id" json:"flow_id" es:"flow_id"`
	FlowTemplateId   int64  `db:"flow_template_id"  gorm:"column:flow_template_id" json:"flow_template_id" es:"flow_template_id"`
	WorkItemStatusId int64  `db:"work_item_status_id"  gorm:"column:work_item_status_id" json:"work_item_status_id" es:"work_item_status_id"`
	WorkItemStatus   string `query:"work_item_status" db:"work_item_status" gorm:"column:work_item_status" json:"work_item_status,omitempty" es:"work_item_status"`
	WorkItemFlowId   int64  `query:"work_item_flow_id" db:"work_item_flow_id" gorm:"column:work_item_flow_id" json:"work_item_flow_id,omitempty" es:"work_item_flow_id"`

	// 以下为doc内字段
	//Doc            datatypes.JSON              `query:"doc" db:"doc" gorm:"column:doc"` //这个字段不能直接读取
	Describe       string                      `db:"doc->>'$.describe'" gorm:"column:describe" dt:"rich-text" json:"describe,omitempty" es:"doc.describe"` //大字段不能用于查询
	Remark         string                      `db:"doc->>'$.remark'" gorm:"column:remark" dt:"rich-text" json:"remark,omitempty" es:"doc.remark"`         //大字段不能用于查询
	ProcessRate    int32                       `query:"process_rate" db:"doc->'$.process_rate'" gorm:"column:process_rate" json:"process_rate,omitempty" es:"doc.process_rate"`
	Priority       string                      `query:"priority" db:"doc->>'$.priority'" gorm:"column:priority" json:"priority,omitempty" es:"doc.priority"`
	PlanStartAt    int64                       `query:"plan_start_at" db:"doc->'$.plan_start_at'" gorm:"column:plan_start_at"  dt:"date" json:"plan_start_at,omitempty" es:"doc.plan_start_at"`
	PlanCompleteAt int64                       `query:"plan_complete_at" db:"doc->'$.plan_complete_at'" gorm:"column:plan_complete_at"  dt:"date" json:"plan_complete_at,omitempty" es:"doc.plan_complete_at"`
	Tags           datatypes.JSONSlice[string] `query:"tags" db:"doc->'$.tags'" gorm:"column:tags" dt:"multi-select" json:"tags,omitempty" es:"doc.tags"`
	Directors      datatypes.JSONSlice[string] `query:"directors" db:"doc->'$.directors'" gorm:"column:directors" dt:"multi-user" json:"directors,omitempty" es:"doc.directors"`                     //当前负责人
	Followers      datatypes.JSONSlice[string] `query:"followers" db:"doc->'$.followers'" gorm:"column:followers" dt:"multi-user" json:"followers,omitempty" es:"doc.followers"`                     //关注人
	Participators  datatypes.JSONSlice[string] `query:"participators" db:"doc->'$.participators'" gorm:"column:participators" dt:"multi-user" json:"participators,omitempty" es:"doc.participators"` //参与人
	NodeDirectors  datatypes.JSONSlice[string] `query:"node_directors" db:"doc->'$.node_directors'" gorm:"column:node_directors" dt:"multi-user" json:"node_directors,omitempty" es:"doc.node_directors"`

	//FlowNode
	NodeStatus int64 `query:"node_status" db:"flow_node_status" json:"-" gorm:"column:-"`

	NodePlanStartAt    int64 `query:"node_plan_start_at" db:"plan_start_at" dt:"date" json:"-" gorm:"column:-"`
	NodePlanCompleteAt int64 `query:"node_plan_complete_at" db:"plan_complete_at" dt:"date" json:"-" gorm:"column:-"`
}

func (m *Model) TableName() string {
	return "space_work_item_v2"
}

func init() {
	typeOf := reflect.TypeOf(Model{})
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		var gormTag = field.Tag.Get("gorm")
		if gormTag == "" {
			panic("tag gorm 不能为空")
		}

		var dbTag = field.Tag.Get("db")
		if dbTag == "" {
			panic("tag db 不能为空")
		}

		if strings.Contains(dbTag, "''") {
			panic(fmt.Sprintf("dbTag %v 存在连续两个单引号", dbTag))
		}

		fieldModelList = append(fieldModelList, FieldModel{
			query: QueryField(field.Tag.Get("query")),
			db:    dbTag,
			gorm:  strings.Split(gormTag, ":")[1],
			dt:    DataType(field.Tag.Get("dt")),
			inObj: strings.Contains(dbTag, "->>") || strings.Contains(dbTag, "->"),
		})
	}

	for _, model := range fieldModelList {
		if v := model.Query(); v != "" {
			query2FieldModelMap[model.Query()] = model
		}
		if v := model.RawGorm(); v != "" {
			column2fieldModelMap[model.RawGorm()] = model
		}
	}
}
