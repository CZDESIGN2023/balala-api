package condition_translater

import (
	"github.com/spf13/cast"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/domain/work_item_status"
	"go-cs/pkg/stream"
)

// FinishedAt 将 FinishedAt 转换为具体的条件
func (ctx *Ctx) FinishedAt(c *v1.Condition) *v1.ConditionGroup {

	var startTime string
	var endTime string
	if len(c.Values) >= 2 {
		startTime = c.Values[0]
		endTime = c.Values[1]
	}

	if endTime == "" {
		endTime = startTime
	}

	spaceIds := ctx.SpaceIds
	if c.SpaceId != 0 {
		spaceIds = []int64{c.SpaceId}
	}

	status, _ := ctx.StatusRepo.GetWorkItemStatusItemsBySpaceIds(ctx.Ctx, spaceIds)

	completeStatus := status.GetCompleteStatus()
	completeStatusIds := stream.Map(completeStatus, func(item *work_item_status.WorkItemStatusItem) string {
		return cast.ToString(item.Id)
	})

	completeStatusCond := &v1.Condition{
		Field:    "work_item_status_id",
		Operator: "IN",
		Values:   completeStatusIds,
		SpaceId:  c.SpaceId,
	}

	field := "last_status_at"

	switch c.Operator {
	case "BETWEEN", "EQ":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
				{Field: field, Operator: "BETWEEN", Values: []string{startTime, endTime}, SpaceId: c.SpaceId},
			},
		}
	case "NOT_EQ":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
			},
			ConditionGroup: []*v1.ConditionGroup{
				{
					Conjunction: "OR",
					Conditions: []*v1.Condition{
						{Field: field, Operator: "LT", Values: []string{startTime}, SpaceId: c.SpaceId},
						{Field: field, Operator: "GT", Values: []string{endTime}, SpaceId: c.SpaceId},
					},
				},
			},
		}
	case "LT":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
				{Field: field, Operator: c.Operator, Values: []string{startTime}, SpaceId: c.SpaceId},
			},
		}
	case "GT":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
				{Field: field, Operator: c.Operator, Values: []string{endTime}, SpaceId: c.SpaceId},
			},
		}
	case "LTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
				{Field: field, Operator: c.Operator, Values: []string{endTime}, SpaceId: c.SpaceId},
			},
		}
	case "GTE":
		return &v1.ConditionGroup{
			Conjunction: "AND",
			Conditions: []*v1.Condition{
				completeStatusCond,
				{Field: field, Operator: c.Operator, Values: []string{startTime}, SpaceId: c.SpaceId},
			},
		}
	}

	return nil
}
