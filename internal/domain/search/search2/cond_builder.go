package search2

import (
	"fmt"
	"strings"
)

type CondType int

const (
	Normal CondType = iota
	JSONArray
)

type Condition struct {
	Field    QueryField        `json:"field"`
	Values   []any             `json:"values"`
	Operator Operator          `json:"operator"`
	SpaceId  int64             `json:"spaceId,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`

	fieldInfo FieldModel
	exp       string
}

type ConditionGroup struct {
	Conjunction    Conjunction       `json:"conjunction,omitempty"`
	Conditions     []*Condition      `json:"conditions,omitempty"`
	ConditionGroup []*ConditionGroup `json:"condition_group,omitempty"`
}

func fillExp(group *ConditionGroup) {
	if group == nil {
		return
	}

	for _, c := range group.Conditions {
		if len(c.Values) == 0 || c.Field == "" {
			continue
		}
		c.Values = ConvertValues(c.fieldInfo.Dt(), c.Values) //将values类型转换为数据库字段对应的类型

		v := convertCond(c) //将条件转换为具体类型
		var exp string      //构建的表达式

		switch c.Operator {
		case IN:
			exp = v.In()
		case NOT_IN:
			exp = v.NotIn()
		case EQ:
			exp = v.Eq()
		case NOT_EQ:
			exp = v.NotEq()
		case GT:
			exp = v.Gt()
		case LT:
			exp = v.Lt()
		case GTE:
			exp = v.Gte()
		case LTE:
			exp = v.Lte()
		case INCLUDE:
			exp = v.Include()
		case NOT_INCLUDE:
			exp = v.Exclude()
		case BETWEEN:
			exp = v.Between()
		default:
			panic("未知的operator: " + c.Operator)
		}

		if c.SpaceId != 0 {
			exp = fmt.Sprintf("space_id = %v AND %v", c.SpaceId, exp)
		}

		c.exp = exp
	}

	for _, condGroup := range group.ConditionGroup {
		fillExp(condGroup)
	}
}

func BuildCondition(group *ConditionGroup) (sql string, values []any) {
	fillFieldInfo(group) //填充field信息
	fillExp(group)

	return buildCondition(group)
}

func buildCondition(group *ConditionGroup) (sql string, values []any) {
	if group == nil {
		return "", nil
	}

	var ret []string

	for _, v := range group.Conditions {
		ret = append(ret, v.exp)
		if v.Values == nil {
			continue
		}
		if v.Operator.Expand() {
			values = append(values, toAny(v.Values)...)
		} else {
			values = append(values, v.Values)
		}
	}

	for _, condGroup := range group.ConditionGroup {
		q, args := buildCondition(condGroup)
		if q != "" {
			ret = append(ret, q)
			values = append(values, toAny(args)...)
		}
	}

	if len(ret) == 0 {
		return "", nil
	}

	sep := fmt.Sprintf(" %v ", group.Conjunction)
	return "(" + strings.Join(ret, sep) + ")", values
}

func convertCond(c *Condition) Operate {
	switch c.fieldInfo.Dt() {
	case MultiSelect, MultiUser:
		return (*CondJSONArray)(c)
	default:
		return (*CondNormal)(c)
	}
}

func fillFieldInfo(g *ConditionGroup) {
	if g == nil {
		return
	}

	for _, v := range g.Conditions {
		info, ok := query2FieldModelMap[v.Field]

		if ok {
			v.fieldInfo = info
		}
	}

	for _, condGroup := range g.ConditionGroup {
		fillFieldInfo(condGroup)
	}
}

func toAny[T any](arr []T) []any {
	if arr == nil {
		return nil
	}
	var ret []any
	for _, v := range arr {
		ret = append(ret, v)
	}
	return ret
}

func (c *Condition) SetAttr(key, value string) string {
	if c.Attrs == nil {
		c.Attrs = make(map[string]string)
	}

	oldVal := c.Attrs[key]
	c.Attrs[key] = value
	return oldVal
}

func (c *Condition) Attr(key string) string {
	if c.Attrs == nil {
		c.Attrs = make(map[string]string)
	}

	return c.Attrs[key]
}
