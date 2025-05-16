package rsp_convert

import (
	"go-cs/internal/bean/vo/rsp"
	domain "go-cs/internal/domain/work_flow"
	tplt_conf "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/utils"

	"github.com/spf13/cast"
)

func StateFlowTpltToRsp(w *domain.WorkFlowTemplate) *rsp.StateFlowTemplate {
	r := &rsp.StateFlowTemplate{
		Id:             w.Id,
		Uuid:           w.Uuid,
		Version:        int64(w.Version),
		SpaceId:        w.SpaceId,
		WorkItemTypeId: w.WorkItemTypeId,
		WorkFlowId:     w.WorkFlowId,
		FlowMode:       string(w.FlowMode),
		Conf:           StateFlowConfToRsp(w.StateFlowConf()),
	}
	return r
}

func StateFlowConfToRsp(s *tplt_conf.StateFlow) *rsp.StateFlowTemplateConf {
	var nodes []*rsp.StateFlowNode
	var connections []*rsp.StateFlowConnection

	for _, v := range s.StateFlowNodes {
		nodes = append(nodes, StateFlowNodeToRsp(v))
	}

	for _, v := range s.StateFlowTransitionRule {
		connections = append(connections, StateFlowConnectionToRsp(v))
	}

	return &rsp.StateFlowTemplateConf{
		Key:                               s.Key,
		Uuid:                              s.Uuid,
		StateFlowNodes:                    nodes,
		StateFlowTransitionRule:           connections,
		TerminatedReasonOptions:           s.TerminatedReasonOptions,
		ResumeReasonOptions:               s.ResumeReasonOptions,
		EnableTerminatedReasonOtherOption: s.EnableTerminatedReasonOtherOption,
		EnableResumeReasonOtherOption:     s.EnableResumeReasonOtherOption,
	}
}

func StateFlowNodeToRsp(node *tplt_conf.StateFlowNode) *rsp.StateFlowNode {
	r := &rsp.StateFlowNode{
		Key:             node.Key,
		SubStateId:      node.SubStateId,
		SubStateKey:     node.SubStateKey,
		SubStateVal:     node.SubStateVal,
		IsArchivedState: node.IsArchivedState,
		IsInitState:     node.IsInitState,
		Name:            node.Name,
		OperationRole:   node.OperationRole,
		Owner:           StateFlowNodeOwnerConfToRsp(node.Owner),
	}

	return r
}

func StateFlowConnectionToRsp(connection *tplt_conf.StateFlowTransitionRule) *rsp.StateFlowConnection {
	r := &rsp.StateFlowConnection{
		ConfirmForm:    connection.ConfirmForm,
		SourceStateKey: connection.SourceStateKey,
		TargetStateKey: connection.TargetStateKey,
	}
	return r
}

func StateFlowNodeOwnerConfToRsp(c *tplt_conf.OwnerConf) *rsp.OwnerConf {

	r := &rsp.OwnerConf{
		ForceOwner: c.ForceOwner,
		UsageMode:  string(c.UsageMode),
		Value:      utils.ToJSON(c.Value),
		OwnerRole:  make([]*rsp.OwnerConf_Role, 0),
	}

	for _, v := range c.OwnerRole {
		r.OwnerRole = append(r.OwnerRole, &rsp.OwnerConf_Role{
			Id:   cast.ToInt64(v.Id),
			Key:  v.Key,
			Uuid: v.Uuid,
		})
	}

	return r
}
