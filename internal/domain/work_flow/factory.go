package work_flow

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

func NewWorkFlow(
	id int64,
	spaceId int64,
	workItemTypeId int64,
	name string,
	key string,
	ranking int64,
	flowMode consts.WorkFlowMode,
	template *WorkFlowTemplate,
	status WorkFlowStatus,
	isSys int32,
	uid int64,
	oper shared.Oper,
) *WorkFlow {

	wf := &WorkFlow{
		Id:             id,
		Uuid:           uuid.NewString(),
		SpaceId:        spaceId,
		WorkItemTypeId: workItemTypeId,
		Name:           name,
		Key:            key,
		Ranking:        ranking,
		FlowMode:       flowMode,
		Status:         status,
		IsSys:          isSys,
		UserId:         uid,
	}

	if template != nil {
		wf.LastTemplateId = template.Id
		wf.Version = template.Version
	}

	wf.AddMessage(oper, &domain_message.CreateWorkFlow{
		SpaceId:      spaceId,
		WorkFlowId:   id,
		WorkFlowName: name,
		FlowMode:     flowMode,
	})

	return wf
}

func NewWorkFlowTemplate(
	id int64,
	spaceId int64,
	workItemTypeId int64,
	workFlowId int64,
	version int32,
	flowMode consts.WorkFlowMode,
	workflowConf *config.WorkFlow,
	stateFlowConf *config.StateFlow,
	status WorkFlowTemplateStatus,
	uid int64,
	oper shared.Oper,
) *WorkFlowTemplate {

	return &WorkFlowTemplate{
		Id:              id,
		UserId:          uid,
		Uuid:            uuid.NewString(),
		SpaceId:         spaceId,
		WorkItemTypeId:  workItemTypeId,
		WorkFlowId:      workFlowId,
		Version:         version,
		FlowMode:        flowMode,
		WorkFLowConfig:  workflowConf,
		StateFlowConfig: stateFlowConf,
		CreatedAt:       time.Now().Unix(),
		Status:          status,
	}
}
