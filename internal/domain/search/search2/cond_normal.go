package search2

import (
	"fmt"
	"github.com/spf13/cast"
)

type CondNormal Condition

func (c *CondNormal) In() string {
	return c.buildExp("IN ?")
}

func (c *CondNormal) NotIn() string {
	return c.buildExp("NOT IN ?")
}

func (c *CondNormal) Eq() string {
	c.keepNValues(1)
	return c.buildExp("= ?")
}

func (c *CondNormal) NotEq() string {
	c.keepNValues(1)
	return c.buildExp("!= ?")
}

func (c *CondNormal) Gt() string {
	c.keepNValues(1)
	return c.buildExp("> ?")
}

func (c *CondNormal) Lt() string {
	c.keepNValues(1)
	return c.buildExp("< ?")
}

func (c *CondNormal) Gte() string {
	c.keepNValues(1)
	return c.buildExp(">= ?")
}

func (c *CondNormal) Lte() string {
	c.keepNValues(1)
	return c.buildExp("<= ?")
}
func (c *CondNormal) Include() string {
	c.keepNValues(1)
	c.Values[0] = "%" + cast.ToString(c.Values[0]) + "%"
	return c.buildExp("LIKE ?")
}

func (c *CondNormal) Exclude() string {
	c.keepNValues(1)
	c.Values[0] = "%" + cast.ToString(c.Values[0]) + "%"
	return c.buildExp("NOT LIKE ?")
}

func (c *CondNormal) Between() string {
	c.keepNValues(2)
	return c.buildExp("BETWEEN ? AND ?")
}

func (c *CondNormal) buildExp(op string) string {
	return fmt.Sprintf("%v %v", c.fieldInfo.DB(), op)
}

func (c *CondNormal) keepNValues(n int) {
	c.Values = c.Values[:n]
	return
}
