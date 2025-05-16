package rsp_convert

import (
	"encoding/json"
	"go-cs/internal/bean/vo/rsp"
	domain "go-cs/internal/domain/work_flow"
	tplt_conf "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/utils"

	"github.com/spf13/cast"
)

func WorkFlowTpltToRsp(w *domain.WorkFlowTemplate) *rsp.WorkFlowTemplate {
	r := &rsp.WorkFlowTemplate{
		Id:             w.Id,
		Uuid:           w.Uuid,
		Version:        int64(w.Version),
		SpaceId:        w.SpaceId,
		WorkItemTypeId: w.WorkItemTypeId,
		WorkFlowId:     w.WorkFlowId,
		FlowMode:       string(w.FlowMode),
		Conf:           WorkFlowConfToRsp(w.WorkFLowConfig),
	}
	return r
}

func WorkFlowDefaultConf(w *tplt_conf.WorkFlowDefaultConf) *rsp.WorkFlowDefaultConf {
	r := &rsp.WorkFlowDefaultConf{
		ReasonOptions: make([]*rsp.WorkFlowDefaultConfReasonOption, 0),
	}

	for _, v := range w.ReasonOptions {
		r.ReasonOptions = append(r.ReasonOptions, &rsp.WorkFlowDefaultConfReasonOption{
			Name:    v.Name,
			Options: v.Options,
		})
	}
	return r
}

func WorkFlowConfToRsp(w *tplt_conf.WorkFlow) *rsp.WorkFlowTemplateConf {
	r := &rsp.WorkFlowTemplateConf{
		Key:                               w.Key,
		Uuid:                              w.Uuid,
		Version:                           0,
		Nodes:                             make([]*rsp.WorkFlowNode, 0),
		TerminatedReasonOptions:           w.TerminatedReasonOptions,
		EnableTerminatedReasonOtherOption: w.EnableRebootReasonOtherOption,
		RebootReasonOptions:               w.RebootReasonOptions,
		EnableRebootReasonOtherOption:     w.EnableRebootReasonOtherOption,
		FormFields:                        make([]*rsp.WorkFlowFormField, 0),
		EnableResumeReasonOtherOption:     w.EnableResumeReasonOtherOption,
		ResumeReasonOptions:               w.ResumeReasonOptions,
	}

	for _, v := range w.Nodes {
		r.Nodes = append(r.Nodes, WorkFlowNodeToRsp(v))
	}

	for _, v := range w.Connections {
		r.Connections = append(r.Connections, WorkFlowConnectionToRsp(v))
	}

	for _, v := range w.FormFields {
		r.FormFields = append(r.FormFields, &rsp.WorkFlowFormField{
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return r
}

func WorkFlowNodeToRsp(n *tplt_conf.WorkFlowNode) *rsp.WorkFlowNode {
	r := &rsp.WorkFlowNode{
		Name:                            n.Name,
		DoneOperationDisplayName:        n.DoneOperationDisplayName,
		Key:                             n.Key,
		StartMode:                       string(n.StartMode),
		NeedDoneOperator:                n.NeedDoneOperator,
		BelongStatus:                    n.BelongStatus,
		DoneOperationRole:               n.DoneOperationRole,
		PassMode:                        string(n.PassMode),
		EnableRollback:                  n.EnableRollback,
		EnableRollbackReasonOtherOption: n.EnableRollbackReasonOtherOption,
		RollbackReasonOptions:           n.RollbackReasonOptions,
		EnableClose:                     n.EnableClose,
		EnableCloseReasonOtherOption:    n.EnableCloseReasonOtherOption,
		CloseReasonOptions:              n.CloseReasonOptions,
		RestartReasonOptions:            n.RestartReasonOptions,
		EnableRestartReasonOtherOption:  n.EnableRestartReasonOtherOption,
		ForcePlanTime:                   n.ForcePlanTime,

		Owner: WorkFlowNodeOwnerConfToRsp(n.Owner),
	}

	for _, v := range n.OnReach {
		r.OnReach = append(r.OnReach, WorkFlowNodeEventToRsp(v))
	}

	for _, v := range n.OnPass {
		r.OnPass = append(r.OnPass, WorkFlowNodeEventToRsp(v))
	}

	return r

}

func WorkFlowNodeEventToRsp(e *tplt_conf.WorkFlowNodeEvent) *rsp.WorkFlowNodeEvent {
	r := &rsp.WorkFlowNodeEvent{}
	r.EventType = e.EventType
	if e.TargetSubState != nil {
		r.TargetSubState = &rsp.WorkFlowSubState{
			Id:   cast.ToInt64(e.TargetSubState.Id),
			Key:  e.TargetSubState.Key,
			Uuid: e.TargetSubState.Uuid,
			Val:  e.TargetSubState.Val,
		}
	}

	if v, err := json.Marshal(e.Condition); err == nil {
		r.Condition = string(v)
	}
	return r
}

func WorkFlowNodeOwnerConfToRsp(c *tplt_conf.OwnerConf) *rsp.OwnerConf {
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

func WorkFlowConnectionToRsp(conn *tplt_conf.WorkFlowConnection) *rsp.WorkFlowConnection {
	r := &rsp.WorkFlowConnection{
		StartNode: conn.StartNode,
		EndNode:   conn.EndNode,
	}
	return r
}
