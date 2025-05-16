package search2

import "fmt"

type Conjunction string

const (
	OR  Conjunction = "OR"
	AND Conjunction = "AND"
)

func (c Conjunction) String() string {
	switch c {
	case AND:
		return "且"
	case OR:
		return "或"
	}

	return fmt.Sprintf("非法的 Conjunction: %v", string(c))
}
