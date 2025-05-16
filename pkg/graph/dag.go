package graph

import (
	"cmp"
	"errors"
	"github.com/dominikbraun/graph"
	"go-cs/pkg/stream"
	"slices"
)

type DAG struct {
	g graph.Graph[string, string]

	startNode string // 开始节点
	endNode   string // 结束节点

	maxDepth    map[string]int      // 每个节点的最大深度
	depth2Nodes map[int][]string    // 深度对应的节点
	outgoing    map[string][]string // 每个节点的出边
	incoming    map[string][]string // 每个节点的入边

	allPathBetween map[string][][]string // 所有路径

	nodes []string // 需要记录原始节点到顺序
	edges [][]string
}

/*
NewDAG 创建一个有向无环图
@param nodes 节点列表 必须是有顺序的，从左到右，从上到下
*/
func NewDAG(nodes []string, edges [][]string) (*DAG, error) {
	// 创建一个空的有向无环图
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	outgoing := map[string][]string{}
	incoming := map[string][]string{}

	for _, node := range nodes {
		err := g.AddVertex(node)
		if err != nil {
			return nil, err
		}

		outgoing[node] = nil
		incoming[node] = nil
	}

	for _, edge := range edges {
		err := g.AddEdge(edge[0], edge[1], graph.EdgeWeight(1))
		if err != nil {
			return nil, err
		}

		outgoing[edge[0]] = append(outgoing[edge[0]], edge[1])
		incoming[edge[1]] = append(incoming[edge[1]], edge[0])
	}

	dag := &DAG{
		g:              g,
		nodes:          nodes,
		edges:          edges,
		maxDepth:       map[string]int{},
		allPathBetween: map[string][][]string{},
		outgoing:       outgoing,
		incoming:       incoming,
	}

	// 确定开始节点和结束节点
	noIncomingEdges := dag.GetNodesWithNoIncomingEdges()
	if len(noIncomingEdges) != 1 {
		return nil, errors.New("multiple start nodes")
	}
	noOutgoingEdges := dag.GetNodesWithNoOutgoingEdges()
	if len(noIncomingEdges) != 1 {
		return nil, errors.New("multiple end nodes")
	}

	dag.startNode = noOutgoingEdges[0]
	dag.endNode = noOutgoingEdges[0]

	// 计算最大深度
	dag.calMaxDepth()

	return dag, nil
}

// GetNodesWithNoIncomingEdges 获取没有入度的节点
func (dag *DAG) GetNodesWithNoIncomingEdges() []string {
	var nodes []string
	for node, incomingEdges := range dag.incoming {
		if len(incomingEdges) == 0 {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetNodesWithNoOutgoingEdges 获取没有出度的节点
func (dag *DAG) GetNodesWithNoOutgoingEdges() []string {
	var nodes []string
	for node, outgoingEdges := range dag.outgoing {
		if len(outgoingEdges) == 0 {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetOutgoingEdges 获取节点的所有出边
func (dag *DAG) GetOutgoingEdges(node string) []string {
	return append([]string(nil), dag.outgoing[node]...)
}

// GetIncomingEdges 获取节点的所有入边
func (dag *DAG) GetIncomingEdges(node string) []string {
	return append([]string(nil), dag.incoming[node]...)
}

// GetMaxDepth 获取每个节点的最大深度
func (dag *DAG) GetMaxDepth() map[string]int {
	cpy := map[string]int{}

	for k, v := range dag.maxDepth {
		cpy[k] = v
	}
	return cpy
}

// CheckNodeOrder 检查节点顺序是否正确
func (dag *DAG) CheckNodeOrder() bool {
	// 获取每个深度对应的节点
	depth2Nodes := stream.ToEntries(dag.GetDepth2Nodes())
	slices.SortFunc(depth2Nodes, func(a, b stream.Entry[int, []string]) int {
		return cmp.Compare(a.Key, b.Key)
	})

	// 检查节点顺序是否正确
	var nodesIdx = 0
	for _, e := range depth2Nodes {
		nodes := e.Val

		for _, node := range nodes {
			if !slices.Contains(dag.nodes[nodesIdx:nodesIdx+len(nodes)], node) {
				return false
			}
		}
		nodesIdx += len(nodes)
	}

	return true
}

// GetDepth2Nodes 获取每个深度对应的节点
func (dag *DAG) GetDepth2Nodes() map[int][]string {
	if dag.depth2Nodes == nil {
		entries := stream.ToEntries(dag.maxDepth)

		// 对节点进行排序
		nodeRank := stream.ToMap(dag.nodes, func(idx int, node string) (string, int) {
			return node, idx
		})

		slices.SortFunc(entries, func(a, b stream.Entry[string, int]) int {
			return cmp.Compare(nodeRank[a.Key], nodeRank[b.Key])
		})

		g := stream.GroupBy(entries, func(entry stream.Entry[string, int]) int {
			return entry.Val
		})

		dag.depth2Nodes = stream.MapValue(g, func(nodes []stream.Entry[string, int]) []string {
			return stream.Map(nodes, func(node stream.Entry[string, int]) string {
				return node.Key
			})
		})

	}

	cpy := map[int][]string{}
	for k, v := range dag.depth2Nodes {
		cpy[k] = append([]string(nil), v...)
	}
	return cpy
}

// calMaxDepth 从开始节点到结束节点计算出每个节点的最大深度
func (dag *DAG) calMaxDepth() {
	// 递归计算每个节点的最大深度，直到结束节点
	var calMaxDepthFromNode func(node string, depth int)
	calMaxDepthFromNode = func(node string, depth int) {
		outNodes := dag.GetOutgoingEdges(node)
		if len(outNodes) == 0 { //没有下一个节点，说明到达终点，停止递归
			return
		}

		for _, nextNode := range outNodes {
			if dag.maxDepth[nextNode] < depth+1 {
				dag.maxDepth[nextNode] = depth + 1
				calMaxDepthFromNode(nextNode, depth+1)
			}
		}
	}

	// 从开始节点开始计算
	startNode := dag.GetNodesWithNoIncomingEdges()
	for _, node := range startNode {
		dag.maxDepth[node] = 0
		calMaxDepthFromNode(node, 0)
	}
}

// GetAllPathBetween 获取从startNode到endNode的所有路径
func (dag *DAG) GetAllPathBetween(startNode, endNode string) [][]string {
	if paths, ok := dag.allPathBetween[startNode+"-"+endNode]; ok {
		return paths
	}

	between, _ := graph.AllPathsBetween(dag.g, startNode, endNode)
	dag.allPathBetween[startNode+"-"+endNode] = between
	return append([][]string(nil), between...)
}
