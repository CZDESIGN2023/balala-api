package facade

import (
	"context"
	witem_status_domain "go-cs/internal/domain/work_item_status"
	witem_status_service "go-cs/internal/domain/work_item_status/service"
)

type WorkItemStatusFacade struct {
	status *witem_status_domain.WorkItemStatusInfo
}

func (w *WorkItemStatusFacade) GetItemByKey(itemKey string) *witem_status_domain.WorkItemStatusItem {
	return w.status.GetItemByKey(itemKey)
}

func (w *WorkItemStatusFacade) HasArchivedItem(itemKey string) bool {
	return w.status.HasArchivedItem(itemKey)
}

func (w *WorkItemStatusFacade) Keyword() *witem_status_domain.WorkItemStatusKeyword {
	return w.status.Keyword()
}

func BuildWorkItemStatusFacade(statusInfo *witem_status_domain.WorkItemStatusInfo) *WorkItemStatusFacade {
	return &WorkItemStatusFacade{
		status: statusInfo,
	}
}

type WorkItemStatusServiceFacade struct {
	service *witem_status_service.WorkItemStatusService
}

func (w *WorkItemStatusServiceFacade) GetWorkItemStatusItem(ctx context.Context, statusId int64) (*witem_status_domain.WorkItemStatusItem, error) {
	return w.service.GetWorkItemStatusItem(ctx, statusId)
}

func BuildWorkItemStatusServiceFacade(service *witem_status_service.WorkItemStatusService) *WorkItemStatusServiceFacade {
	return &WorkItemStatusServiceFacade{
		service: service,
	}
}
