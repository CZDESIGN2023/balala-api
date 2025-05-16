package group

import (
	"cmp"
	"slices"
	"sort"
)

type Group[T any] struct {
	FieldName string
	Key       string
	Values    []T

	NextNodes      []*Group[T]
	prevLevelNodes []*Group[T]
	curLevelNodes  []*Group[T]
}

func New[T any](values []T) *Group[T] {
	g := &Group[T]{
		Values: values,
	}
	g.curLevelNodes = []*Group[T]{g}
	return g
}

func (g *Group[T]) GroupBy(fieldName string, fn func(T) string) {
	f := func(t T) []string {
		return []string{fn(t)}
	}

	g.GroupByToMulti(fieldName, f)
}

// GroupByToMulti 将一个对象分配到多个组
func (g *Group[T]) GroupByToMulti(fieldName string, fn func(T) []string) {
	var nextLevelGroups []*Group[T]

	for _, group := range g.curLevelNodes {
		nextLevelGroups = append(nextLevelGroups, GroupBy(group, fn, fieldName)...)
	}
	g.prevLevelNodes = g.curLevelNodes
	g.curLevelNodes = nextLevelGroups
}

func (g *Group[T]) Sort(fn func(ig, jg *Group[T]) bool) {
	if g.prevLevelNodes == nil {
		return
	}

	for _, group := range g.prevLevelNodes {
		ng := group.NextNodes
		sort.Slice(ng, func(i, j int) bool {
			v := fn(ng[i], ng[j])
			return v
		})
	}
}

func GroupBy[T any](g *Group[T], fn func(u T) []string, fieldName string) []*Group[T] {
	var m = map[string][]T{}
	var rankMap = map[string]int{}
	for i := 0; i < len(g.Values); i++ {
		for _, k := range fn(g.Values[i]) {
			if _, ok := rankMap[k]; !ok {
				rankMap[k] = i
			}
			m[k] = append(m[k], g.Values[i])
		}
	}

	var ret []*Group[T]
	for k, values := range m {
		ret = append(ret, &Group[T]{
			FieldName: fieldName,
			Key:       k,
			Values:    values,
		})
	}

	// 保证最终顺序
	slices.SortFunc(ret, func(a, b *Group[T]) int {
		return cmp.Compare(rankMap[a.Key], rankMap[b.Key])
	})

	g.NextNodes = ret

	return ret
}
