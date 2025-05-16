package search2

import (
	"testing"
)

func TestBuildCondition(t *testing.T) {
	group := &ConditionGroup{
		Conjunction: AND,
		ConditionGroup: []*ConditionGroup{
			{
				Conjunction: OR,
				Conditions: []*Condition{
					{Field: "user_id", Values: []any{"42"}, Operator: EQ},
					{Field: "directors", Values: []any{"42"}, Operator: IN},
				},
			},
			{
				Conjunction: AND,
				Conditions: []*Condition{
					{Field: "user_id", Values: []any{"42"}, Operator: EQ},
					{Field: "plan_start_at", Values: []any{"2023/01/01 00:00:00"}, Operator: GT},
				},
			},
		},
	}

	q, args := BuildCondition(group)
	t.Log(q)
	t.Log(args)
}
