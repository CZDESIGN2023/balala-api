package stream

import (
	"go-cs/pkg/pprint"
	"slices"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	strArr := []string{"1", "2"}
	Map(strArr, func(v string) int {
		i, _ := strconv.Atoi(v)
		return i
	})

	t.Log(strArr)
}

func TestFlat(t *testing.T) {
	l1 := Map([]string{"1", "2"}, func(v string) []int64 {
		return []int64{1, 2}
	})

	t.Log(l1)

	flat := Flat(l1)

	t.Log(flat)
}

func TestUnique(t *testing.T) {
	list := Unique([]int64{1, 2, 3, 1, 1, 2})
	pprint.Println(list)
}

func TestConcat(t *testing.T) {
	a1 := []int{1, 2}
	a2 := []int{3, 4}
	a3 := []int{4, 5}

	merge := Concat(a1, a2, a3)
	t.Log(merge)
}

func TestCompact(t *testing.T) {
	compact := slices.Compact([]int{1, 3, 2, 3, 3})
	t.Log(compact)
}

func TestContainsArr(t *testing.T) {
	t.Log(ContainsArr([]int{1, 2, 3}, []int{3, 2, 4}))
}

func TestUnion(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 2, 3, 4}
	t.Log(Union(a, b))
	t.Log(Diff(a, b))
	t.Log(Diff(b, a))
}

func TestMove(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	t.Log(a)
	Move(a, 3, 1)
	t.Log(a)
}
