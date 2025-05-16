package convert

import (
	"encoding/json"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_flow"
	flow_config "go-cs/internal/domain/work_flow/flow_tplt_config"
)

func WorkFlowTemplateEntityToPo(template *domain.WorkFlowTemplate) *db.WorkFlowTemplate {
	po := &db.WorkFlowTemplate{
		Id:             template.Id,
		Uuid:           template.Uuid,
		SpaceId:        template.SpaceId,
		UserId:         template.UserId,
		WorkItemTypeId: template.WorkItemTypeId,
		WorkFlowId:     template.WorkFlowId,
		FlowMode:       string(template.FlowMode),
		Version:        template.Version,
		Status:         int64(template.Status),
		CreatedAt:      template.CreatedAt,
		UpdatedAt:      template.UpdatedAt,
		DeletedAt:      template.DeletedAt,
	}

	if template.WorkFLowConfig != nil {
		po.Setting = template.WorkFLowConfig.ToString()
	}

	if template.StateFlowConfig != nil {
		po.Setting = template.StateFlowConfig.ToString()
	}

	return po
}

func WorkFlowTemplatePoToEntity(po *db.WorkFlowTemplate) *domain.WorkFlowTemplate {
	ent := &domain.WorkFlowTemplate{
		Id:             po.Id,
		Uuid:           po.Uuid,
		UserId:         po.UserId,
		SpaceId:        po.SpaceId,
		WorkItemTypeId: po.WorkItemTypeId,
		WorkFlowId:     po.WorkFlowId,
		FlowMode:       consts.WorkFlowMode(po.FlowMode),
		Version:        po.Version,
		Status:         domain.WorkFlowTemplateStatus(po.Status),
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
		DeletedAt:      po.DeletedAt,
		//WorkFLowConfig:  domain.WorkFLowConfig{},
	}

	if po.Setting != "" {
		switch po.FlowMode {
		case string(consts.FlowMode_WorkFlow):
			ent.WorkFLowConfig = &flow_config.WorkFlow{}
			json.Unmarshal([]byte(po.Setting), ent.WorkFLowConfig)
		case string(consts.FlowMode_StateFlow):
			ent.StateFlowConfig = &flow_config.StateFlow{}
			json.Unmarshal([]byte(po.Setting), ent.StateFlowConfig)
		}
	}

	return ent
}

func WorkFlowTemplatePoToEntities(po []*db.WorkFlowTemplate) []*domain.WorkFlowTemplate {
	list := make([]*domain.WorkFlowTemplate, 0, len(po))
	for _, v := range po {
		list = append(list, WorkFlowTemplatePoToEntity(v))
	}
	return list
}
