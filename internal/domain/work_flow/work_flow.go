package work_flow

import (
	"errors"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

type WorkFlowStatus int32

var (
	WorkFlowStatus_Enable  WorkFlowStatus = 1
	WorkFlowStatus_Disable WorkFlowStatus = 0
	WorkFlowStatus_Hide    WorkFlowStatus = -1
)

type WorkFlows []*WorkFlow

type WorkFlow struct {
	shared.AggregateRoot

	Id             int64               `json:"id"`
	Uuid           string              `json:"uuid"`
	UserId         int64               `json:"user_id"`
	SpaceId        int64               `json:"space_id"`
	WorkItemTypeId int64               `json:"work_item_type_id"`
	Name           string              `json:"name"`
	FlowMode       consts.WorkFlowMode `json:"flow_mode"`
	Version        int32               `json:"version"`
	Key            string              `json:"key"`
	Ranking        int64               `json:"ranking"`
	LastTemplateId int64               `json:"last_template_id"`
	IsSys          int32               `json:"is_sys"`

	Status    WorkFlowStatus `json:"status"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt int64          `json:"deleted_at"`
}

func (ws *WorkFlows) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)

	for _, w := range *ws {
		logs = append(logs, w.GetMessages()...)
	}

	return logs
}

func (w *WorkFlow) UpdateLastTemplate(template *WorkFlowTemplate) {
	w.LastTemplateId = template.Id
	w.Version = template.Version

	w.AddDiff(Diff_LastTemplate, Diff_Version)
}

func (w *WorkFlow) IsSameSpace(spaceId int64) bool {
	return w.SpaceId == spaceId
}

func (w *WorkFlow) IsLastTemplate(template *WorkFlowTemplate) bool {
	return w.LastTemplateId == template.Id
}

func (w *WorkFlow) UpdateRanking(ranking int64) {
	if w.Ranking == ranking {
		return
	}

	w.Ranking = ranking
	w.AddDiff(Diff_Ranking)
}

func (w *WorkFlow) ChangeRanking(ranking int64, oper shared.Oper) {
	if w.Ranking == ranking {
		return
	}

	oldValue := w.Ranking
	w.UpdateRanking(ranking)

	w.AddMessage(oper, &domain_message.ModifyWorkFlow{
		SpaceId:      w.SpaceId,
		WorkFlowName: w.Name,
		WorkFlowId:   w.Id,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "ranking",
				OldValue: oldValue,
				NewValue: ranking,
			},
		},
	})

}

func (w *WorkFlow) ChangeStatus(status int64, oper shared.Oper) error {

	if WorkFlowStatus(status) != WorkFlowStatus_Enable &&
		WorkFlowStatus(status) != WorkFlowStatus_Disable &&
		WorkFlowStatus(status) != WorkFlowStatus_Hide {
		return errors.New("状态错误")
	}

	oldValue := w.Status

	w.UpdateStatus(status)

	w.AddMessage(oper, &domain_message.ModifyWorkFlow{
		SpaceId:      w.SpaceId,
		WorkFlowName: w.Name,
		WorkFlowId:   w.Id,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "status",
				OldValue: oldValue,
				NewValue: w.Status,
			},
		},
	})

	return nil
}

func (w *WorkFlow) UpdateStatus(status int64) {

	if w.Status == WorkFlowStatus(status) {
		return
	}

	w.Status = WorkFlowStatus(status)
	w.AddDiff(Diff_Status)

}

func (w *WorkFlow) IsEnable() bool {
	return w.Status == WorkFlowStatus_Enable
}

func (w *WorkFlow) IsSyPreset() bool {
	return w.IsSys == 1
}

func (w *WorkFlow) UpdateName(name string) {

	if w.Name == name {
		return
	}

	w.Name = name
	w.AddDiff(Diff_Name)
}

func (w *WorkFlow) ChangeName(name string, oper shared.Oper) {
	oldValue := w.Name

	w.UpdateName(name)

	w.AddMessage(oper, &domain_message.ModifyWorkFlow{
		SpaceId:      w.SpaceId,
		WorkFlowName: w.Name,
		WorkFlowId:   w.Id,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: w.Name,
			},
		},
	})
}

func (w *WorkFlow) IsWorkFlowMode() bool {
	return w.FlowMode == consts.FlowMode_WorkFlow
}

func (w *WorkFlow) IsStateFlowMode() bool {
	return w.FlowMode == consts.FlowMode_StateFlow
}

func (w *WorkFlow) OnDelete(oper shared.Oper) {
	w.UpdateDeleted()
	w.AddMessage(oper, &domain_message.DeleteWorkFlow{
		SpaceId:      w.SpaceId,
		WorkFlowId:   w.Id,
		WorkFlowName: w.Name,
	})
}

func (w *WorkFlow) UpdateDeleted() {
	w.DeletedAt = time.Now().Unix()
	w.AddDiff(Diff_DeletedAt)
}
