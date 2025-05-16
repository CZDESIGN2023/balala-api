package search_es

import (
	"fmt"
)

var fieldModelList []FieldModel

func FieldByQuery(query string) FieldModel {
	return query2FieldModelMap[QueryField(query)]
}

type FieldModel struct {
	query     string //外部字段名
	es        string
	esKeyword string     //es关键字
	dt        DataType   //数据类型
	esDt      EsDataType //es对应数据类型
}

func NewFieldModel(query string, es string, dt DataType) *FieldModel {
	m := &FieldModel{
		query: query,
		es:    es,
		dt:    dt,
	}
	return m
}

func (f *FieldModel) Query() string {
	return f.query
}

func (f *FieldModel) Dt() DataType {
	return f.dt
}

func (f *FieldModel) ES() string {
	return f.es
}

func (f *FieldModel) EsKeyword() string {
	if f.esKeyword != "" {
		return f.esKeyword
	}
	return f.es
}

func quote(s string) string {
	return fmt.Sprintf("`%s`", s)
}
