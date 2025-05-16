package task

import (
	"context"
	"go-cs/api/notify"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/utils/date"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"
)

// 项目异常检查，逾期任务太多
func (task *SpaceTask) CheckSpaceAbnormal() {
	begin := time.Now()
	task.log.Info("CheckSpaceAbnormal start")
	defer task.log.Infof("CheckSpaceAbnormal end %v", time.Now().Sub(begin))

	expiredNumMap := task.CountProgressingExpiredWorkItem()
	TotalNumMap := task.CountTotalWorkItem()

	var needNotifySpaceIds []int64

	for spaceId, total := range TotalNumMap {
		if total == 0 {
			continue
		}

		expiredNum := expiredNumMap[spaceId]

		if expiredNum*100/total >= 40 {
			needNotifySpaceIds = append(needNotifySpaceIds, spaceId)
		}
	}

	spaceMap, err := task.spaceRepo.SpaceMap(context.Background(), needNotifySpaceIds)
	if err != nil {
		return
	}

	for _, spaceId := range needNotifySpaceIds {
		expiredNum := expiredNumMap[spaceId]
		space := spaceMap[spaceId]
		if space == nil {
			continue
		}

		bus.Emit(notify.Event_SpaceAbnormal, &event.SpaceAbnormal{
			Event:      notify.Event_SpaceAbnormal,
			Space:      space,
			ExpiredNum: expiredNum,
		})
	}
}

func (task *SpaceTask) CountProgressingExpiredWorkItem() map[int64]int64 {
	const sql = `
SELECT
	t.space_id,
	COUNT(*) AS total
FROM
	space_work_item_v2 t INNER JOIN work_item_status s ON t.work_item_status_id = s.id
WHERE 
    s.status_type != 3 and 
    doc->'$.plan_complete_at' <  ?
GROUP BY
	space_id;
`
	type item struct {
		SpaceId int64
		Total   int64
	}
	var list []item

	err := task.data.RoDB(context.Background()).Model(&db.SpaceWorkItemV2{}).
		Raw(sql, date.TodayBegin().Unix()).
		Find(&list).Error
	if err != nil {
		task.log.Error(err)
		return nil
	}

	return stream.ToMap(list, func(i int, t item) (int64, int64) {
		return t.SpaceId, t.Total
	})
}

func (task *SpaceTask) CountTotalWorkItem() map[int64]int64 {
	const sql = `
SELECT
	space_id,
	count(*) AS total
FROM
	space_work_item_v2
GROUP BY
	space_id;
`
	type item struct {
		SpaceId int64
		Total   int64
	}
	var list []item

	err := task.data.RoDB(context.Background()).Model(&db.SpaceWorkItemV2{}).
		Raw(sql).
		Find(&list).Error
	if err != nil {
		task.log.Error(err)
		return nil
	}

	return stream.ToMap(list, func(i int, t item) (int64, int64) {
		return t.SpaceId, t.Total
	})
}
