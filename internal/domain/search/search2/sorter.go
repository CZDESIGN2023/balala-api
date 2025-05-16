package search2

import (
	"slices"
)

// MergeCmpFunc 合并比较函数，实现 order by a, b, c
func MergeCmpFunc[T any](fns ...func(T, T) int) func(a, b T) int {
	return func(a, b T) int {
		for _, fn := range fns {
			if v := fn(a, b); v != 0 {
				return v
			}
		}
		return 0
	}
}

// ReverseCmpFunc 反转顺序
func ReverseCmpFunc[T any](fn func(a, b T) int) func(a, b T) int {
	return func(a, b T) int {
		return fn(b, a)
	}
}

// Sort 排序
func Sort[T any](list []T, fn ...func(a, b T) int) {
	slices.SortStableFunc(list, MergeCmpFunc(fn...))
}
