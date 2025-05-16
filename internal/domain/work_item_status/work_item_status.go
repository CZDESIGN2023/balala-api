package work_item_status

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
)

type WorkItemStatusKeyword struct {
	Terminated *WorkItemStatusItem
}

type WorkItemStatusInfo struct {
	shared.AggregateRoot

	SpaceId        int64
	WorkItemTypeId int64

	Items       WorkItemStatusItems
	keywordItem *WorkItemStatusKeyword
	itemMap     map[string]*WorkItemStatusItem
}

func (w *WorkItemStatusInfo) Init() {

	w.itemMap = map[string]*WorkItemStatusItem{}
	for _, v := range w.Items {
		w.itemMap[v.Key] = v
	}

	w.keywordItem = &WorkItemStatusKeyword{
		Terminated: w.itemMap[string(consts.WorkItemStatus_TerminatedKey)],
	}
}

func (w *WorkItemStatusInfo) Keyword() *WorkItemStatusKeyword {
	return w.keywordItem
}

func (w *WorkItemStatusInfo) GetItemById(id int64) *WorkItemStatusItem {
	for _, v := range w.Items {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func (w *WorkItemStatusInfo) GetItemByKey(itemKey string) *WorkItemStatusItem {
	return w.itemMap[itemKey]
}

func (w *WorkItemStatusInfo) HasArchivedItem(key string) bool {
	item := w.GetItemByKey(key)
	return item != nil && item.IsArchivedTypeState()
}

func (w *WorkItemStatusInfo) IsSameSpace(spaceId int64) bool {
	return w.SpaceId == spaceId
}

func (w *WorkItemStatusInfo) GetProcessTypeItems(flowScope consts.FlowScope) WorkItemStatusItems {
	items := make(WorkItemStatusItems, 0)
	for _, v := range w.Items {
		if w.isInFlowScope(v, flowScope) {
			if !v.IsArchivedTypeState() {
				items = append(items, v)
			}
		}
	}
	return items
}

func (w *WorkItemStatusInfo) isInFlowScope(statusItem *WorkItemStatusItem, flowScope consts.FlowScope) bool {
	if flowScope == consts.FlowScope_All || statusItem.FlowScope == flowScope || statusItem.FlowScope == consts.FlowScope_All {
		return true
	}
	return false
}

func (w *WorkItemStatusInfo) GetArchivedTypeItems(flowScope consts.FlowScope) WorkItemStatusItems {
	items := make(WorkItemStatusItems, 0)
	for _, v := range w.Items {
		if w.isInFlowScope(v, flowScope) {
			if v.IsArchivedTypeState() {
				items = append(items, v)
			}
		}
	}
	return items
}

func (w *WorkItemStatusInfo) GetItemsByFlowScope(flowScope consts.FlowScope) WorkItemStatusItems {

	items := make(WorkItemStatusItems, 0)
	for _, v := range w.Items {
		if w.isInFlowScope(v, flowScope) {
			items = append(items, v)
		}
	}
	return items
}
