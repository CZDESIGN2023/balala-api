package query

import domain "go-cs/internal/domain/work_flow"

type TaskWorkFlowListQuery struct {
	//空间id
	SpaceId int64
	//主任务类型id
	WorkItemTypeIds []int64
}

type TaskWorkFlowListQueryResult struct {
	List []*TaskWorkFlowListQueryResult_Item `json:"list"`
}

type TaskWorkFlowListQueryResult_Item struct {
	WorkFlow         *domain.WorkFlow         `json:"work_flow"`
	WorkFlowTemplate *domain.WorkFlowTemplate `json:"work_flow_template"`
}
