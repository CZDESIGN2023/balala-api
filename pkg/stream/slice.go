package stream

import (
	"cmp"
	"slices"
)

func GroupBy[K comparable, V any](arr []V, fn func(v V) K) map[K][]V {
	var m = map[K][]V{}
	for _, e := range arr {
		k := fn(e)
		m[k] = append(m[k], e)
	}

	return m
}

func Map[T any, V any](list []T, fn func(v T) V) []V {
	ret := make([]V, 0, len(list))
	for _, v := range list {
		ret = append(ret, fn(v))
	}

	return ret
}

// Flat 压平
func Flat[S []T, T any](list []S) S {
	var size int
	for _, v := range list {
		size += len(v)
	}

	ret := make([]T, 0, size)
	for _, v := range list {
		ret = append(ret, v...)
	}

	return ret
}

// FlatMap 对 Map 结果压平
func FlatMap[T, E any, V []E](list []T, fn func(v T) V) []E {
	return Flat(Map(list, fn))
}

// Concat 合并
func Concat[S []T, T any](list ...S) S {
	var size int
	for _, v := range list {
		size += len(v)
	}

	ret := make([]T, 0, size)

	for _, v := range list {
		ret = append(ret, v...)
	}

	return ret
}

// Unique 去重
func Unique[T comparable](list []T) []T {
	ret := make([]T, 0, len(list))

	for _, v := range list {
		if Contains(ret, v) {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}

// UniqueFunc 去重
func UniqueFunc[T, K comparable](list []T, f func(T) K) []T {
	ret := make([]T, 0, len(list))
	m := make(map[K]struct{}, len(list))

	for _, v := range list {
		k := f(v)
		if _, ok := m[k]; !ok {
			m[k] = struct{}{}
			ret = append(ret, v)
		}
	}
	return ret
}

func Filter[T any](list []T, fn func(e T) bool) []T {
	var ret []T
	for _, v := range list {
		if fn(v) {
			ret = append(ret, v)
		}
	}

	return ret
}

func FilterOne[T any](list []T, fn func(e T) bool) T {
	var t T
	for _, v := range list {
		if fn(v) {
			return v
		}
	}
	return t
}

func Each[T any](list []T, fn func(v T)) {
	for _, t := range list {
		fn(t)
	}
}

func Sum[T cmp.Ordered](list []T) T {
	var ret T
	for _, v := range list {
		ret += v
	}

	return ret
}

func Contains[T comparable](list []T, in T) bool {
	for _, v := range list {
		if v == in {
			return true
		}
	}

	return false
}

func ContainsFunc[T any](list []T, fn func(T) bool) bool {
	for _, v := range list {
		if fn(v) {
			return true
		}
	}

	return false
}

func ToMap[K comparable, T, V any](arr []T, fn func(int, T) (K, V)) map[K]V {
	m := make(map[K]V, len(arr))

	for i, v := range arr {
		k, v := fn(i, v)
		m[k] = v
	}
	return m
}

func ContainsArr[T comparable](arr1, arr2 []T) bool {
	if len(arr2) > len(arr1) {
		return false
	}

	m1 := ToMap(arr1, func(i int, t T) (T, struct{}) {
		return t, struct{}{}
	})

	for _, t := range arr2 {
		if _, ok := m1[t]; !ok {
			return false
		}
	}
	return true
}

// All 所有成员满足 f 条件
func All[T any](arr []T, f func(t T) bool) bool {
	for _, t := range arr {
		if !f(t) {
			return false
		}
	}

	return true
}

// Some 至少一个成员满足 f 条件
func Some[T any](arr []T, f func(t T) bool) bool {
	for _, t := range arr {
		if f(t) {
			return true
		}
	}

	return false
}

func Remove[E comparable](a []E, e E) []E {
	index := slices.Index(a, e)

	if index < 0 {
		return a
	}

	return slices.Delete(a, index, index+1)
}

// Union 定义一个泛型函数，用于实现两个集合的并集
func Union[E comparable](a, b []E) []E {
	m := make(map[E]struct{})
	for _, item := range a {
		m[item] = struct{}{}
	}
	for _, item := range b {
		m[item] = struct{}{}
	}
	result := make([]E, 0, len(m))
	for item := range m {
		result = append(result, item)
	}
	return result
}

// Intersect 定义一个泛型函数，用于实现两个集合的交集
func Intersect[E comparable](a, b []E) []E {
	m := make(map[E]bool)
	result := make([]E, 0)
	for _, item := range a {
		if Contains(b, item) {
			m[item] = true
		}
	}
	for item := range m {
		result = append(result, item)
	}
	return result
}

// Diff 定义一个泛型函数，用于实现两个集合的差集
func Diff[E comparable](a, b []E) []E {
	m := make(map[E]bool)
	result := make([]E, 0)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, found := m[item]; !found {
			result = append(result, item)
		}
	}
	return result
}

func ForEach[E any](list []E, fn func(idx int, item E)) {
	for i, item := range list {
		fn(i, item)
	}
}

func Zip[K comparable, V any](keys []K, values []V) map[K]V {
	var m = make(map[K]V, len(keys))
	for i, k := range keys {
		var v V
		if i < len(values) {
			v = values[i]
		}

		m[k] = v
	}

	return m
}

func Move[S []E, E comparable](s S, srcIdx, dstIdx int) {
	e := s[srcIdx]

	if dstIdx > srcIdx {
		copy(s[srcIdx:], s[srcIdx+1:dstIdx+1])
		s[dstIdx] = e
	} else {
		copy(s[dstIdx+1:], s[dstIdx:srcIdx])
		s[dstIdx] = e
	}
}

func Reduce[E cmp.Ordered](list []E, fn func(acc E, v E) E) E {
	var acc E
	for _, v := range list {
		acc = fn(acc, v)
	}

	return acc
}
