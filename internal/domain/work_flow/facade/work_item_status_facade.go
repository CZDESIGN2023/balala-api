package facade

import (
	witem_status_domain "go-cs/internal/domain/work_item_status"
)

type WorkItemStatusInfo struct {
	status *witem_status_domain.WorkItemStatusInfo
}

func (w *WorkItemStatusInfo) GetItemByKey(itemKey string) *witem_status_domain.WorkItemStatusItem {
	return w.status.GetItemByKey(itemKey)
}

func BuildWorkItemStatusInfo(statusInfo *witem_status_domain.WorkItemStatusInfo) *WorkItemStatusInfo {
	return &WorkItemStatusInfo{
		status: statusInfo,
	}
}
