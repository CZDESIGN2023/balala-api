package condition_translater

import (
	"context"
	status_repo "go-cs/internal/domain/work_item_status/repo"
)

type Ctx struct {
	Ctx         context.Context
	StatusRepo  status_repo.WorkItemStatusRepo
	SpaceIds    []int64
	IsWorkBench bool
}
