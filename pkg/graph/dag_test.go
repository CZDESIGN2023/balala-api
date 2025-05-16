package graph

import "testing"

func TestNewDAG(t *testing.T) {
	nodes := []string{"1", "3", "2", "5", "6", "4"}
	edges := [][]string{
		{"1", "2"},
		{"1", "3"},
		{"2", "4"},
		{"3", "4"},
		{"3", "5"},
		{"5", "6"},
		{"6", "4"},
	}
	dag, err := NewDAG(nodes, edges)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(dag.GetAllPathBetween("1", "4"))
	t.Log(dag.GetMaxDepth())
	t.Log(dag.GetDepth2Nodes())
	t.Log(dag.CheckNodeOrder())
}

func TestDAG_GetAllPathBetween(t *testing.T) {
}
