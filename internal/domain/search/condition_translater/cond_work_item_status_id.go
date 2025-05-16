package condition_translater

import (
	v1 "go-cs/api/search/v1"
	"go-cs/internal/consts"
)

// WorkItemStatusId 将 WorkItemStatusId 转换为具体的条件
func (ctx *Ctx) WorkItemStatusId(c *v1.Condition) *v1.ConditionGroup {
	spaceId := c.SpaceId

	var flowScope string
	if c.Attrs != nil {
		flowScope = c.Attrs["flow_scope"]
	}

	conditions := []*v1.Condition{c}

	if flowScope != "" {
		var values []string
		switch consts.FlowScope(flowScope) {
		case consts.FlowScope_Workflow:
			values = []string{string(consts.WorkItemTypeKey_Task), string(consts.WorkItemTypeKey_SubTask)}
		case consts.FlowScope_Stateflow:
			values = []string{string(consts.WorkItemTypeKey_StateTask)}
		}

		if len(values) != 0 {
			conditions = append(conditions, &v1.Condition{
				Field:    "work_item_type_key",
				Operator: "IN",
				Values:   values,
				SpaceId:  spaceId,
			})
		}
	}

	return &v1.ConditionGroup{
		Conjunction: "AND",
		Conditions:  conditions,
	}
}
