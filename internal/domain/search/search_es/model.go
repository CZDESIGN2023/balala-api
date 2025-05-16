package search_es

import (
	"reflect"

	"gorm.io/datatypes"
)

var query2FieldModelMap = map[QueryField]FieldModel{} // query: FieldModel

// Model
// query 标签：【查询】参数字段明
// dt 标签：【查询】参数字段类型
// es 标签: es字段
// esDt 标签: es字段类型
type Model struct {
	Id              int64  `query:"work_item_id" json:"id,omitempty" es:"id" esDt:"int"`
	Pid             int64  `query:"pid" json:"pid,omitempty" es:"pid" esDt:"int"`
	SpaceId         int64  `query:"space_id" json:"space_id,omitempty" es:"space_id" esDt:"int"`
	UserId          int64  `query:"user_id" json:"user_id,omitempty" es:"user_id" esDt:"int"`
	WorkItemType    int64  `query:"work_item_type" json:"work_item_type,omitempty" es:"work_item_type" esDt:"int"`
	WorkItemTypeKey int64  `query:"work_item_type_key" json:"work_item_type_key,omitempty" es:"work_item_type_key" esDt:"string"`
	WorkObjectId    int64  `query:"work_object_id" json:"work_object_id,omitempty" es:"work_object_id" esDt:"int"`
	WorkItemGuid    string `query:"work_item_guid" json:"work_item_guid,omitempty" es:"work_item_guid"  esKeyword:"work_item_guid.keyword" esDt:"string"`
	WorkItemName    string `query:"work_item_name" json:"work_item_name,omitempty" es:"work_item_name" esKeyword:"work_item_name.keyword" esDt:"string"`
	CreatedAt       int64  `query:"created_at"  json:"created_at,omitempty" queryDt:"date" es:"created_at" esDt:"int"`
	UpdatedAt       int64  `query:"updated_at" json:"updated_at,omitempty" queryDt:"date" es:"updated_at" esDt:"int"`
	DeletedAt       int64  `query:"deleted_at" json:"deleted_at,omitempty" queryDt:"date" es:"deleted_at" esDt:"int"`
	VersionId       int64  `query:"version_id" json:"version_id" es:"version_id" esDt:"int"`
	FlowMode        string `query:"flow_mode" json:"flow_mode" es:"flow_mode" esKeyword:"flow_mode.keyword" esDt:"string"`
	FlowModeVersion string `query:"flow_mode_version" json:"flow_mode_version" es:"flow_mode_version" esKeyword:"flow_mode_version.keyword" esDt:"string"`
	FlowModeCode    string `query:"flow_mode_code" json:"flow_mode_code" es:"flow_mode_code" esKeyword:"flow_mode_code.keyword" esDt:"string"`
	LastStatusAt    int64  `query:"last_status_at" json:"last_status_at" queryDt:"date" es:"last_status_at" esDt:"int"`
	LastStatus      int64  `query:"last_status" json:"last_status" es:"last_status" esDt:"int"`
	IsRestart       int32  `query:"is_restart" json:"is_restart" es:"is_restart" esDt:"int"`
	RestartAt       int64  `query:"restart_at" json:"restart_at" es:"restart_at" esDt:"int"`
	RestartUserId   int64  `query:"restart_user_id" json:"restart_user_id" es:"restart_user_id" esDt:"int"`
	IconFlags       uint32 `query:"icon_flags" json:"icon_flags" es:"icon_flags" esDt:"int"`
	CommentNum      uint32 `query:"comment_num" json:"comment_num" es:"comment_num" esDt:"int"`
	ResumeAt        int64  `query:"resume_at" json:"resume_at" queryDt:"date" es:"resume_at" esDt:"int"`
	ChildNum        int64  `query:"child_num" db:"child_num"  gorm:"column:child_num" json:"child_num" es:"child_num"`

	WorkItemTypeId   int64  `db:"work_item_type_id"  gorm:"column:work_item_type_id" json:"work_item_type_id" es:"work_item_type_id"`
	FlowId           int64  `query:"flow_id" db:"flow_id"  gorm:"column:flow_id" json:"flow_id" es:"flow_id"`
	FlowTemplateId   int64  `db:"flow_template_id"  gorm:"column:flow_template_id" json:"flow_template_id" es:"flow_template_id"`
	WorkItemStatusId int64  `query:"work_item_status_id" db:"work_item_status_id"  gorm:"column:work_item_status_id" json:"work_item_status_id" es:"work_item_status_id" esKeyword:"work_item_status_id"`
	WorkItemStatus   string `query:"work_item_status" db:"work_item_status" gorm:"column:work_item_status" json:"work_item_status,omitempty" es:"work_item_status" esKeyword:"work_item_status.keyword"`
	WorkItemFlowId   int64  `query:"work_item_flow_id" db:"work_item_flow_id" gorm:"column:work_item_flow_id" json:"work_item_flow_id,omitempty" es:"work_item_flow_id" esKeyword:"work_item_flow_id"`

	// 以下为doc内字段
	//Doc            datatypes.JSON            `query:"doc" gorm:"column:doc"` //这个字段不能直接读取
	Describe       string                      `query:"describe" json:"describe,omitempty" es:"doc.describe" esKeyword:"doc.describe.keyword" esDt:"string"`
	Remark         string                      `query:"remark" json:"remark,omitempty" es:"doc.remark"  esKeyword:"doc.remark.keyword" esDt:"string"`
	ProcessRate    int32                       `query:"process_rate" json:"process_rate,omitempty" es:"doc.process_rate" esDt:"int"`
	Priority       string                      `query:"priority" json:"priority,omitempty" es:"doc.priority"  esKeyword:"doc.priority.keyword" esDt:"string"`
	PlanStartAt    int64                       `query:"plan_start_at" queryDt:"date" json:"plan_start_at,omitempty" es:"doc.plan_start_at" esDt:"int"`
	PlanCompleteAt int64                       `query:"plan_complete_at" queryDt:"date" json:"plan_complete_at,omitempty" es:"doc.plan_complete_at" esDt:"int"`
	Tags           datatypes.JSONSlice[string] `query:"tags" json:"tags,omitempty" es:"doc.tags"  esKeyword:"doc.tags.keyword" esDt:"array[string]"`
	Directors      datatypes.JSONSlice[string] `query:"directors" json:"directors,omitempty" es:"doc.directors"  esKeyword:"doc.directors.keyword" esDt:"array[string]"`
	Followers      datatypes.JSONSlice[string] `query:"followers" json:"followers,omitempty" es:"doc.followers"  esKeyword:"doc.followers.keyword" esDt:"array[string]"`
	Participators  datatypes.JSONSlice[string] `query:"participators" json:"participators,omitempty" es:"doc.participators"  esKeyword:"doc.participators.keyword" esDt:"array[string]"`
	NodeDirectors  datatypes.JSONSlice[string] `query:"node_directors" json:"node_directors,omitempty" es:"doc.node_directors"  esKeyword:"doc.node_directors.keyword" esDt:"array[string]"`

	//Reason
	TerminateReason string `query:"terminate_reason" json:"terminate_reason,omitempty" es:"reason.terminate" esKeyword:"reason.terminate.keyword" esDt:"string"`
	ResumeReason    string `query:"resume_reason" json:"resume_reason,omitempty" es:"reason.resume" esKeyword:"reason.resume.keyword" esDt:"string"`
	RollbackReason  string `query:"rollback_reason" json:"rollback_reason,omitempty" es:"reason.rollback" esKeyword:"reason.rollback.keyword" esDt:"string"`
	CloseReason     string `query:"close_reason" json:"close_reason,omitempty" es:"reason.close" esKeyword:"reason.close.keyword" esDt:"string"`

	//FlowNode
	NodeStatus         int64 `query:"node_status" db:"flow_node_status" json:"-" gorm:"column:-" es:"-"`
	NodePlanStartAt    int64 `query:"node_plan_start_at" db:"plan_start_at" dt:"date" json:"-" gorm:"column:-" es:"-"`
	NodePlanCompleteAt int64 `query:"node_plan_complete_at" db:"plan_complete_at" dt:"date" json:"-" gorm:"column:-" es:"-"`

	StateDirectors []string `query:"state_directors" json:"-" gorm:"column:-" es:"-"`
}

