package work_item_status

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

func BuildWorkItemStatusInfo(spaceId int64, items []*WorkItemStatusItem) *WorkItemStatusInfo {
	info := &WorkItemStatusInfo{
		SpaceId:        spaceId,
		WorkItemTypeId: 0,
		Items:          items,
	}

	info.Init()
	return info
}

func NewWorkItemStatusItem(id int64, spaceId int64, name string, key string, val string, isSys int32, ranking int64, statusType consts.WorkItemStatusType, uid int64, flowScope consts.FlowScope, oper shared.Oper) *WorkItemStatusItem {

	ins := &WorkItemStatusItem{
		Id:             id,
		SpaceId:        spaceId,
		UserId:         uid,
		WorkItemTypeId: 0,
		Name:           name,
		Key:            key,
		Val:            val,
		IsSys:          isSys,
		Status:         1,
		StatusType:     statusType,
		Uuid:           uuid.NewString(),
		CreatedAt:      time.Now().Unix(),
		Ranking:        ranking,
		FlowScope:      flowScope,
	}

	ins.AddMessage(oper, &domain_message.CreateWorkItemStatus{
		SpaceId:            spaceId,
		WorkItemStatusName: name,
		WorkItemStatusId:   id,
		FlowScope:          flowScope,
	})
	return ins
}
