package search_es

type QueryField string

const (
	WorkItemIdField       QueryField = "work_item_id"
	PidField              QueryField = "pid"
	SpaceIdField          QueryField = "space_id"
	UserIdField           QueryField = "user_id"
	WorkItemTypeField     QueryField = "work_item_type"
	WorkObjectIdField     QueryField = "work_object_id"
	WorkItemFlowIdField   QueryField = "work_item_flow_id"
	WorkItemGuidField     QueryField = "work_item_guid"
	WorkItemNameField     QueryField = "work_item_name"
	WorkItemStatusField   QueryField = "work_item_status"
	WorkItemStatusIdField QueryField = "work_item_status_id"
	CreatedAtField        QueryField = "created_at"
	UpdatedAtField        QueryField = "updated_at"
	DeletedAtField        QueryField = "deleted_at"
	VersionIdField        QueryField = "version_id"
	ChildNumField         QueryField = "child_num"
	FlowModeField         QueryField = "flow_mode"
	FlowModeVersionField  QueryField = "flow_mode_version"
	FlowModeCodeField     QueryField = "flow_mode_code"
	LastStatusAtField     QueryField = "last_status_at"
	LastStatusField       QueryField = "last_status"
	IsRestartField        QueryField = "is_restart"
	RestartAtField        QueryField = "restart_at"
	RestartUserIdField    QueryField = "restart_user_id"
	IconFlagsField        QueryField = "icon_flags"
	CommentNumField       QueryField = "comment_num"
	ResumeAtField         QueryField = "resume_at"
	DescribeField         QueryField = "describe"
	RemarkField           QueryField = "remark"
	ProcessRateField      QueryField = "process_rate"
	PriorityField         QueryField = "priority"
	PlanStartAtField      QueryField = "plan_start_at"
	PlanCompleteAtField   QueryField = "plan_complete_at"
	TagsField             QueryField = "tags"
	DirectorsField        QueryField = "directors"
	FollowersField        QueryField = "followers"
	ParticipatorsField    QueryField = "participators"
	FlowMode              QueryField = "flow_mode"

	// 混合字段
	PlanTimeField       QueryField = "plan_time"
	FinishedAtField     QueryField = "finished_at"
	NodeDirectorsField  QueryField = "node_directors"
	StateDirectorsField QueryField = "state_directors"
)

var queryFieldStrMap = map[QueryField]string{
	WorkItemIdField:       "任务ID",
	WorkItemFlowIdField:   "任务流程",
	PriorityField:         "优先级",
	SpaceIdField:          "空间",
	WorkObjectIdField:     "模块",
	UserIdField:           "创建人",
	PlanTimeField:         "总排期",
	PlanStartAtField:      "排期开始时间",
	PlanCompleteAtField:   "排期结束时间",
	DirectorsField:        "当前负责人",
	FollowersField:        "关注人",
	WorkItemStatusField:   "任务状态",
	DescribeField:         "描述",
	WorkItemNameField:     "任务名",
	VersionIdField:        "版本",
	WorkItemStatusIdField: "任务状态",
	FinishedAtField:       "任务完成时间",
	NodeDirectorsField:    "节点负责人",
	StateDirectorsField:   "状态负责人",
}

func (q QueryField) String() string {
	v, ok := queryFieldStrMap[q]

	if !ok {
		return string(q)
	}

	return v
}

func (q QueryField) Es() string {
	if info, isOk := query2FieldModelMap[q]; isOk {
		return info.ES()
	}
	return ""
}
func (q QueryField) EsKeyword() string {
	if info, isOk := query2FieldModelMap[q]; isOk {
		return info.EsKeyword()
	}
	return ""
}
