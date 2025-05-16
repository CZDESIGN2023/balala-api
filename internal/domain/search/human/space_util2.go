package human

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"go-cs/internal/biz/space_temp_config"
	"go-cs/internal/domain/search/search_es"
	wFlow_repo "go-cs/internal/domain/work_flow/repo"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"go-cs/pkg/stream/tuple"
	"slices"
	"strings"
)

/*
转换为人类可读内容
*/

type ConditionGroupParser struct {
	spaceId    int64
	statusRepo witem_status_repo.WorkItemStatusRepo
	flowRepo   wFlow_repo.WorkFlowRepo
	ctx        context.Context
}

func NewConditionGroupParser(spaceId int64, statusRepo witem_status_repo.WorkItemStatusRepo, flowRepo wFlow_repo.WorkFlowRepo) *ConditionGroupParser {
	return &ConditionGroupParser{
		spaceId:    spaceId,
		statusRepo: statusRepo,
		flowRepo:   flowRepo,
		ctx:        context.Background(),
	}
}

func (p *ConditionGroupParser) Parse(condGroup search_es.ConditionGroup) string {
	return p.parse(condGroup)
}

func (p *ConditionGroupParser) parse(condGroup search_es.ConditionGroup) string {
	parsedConditions := stream.Map(condGroup.Conditions, func(v *search_es.Condition) string {
		return p.parseCondition(v)
	})

	if len(parsedConditions) == 1 {
		return parsedConditions[0]
	}

	return strings.Join(parsedConditions, " "+condGroup.Conjunction.String()+" ")
}

func parseVariable(val string) (string, bool) {
	switch val {
	case "${TODAY}":
		return "今天", true
	case "${TOMORROW}":
		return "明天", true
	case "${YESTERDAY}":
		return "昨天", true
	// 本周
	case "${THIS_WEEK}":
		return "本周", true
	// 下周
	case "${NEXT_WEEK}":
		return "下周", true
	// 上周
	case "${LAST_WEEK}":
		return "上周", true
	// 本月
	case "${THIS_MONTH}":
		return "本月", true
	// 上月
	case "${LAST_MONTH}":
		return "上月", true
	// 下月
	case "${NEXT_MONTH}":
		return "下月", true
	//今年
	case "${THIS_YEAR}":
		return "今年", true
	//明年
	case "${NEXT_YEAR}":
		return "明年", true
	//去年
	case "${LAST_YEAR}":
		return "去年", true
	}

	return val, false
}

func (p *ConditionGroupParser) parseCondition(condition *search_es.Condition) string {
	var field = condition.Field.String()
	var operator = condition.Operator.String()
	var value string

	conditionValues := utils.ToStrArray(condition.Values)

	switch condition.Field {
	case search_es.WorkItemStatusIdField:
		statusMap, err := p.statusRepo.WorkItemStatusMap(p.ctx, p.spaceId)
		if err != nil {
			return ""
		}
		values := stream.Map(conditionValues, func(v string) string {
			status := statusMap[cast.ToInt64(v)]
			if status != nil {
				return status.Name
			}

			return ""
		})

		value = strings.Join(values, "、")

	case search_es.PlanTimeField:
		var isVariable bool
		values := stream.Map(conditionValues, func(v string) string {
			if val, ok := parseVariable(v); ok {
				isVariable = true
				return val
			}

			return strings.Split(v, " ")[0]
		})

		if isVariable {
			value = values[0]
		} else {
			value = strings.Join(values, " ～ ")
		}

	case search_es.WorkItemFlowIdField:
		flowMap, err := p.flowRepo.WorkFlowMap(p.ctx, utils.ToInt64Array(conditionValues))
		if err != nil {
			return ""
		}

		values := stream.Map(conditionValues, func(v string) string {
			id := cast.ToInt64(v)
			return flowMap[id].Name
		})

		value = strings.Join(values, "、")
	case search_es.FinishedAtField:
		values := stream.Map(conditionValues, func(v string) string {
			if val, ok := parseVariable(v); ok {
				return val
			}

			return strings.Split(v, " ")[0]
		})

		if condition.Operator != search_es.BETWEEN {
			values = values[:1]
		}

		value = strings.Join(values, " ～ ")
	}

	return field + " " + operator + " " + value
}

func GetUpdatedConditionGroup(oldVal, newVal string) (res []tuple.Pair[space_temp_config.CondItem, space_temp_config.CondItem], orderChanged bool) {
	var oldGroupList []space_temp_config.CondItem
	var newGroupList []space_temp_config.CondItem

	json.Unmarshal([]byte(oldVal), &oldGroupList)
	json.Unmarshal([]byte(newVal), &newGroupList)

	oldGroupMap := stream.ToMap(oldGroupList, func(_ int, v space_temp_config.CondItem) (string, space_temp_config.CondItem) {
		return v.Type, v
	})

	newGroupMap := stream.ToMap(newGroupList, func(_ int, v space_temp_config.CondItem) (string, space_temp_config.CondItem) {
		return v.Type, v
	})

	delete(oldGroupMap, "followed")
	delete(newGroupMap, "followed")

	var updatedCondList []tuple.Pair[space_temp_config.CondItem, space_temp_config.CondItem]

	for k, newVal := range newGroupMap {
		oldVal := oldGroupMap[k]
		if oldVal.Desc != newVal.Desc || !oldVal.Value.Equal(&newVal.Value) {
			updatedCondList = append(updatedCondList, tuple.T2(oldVal, newVal))
		}
	}

	if updatedCondList != nil {
		return updatedCondList, false
	}

	//检查顺序是否变化了
	oldKeys := stream.Map(oldGroupList, func(v space_temp_config.CondItem) string {
		return v.Type
	})
	newKeys := stream.Map(newGroupList, func(v space_temp_config.CondItem) string {
		return v.Type
	})

	return nil, !slices.Equal(oldKeys, newKeys)
}

type Update struct {
	Field    string
	OldValue string
	NewValue string
}

func ParseToChangeLog(cgParser *ConditionGroupParser, t tuple.Pair[space_temp_config.CondItem, space_temp_config.CondItem]) []Update {
	var updates []Update

	if t.First.Desc != t.Second.Desc {
		updates = append(updates, Update{
			Field:    "desc",
			OldValue: t.First.Desc,
			NewValue: t.Second.Desc,
		})
	}

	oldCondVal := cgParser.Parse(t.First.Value)
	newCondVal := cgParser.Parse(t.Second.Value)
	if oldCondVal != newCondVal {
		updates = append(updates, Update{
			Field:    "condition",
			OldValue: oldCondVal,
			NewValue: newCondVal,
		})
	}

	return updates
}
