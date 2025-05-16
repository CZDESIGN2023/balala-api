package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Int64ArrayToString Int64 array轉字串, 使用逗號,串接
func Int64ArrayToString(nums []int64) string {
	strNums := make([]string, len(nums))
	for i, num := range nums {
		strNums[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strNums, ",")
}

// Int64ArrayToStringWithSeparator Int64 array轉字串, 使用自訂分隔符號串接
func Int64ArrayToStringWithSeparator(nums []int64, separator string) string {
	strNums := make([]string, len(nums))
	for i, num := range nums {
		strNums[i] = strconv.FormatInt(num, 10)
	}
	return strings.Join(strNums, separator)
}

func JSONString2StringArray(s string) []string {
	var list []string
	if s == "" {
		return list
	}

	json.Unmarshal([]byte(s), &list)

	return list
}

// StringToInt64Array 逗號,串接的字串轉Int64 array
func StringToInt64Array(s string) []int64 {
	if s == "" {
		return []int64{}
	}

	strNums := strings.Split(s, ",")
	nums := make([]int64, len(strNums))

	for i, strNum := range strNums {
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return []int64{}
		}
		nums[i] = num
	}

	return nums
}

func StringArrToInt64Arr(s []string) []int64 {
	if s == nil {
		return []int64{}
	}
	nums := make([]int64, len(s))
	for i, strNum := range s {
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return []int64{}
		}
		nums[i] = num
	}
	return nums
}

// StringToInt64ArrayWithSeparator 自訂分隔符號串接的字串轉Int64 array
func StringToInt64ArrayWithSeparator(s string, separator string) []int64 {
	if s == "" {
		return []int64{}
	}

	strNums := strings.Split(s, separator)
	nums := make([]int64, len(strNums))

	for i, strNum := range strNums {
		num, err := strconv.ParseInt(strNum, 10, 64)
		if err != nil {
			return []int64{}
		}
		nums[i] = num
	}

	return nums
}
