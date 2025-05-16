package search2

import (
	"fmt"
	"strings"
)

var fieldModelList []FieldModel
var query2FieldModelMap = map[QueryField]FieldModel{} // query: FieldModel
var column2fieldModelMap = map[string]FieldModel{}    // query: FieldModel
var es2fieldModelMap = map[string]FieldModel{}        // query: FieldModel

func FieldByQuery(query QueryField) FieldModel {
	return query2FieldModelMap[query]
}

func FieldByColumn(column string) FieldModel {
	return column2fieldModelMap[column]
}

type FieldModel struct {
	query QueryField //外部字段名
	db    string
	gorm  string
	es    string
	dt    DataType //数据类型
	inObj bool     //是否在json对象中
}

func (f *FieldModel) DB() string {
	if f.db == "" {
		return ""
	}

	if f.inObj {
		split := strings.Split(f.db, "->>")
		if len(split) == 2 {
			return fmt.Sprintf("%v->>%v", quote(split[0]), split[1])
		}
		split = strings.Split(f.db, "->")
		if len(split) == 2 {
			return fmt.Sprintf("%v->%v", quote(split[0]), split[1])
		}
	}

	return quote(f.db)
}

func (f *FieldModel) Query() QueryField {
	return f.query
}

func (f *FieldModel) Dt() DataType {
	return f.dt
}

func (f *FieldModel) Gorm() string {
	if f.gorm == "" {
		return ""
	}
	return quote(f.gorm)
}

func (f *FieldModel) RawGorm() string {
	return f.gorm
}

func (f *FieldModel) RawDb() string {
	return f.db
}

func (f *FieldModel) ES() string {
	return f.es
}

func (f *FieldModel) IsJSONField() bool {
	return f.inObj
}

func quote(s string) string {
	return fmt.Sprintf("`%s`", s)
}
