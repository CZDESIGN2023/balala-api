package condition_translater

import (
	v1 "go-cs/api/search/v1"
	"go-cs/internal/domain/search/search_es"
)

func (ctx *Ctx) TranslateCond(cg *v1.ConditionGroup) *v1.ConditionGroup {
	if cg == nil {
		return nil
	}

	var groups []*v1.ConditionGroup
	for _, group := range cg.ConditionGroup {
		g := ctx.TranslateCond(group)
		groups = append(groups, g)
	}

	var conditions []*v1.Condition
	for _, condition := range cg.Conditions {
		// 处理变量
		ctx.handleVariable(condition)
		switch search_es.QueryField(condition.Field) {
		case "created_at":
			groups = append(groups, ctx.CreatedAt(condition))
		case "finished_at":
			groups = append(groups, ctx.FinishedAt(condition))
		case "node_directors":
			groups = append(groups, ctx.NodeDirectors(condition))
		case "state_directors":
			groups = append(groups, ctx.StateDirectors(condition))
		case "plan_time":
			groups = append(groups, ctx.PlanTime(condition))
		case "work_item_status_id":
			groups = append(groups, ctx.WorkItemStatusId(condition))
		default:
			conditions = append(conditions, condition)
		}
	}

	return &v1.ConditionGroup{
		Conjunction:    cg.Conjunction,
		Conditions:     conditions,
		ConditionGroup: groups,
	}
}
