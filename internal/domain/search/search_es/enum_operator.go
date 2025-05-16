package search_es

import "slices"

type Operator string

const (
	IN          Operator = "IN"
	NOT_IN      Operator = "NOT_IN"
	EQ          Operator = "EQ"
	NOT_EQ      Operator = "NOT_EQ"
	GT          Operator = "GT"
	LT          Operator = "LT"
	GTE         Operator = "GTE"
	LTE         Operator = "LTE"
	INCLUDE     Operator = "INCLUDE"
	NOT_INCLUDE Operator = "NOT_INCLUDE"
	BETWEEN     Operator = "BETWEEN"
)

var allOperator = []Operator{IN, NOT_IN, EQ, NOT_EQ, GT, LT, GTE, LTE, INCLUDE, NOT_INCLUDE, BETWEEN}

func IsValidOperator(op string) bool {
	return slices.Contains(allOperator, Operator(op))
}

// Expand 是否需要展开参数
func (o Operator) Expand() bool {
	switch o {
	case IN, NOT_IN:
		return false
	}
	return true
}

func (o Operator) String() string {
	switch o {
	case INCLUDE:
		return "包含"
	case NOT_INCLUDE:
		return "不包含"
	case IN:
		return "存在选项属于"
	case NOT_IN:
		return "全部选项均不属于"
	case EQ:
		return "等于"
	case NOT_EQ:
		return "不等于"
	case GT:
		return "大于"
	case LT:
		return "小于"
	case GTE:
		return "大于等于"
	case LTE:
		return "小于等于"
	case BETWEEN:
		return "在区间内"
	default:
		return string(o)
	}
}
