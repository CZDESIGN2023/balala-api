package stream

import "go-cs/pkg/stream/tuple"

func MapValue[K comparable, V, T any](m map[K]V, fn func(V) T) map[K]T {
	ms := map[K]T{}
	for k, v := range m {
		ms[k] = fn(v)
	}
	return ms
}

func MapValue2[K comparable, V, T any](m map[K]V, fn func(V) (T, error)) map[K]T {
	ms := map[K]T{}
	for k, v := range m {
		fV, err := fn(v)
		if err != nil {
			ms[k] = fV
		}
	}
	return ms
}

func MapKey[K, K2 comparable, V any](m map[K]V, fn func(K) K2) map[K2]V {
	ms := map[K2]V{}
	for k, v := range m {
		ms[fn(k)] = v
	}

	return ms
}

func MapKV[K, K2 comparable, V, V2 any](m map[K]V, fn func(K, V) (K2, V2)) map[K2]V2 {
	ms := map[K2]V2{}
	for k, v := range m {
		k2, v2 := fn(k, v)
		ms[k2] = v2
	}

	return ms
}

func FilterByKey[K comparable, V any](m map[K]V, fn func(K) bool) map[K]V {
	ms := map[K]V{}
	for k, v := range m {
		if fn(k) {
			ms[k] = v
		}
	}

	return ms
}

func FilterByValue[K comparable, V any](m map[K]V, fn func(V) bool) map[K]V {
	ms := map[K]V{}
	for k, v := range m {
		if fn(v) {
			ms[k] = v
		}
	}

	return ms
}

func ContainsKey[K comparable, V any](m map[K]V, fn func(K) bool) bool {
	for k, _ := range m {
		if fn(k) {
			return true
		}
	}

	return false
}

func GroupByKey[K, K2 comparable, V any](m map[K]V, fn func(K) K2) map[K2][]V {
	ms := map[K2][]V{}
	for k, v := range m {
		ms[fn(k)] = append(ms[fn(k)], v)
	}

	return ms
}

func FlatValue[K comparable, V any](m map[K][][]V) map[K][]V {
	ret := make(map[K][]V, len(m))
	for k, v := range m {
		ret[k] = append(ret[k], Flat(v)...)
	}

	return ret
}

func Values[K comparable, V any](m map[K]V) []V {
	var arr []V
	for _, v := range m {
		arr = append(arr, v)
	}

	return arr
}

func Keys[K comparable, V any](m map[K]V) []K {
	var arr []K
	for k := range m {
		arr = append(arr, k)
	}

	return arr
}

func PickBy[K comparable, V any](in map[K]V, fn func(K, V) bool) map[K]V {
	var ret = make(map[K]V)

	for k, v := range in {
		if fn(k, v) {
			ret[k] = v
		}
	}

	return ret
}

func ReverseMap[K comparable, V any, T comparable](m map[K]V, fn func(V) T) map[T]K {
	ms := map[T]K{}
	for k, v := range m {
		ms[fn(v)] = k
	}

	return ms
}

// m1 1:N
// m2 1:N
// left join
func Join1NAnd1N[K comparable, V comparable, T any](m1 map[K][]V, m2 map[V][]T) map[K][]T {
	ms := map[K][]T{}
	for k, v := range m1 {
		for _, v2 := range v {
			ms[k] = append(ms[k], m2[v2]...)
		}
	}

	return ms
}

func Join[T1, T2 any](s1 []T1, s2 []T2, fn func(T1, T2) bool) []tuple.Pair[T1, T2] {
	var ret []tuple.Pair[T1, T2]

	for _, t1 := range s1 {
		for _, t2 := range s2 {
			if fn(t1, t2) {
				ret = append(ret, tuple.Pair[T1, T2]{
					First:  t1,
					Second: t2,
				})
			}
		}
	}

	return ret
}

func InnerJoin2Map[T1, T2 comparable](s1 []T1, s2 []T2, fn func(T1, T2) bool) map[T1]T2 {
	var ret = make(map[T1]T2)

	for _, t1 := range s1 {
		for _, t2 := range s2 {
			if fn(t1, t2) {
				ret[t1] = t2
			}
		}
	}

	return ret
}

func LeftJoin2Map[T1, T2 comparable](s1 []T1, s2 []T2, fn func(T1, T2) bool) map[T1]T2 {
	var ret = make(map[T1]T2)

	for _, t1 := range s1 {
		var zero T2
		for _, t2 := range s2 {
			if fn(t1, t2) {
				ret[t1] = t2
			}
		}
		if _, ok := ret[t1]; !ok {
			ret[t1] = zero
		}
	}

	return ret
}

func ToSlice[K comparable, V any, T any](m map[K]V, f func(K, V) T) []T {
	list := make([]T, 0, len(m))
	for k, v := range m {
		list = append(list, f(k, v))
	}

	return list
}

type Entry[K any, V any] struct {
	Key K `json:"key"`
	Val V `json:"val"`
}

func ToEntries[K comparable, V any](m map[K]V) []Entry[K, V] {
	return ToSlice(m, func(k K, v V) Entry[K, V] {
		return Entry[K, V]{
			Key: k,
			Val: v,
		}
	})
}

func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	return ToMap(entries, func(_ int, e Entry[K, V]) (K, V) {
		return e.Key, e.Val
	})
}

// MergeMap 合并map，覆盖相同key的值
func MergeMap[K comparable, V any](ms ...map[K]V) map[K]V {
	ret := make(map[K]V)

	for _, m := range ms {
		for k, v := range m {
			ret[k] = v
		}
	}

	return ret
}
