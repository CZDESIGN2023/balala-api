package condition_translater

import (
	v1 "go-cs/api/search/v1"
)

// FinishedAt 将 FinishedAt 转换为具体的条件
func (ctx *Ctx) CreatedAt(c *v1.Condition) *v1.ConditionGroup {

	var startTime string
	var endTime string
	if len(c.Values) >= 2 {
		startTime = c.Values[0]
		endTime = c.Values[1]
	}

	if endTime == "" {
		endTime = startTime
	}

	switch c.Operator {
	case "BETWEEN", "EQ":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: "BETWEEN", Values: []string{startTime, endTime}, SpaceId: c.SpaceId},
			},
		}
	case "NOT_EQ":
		return &v1.ConditionGroup{
			Conjunction: "OR",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: "LT", Values: []string{startTime}, SpaceId: c.SpaceId},
				{Field: c.Field, Operator: "GT", Values: []string{endTime}, SpaceId: c.SpaceId},
			},
		}
	case "LT":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: c.Operator, Values: []string{startTime}, SpaceId: c.SpaceId},
			},
		}
	case "GT":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: c.Operator, Values: []string{endTime}, SpaceId: c.SpaceId},
			},
		}
	case "LTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: c.Operator, Values: []string{endTime}, SpaceId: c.SpaceId},
			},
		}
	case "GTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: c.Field, Operator: c.Operator, Values: []string{startTime}, SpaceId: c.SpaceId},
			},
		}
	}

	return nil
}
