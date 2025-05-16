package search_es

import (
	"fmt"
	esV8 "go-cs/internal/utils/es/v8"

	"github.com/spf13/cast"
)

type SearchCondExp interface {
	In() esV8.Query
	NotIn() esV8.Query
	Eq() esV8.Query
	NotEq() esV8.Query
	Gt() esV8.Query
	Lt() esV8.Query
	Gte() esV8.Query
	Lte() esV8.Query
	Include() esV8.Query
	Exclude() esV8.Query
	Between() esV8.Query
}

func NewSearchCondExp(model FieldModel, values ...interface{}) SearchCondExp {
	if model.esDt == EsArrayInt || model.esDt == EsArrayString {
		return &ArraySearchCondExp{model: model, values: values}
	}
	return &NormalSearchCondExp{model: model, values: values}
}

type NormalSearchCondExp struct {
	model  FieldModel
	values []interface{}
}

func (exp *NormalSearchCondExp) In() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}
	clause := esV8.NewTermsQuery(exp.model.EsKeyword(), exp.values...)
	return clause
}

func (exp *NormalSearchCondExp) NotIn() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}

	clause := esV8.NewBoolQuery()
	for _, v := range exp.values {
		clause.MustNot(esV8.NewTermQuery(exp.model.EsKeyword(), v))
	}
	return clause
}

func (exp *NormalSearchCondExp) Eq() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}

	if len(exp.values) == 1 {
		return esV8.NewTermQuery(exp.model.EsKeyword(), exp.values[0])
	}

	boolQuery := esV8.NewBoolQuery()
	for _, v := range exp.values {
		boolQuery.Must(esV8.NewTermQuery(exp.model.EsKeyword(), v))
	}
	return boolQuery
}

func (exp *NormalSearchCondExp) NotEq() esV8.Query {

	if len(exp.values) == 0 {
		return nil
	}

	boolQuery := esV8.NewBoolQuery()
	for _, v := range exp.values {
		boolQuery.MustNot(esV8.NewTermQuery(exp.model.EsKeyword(), v))
	}
	return boolQuery
}

func (exp *NormalSearchCondExp) Gt() esV8.Query {
	clause := esV8.NewRangeQuery(exp.model.EsKeyword())
	clause.Gt(exp.values[0])
	return clause
}

func (exp *NormalSearchCondExp) Lt() esV8.Query {
	clause := esV8.NewRangeQuery(exp.model.EsKeyword())
	clause.Lt(exp.values[0])
	return clause
}

func (exp *NormalSearchCondExp) Gte() esV8.Query {
	clause := esV8.NewRangeQuery(exp.model.EsKeyword())
	clause.Gte(exp.values[0])
	return clause
}

func (exp *NormalSearchCondExp) Lte() esV8.Query {
	clause := esV8.NewRangeQuery(exp.model.EsKeyword())
	clause.Lte(exp.values[0])
	return clause
}
func (exp *NormalSearchCondExp) Include() esV8.Query {
	clause := esV8.NewWildcardQuery(exp.model.EsKeyword(), "*"+cast.ToString(exp.values[0])+"*").
		CaseInsensitive(true)
	return clause
}

func (exp *NormalSearchCondExp) Exclude() esV8.Query {
	clause := esV8.NewWildcardQuery(exp.model.EsKeyword(), "NOT *"+cast.ToString(exp.values[0])+"*")
	return clause
}

func (exp *NormalSearchCondExp) Between() esV8.Query {
	clause := esV8.NewRangeQuery(exp.model.EsKeyword())
	clause.Gte(exp.values[0])
	clause.Lte(exp.values[1])
	return clause
}

// 数组类型
type ArraySearchCondExp struct {
	model  FieldModel
	values []interface{}
}

func (exp *ArraySearchCondExp) In() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}
	clause := esV8.NewTermsQuery(exp.model.EsKeyword(), exp.values...)
	return clause
}

func (exp *ArraySearchCondExp) NotIn() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}

	condition := esV8.NewTermsQuery(exp.model.EsKeyword(), exp.values...)
	clause := esV8.NewBoolQuery().MustNot(condition)
	return clause
}

func (exp *ArraySearchCondExp) Eq() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}

	//如果字段是数组类的
	script := fmt.Sprintf(`if (doc['%s'].size() == params.vals.size() && doc['%s'].containsAll(params.vals)) { return true } else { return false }`,
		exp.model.EsKeyword(), exp.model.EsKeyword())
	clause := esV8.NewScriptWrap(esV8.NewScript(script).Lang("painless").Params(map[string]interface{}{"vals": exp.values}))
	return clause
}

func (exp *ArraySearchCondExp) NotEq() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}

	//如果字段是数组类的
	script := fmt.Sprintf(`if (doc['%s'].size() == params.vals.size() && doc['%s'].containsAll(params.vals)) { return false } else { return true }`,
		exp.model.EsKeyword(), exp.model.EsKeyword())
	clause := esV8.NewScriptWrap(esV8.NewScript(script).Lang("painless").Params(map[string]interface{}{"vals": exp.values}))
	return clause
}

func (exp *ArraySearchCondExp) Gt() esV8.Query {
	return nil
}

func (exp *ArraySearchCondExp) Lt() esV8.Query {
	return nil
}

func (exp *ArraySearchCondExp) Gte() esV8.Query {
	return nil
}

func (exp *ArraySearchCondExp) Lte() esV8.Query {
	return nil
}
func (exp *ArraySearchCondExp) Include() esV8.Query {
	if len(exp.values) == 0 {
		return nil
	}
	//如果字段是数组类的
	script := fmt.Sprintf(`if (doc['%s'].containsAll(params.vals)) { return true } else { return false }`,
		exp.model.EsKeyword())
	clause := esV8.NewScriptWrap(esV8.NewScript(script).Lang("painless").Params(map[string]interface{}{"vals": exp.values}))
	return clause
}

func (exp *ArraySearchCondExp) Exclude() esV8.Query {
	//如果字段是数组类的
	clause := esV8.NewBoolQuery()
	for _, v := range exp.values {
		clause.MustNot(esV8.NewTermQuery(exp.model.EsKeyword(), v))
	}

	return clause
}

func (exp *ArraySearchCondExp) Between() esV8.Query {
	return nil
}
