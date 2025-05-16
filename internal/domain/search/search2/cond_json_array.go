package search2

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Operate interface {
	In() string
	NotIn() string
	Eq() string
	NotEq() string
	Gt() string
	Lt() string
	Gte() string
	Lte() string
	Include() string
	Exclude() string
	Between() string
}

type CondJSONArray Condition

func (c *CondJSONArray) In() string {
	values := c.Values
	fieldArr := c.GetField()
	c.ClearFieldAndValues()

	var exps []string
	for _, v := range values {
		exps = append(exps, fmt.Sprintf("JSON_CONTAINS(%v, '%q')", fieldArr, v))
	}

	return "(" + strings.Join(exps, " OR ") + ")"
}

func (c *CondJSONArray) NotIn() string {
	values := c.Values
	fieldArr := c.GetField()
	c.ClearFieldAndValues()

	var exps []string
	for _, v := range values {
		exps = append(exps, fmt.Sprintf("NOT JSON_CONTAINS(%v, '%q')", fieldArr, v))
	}

	return "(" + strings.Join(exps, " AND ") + ")"
}

func (c *CondJSONArray) Eq() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(JSON_CONTAINS(%v, %v) AND JSON_LENGTH(%v) = JSON_LENGTH(%v))", fieldArr, valueArr, fieldArr, valueArr)
}

func (c *CondJSONArray) NotEq() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(NOT JSON_CONTAINS(%v, %v) OR JSON_LENGTH(%v) != JSON_LENGTH(%v))", fieldArr, valueArr, fieldArr, valueArr)
}

func (c *CondJSONArray) Gt() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(JSON_CONTAINS(%v, %v) AND JSON_LENGTH(%v) > JSON_LENGTH(%v))", fieldArr, valueArr, fieldArr, valueArr)
}

func (c *CondJSONArray) Lt() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(JSON_CONTAINS(%v, %v) AND JSON_LENGTH(%v) < JSON_LENGTH(%v))", valueArr, fieldArr, fieldArr, valueArr)
}

func (c *CondJSONArray) Gte() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(JSON_CONTAINS(%v, %v) AND JSON_LENGTH(%v) >= JSON_LENGTH(%v))", fieldArr, valueArr, fieldArr, valueArr)
}

func (c *CondJSONArray) Lte() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("(JSON_CONTAINS(%v, %v) AND JSON_LENGTH(%v) <= JSON_LENGTH(%v))", valueArr, fieldArr, fieldArr, valueArr)
}

func (c *CondJSONArray) Include() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("JSON_CONTAINS(%v, %v)", fieldArr, valueArr)
}

func (c *CondJSONArray) Exclude() string {
	fieldArr := c.GetField()
	valueArr := c.GetValuesJSON()
	c.ClearFieldAndValues()

	return fmt.Sprintf("NOT JSON_CONTAINS(%v, %v)", fieldArr, valueArr)
}

func (c *CondJSONArray) Between() string {
	// 数组没有 between 运算
	c.ClearFieldAndValues()
	return ""
}

func (c *CondJSONArray) GetField() string {
	f := c.fieldInfo.DB()
	//f = strings.Replace(f, "->>", "->", 1)
	return f
}

func (c *CondJSONArray) GetValuesJSON() string {
	marshal, _ := json.Marshal(c.Values)
	return fmt.Sprintf("'%s'", marshal)
}

func (c *CondJSONArray) ClearFieldAndValues() {
	c.Values = nil
	//c.Field = ""
}
