package work_item

import (
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"
)

type WorkItemFlowRoles []*WorkItemFlowRole

func (w *WorkItemFlowRoles) GetIds() []int64 {
	ids := make([]int64, 0, len(*w))
	for _, item := range *w {
		ids = append(ids, item.Id)
	}
	return ids
}

func (w *WorkItemFlowRoles) GetByRoleKey(roleKey string) *WorkItemFlowRole {
	for _, item := range *w {
		if item.WorkItemRoleKey == roleKey {
			return item
		}
	}

	return nil
}

type WorkItemFlowRole struct {
	shared.DomainEntity

	Id              int64     `json:"id"`
	SpaceId         int64     `json:"space_id"`
	WorkItemId      int64     `json:"work_item_id"`
	FlowId          int64     `json:"flow_id"`
	FlowTemplateId  int64     `json:"flow_template_id"`
	WorkItemRoleId  int64     `json:"work_item_role_id"`
	WorkItemRoleKey string    `json:"work_item_role_key"`
	Directors       Directors `json:"directors"`
	CreatedAt       int64     `json:"created_at"`
	UpdatedAt       int64     `json:"updated_at"`
	DeletedAt       int64     `json:"deleted_at"`
}

func (role *WorkItemFlowRole) UpdateDirectors(directors Directors) {
	role.Directors = stream.Unique(directors)
	role.AddDiff(Diff_WorkItemFlowRole_Directors)
}

func (role *WorkItemFlowRole) AddDirectors(directors Directors) {

	addDirectors := stream.Filter(directors, func(id string) bool {
		return !stream.Contains(role.Directors, id)
	})

	if len(addDirectors) == 0 {
		return
	}

	newDirectors := append(role.Directors, addDirectors...)
	role.UpdateDirectors(newDirectors)
}

func (role *WorkItemFlowRole) RemoveDirectors(directors Directors) {

	removeDirectors := stream.Filter(directors, func(id string) bool {
		return stream.Contains(role.Directors, id)
	})

	if len(removeDirectors) == 0 {
		return
	}

	newDirectors := stream.Filter(role.Directors, func(id string) bool {
		return !stream.Contains(removeDirectors, id)
	})

	role.UpdateDirectors(newDirectors)
}
