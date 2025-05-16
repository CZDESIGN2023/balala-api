package data

import (
	"context"
	search22 "go-cs/internal/domain/search/search2"
	"go-cs/internal/utils/date"
	"go-cs/pkg/pprint"
	"go-cs/pkg/stream"
	"testing"

	"github.com/spf13/cast"
)

func TestSearchRepo_SearchWorkItem(t *testing.T) {
	userId := 42
	uidStr := cast.ToString(userId)

	begin, end := date.ThisWeekBeginEnd()

	var weekStart = date.Format(begin)
	var weekEnd = date.Format(end)

	weekCond := &search22.ConditionGroup{
		Conjunction: "OR",
		Conditions: []*search22.Condition{
			{Field: "plan_start_at", Operator: "BETWEEN", Values: []any{weekStart, weekEnd}},
			{Field: "plan_complete_at", Operator: "BETWEEN", Values: []any{weekStart, weekEnd}},
		},
		ConditionGroup: []*search22.ConditionGroup{
			{
				Conjunction: "AND",
				Conditions: []*search22.Condition{
					{Field: "plan_start_at", Operator: "LT", Values: []any{weekStart}},
					{Field: "plan_complete_at", Operator: "GT", Values: []any{weekEnd}},
				},
			},
		},
	}

	userCond :=
		&search22.ConditionGroup{
			Conjunction: search22.OR,
			Conditions: []*search22.Condition{
				{
					Field:    "user_id",
					Operator: search22.EQ,
					Values:   []any{uidStr},
				},
				{
					Field:    "directors",
					Operator: search22.IN,
					Values:   []any{uidStr},
				},
			},
		}

	processingStatus := stream.Map([]int64{1, 3, 4, 5, 9, 10}, func(v int64) any {
		return v
	})
	group := &search22.ConditionGroup{
		Conjunction: search22.AND,
		Conditions: []*search22.Condition{
			{
				Field:    "work_item_status",
				Operator: search22.IN,
				Values:   processingStatus,
			},
		},
		ConditionGroup: []*search22.ConditionGroup{
			userCond,
			weekCond,
		},
	}

	ids, err := SpaceRepo.GetUserSpaceIds(context.Background(), 21)
	if err != nil {
		t.Error(err)
		return
	}

	item, err := SearchRepo.QueryWorkItem(context.Background(), ids, group, search22.SelectByQuery("id", "directors"))
	if err != nil {
		t.Error(err)
	}
	pprint.Println(item)
}

func TestSearchRepo_WorkBenchSearchWorkItem(t *testing.T) {
	userId := int64(152)
	uidStr := cast.ToString(userId)

	userCond := &search22.Condition{
		Field:    "directors",
		Operator: search22.IN,
		Values:   []any{uidStr},
	}

	// endStatusList := stream.Map(consts.WorkItemEndStatusList(), func(v int64) any {
	// 	return cast.ToString(v)
	// })

	endStatusList := []any{"2", "3", "4", "5", "9", "10"}
	statusCond := &search22.Condition{
		Field:    "work_item_status",
		Operator: search22.NOT_IN,
		Values:   endStatusList,
	}

	group := &search22.ConditionGroup{
		Conjunction: search22.AND,

		ConditionGroup: []*search22.ConditionGroup{
			{
				Conjunction: search22.AND,
				Conditions: []*search22.Condition{
					userCond,
					statusCond,
				},
			},
			{
				Conditions: []*search22.Condition{
					{
						SpaceId:  87,
						Field:    "directors",
						Operator: search22.IN,
						Values:   []any{"152"},
					},
				},
			},
		},
	}

	ids, err := SpaceRepo.GetUserSpaceIds(context.Background(), userId)
	if err != nil {
		t.Error(err)
		return
	}

	item, err := SearchRepo.QueryWorkItem(context.Background(), ids, group, search22.SelectByQuery("id", "directors"))
	if err != nil {
		t.Error(err)
	}
	pprint.Println(item)
}

func TestSearchRepo_PendingWorkItem(t *testing.T) {
	ids, err := SpaceRepo.GetUserSpaceIds(context.Background(), 42)
	if err != nil {
		t.Error(err)
		return
	}

	item, err := SearchRepo.PendingWorkItem(context.Background(), 42, ids)
	if err != nil {
		t.Error(err)
	}
	pprint.Println(item)
}

func TestSearchRepo_QueryWorkItemEsByPid(t *testing.T) {
	item, err := SearchRepo.QueryWorkItemEsByPid(context.Background(), []int64{618})
	if err != nil {
		t.Error(err)
	}
	pprint.Println(item)
}
