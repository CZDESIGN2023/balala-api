package search2

import "testing"

func TestSelectAll(t *testing.T) {
	all := SelectAll()
	t.Log(all)
}

func TestSelectByQuery(t *testing.T) {
	query := SelectByQuery("id")
	t.Log(query)
}
