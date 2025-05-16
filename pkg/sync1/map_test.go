package sync1

import "testing"

func Test(t *testing.T) {
	m := Map[int64, string]{}
	m.Store(1, "sdrf")

	//t.Log(m.Load(1))

	m.Range(func(k int64, v string) {
		t.Log(k, v)
	})
}
