package stream

import "testing"

func Test(t *testing.T) {
	s := Of([]int{1, 2, 3, 4, 5})

	s.Diff(1).Diff(1, 2).Concat(10, 11).Concat(1, 1, 1).Unique()

	t.Log(s.List())
}
