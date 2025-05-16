package service

import (
	v1 "go-cs/api/search/v1"
	"testing"
)

func Test(t *testing.T) {
	g := clearCondGroup(&v1.ConditionGroup{
		ConditionGroup: []*v1.ConditionGroup{
			{
				Conjunction: "AND",
				Conditions: []*v1.Condition{
					{
						Field:    "1",
						Values:   []string{"1"},
						Operator: "1",
						SpaceId:  0,
					},
				},
			},
		},
	})

	t.Log(g)
}
