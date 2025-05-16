package space_temp_config

import (
	"go-cs/internal/consts"
	"go-cs/internal/domain/search/search_es"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
)

func buildProcessingCond(processingStatusIds []int64) search_es.ConditionGroup {
	return search_es.ConditionGroup{
		Conjunction: search_es.AND,
		Conditions: []*search_es.Condition{
			{
				Field:    search_es.WorkItemStatusIdField,
				Operator: search_es.IN,
				Values:   utils.ToAnyArray(utils.ToStrArray(processingStatusIds)),
			},
		},
	}
}

func buildDoneCond(completeStatusIds []int64) search_es.ConditionGroup {
	return search_es.ConditionGroup{
		Conjunction: search_es.AND,
		Conditions: []*search_es.Condition{
			{Field: search_es.WorkItemStatusIdField, Operator: search_es.IN, Values: utils.ToAnyArray(utils.ToStrArray(completeStatusIds))},
		},
	}
}

func buildExpiredCond(processingStatusIds []int64) search_es.ConditionGroup {
	return search_es.ConditionGroup{
		Conjunction: search_es.AND,
		Conditions: []*search_es.Condition{
			{Field: search_es.PlanTimeField, Operator: search_es.LT, Values: []any{"${TODAY}", "${TODAY}"}},
			{Field: search_es.WorkItemStatusIdField, Operator: search_es.IN, Values: utils.ToAnyArray(utils.ToStrArray(processingStatusIds))},
		},
	}
}

// 总数
var totalCond = search_es.ConditionGroup{
	Conjunction: search_es.AND,
}

type CondItem struct {
	Desc    string                   `json:"desc"`
	Value   search_es.ConditionGroup `json:"value"`
	Type    string                   `json:"type"`
	Version string                   `json:"version"`
}

var validKeys = []string{
	consts.SpaceTempConfigKey_spaceWorkbenchCountConditions,
	consts.SpaceTempConfigKey_overviewOrder,
}

func IsValidKey(key string) bool {
	return stream.Contains(validKeys, key)
}

func GetDefaultByKey(key string, statusList work_item_status.WorkItemStatusItems) string {
	completeStatusItems := statusList.GetCompleteStatus()
	processingStatusItems := statusList.GetProcessingStatus()

	completeStatusList := stream.Map(completeStatusItems, func(item *work_item_status.WorkItemStatusItem) int64 {
		return item.Id
	})
	processingStatusList := stream.Map(processingStatusItems, func(item *work_item_status.WorkItemStatusItem) int64 {
		return item.Id
	})

	switch key {
	case consts.SpaceTempConfigKey_spaceWorkbenchCountConditions:
		return utils.ToJSON([]CondItem{
			{Desc: "待办任务", Value: buildProcessingCond(processingStatusList), Type: "processing"},
			{Desc: "完成任务", Value: buildDoneCond(completeStatusList), Type: "done"},
			{Desc: "已逾期任务", Value: buildExpiredCond(processingStatusList), Type: "expired"},
			{Desc: "任务总数", Value: totalCond, Type: "total"},
		})
	case consts.SpaceTempConfigKey_overviewOrder:
		return ""
	}

	return ""
}
