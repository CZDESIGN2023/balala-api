package condition_translater

import (
	"fmt"
	"regexp"
)

var isVariableExp = regexp.MustCompile(`^\$\{(.+)\}$`)

func isVariable(s string) bool {
	fmt.Println(isVariableExp.FindStringSubmatch(s))

	return isVariableExp.MatchString(s)
}

func extractVariable(s string) string {
	match := isVariableExp.FindStringSubmatch(s)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}
