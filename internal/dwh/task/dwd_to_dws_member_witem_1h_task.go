package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-cs/internal/dwh/data"
	"go-cs/internal/dwh/model/dws"
	"go-cs/internal/utils/date"
	"go-cs/pkg/stream"
	"time"

	"go-cs/internal/dwh/pkg"

	"gorm.io/gorm"
)

type DwdToDwsMemberWitem1hTaskStateValue struct {
	LastDwdDate string `json:"last_dwd_date"`
}

func (t *DwdToDwsMemberWitem1hTaskStateValue) String() string {
	v, _ := json.Marshal(t)
	return string(v)
}

// 每小时报表 空间版本任务轻度汇总
// 任务状态变量: 最后一次生成dws的时间 last_dws_date 任务被执行的时间
type DwdToDwsMemberWitem1hTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewDwdToDwsMemberWitem1hTask(
	id string,
	ctx *pkg.TaskContext,
) *DwdToDwsMemberWitem1hTask {
	return &DwdToDwsMemberWitem1hTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *DwdToDwsMemberWitem1hTask) Id() string {
	return t.id
}

func (t *DwdToDwsMemberWitem1hTask) Name() string {
	return "dwd_to_dws_member_witem_1h"
}

func (t *DwdToDwsMemberWitem1hTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *DwdToDwsMemberWitem1hTask) Status() string {
	return t.status
}

type memberAggSearchResult struct {
	Num         int64 `gorm:"column:num"`
	ExpireNum   int64 `gorm:"column:expire_num"`
	TodoNum     int64 `gorm:"column:todo_num"`
	CompleteNum int64 `gorm:"column:complete_num"`
	CloseNum    int64 `gorm:"column:close_num"`
	AbortNum    int64 `gorm:"column:abort_num"`
}

type spaceMemberSearchResult struct {
	SpaceId  int64 `gorm:"column:space_id"`
	MemberId int64 `gorm:"column:member_id"`
	UserId   int64 `gorm:"column:user_id"`
}

