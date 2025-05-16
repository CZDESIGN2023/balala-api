package search_es

import (
	esV8 "go-cs/internal/utils/es/v8"
	"reflect"
)

type CondType int

const (
	Normal = CondType(iota)
	JSONArray
)

type Condition struct {
	Field     QueryField `json:"field"`
	Values    []any      `json:"values"`
	Operator  Operator   `json:"operator"`
	fieldInfo FieldModel
	spaceId   int64
	Attrs     map[string]string `json:"attrs,omitempty"`
}

func (c *Condition) Equal(b *Condition) bool {
	if c == nil && b == nil {
		return true
	}
	if c == nil || b == nil {
		return false
	}
	return c.Field == b.Field &&
		c.Operator == b.Operator &&
		c.spaceId == b.spaceId &&
		reflect.DeepEqual(c.Values, b.Values)
}

type ConditionGroup struct {
	Conjunction    Conjunction       `json:"conjunction"`
	Conditions     []*Condition      `json:"conditions"`
	ConditionGroup []*ConditionGroup `json:"conditionGroup"`
}

func (g *ConditionGroup) Equal(b *ConditionGroup) bool {
	if g == nil && b == nil {
		return true
	}
	if g == nil || b == nil {
		return false
	}
	if g.Conjunction != b.Conjunction {
		return false
	}
	if len(g.Conditions) != len(b.Conditions) {
		return false
	}

	for i, v := range g.Conditions {
		if !v.Equal(b.Conditions[i]) {
			return false
		}
	}

	for i, v := range g.ConditionGroup {
		if !v.Equal(b.ConditionGroup[i]) {
			return false
		}
	}

	return true
}

func BuildCondition(group *ConditionGroup) *esV8.BoolQuery {
	fillFieldInfo(group) //填充field信息
	return buildCondition(group)
}

func buildCondition(group *ConditionGroup) *esV8.BoolQuery {

	if group == nil {
		return nil
	}

	if group.Conjunction == "" {
		group.Conjunction = AND
	}

	var clauses []esV8.Query

	for _, v := range group.Conditions {
		if len(v.Values) == 0 || v.Field == "" {
			continue
		}
		v.Values = ConvertValues(v.fieldInfo.Dt(), v.Values) //将values类型转换为数据库字段对应的类型

		exp := NewSearchCondExp(v.fieldInfo, v.Values...)

		var c esV8.Query
		//这里确实可以对各个具体的FiledModel进行exp， 之后回头再来优化
		//比如 user 这个字段可能是单选，也可能是需要二次解释的多选条件 需要进行转换成对应的搜索条件
		switch v.Operator {
		case IN:
			c = exp.In()
		case NOT_IN:
			c = exp.NotIn()
		case EQ:
			c = exp.Eq()
		case NOT_EQ:
			c = exp.NotEq()
		case GT:
			c = exp.Gt()
		case LT:
			c = exp.Lt()
		case GTE:
			c = exp.Gte()
		case LTE:
			c = exp.Lte()
		case INCLUDE:
			c = exp.Include()
		case NOT_INCLUDE:
			c = exp.Exclude()
		case BETWEEN:
			c = exp.Between()
		default:
			panic("未知的operator: " + v.Operator)
		}

		if c != nil {
			if v.spaceId != 0 {
				clause := esV8.NewBoolQuery().Must(c, esV8.NewTermQuery("space_id", v.spaceId))
				clauses = append(clauses, clause)
			} else {
				clauses = append(clauses, c)
			}
		}
	}

	query := esV8.NewBoolQuery()

	switch group.Conjunction {
	case AND:
		if len(clauses) != 0 { // 如果为空，则不添加该条件
			query.Must(clauses...)
		}
	case OR:
		if len(clauses) != 0 { // 如果为空，则不添加该条件
			query.Should(clauses...).MinimumShouldMatch("1")
		}
	}

	for _, v := range group.ConditionGroup {
		clauseCondGroup := buildCondition(v)
		switch group.Conjunction {
		case AND:
			query.Must(clauseCondGroup)
		case OR:
			query.Should(clauseCondGroup).MinimumShouldMatch("1")
		}
	}
	return query
}

func fillFieldInfo(g *ConditionGroup) {
	if g == nil {
		return
	}

	for _, v := range g.Conditions {
		if info, ok := query2FieldModelMap[v.Field]; ok {
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
