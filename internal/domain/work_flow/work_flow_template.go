package work_flow

import (
	"go-cs/internal/consts"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"

	"github.com/spf13/cast"
)

type WorkFlowTemplateStatus int32

var (
	WorkFlowTemplateStatus_Enable  WorkFlowTemplateStatus = 1
	WorkFlowTemplateStatus_Disable WorkFlowTemplateStatus = 1
)

type WorkFlowTemplates []*WorkFlowTemplate

type WorkFlowTemplate struct {
	shared.AggregateRoot
	Id             int64               `json:"id"`
	UserId         int64               `json:"user_id"`
	Uuid           string              `json:"uuid"`
	SpaceId        int64               `json:"space_id"`
	WorkItemTypeId int64               `json:"work_item_type_id"`
	WorkFlowId     int64               `json:"work_flow_id"`
	Name           string              `json:"name"`
	FlowMode       consts.WorkFlowMode `json:"flow_mode"`
	Version        int32               `json:"version"`
	Key            string              `json:"key"`
	Ranking        int64               `json:"ranking"`

	Status    WorkFlowTemplateStatus `json:"status"`
	CreatedAt int64                  `json:"created_at"`
	UpdatedAt int64                  `json:"updated_at"`
	DeletedAt int64                  `json:"deleted_at"`

	WorkFLowConfig  *config.WorkFlow  `json:"work_flow_conf"`
	StateFlowConfig *config.StateFlow `json:"state_flow_conf"`
}

func (w *WorkFlowTemplate) WorkFlowConf() *config.WorkFlow {
	return w.WorkFLowConfig
}

func (w *WorkFlowTemplate) StateFlowConf() *config.StateFlow {
	return w.StateFlowConfig
}

func (w *WorkFlowTemplate) IsStateFlow() bool {
	return w.FlowMode == consts.FlowMode_StateFlow
}

func (w *WorkFlowTemplate) IsWorkFlow() bool {
	return w.FlowMode == consts.FlowMode_WorkFlow
}

func (w *WorkFlowTemplate) RemoveOwner(userId int64) {
	userIdStr := cast.ToString(userId)
	switch w.FlowMode {
	case consts.FlowMode_WorkFlow:
		for _, node := range w.WorkFLowConfig.Nodes {
			if node.Owner == nil {
				continue
			}

			switch node.Owner.UsageMode {
			case config.UsageMode_Appointed:
				val := node.Owner.GetAppointedUsageModeVal()

				// 如果用户在 owner 列表中，则设置为 none
				if val.Contains(config.FillOwnerType_User, userIdStr) {
					node.Owner.UsageMode = config.UsageMode_None
					node.Owner.Value = &config.OwnerConf_UsageMode_None{
						FillOwner: []*config.OwnerConf_UsageMode_FillOwner{}, //必须为空数组
					}
				}
			case config.UsageMode_None:
				val := node.Owner.GetNoneUsageModeVal()

				// 如果用户在 owner 列表中，则移除
				node.Owner.Value = &config.OwnerConf_UsageMode_None{
					FillOwner: stream.Filter(val.FillOwner, func(e *config.OwnerConf_UsageMode_FillOwner) bool {
						return e.Value != userIdStr
					}),
				}
			}
		}
	case consts.FlowMode_StateFlow:
		for _, node := range w.StateFlowConfig.StateFlowNodes {
			if node.Owner == nil {
				continue
			}

			switch node.Owner.UsageMode {
			case config.UsageMode_Appointed:
				val := node.Owner.GetAppointedUsageModeVal()

				// 如果用户在 owner 列表中，则设置为 none
				if val.Contains(config.FillOwnerType_User, userIdStr) {
					node.Owner.UsageMode = config.UsageMode_None
					node.Owner.Value = &config.OwnerConf_UsageMode_None{
						FillOwner: []*config.OwnerConf_UsageMode_FillOwner{}, //必须为空数组
					}
				}
			case config.UsageMode_None:
				val := node.Owner.GetNoneUsageModeVal()

				// 如果用户在 owner 列表中，则移除
				node.Owner.Value = &config.OwnerConf_UsageMode_None{
					FillOwner: stream.Filter(val.FillOwner, func(e *config.OwnerConf_UsageMode_FillOwner) bool {
						return e.Value != userIdStr
					}),
				}
			}
		}
	}

	w.AddDiff(Diff_TemplateSetting)
}
