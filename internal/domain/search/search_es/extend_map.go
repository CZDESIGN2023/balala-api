package search_es

import "go-cs/pkg/stream"

type ExtendMap map[string]any

func (e ExtendMap) Int64(k string) int64 {
	return e[k].(int64)
}

func (e ExtendMap) Float64(k string) float64 {
	return e[k].(float64)
}

func (e ExtendMap) String(k string) string {
	return e[k].(string)
}

func (e ExtendMap) Int64s(k string) []int64 {
	return stream.Map(e[k].([]any), func(v any) int64 {
		return v.(int64)
	})
}

func (e ExtendMap) Float64s(k string) []float64 {
	return stream.Map(e[k].([]any), func(v any) float64 {
		return v.(float64)
	})
}

func (e ExtendMap) Strings(k string) []string {
	return stream.Map(e[k].([]any), func(v any) string {
		return v.(string)
	})
}
