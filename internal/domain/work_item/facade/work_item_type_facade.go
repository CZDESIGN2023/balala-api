package facade

import (
	"go-cs/internal/bean/vo/query"
)

type WorkItemTypeFacade struct {
	info *query.WorkItemTypeInfoQueryResult
}

func BuildWorkItemTypeFacade(info *query.WorkItemTypeInfoQueryResult) *WorkItemTypeFacade {
	return &WorkItemTypeFacade{
		info: info,
	}
}
