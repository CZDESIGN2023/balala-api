package search2

import (
	"fmt"
	"go-cs/pkg/stream"
	"strings"
)

func SelectByQuery(query ...string) string {
	var ret []string
	for _, q := range query {
		info := query2FieldModelMap[QueryField(q)]

		var s = buildSelectField(info)
		if s == "" {
			s = q
		}

		ret = append(ret, s)
	}

	return strings.Join(ret, ", ")
}

func SelectByColumn(columnName ...string) string {
	var ret []string
	for _, q := range columnName {
		info := column2fieldModelMap[q]

		var s = buildSelectField(info)
		if s == "" {
			s = q
		}

		ret = append(ret, s)
	}

	return strings.Join(ret, ", ")
}

func SelectAll() string {
	var ret []string
	for _, info := range fieldModelList {
		if info.Dt().IsBig() { //查询排除大字段
			continue
		}

		var s = buildSelectField(info)

		ret = append(ret, s)
	}

	return strings.Join(ret, ", ")
}

func SelectExt(ext ...string) string {
	var ret []string
	for _, info := range fieldModelList {
		if info.Dt().IsBig() { //查询排除大字段
			continue
		}

		if info.query != "" && stream.Contains(ext, string(info.query)) {
			continue
		}

		var s = buildSelectField(info)

		ret = append(ret, s)
	}

	return strings.Join(ret, ", ")
}

// SelectAllWithBig 不排除大字段
func SelectAllWithBig() string {
	var ret []string
	for _, info := range fieldModelList {
		var s = buildSelectField(info)

		ret = append(ret, s)
	}

	return strings.Join(ret, ", ")
}

func buildSelectField(info FieldModel) string {
	var s = info.DB()
	if info.DB() != info.Gorm() {
		s = fmt.Sprintf("%v AS %v", info.DB(), info.Gorm())
	}

	return s
}
