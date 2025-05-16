package query

import (
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_type"
	"go-cs/pkg/stream"
)

type WorkItemTypeInfoQuery struct {
	SpaceId int64
}

func BuildWorkItemTypeInfoQueryResult(spaceId int64, list []*domain.WorkItemType) *WorkItemTypeInfoQueryResult {
	return &WorkItemTypeInfoQueryResult{
		SpaceId: spaceId,
		list:    list,
	}
}

type WorkItemTypeInfoQueryResult struct {
	SpaceId int64
	list    []*domain.WorkItemType
}

func (s *WorkItemTypeInfoQueryResult) HasType(id int64) bool {
	for _, v := range s.list {
		if v.Id == id {
			return true
		}
	}
	return false
}

func (s *WorkItemTypeInfoQueryResult) GetFlowTaskTypeByFlowMode(flowMode consts.WorkFlowMode) *domain.WorkItemType {
	for _, v := range s.list {
		if v.IsWorkFlowTaskType() {
			return v
		}
	}
	return nil
}

func (s *WorkItemTypeInfoQueryResult) GetWorkFlowTaskType() *domain.WorkItemType {
	for _, v := range s.list {
		if v.IsWorkFlowTaskType() {
			return v
		}
	}
	return nil
}

func (s *WorkItemTypeInfoQueryResult) GetStateFlowTaskType() *domain.WorkItemType {
	for _, v := range s.list {
		if v.IsStateFlowTaskType() {
			return v
		}
	}
	return nil
}

func (s *WorkItemTypeInfoQueryResult) GetSubTaskType() *domain.WorkItemType {
	for _, v := range s.list {
		if v.IsSubTaskType() {
			return v
		}
	}
	return nil
}

// GetMainTaskTypeIds 获取主任务类型id，不包含子任务类型
func (s *WorkItemTypeInfoQueryResult) GetMainTaskTypeIds() []int64 {
	var ids []int64
	for _, v := range s.list {
		if !v.IsSubTaskType() {
			ids = append(ids, v.Id)
		}
	}
	return ids
}

func (s *WorkItemTypeInfoQueryResult) Types() []*domain.WorkItemType {
	return s.list
}

func (s *WorkItemTypeInfoQueryResult) Maps() map[int64]*domain.WorkItemType {
	return stream.ToMap(s.list, func(i int, t *domain.WorkItemType) (int64, *domain.WorkItemType) {
		return t.Id, t
	})
}

func (s *WorkItemTypeInfoQueryResult) GetByTypeId(workItemTypeId int64) *domain.WorkItemType {
	for _, v := range s.list {
		if v.Id == workItemTypeId {
			return v
		}
	}
	return nil
}
