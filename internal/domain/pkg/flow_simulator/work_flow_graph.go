package flow_simulator

import (
	"cmp"
	"github.com/spf13/cast"
	"go-cs/internal/consts"
	"go-cs/internal/domain/work_flow"
	flow_config "go-cs/internal/domain/work_flow/flow_tplt_config"
	domain "go-cs/internal/domain/work_item"
	"go-cs/pkg/graph"
	"go-cs/pkg/stream"
	"slices"
)

type WorkFlowGraph struct {
	g        *graph.DAG
	workItem *domain.WorkItem
	template *work_flow.WorkFlowTemplate

	flowNodeMap map[string]*domain.WorkItemFlowNode
}

func NewWorkFlowGraph(workItem *domain.WorkItem, template *work_flow.WorkFlowTemplate) (*WorkFlowGraph, error) {
	nodes := stream.Map(workItem.WorkItemFlowNodes, func(v *domain.WorkItemFlowNode) string {
		return v.FlowNodeCode
	})
	edges := stream.Map(template.WorkFLowConfig.Connections, func(v *flow_config.WorkFlowConnection) []string {
		return []string{v.StartNode, v.EndNode}
	})

	dag, err := graph.NewDAG(nodes, edges)
	if err != nil {
		return nil, err
	}

	flowNodeMap := stream.ToMap(workItem.WorkItemFlowNodes, func(_ int, v *domain.WorkItemFlowNode) (string, *domain.WorkItemFlowNode) {
		return v.FlowNodeCode, v
	})

	return &WorkFlowGraph{
		g:           dag,
		workItem:    workItem,
		template:    template,
		flowNodeMap: flowNodeMap,
	}, nil
}

func (g *WorkFlowGraph) GetNextNodesInfo(nodeCode string) []*domain.WorkItemFlowNode {
	nextNodes := g.g.GetOutgoingEdges(nodeCode)
	return stream.Map(nextNodes, func(v string) *domain.WorkItemFlowNode {
		return g.flowNodeMap[v]
	})
}

func (g *WorkFlowGraph) GetPrevNodesInfo(nodeCode string) []*domain.WorkItemFlowNode {
	prevNodes := g.g.GetIncomingEdges(nodeCode)
	return stream.Map(prevNodes, func(v string) *domain.WorkItemFlowNode {
		return g.flowNodeMap[v]
	})
}

func (g *WorkFlowGraph) GetNodeConfig(nodeCode string) *flow_config.WorkFlowNode {
	return g.template.WorkFLowConfig.GetNode(nodeCode)
}

func (g *WorkFlowGraph) GetNodeInfo(nodeCode string) *domain.WorkItemFlowNode {
	return g.flowNodeMap[nodeCode]
}

// ConfirmNode 确认节点
func (g *WorkFlowGraph) ConfirmNode(nodeCode string) (affectedNodes []*domain.WorkItemFlowNode) {
	nodeInfo := g.flowNodeMap[nodeCode]

	if nodeInfo.IsCompleted() {
		return
	}

	nodeInfo.UpdateStateOnNodePass()

	// 处理后继节点
	nextNodesInfo := g.GetNextNodesInfo(nodeCode)
	for i := 0; i < len(nextNodesInfo); i++ {
		node := nextNodesInfo[i]
		if node.IsCompleted() {
			continue
		}

		nodeConfig := g.GetNodeConfig(node.FlowNodeCode)
		switch nodeConfig.StartMode {
		case flow_config.WorkflowNodeStartMode_PreAllDone:
			// 判断所有前置节点是否完成
			allDone := stream.All(g.GetPrevNodesInfo(node.FlowNodeCode), func(v *domain.WorkItemFlowNode) bool {
				return v.IsCompleted()
			})

			// 如果所有前置节点都完成，则设置当前节点为进行中状态
			if allDone {
				node.UpdateStateOnNodeReach()
				affectedNodes = append(affectedNodes, node)
			}

			// 如果当前节点处于进行中状态，且节点配置为自动通过，则设置当前节点为完成状态
			if node.IsInProcess() && nodeConfig.IsAutoPass() {
				node.UpdateStateOnNodePass()

				// 添加当前节点的后继节点到待处理列表中
				nextNodesInfo = append(nextNodesInfo, g.GetNextNodesInfo(node.FlowNodeCode)...)
			}
		}
	}

	return
}

