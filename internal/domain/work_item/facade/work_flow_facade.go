package facade

import (
	"go-cs/internal/domain/work_flow"
)

type WorkFlowFacade struct {
	flow *work_flow.WorkFlow
}

func (w *WorkFlowFacade) Flow() *work_flow.WorkFlow {
	return w.flow
}

func BuildWorkFlowFacade(flow *work_flow.WorkFlow) *WorkFlowFacade {
	return &WorkFlowFacade{
		flow: flow,
	}
}
