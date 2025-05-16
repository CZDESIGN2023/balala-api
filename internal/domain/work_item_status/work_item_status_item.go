package work_item_status

import (
	"cmp"
	"errors"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"
	"slices"
	"time"
)

type WorkItemStatusItems []*WorkItemStatusItem

type WorkItemStatusItem struct {
	shared.DomainEntity

	Id             int64
	Uuid           string
	UserId         int64
	SpaceId        int64
	WorkItemTypeId int64
	Name           string
	Key            string
	Val            string
	StatusType     consts.WorkItemStatusType
	Ranking        int64
	Status         int64
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      int64
	IsSys          int32
	FlowScope      consts.FlowScope
}

func (wis *WorkItemStatusItems) GetMessages() shared.DomainMessages {
	var opsLogs shared.DomainMessages
	for _, wi := range *wis {
		opsLogs = append(opsLogs, wi.GetMessages()...)
	}
	return opsLogs
}

func (wis *WorkItemStatusItems) GetStatusKeys() []string {
	statusKeys := make([]string, 0)
	for _, v := range *wis {
		statusKeys = append(statusKeys, v.Key)
	}
	return statusKeys
}

func (wis *WorkItemStatusItems) GetStatusByVal(val string) *WorkItemStatusItem {
	for _, v := range *wis {
		if v.Val == val {
			return v
		}
	}

	return nil
}

func (wis *WorkItemStatusItems) GetCompleteStatus() []*WorkItemStatusItem {
	return stream.Filter(*wis, func(wi *WorkItemStatusItem) bool {
		return wi.IsCompleted() ||
			wi.FlowScope == consts.FlowScope_Stateflow && wi.StatusType == consts.WorkItemStatusType_Archived
	})
}

func (wis *WorkItemStatusItems) GetArchivedStatus() []*WorkItemStatusItem {
	return stream.Filter(*wis, func(wi *WorkItemStatusItem) bool {
		return wi.StatusType == consts.WorkItemStatusType_Archived
	})
}

func (wis *WorkItemStatusItems) GetProcessingStatus() []*WorkItemStatusItem {
	return stream.Filter(*wis, func(wi *WorkItemStatusItem) bool {
		return wi.StatusType == consts.WorkItemStatusType_Process
	})
}

func (w *WorkItemStatusItem) IsArchivedTypeState() bool {
	return w.StatusType == consts.WorkItemStatusType_Archived
}

func (w *WorkItemStatusItem) IsProcessingTypeState() bool {
	return w.StatusType == consts.WorkItemStatusType_Process
}

func (w *WorkItemStatusItem) IsTerminated() bool {
	return w.Key == string(consts.WorkItemStatus_TerminatedKey)
}

func (w *WorkItemStatusItem) IsCompleted() bool {
	return w.Key == string(consts.WorkItemStatus_WorkFlowCompletedDefaultKey)
}

func (w *WorkItemStatusItem) IsClose() bool {
	return w.Key == string(consts.WorkItemStatus_WorkFlowCloseDefaultKey)
}

// IsSysDefaultStatus 判断是否是预设任务状态
func (w *WorkItemStatusItem) IsSysDefaultStatus() bool {
	return w.IsSys == 1
}

// IsSysFixStatus 判断是否是预设固定任务状态，不可编辑，不可删除，不可排序
func (w *WorkItemStatusItem) IsSysFixStatus() bool {
	return slices.Contains(consts.SysStatusKeyList(), consts.WorkItemStatusDefaultKey(w.Key))
}

func (w *WorkItemStatusItem) ChangeRanking(ranking int64, oper shared.Oper) {

	if w.Ranking == ranking {
		return
	}

	oldValue := w.Ranking
	w.UpdateRanking(ranking)

	w.AddMessage(oper, &domain_message.ModifyWorkItemStatus{
		SpaceId:            w.SpaceId,
		WorkItemStatusName: w.Name,
		WorkItemStatusId:   w.Id,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "ranking",
				OldValue: oldValue,
				NewValue: w.Ranking,
			},
		},
	})
}

func (w *WorkItemStatusItem) UpdateRanking(ranking int64) {
	if w.Ranking == ranking {
		return
	}

	w.Ranking = ranking
	w.AddDiff(Diff_Ranking)
}

func (w *WorkItemStatusItem) ChangeName(name string, oper shared.Oper) error {
	oldValue := w.Name
	w.UpdateName(name)

	w.AddMessage(oper, &domain_message.ModifyWorkItemStatus{
		SpaceId:            w.SpaceId,
		WorkItemStatusName: w.Name,
		WorkItemStatusId:   w.Id,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: w.Name,
			},
		},
	})

	return nil
}

func (w *WorkItemStatusItem) UpdateName(name string) {

	if w.Name == name {
		return
	}

	w.Name = name
	w.AddDiff(Diff_Name)
}

func (w *WorkItemStatusItem) OnDelete(oper shared.Oper) error {
	if w.IsSysDefaultStatus() {
		return errors.New("不能删除预设任务状态")
	}

	w.DeletedAt = time.Now().Unix()
	w.AddDiff(Diff_DeletedAt)
	w.AddMessage(oper, &domain_message.DeleteWorkItemStatus{
		SpaceId:            w.SpaceId,
		WorkItemStatusName: w.Name,
		WorkItemStatusId:   w.Id,
		FlowScope:          w.FlowScope,
	})

	return nil
}

func (w *WorkItemStatusItem) IsSameSpace(spaceId int64) bool {
	return w.SpaceId == spaceId
}

func (w *WorkItemStatusItem) IsForWorkflow() bool {
	return w.FlowScope == consts.FlowScope_Workflow || w.FlowScope == consts.FlowScope_All
}

func (w *WorkItemStatusItem) IsForStateFlow() bool {
	return w.FlowScope == consts.FlowScope_Stateflow || w.FlowScope == consts.FlowScope_All
}

func (w *WorkItemStatusItem) Compare(b *WorkItemStatusItem) int {
	archivedStatusKeyList := consts.SysStatusKeyList()

	aKey := consts.WorkItemStatusDefaultKey(w.Key)
	bKey := consts.WorkItemStatusDefaultKey(b.Key)

	if slices.Contains(archivedStatusKeyList, aKey) && slices.Contains(archivedStatusKeyList, bKey) {
		aIdx := slices.Index(archivedStatusKeyList, aKey)
		bIdx := slices.Index(archivedStatusKeyList, bKey)
		return cmp.Compare(aIdx, bIdx)
	}

	if slices.Contains(archivedStatusKeyList, aKey) {
		return 1
	}

	if slices.Contains(archivedStatusKeyList, bKey) {
		return -1
	}

	if w.FlowScope == b.FlowScope {
		return -cmp.Compare(w.Ranking, b.Ranking)
	}

	return cmp.Compare(w.FlowScope, b.FlowScope)
}
