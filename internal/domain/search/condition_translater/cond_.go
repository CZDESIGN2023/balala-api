package condition_translater

import (
	v1 "go-cs/api/search/v1"
	"go-cs/internal/consts"
)

// PlanTime 将 plan_time 转换为具体的条件
func (ctx *Ctx) StateDirectors(c *v1.Condition) *v1.ConditionGroup {
	spaceId := c.SpaceId

	return &v1.ConditionGroup{
		Conjunction: "AND",
		Conditions: []*v1.Condition{
			{Field: "work_item_type_key", Operator: "IN", Values: []string{string(consts.WorkItemTypeKey_StateTask)}, SpaceId: spaceId},
			{Field: "node_directors", Operator: c.Operator, Values: c.Values},
		},
	}
}