func init() {
	typeOf := reflect.TypeOf(Model{})
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		var esTag = field.Tag.Get("es")
		if esTag == "" {
			panic("tag es 不能为空")
		}

		var esKeyword = field.Tag.Get("esKeyword")

		fieldModelList = append(fieldModelList, FieldModel{
			query:     field.Tag.Get("query"),
			dt:        DataType(field.Tag.Get("queryDt")),
			es:        esTag,
			esKeyword: esKeyword,
			esDt:      EsDataType(field.Tag.Get("esDt")),
		})
	}

	for _, model := range fieldModelList {
		if v := model.Query(); v != "" {
			query2FieldModelMap[QueryField(model.Query())] = model
		}
	}
}

func SelectEsFiledByQuery(query ...string) []string {
	var ret []string
	for _, q := range query {
		if info, isOk := query2FieldModelMap[QueryField(q)]; isOk {
			ret = append(ret, info.ES())
		}
	}
	return ret
}

func GetFieldEs(query string) string {
	if info, isOk := query2FieldModelMap[QueryField(query)]; isOk {
		return info.ES()
	}
	return ""
}

func GetFieldEsKeyword(query string) string {
	if info, isOk := query2FieldModelMap[QueryField(query)]; isOk {
		return info.EsKeyword()
	} else {
		return info.ES()
	}
}