func (t *DwdToDwsMemberWitem1hTask) Run() {

	if t.status == pkg.TASK_STATUS_RUNNING {
		return
	}

	defer func() {
		t.status = pkg.TASK_STATUS_READY
	}()

	t.status = pkg.TASK_STATUS_RUNNING

	//获取任务状态变量
	jobVar, err := t.variablesRepo.GetVariablesByName(t.FullName(), "val")
	if err != nil {
		fmt.Println(err)
		return
	}

	var stateVal *DwdToDwsMemberWitem1hTaskStateValue
	json.Unmarshal([]byte(jobVar.VariableValue), &stateVal)
	if stateVal == nil {
		//如果重来没有生成过，那么就从前一天的开始生成
		preTime := time.Now().Add(-24 * time.Hour)
		lastDwdDate := time.Date(preTime.Year(), preTime.Month(), preTime.Day(), 0, 0, 0, 0, time.Local)
		stateVal = &DwdToDwsMemberWitem1hTaskStateValue{
			LastDwdDate: lastDwdDate.Format("2006-01-02 15:04:05"),
		}
		//保存一下
		jobVar.VariableValue = stateVal.String()
		err = t.variablesRepo.SaveVariables(jobVar)
		if err != nil {
			fmt.Println(err)
		}
	}

	//距离最后一次时间 到 (现在的时间-1小时）段 未被生成的数据
	bizNowTime := time.Now()
	lastDwdDate := date.ParseInLocation("2006-01-02 15:04:05", stateVal.LastDwdDate)
	subHours := bizNowTime.Sub(lastDwdDate).Hours()

	var lastEndTime time.Time

	for i := 0; i < int(subHours); i++ {

		bizStartTime := lastDwdDate.Add(time.Duration(i) * time.Hour)
		bizEndTime := bizStartTime.Add(time.Hour)

		err := t.run(bizStartTime, bizEndTime)
		if err != nil {
			fmt.Println(err)
			return
		}

		//最后一次的时间
		lastEndTime = bizEndTime
		//保存任务状态
		if !lastEndTime.IsZero() {
			stateVal.LastDwdDate = lastEndTime.Format("2006-01-02 15:04:05")
			jobVar.VariableValue = stateVal.String()
			err = t.variablesRepo.SaveVariables(jobVar)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func (t *DwdToDwsMemberWitem1hTask) run(startTime, endTime time.Time) error {
	spaceMemberMap, err := t.SpaceMemberMap(endTime)
	if err != nil {
		return err
	}

	aggResult, err := t.AggResult(endTime)
	if err != nil {
		return err
	}

	aggResultMap := stream.ToMap(aggResult, func(_ int, v *dws.DwsMbrWitem1h) (string, *dws.DwsMbrWitem1h) {
		return fmt.Sprintf("%d_%d", v.SpaceId, v.UserId), v
	})

	var finalResult []*dws.DwsMbrWitem1h
	//需要补0
	for _, v := range aggResultMap {
		finalResult = append(finalResult, &dws.DwsMbrWitem1h{
			SpaceId:     v.SpaceId,
			UserId:      v.UserId,
			Num:         v.Num,
			ExpireNum:   v.ExpireNum,
			TodoNum:     v.TodoNum,
			CloseNum:    v.CloseNum,
			AbortNum:    v.AbortNum,
			CompleteNum: v.CompleteNum,
			StartDate:   startTime,
			EndDate:     endTime,
		})
	}
	for _, members := range spaceMemberMap {
		for _, v := range members {
			key := fmt.Sprintf("%d_%d", v.SpaceId, v.UserId)
			if _, ok := aggResultMap[key]; ok {
				continue
			}

			finalResult = append(finalResult, &dws.DwsMbrWitem1h{
				SpaceId:   v.SpaceId,
				UserId:    v.UserId,
				StartDate: startTime,
				EndDate:   endTime,
			})
		}
	}

	//写入dws汇总表
	if len(finalResult) > 0 {
		err = t.data.Db().Transaction(func(tx *gorm.DB) error {
			//保存dws报表
			for _, v := range finalResult {
				res := tx.Model(v).
					Where("space_id=? AND user_id = ? AND end_date = ?", v.SpaceId, v.UserId, v.StartDate.Format("2006-01-02 15:04:05")).
					Where("num = ? AND expire_num = ? AND todo_num = ? AND close_num = ? AND abort_num = ? AND complete_num = ?", v.Num, v.ExpireNum, v.TodoNum, v.CloseNum, v.AbortNum, v.CompleteNum).
					Order("_id asc").Limit(1).
					Update("end_date", v.EndDate.Format("2006-01-02 15:04:05"))
				if res.Error != nil {
					return err
				}

				if res.RowsAffected == 0 {
					err = tx.Create(v).Error
					if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
						return err
					}
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return err
}

func (t *DwdToDwsMemberWitem1hTask) SpaceMemberMap(endDate time.Time) (map[int64][]*spaceMemberSearchResult, error) {
	end := endDate.Format("2006-01-02 15:04:05")

	const sql = `
SELECT * FROM (
	SELECT 
		*,
		ROW_NUMBER() over ( PARTITION BY space_id, user_id ORDER BY _id DESC ) AS ranking 
	FROM dwd_member
	WHERE start_date < ? AND end_date >= ? 
) AS t 
WHERE t.ranking = 1
`
	//获取符合时间段的空间成员
	var spaceMemberSearchResults []*spaceMemberSearchResult
	err := t.data.Db().Raw(sql, end, end).Scan(&spaceMemberSearchResults).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//按空间分组
	spaceMap := stream.GroupBy(spaceMemberSearchResults, func(v *spaceMemberSearchResult) int64 {
		return v.SpaceId
	})

	return spaceMap, nil
}

func (t *DwdToDwsMemberWitem1hTask) AggResult(endDate time.Time) ([]*dws.DwsMbrWitem1h, error) {
	end := endDate.Format("2006-01-02 15:04:05")

	const sql = `
SELECT 
    jt.order_id user_id,
	w.space_id,
    COUNT(w.work_item_id) AS num, -- 总任务数
    SUM(CASE WHEN stu.status_type = 2 AND w.plan_complete_at > 0 AND ? > w.plan_complete_at THEN 1 ELSE 0 END) AS expire_num, -- 过期任务数
    SUM(CASE WHEN stu.status_type = 2 AND JSON_CONTAINS(w.directors, CAST(CONCAT('"', jt.order_id, '"') AS JSON)) THEN 1 ELSE 0 END) AS todo_num, -- 待办任务数
    SUM(CASE WHEN stu.status_type = 3 AND stu.status_key NOT IN ('close', 'terminated') THEN 1 ELSE 0 END) AS complete_num, -- 已完成任务数
    SUM(CASE WHEN stu.status_key = 'close' THEN 1 ELSE 0 END) AS close_num, -- 已关闭任务数
    SUM(CASE WHEN stu.status_key = 'terminated' THEN 1 ELSE 0 END) AS abort_num -- 已终止任务数
FROM 
    dwd_witem w
JOIN JSON_TABLE(
    w.participators,            -- JSON 数组字段
    "$[*]"                      -- 遍历数组所有元素
    COLUMNS (
        order_id INT PATH "$"   -- 提取每个元素作为 order_id 列
    )
) AS jt
LEFT JOIN dim_witem_status stu ON w.status_id = stu.status_id
WHERE 
    w.start_date < ? 
    AND w.end_date >= ?
GROUP BY 
    jt.order_id, w.space_id
`

	var result []*dws.DwsMbrWitem1h
	err := t.data.Db().Raw(sql, endDate.Unix(), end, end).Scan(&result).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return result, nil
}

func (t *DwdToDwsMemberWitem1hTask) Stop() {
}
