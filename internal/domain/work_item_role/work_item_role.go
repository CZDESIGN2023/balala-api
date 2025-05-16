package work_item_role

import (
	"errors"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

type WorkItemRoles []*WorkItemRole

type WorkItemRole struct {
	shared.AggregateRoot

	Id             int64            `json:"id"`
	Uuid           string           `json:"uuid"`
	UserId         int64            `json:"user_id"`
	SpaceId        int64            `json:"space_id"`
	WorkItemTypeId int64            `json:"work_item_type_id"`
	Key            string           `json:"key"`
	Name           string           `json:"name"`
	IsSys          int32            `json:"is_sys"`
	Status         int64            `json:"status"`
	Ranking        int64            `json:"ranking"`
	CreatedAt      int64            `json:"created_at"`
	UpdatedAt      int64            `json:"updated_at"`
	DeletedAt      int64            `json:"deleted_at"`
	FlowScope      consts.FlowScope `json:"flow_scope"`
}

func (wr *WorkItemRoles) GetRoleById(id int64) *WorkItemRole {
	for _, role := range *wr {
		if role.Id == id {
			return role
		}
	}
	return nil
}

func (wr *WorkItemRoles) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)
	for _, role := range *wr {
		logs = append(logs, role.GetMessages()...)
	}
	return logs
}

func (w *WorkItemRole) UpdateName(name string) {

	if w.Name == name {
		return
	}
	w.Name = name
	w.AddDiff(Diff_Name)
}

func (w *WorkItemRole) ChangeName(name string, oper shared.Oper) error {
	if w.Name == name {
		return errors.New("角色名称未修改")
	}

	oldValue := w.Name
	w.UpdateName(name)

	w.AddMessage(oper, &domain_message.ModifyWorkItemRole{
		SpaceId:          w.SpaceId,
		WorkItemRoleId:   w.Id,
		WorkItemRoleName: w.Name,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: name,
			},
		},
	})

	return nil
}

func (w *WorkItemRole) OnDelete(oper shared.Oper) error {
	if w.IsSysDefaultRole() {
		return errors.New("不能删除预设角色")
	}

	w.AddMessage(oper, &domain_message.DeleteWorkItemRole{
		SpaceId:          w.SpaceId,
		WorkItemRoleId:   w.Id,
		WorkItemRoleName: w.Name,
		FlowScope:        w.FlowScope,
	})

	return nil
}

func (w *WorkItemRole) Enable() {
	w.Status = 1
}

func (w *WorkItemRole) Disable() {
	w.Status = 0
}

func (w *WorkItemRole) ChangeRanking(ranking int64, oper shared.Oper) {
	if w.Ranking == ranking {
		return
	}

	oldValue := w.Ranking
	w.UpdateRanking(ranking)

	w.AddMessage(oper, &domain_message.ModifyWorkItemRole{
		SpaceId:          w.SpaceId,
		WorkItemRoleId:   w.Id,
		WorkItemRoleName: w.Name,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "ranking",
				OldValue: oldValue,
				NewValue: w.Ranking,
			},
		},
	})

}

func (w *WorkItemRole) UpdateRanking(ranking int64) {
	if w.Ranking == ranking {
		return
	}

	w.Ranking = ranking

	w.AddDiff(Diff_Ranking)
}

func (w *WorkItemRole) IsSysDefaultRole() bool {
	return w.IsSys == 1
}

func (w *WorkItemRole) IsSameSpace(spaceId int64) bool {
	return w.SpaceId == spaceId
}

func (w *WorkItemRole) IsForWorkflow() bool {
	return w.FlowScope == consts.FlowScope_Workflow || w.FlowScope == consts.FlowScope_All
}

func (w *WorkItemRole) IsForStateFlow() bool {
	return w.FlowScope == consts.FlowScope_Stateflow || w.FlowScope == consts.FlowScope_All
}
