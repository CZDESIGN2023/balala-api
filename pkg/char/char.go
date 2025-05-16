package char

import "unicode"

// Filter 移除所有不可打印的字符
func Filter(s string) string {
	var ret []rune
	r := []rune(s)
	for _, v := range r {
		if unicode.IsPrint(v) && !unicode.IsMark(v) {
			ret = append(ret, v)
		}
	}

	return string(ret)
}
