package biz

import (
	"context"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/test/dbt"
	"go-cs/pkg/pprint"
	"testing"
)

func TestGetSpaceWorkbenchCount2(t *testing.T) {
	cgs := []*v1.ConditionGroup{
		{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				{Field: "user_id", Values: []string{"42"}, Operator: "IN"},
			},
		},
	}

	info, err := dbt.UC.StaticsUsecase.GetSpaceWorkbenchCount2(context.Background(), 42, 87, cgs)
	if err != nil {
		t.Error(err)
	}

	pprint.Println(info)
}
