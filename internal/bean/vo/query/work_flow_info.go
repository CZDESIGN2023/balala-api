package query

import domain "go-cs/internal/domain/work_flow"

type WorkFlowInfoQuery struct {
	FlowId int64 `json:"flowId"`
}

type WorkFlowInfoQueryResult struct {
	WorkFlow         *domain.WorkFlow         `json:"work_flow"`
	WorkFlowTemplate *domain.WorkFlowTemplate `json:"work_flow_template"`
}
