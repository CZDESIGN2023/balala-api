package condition_translater

import v1 "go-cs/api/search/v1"

// PlanTime 将 plan_time 转换为具体的条件
func (ctx *Ctx) PlanTime(c *v1.Condition) *v1.ConditionGroup {

	var startTime string
	var endTime string
	if len(c.Values) >= 2 {
		startTime = c.Values[0]
		endTime = c.Values[1]
	}

	if endTime == "" {
		endTime = startTime
	}

	spaceId := c.SpaceId

	switch c.Operator {
	case "BETWEEN":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: "plan_start_at", Operator: "LT", Values: []string{endTime}, SpaceId: spaceId},
				{Field: "plan_complete_at", Operator: "GT", Values: []string{startTime}, SpaceId: spaceId},
			},
		}
	case "EQ":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: "plan_start_at", Operator: c.Operator, Values: []string{startTime}, SpaceId: spaceId},
				{Field: "plan_complete_at", Operator: c.Operator, Values: []string{endTime}, SpaceId: spaceId},
			},
		}
	case "NOT_EQ":
		return &v1.ConditionGroup{
			Conjunction: "OR",
			Conditions: []*v1.Condition{
				{Field: "plan_start_at", Operator: c.Operator, Values: []string{startTime}, SpaceId: spaceId},
				{Field: "plan_complete_at", Operator: c.Operator, Values: []string{endTime}, SpaceId: spaceId},
			},
		}
	case "LT", "LTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: "plan_complete_at", Operator: c.Operator, Values: []string{endTime}, SpaceId: spaceId},
			},
		}
	case "GT", "GTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: "plan_complete_at", Operator: c.Operator, Values: []string{endTime}, SpaceId: spaceId},
			},
		}
	}

	return nil
}
