package stream

import (
	"go-cs/pkg/pprint"
	"testing"
)

func TestGroupBy(t *testing.T) {
	v := GroupBy([]int64{1, 3, 4, 5, 6}, func(v int64) int64 {
		return v % 2
	})

	t.Log(v)
}

func TestMapMap(t *testing.T) {
	m := map[int64]string{
		1: "1",
		2: "2",
		3: "3",
	}

	m2 := MapValue(m, func(v string) string {
		return "%" + v
	})

	t.Log(m2)
}

func TestKeys(t *testing.T) {
	m := map[int64]string{
		1: "A",
		2: "B",
		3: "C",
	}

	t.Log(Keys(m))
	t.Log(Values(m))
}

func TestJoin(t *testing.T) {
	m1 := map[string][]int{
		"A": {1, 2, 3},
		"B": {},
		"C": {},
	}

	m2 := map[int][]string{
		1: {"I"},
		2: {"II"},
		3: {"III"},
	}

	t.Log(Join1NAnd1N(m1, m2))
}

func TestToEntries(t *testing.T) {
	m := map[int64]string{
		1: "1",
		2: "2",
		3: "3",
	}

	entries := ToEntries(m)

	pprint.Println(entries)
}
