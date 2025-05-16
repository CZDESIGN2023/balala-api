package workitem_type

import (
	"go-cs/internal/consts"
	"time"

	"github.com/google/uuid"
)

func NewWorkItemType(id int64, spaceId SpaceId, name string, key string, flowMode consts.WorkFlowMode, isSys int32, uid int64) *WorkItemType {
	return &WorkItemType{
		Id:        id,
		UserId:    uid,
		Uuid:      uuid.NewString(),
		SpaceId:   spaceId,
		Name:      name,
		Key:       key,
		FlowMode:  flowMode,
		Ranking:   time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
		IsSys:     isSys,
		Status:    1,
	}
}
