package event

import (
	"go-cs/api/notify"

	"go-cs/internal/domain/space"
	"go-cs/internal/domain/work_item"
)

type RemindWork struct {
	Event     notify.Event
	Space     *space.Space
	WorkItem  *work_item.WorkItem
	Operator  int64
	TargetIds []int64
}