// ReCalculateWorkItemStatus 重新计算工作项状态
func (g *WorkFlowGraph) ReCalculateWorkItemStatus() {
	nodeCode := g.FindRightFirstInProcessNode()
	if nodeCode == "" {
		// 如果没有进行中的节点，则说明所有节点都完成了
		nodeCode = flow_config.WorkflowNodeCode_Ended
	}

	nodeConfig := g.GetNodeConfig(nodeCode)

	var status domain.WorkItemStatus
	for _, v := range nodeConfig.OnReach {
		status = domain.WorkItemStatus{
			Val: v.TargetSubState.Val,
			Key: v.TargetSubState.Key,
			Id:  cast.ToInt64(v.TargetSubState.Id),
		}
		_ = g.workItem.UpdateStatus(status)

	}
}

// 找到第一个未完成的节点
func (g *WorkFlowGraph) FindRightFirstInProcessNode() string {
	return g.findFirstInProcessNode("right")
}

func (g *WorkFlowGraph) FindLeftFirstInProcessNode() string {
	return g.findFirstInProcessNode("left")
}

func (g *WorkFlowGraph) findFirstInProcessNode(direction string) string {
	// 按深度遍历
	depth2Nodes := stream.ToEntries(g.g.GetDepth2Nodes())
	slices.SortFunc(depth2Nodes, func(a, b stream.Entry[int, []string]) int {
		return cmp.Compare(a.Key, b.Key)
	})

	iter := slices.All(depth2Nodes)
	if direction == "right" {
		iter = slices.Backward(depth2Nodes)
	}

	for _, depth2Node := range iter {
		for _, nodeCode := range depth2Node.Val {
			if g.GetNodeInfo(nodeCode).IsInProcess() {
				return nodeCode
			}
		}
	}

	return ""
}

// 重启至节点
func (g *WorkFlowGraph) RebootToNode(nodeCode string) {
	nodeInfo := g.GetNodeInfo(nodeCode)
	nodeInfo.ResetProgressStatus()

	g.reCalcNodeStatusChain(nodeCode)
}

// 重置节点链路状态
func (g *WorkFlowGraph) ReCalcAllNodeStatus() {
	// 如果工作项状态为归档状态，则将所有节点状态重置为完成状态
	if g.workItem.WorkItemStatus.Key == string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey) {
		for _, node := range g.workItem.WorkItemFlowNodes {
			node.UpdateStateOnNodePass()
		}

		return
	}

	g.reCalcNodeStatusChain(flow_config.WorkflowNodeCode_Started)
}

// 重新计算节点链路状态，将节点的错误状态修正
// startNodeCode 节点为起始节点，从该节点开始，依次遍历所有节点，如果节点处于错误状态，则修正节点状态
// 起始节点的状态是正确的
func (g *WorkFlowGraph) reCalcNodeStatusChain(startNodeCode string) {
	nextNodesInfo := g.GetNextNodesInfo(startNodeCode)

	visited := map[string]struct{}{}

	for i := 0; i < len(nextNodesInfo); i++ {
		node := nextNodesInfo[i]
		if _, ok := visited[node.FlowNodeCode]; ok {
			continue
		}
		visited[node.FlowNodeCode] = struct{}{}

		nodeConfig := g.GetNodeConfig(node.FlowNodeCode)

		switch nodeConfig.StartMode {
		case flow_config.WorkflowNodeStartMode_PreAllDone:
			// 判断所有前置节点是否完成
			allDone := stream.All(g.GetPrevNodesInfo(node.FlowNodeCode), func(v *domain.WorkItemFlowNode) bool {
				return v.IsCompleted()
			})

			if allDone { // 如果所有前置节点都完成，则设置当前节点为进行中状态
				if node.IsUnstart() {
					node.ResetProgressStatus()
				}
			} else { // 如果所有前置节点未完成，则重置当前节点状态
				node.ResetStatus()
			}

			// 如果当前节点处于进行中状态，且节点配置为自动通过，则设置当前节点为完成状态
			if node.IsInProcess() && nodeConfig.IsAutoPass() {
				node.UpdateStateOnNodePass()
			}

			// 添加当前节点的后继节点到待处理列表中
			nextNodesInfo = append(nextNodesInfo, g.GetNextNodesInfo(node.FlowNodeCode)...)
		}
	}
}
