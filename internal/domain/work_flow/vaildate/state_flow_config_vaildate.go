package vaildate

import (
	"context"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/utils/errs"
)

// --各种检查配置是否正确的
type StateFlowConfigValidate struct {
	conf *config.StateFlow
	ctx  context.Context
}

func NewStateFlowConfigValidate(conf *config.StateFlow) *StateFlowConfigValidate {
	return &StateFlowConfigValidate{
		conf: conf,
	}
}

func (w *StateFlowConfigValidate) WithCtx(ctx context.Context) *StateFlowConfigValidate {
	w.ctx = ctx
	return w
}

func (w *StateFlowConfigValidate) Ctx() context.Context {
	if w.ctx == nil {
		w.ctx = context.Background()
	}
	return w.ctx
}

func (w *StateFlowConfigValidate) Valid() error {

	//if err := w.validStateFlowConnection(); err != nil {
	//	return err
	//}

	if err := w.validNodes(); err != nil {
		return err
	}

	return nil
}

func (w *StateFlowConfigValidate) validStateFlowConnection() error {
	// 每个节点至少有一个度
	degreeCount := map[string]int{}

	for _, v := range w.conf.StateFlowTransitionRule {
		degreeCount[v.SourceStateKey]++
		degreeCount[v.TargetStateKey]++
	}

	for _, v := range w.conf.StateFlowNodes {
		if degreeCount[v.Key] == 0 {
			return errs.Business(w.Ctx(), "状态节点至少有一个度")
		}
	}

	return nil
}

func (w *StateFlowConfigValidate) validNodes() error {
	nodeLen := len(w.conf.StateFlowNodes)
	if nodeLen > 40 {
		return errs.Business(w.Ctx(), "状态数量已达最大数量")
	}

	initStateNodeLength := 0
	for _, v := range w.conf.StateFlowNodes {
		if v.IsInitState {
			initStateNodeLength++
		}
	}

	if initStateNodeLength == 0 || initStateNodeLength > 1 {
		return errs.Business(w.Ctx(), "必须设置且只需要一个初始状态")
	}

	//检查节点的基本配置
	for _, v := range w.conf.StateFlowNodes {
		if v.Name == "" || v.Key == "" {
			return errs.Business(w.Ctx(), "状态名称错误")
		}

		if v.Key != v.SubStateKey {
			return errs.Business(w.Ctx(), "key != subStateKey")
		}

		if v.SubStateId == "" || v.SubStateKey == "" || v.SubStateVal == "" {
			return errs.Business(w.Ctx(), "节点subState配置错误")
		}

		err := w.ValidNodeOwnerConf(v.Owner)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *StateFlowConfigValidate) ValidNodeOwnerConf(conf *config.OwnerConf) error {
	if conf == nil {
		return errs.Business(w.Ctx(), "状态负责人配置错误-未设置关联角色")
	}

	if len(conf.OwnerRole) == 0 {
		return errs.Business(w.Ctx(), "状态负责人配置错误-未设置关联角色")
	}

	if !conf.ForceOwner {
		return errs.Business(w.Ctx(), "状态负责人配置错误-未设置强制关联角色")
	}

	switch conf.UsageMode {
	case config.UsageMode_None:
		return nil
	case config.UsageMode_Appointed:
		ruleVal := conf.GetUsageModeVal()
		rule, isOk := ruleVal.(*config.OwnerConf_UsageMode_Appointed)
		if !isOk {
			return errs.Business(w.Ctx(), "状态负责人配置错误-为设置正确的分配规则")
		}

		if rule.AppointedOwner == nil || len(rule.AppointedOwner) == 0 {
			return errs.Business(w.Ctx(), "状态负责人配置错误-未设置指定的负责人")
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
