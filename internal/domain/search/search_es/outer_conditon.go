package search_es

import (
	"fmt"
	v1 "go-cs/api/search/v1"
)

func ConvertCondition(cg *v1.ConditionGroup) (*ConditionGroup, error) {
	group, err := convertCondition(cg)
	if err != nil {
		return nil, err
	}

	return group, err
}

func convertCondition(cg *v1.ConditionGroup) (*ConditionGroup, error) {
	if cg == nil {
		return nil, nil
	}

	var conditions []*Condition
	for _, v := range cg.Conditions {
		if v.Field == "" {
			return nil, fmt.Errorf("条件 field 不能为空")
		}
		if v.Operator == "" {
			return nil, fmt.Errorf("条件 operator 不能为空")
		}
		if len(v.Values) == 0 {
			return nil, fmt.Errorf("条件 values 不能为空")
		}

		if !IsValidOperator(v.Operator) {
			return nil, fmt.Errorf("条件 operatro: %v 不存在", v.Operator)
		}

		f, ok := query2FieldModelMap[QueryField(v.Field)]
		if !ok {
			return nil, fmt.Errorf("条件参数不支持 %q", v.Field)
		}

		conditions = append(conditions, &Condition{
			Field:     QueryField(v.Field),
			Values:    toAny(v.Values),
			Operator:  Operator(v.Operator),
			fieldInfo: f,
			spaceId:   v.SpaceId,
		})
	}

	var groups []*ConditionGroup
	for _, group := range cg.ConditionGroup {
		g, err := convertCondition(group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return &ConditionGroup{
		Conjunction:    Conjunction(cg.Conjunction),
		Conditions:     conditions,
		ConditionGroup: groups,
	}, nil
}
