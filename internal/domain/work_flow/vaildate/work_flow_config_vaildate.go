package vaildate

import (
	"context"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/graph"
	"go-cs/pkg/stream"
)

// --各种检查配置是否正确的
type WorkFlowConfigValidate struct {
	conf *config.WorkFlow
	ctx  context.Context
}

func NewWorkFlowConfigValidate(conf *config.WorkFlow) *WorkFlowConfigValidate {
	return &WorkFlowConfigValidate{
		conf: conf,
	}
}

func (w *WorkFlowConfigValidate) WithCtx(ctx context.Context) *WorkFlowConfigValidate {
	w.ctx = ctx
	return w
}

func (w *WorkFlowConfigValidate) Ctx() context.Context {
	if w.ctx == nil {
		w.ctx = context.Background()
	}
	return w.ctx
}

func (w *WorkFlowConfigValidate) Valid() error {

	if err := w.validConnections(); err != nil {
		return err
	}

	if err := w.validNodes(); err != nil {
		return err
	}

	return nil
}

func (w *WorkFlowConfigValidate) validConnections() error {
	nodes := stream.Map(w.conf.Nodes, func(v *config.WorkFlowNode) string {
		return v.Key
	})

	edges := stream.Map(w.conf.Connections, func(v *config.WorkFlowConnection) []string {
		return []string{v.StartNode, v.EndNode}
	})

	dag, err := graph.NewDAG(nodes, edges)
	if err != nil {
		return err
	}

	startNode := dag.GetNodesWithNoIncomingEdges()
	endNode := dag.GetNodesWithNoOutgoingEdges()

	if startNode[0] != config.WorkflowNodeCode_Started {
		return errs.Business(w.Ctx(), "流程节点配置错误-未设置开始节点")
	}

	if endNode[0] != config.WorkflowNodeCode_Ended {
		return errs.Business(w.Ctx(), "流程节点配置错误-未设置结束节点")
	}

	if !dag.CheckNodeOrder() {
		return errs.Business(w.Ctx(), "流程节点配置错误-节点顺序错误")
	}

	return nil
}

func (w *WorkFlowConfigValidate) validNodes() error {
	nodeLen := len(w.conf.Nodes)
	if nodeLen < 3 {
		return errs.Business(w.Ctx(), "流程节点数量不足")
	}

	if nodeLen > 40 {
		return errs.Business(w.Ctx(), "流程节点已达最大数量")
	}

	if w.conf.GetNode(config.WorkflowNodeCode_Started) == nil {
		return errs.Business(w.Ctx(), "未配置开始节点")
	}

	if w.conf.GetNode(config.WorkflowNodeCode_Ended) == nil {
		return errs.Business(w.Ctx(), "未配置结束节点")
	}

	//检查节点的基本配置
	for _, v := range w.conf.Nodes {
		if v.Name == "" || v.Key == "" {
			return errs.Business(w.Ctx(), "节点配置名称错误")
		}

		if v.StartMode != config.WorkflowNodeStartMode_PreAllDone {
			return errs.Business(w.Ctx(), "节点配置启动模式错误")
		}

		if v.BelongStatus != config.WorkflowNodeCode_Started {
			return errs.Business(w.Ctx(), "节点配置状态错误")
		}

		if v.PassMode != config.WorkflowNodePassMode_Auto {
			if !v.NeedDoneOperator {
				return errs.Business(w.Ctx(), "节点配置是否需要操作人错误")
			}

			if v.DoneOperationDisplayName == "" {
				return errs.Business(w.Ctx(), "节点配置完成操作人显示名称错误")
			}

			err := w.ValidNodeOwnerConf(v.Owner)
			if err != nil {
				return err
			}
		}

		var nodeEvt []*config.WorkFlowNodeEvent
		nodeEvt = append(nodeEvt, v.OnReach...)
		nodeEvt = append(nodeEvt, v.OnPass...)
		for _, evt := range nodeEvt {
			if evt.TargetSubState != nil {
				if evt.TargetSubState.Id == "" || evt.TargetSubState.Key == "" {
					return errs.Business(w.Ctx(), "节点配置事件错误")
				}
			}
		}
	}

	return nil
}

func (w *WorkFlowConfigValidate) ValidNodeOwnerConf(conf *config.OwnerConf) error {
	if conf == nil {
		return errs.Business(w.Ctx(), "节点负责人配置错误-未设置关联角色")
	}

	if len(conf.OwnerRole) == 0 {
		return errs.Business(w.Ctx(), "节点负责人配置错误-未设置关联角色")
	}

	if !conf.ForceOwner {
		return errs.Business(w.Ctx(), "节点负责人配置错误-未设置强制关联角色")
	}

	switch conf.UsageMode {
	case config.UsageMode_None:
		return nil
	case config.UsageMode_Appointed:
		ruleVal := conf.GetUsageModeVal()
		rule, isOk := ruleVal.(*config.OwnerConf_UsageMode_Appointed)
		if !isOk {
			return errs.Business(w.Ctx(), "节点负责人配置错误-为设置正确的分配规则")
		}

		if rule.AppointedOwner == nil || len(rule.AppointedOwner) == 0 {
			return errs.Business(w.Ctx(), "节点负责人配置错误-未设置指定的负责人")
		}

		// if len(rule.FillOwner) > 0 {
		// 	if !stream.ContainsArr(rule.AppointedOwner, rule.FillOwner) {
		// 		return errs.Business(w.Ctx(), "节点负责人配置错误-指定的负责人必须包含在指定的负责人中")
		// 	}
		// }

	default:
		return errs.Business(w.Ctx(), "节点负责人配置错误-未设置正确的使用模式")
	}

	return nil
}
