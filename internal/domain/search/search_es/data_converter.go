package search_es

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

func ConvertValues(typ DataType, values []any) []any {
	f, ok := convertFuncMap[typ]
	if ok {
		return f(values)
	}

	return toAny(values)
}

var convertFuncMap = map[DataType]func([]any) []any{
	Date:      DateFunc,
	DateRange: DateFunc,
	Integer:   IntegerFunc,
}

func DateFunc(values []any) []any {
	var ret []any
	for _, value := range values {
		s := value
		parse, err := time.ParseInLocation("2006/01/02 15:04:05", s.(string), time.Local)
		if err != nil {
			fmt.Println(err)
		}
		ret = append(ret, parse.Unix())
	}
	return ret
}

func IntegerFunc(values []any) []any {
	var ret []any
	for _, value := range values {
		ret = append(ret, cast.ToInt64(value))
	}
	return ret
}
