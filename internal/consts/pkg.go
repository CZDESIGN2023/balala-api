package consts

type FlowScope string

const (
	FlowScope_All       FlowScope = ""
	FlowScope_Workflow  FlowScope = "work_flow"
	FlowScope_Stateflow FlowScope = "state_flow"
)

var FlowScopeList = []FlowScope{
	FlowScope_Workflow,
	FlowScope_Stateflow,
}

func ConvertFlowModeToFlowScope(flowMode WorkFlowMode) FlowScope {
	switch flowMode {
	case FlowMode_WorkFlow:
		return FlowScope_Workflow
	case FlowMode_StateFlow:
		return FlowScope_Stateflow
	default:
		return FlowScope_All
	}
}
