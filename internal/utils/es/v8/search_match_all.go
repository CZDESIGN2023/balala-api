package v8

import "encoding/json"

type MatchAllQuery struct {
	boost     *float64
	queryName string
}

func NewMatchAllQuery() *MatchAllQuery {
	return &MatchAllQuery{}
}

func (q *MatchAllQuery) Boost(boost float64) *MatchAllQuery {
	q.boost = &boost
	return q
}

func (q *MatchAllQuery) QueryName(name string) *MatchAllQuery {
	q.queryName = name
	return q
}

func (q *MatchAllQuery) Source() (interface{}, error) {
	// {
	//   "match_all" : { ... }
	// }
	source := make(map[string]interface{})
	params := make(map[string]interface{})
	source["match_all"] = params
	if q.boost != nil {
		params["boost"] = *q.boost
	}
	if q.queryName != "" {
		params["_name"] = q.queryName
	}
	return source, nil
}

func (q *MatchAllQuery) MarshalJSON() ([]byte, error) {
	if q == nil {
		return nil, nil
	}
	src, err := q.Source()
	if err != nil {
		return nil, err
	}
	return json.Marshal(src)
}
