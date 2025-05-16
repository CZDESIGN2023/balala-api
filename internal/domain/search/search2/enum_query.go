package search2

type QueryField string

const (
	WorkItemId       QueryField = "work_item_id"
	WorkItemFlowId   QueryField = "work_item_flow_id"
	Priority         QueryField = "priority"
	SpaceId          QueryField = "space_id"
	WorkObjectId     QueryField = "work_object_id"
	UserId           QueryField = "user_id"
	PlanTime         QueryField = "plan_time"
	PlanStartAt      QueryField = "plan_start_at"
	PlanCompleteAt   QueryField = "plan_complete_at"
	Directors        QueryField = "directors"
	Followers        QueryField = "followers"
	WorkItemStatus   QueryField = "work_item_status"
	Describe         QueryField = "describe"
	WorkItemName     QueryField = "work_item_name"
	VersionId        QueryField = "version_id"
	WorkItemStatusId QueryField = "work_item_status_id"
	FinishedAt       QueryField = "finished_at"
	NodeDirectors    QueryField = "node_directors"
	StateDirectors   QueryField = "state_directors"
)

var queryFieldStrMap = map[QueryField]string{
	WorkItemId:       "任务ID",
	WorkItemFlowId:   "任务流程",
	Priority:         "优先级",
	SpaceId:          "空间",
	WorkObjectId:     "模块",
	UserId:           "创建人",
	PlanTime:         "总排期",
	PlanStartAt:      "排期开始时间",
	PlanCompleteAt:   "排期结束时间",
	Directors:        "当前负责人",
	Followers:        "关注人",
	WorkItemStatus:   "任务状态",
	Describe:         "描述",
	WorkItemName:     "任务名",
	VersionId:        "版本",
	WorkItemStatusId: "任务状态",
	FinishedAt:       "任务完成时间",
	NodeDirectors:    "节点负责人",
	StateDirectors:   "状态负责人",
}

func (q QueryField) String() string {
	v, ok := queryFieldStrMap[q]

	if !ok {
		return string(q)
	}

	return v
}
