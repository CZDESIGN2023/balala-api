package facade

import "go-cs/internal/domain/work_item"

type WorkItemFacade struct {
	workItem *work_item.WorkItem
}

func (w *WorkItemFacade) GetWorkItem() *work_item.WorkItem {
	return w.workItem
}

func BuildWorkItemFacade(workItem *work_item.WorkItem) *WorkItemFacade {
	return &WorkItemFacade{
		workItem: workItem,
	}
}
