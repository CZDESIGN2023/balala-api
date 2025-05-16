package utils

import (
	"encoding/json"
	"go-cs/pkg/stream"
	"strconv"

	"github.com/spf13/cast"
)

func ToJSON(v any) string {
	return string(ToJSONBytes(v))
}

func ToJSONBytes(v any) []byte {
	marshal, _ := json.Marshal(v)
	return marshal
}

func ToStrArray[T any](list []T) []string {
	strings := stream.Map(list, func(v T) string {
		return cast.ToString(v)
	})

	if strings == nil {
		return []string{}
	}

	return strings
}

// ToInt64Array ["1","2","3"] --> [1,2,3]
func ToInt64Array(arr []string) []int64 {
	strings := stream.Map(arr, func(v string) int64 {
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	})

	if strings == nil {
		return []int64{}
	}

	return strings
}

// StrToStrArray '["1","2","3"]' --> ["1","2","3"]
func StrToStrArray(str string) []string {
	if str == "" {
		return nil
	}

	var list []string
	json.Unmarshal([]byte(str), &list)

	return list
}

// ToJSONArrayStr [1,2,3] --> '["1","2","3"]'
func ToJSONArrayStr(arr []int64) string {
	return ToJSON(ToStrArray(arr))
}

func ToAnyArray[E any](arr []E) []any {
	return stream.Map(arr, func(v E) any {
		return v
	})
}
