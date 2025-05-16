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

type DwdToDwsVersionWitem1hTaskStateValue struct {
	LastDwdDate string `json:"last_dwd_date"`
}

func (t *DwdToDwsVersionWitem1hTaskStateValue) String() string {
	v, _ := json.Marshal(t)
	return string(v)
}

// 每小时报表 空间版本任务轻度汇总
// 任务状态变量: 最后一次生成dws的时间 last_dws_date 任务被执行的时间
type DwdToDwsVersionWitem1hTask struct {
	id     string
	job    pkg.Job
	status string

	data          *data.DwhData
	variablesRepo *data.JobVariablesRepo
}

func NewDwdToDwsVersionWitem1hTask(
	id string,
	ctx *pkg.TaskContext,
) *DwdToDwsVersionWitem1hTask {
	return &DwdToDwsVersionWitem1hTask{
		id:            id,
		job:           ctx.Job,
		data:          ctx.Data,
		variablesRepo: ctx.JobVariablesRepo,
		status:        pkg.TASK_STATUS_READY,
	}
}

func (t *DwdToDwsVersionWitem1hTask) Id() string {
	return t.id
}

func (t *DwdToDwsVersionWitem1hTask) Name() string {
	return "dwd_to_dws_version_witem_1h"
}

func (t *DwdToDwsVersionWitem1hTask) FullName() string {
	if t.job != nil {
		return t.job.FullName() + ":" + t.Name() + ":" + t.Id()
	}
	return t.Name() + ":" + t.Id()
}

func (t *DwdToDwsVersionWitem1hTask) Status() string {
	return t.status
}

type versionAggSearchResult struct {
	SpaceId     int64 `gorm:"column:space_id" json:"space_id"`
	VersionId   int64 `gorm:"column:version_id" json:"version_id"`
	Num         int64 `gorm:"column:num" json:"num"`
	ExpireNum   int64 `gorm:"column:expire_num" json:"expire_num"`
	TodoNum     int64 `gorm:"column:todo_num" json:"todo_num"`
	CompleteNum int64 `gorm:"column:complete_num" json:"complete_num"`
	CloseNum    int64 `gorm:"column:close_num" json:"close_num"`
	AbortNum    int64 `gorm:"column:abort_num" json:"abort_num"`
}

//--统计数据

type versionSearchResult struct {
	SpaceId   int64 `gorm:"column:space_id" json:"space_id"`
	VersionId int64 `gorm:"column:version_id" json:"version_id"`
}

func (t *DwdToDwsVersionWitem1hTask) Run() {

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

	var stateVal *DwdToDwsVersionWitem1hTaskStateValue
	json.Unmarshal([]byte(jobVar.VariableValue), &stateVal)
	if stateVal == nil {
		//如果重来没有生成过，那么就从前一天的开始生成
		preTime := time.Now().Add(-24 * time.Hour)
		lastDwdDate := time.Date(preTime.Year(), preTime.Month(), preTime.Day(), 0, 0, 0, 0, time.Local)
		stateVal = &DwdToDwsVersionWitem1hTaskStateValue{
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

		err = t.run(bizStartTime, bizEndTime)
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

func (t *DwdToDwsVersionWitem1hTask) run(startTime, endTime time.Time) error {
	spaceVersionMap, err := t.SpaceVersionMap(endTime)
	if err != nil {
		return err
	}

	aggResultMap, err := t.VersionAggResultMap(endTime)
	if err != nil {
		return err
	}

	//需要补0
	var result []*dws.DwsVersWitem1h
	for _, v := range aggResultMap {
		result = append(result, &dws.DwsVersWitem1h{
			SpaceId:     v.SpaceId,
			VersionId:   v.VersionId,
			StartDate:   startTime,
			EndDate:     endTime,
			Num:         v.Num,
			ExpireNum:   v.ExpireNum,
			TodoNum:     v.TodoNum,
			CloseNum:    v.CloseNum,
			AbortNum:    v.AbortNum,
			CompleteNum: v.CompleteNum,
		})
	}
	for _, versions := range spaceVersionMap {
		for _, version := range versions {
			if _, ok := aggResultMap[version.VersionId]; ok {
				continue
			}

			result = append(result, &dws.DwsVersWitem1h{
				SpaceId:   version.SpaceId,
				VersionId: version.VersionId,
				StartDate: startTime,
				EndDate:   endTime,
			})
		}
	}

	// 统计项目数据
	spaceAggResultMap := stream.GroupBy(result, func(v *dws.DwsVersWitem1h) int64 {
		return v.SpaceId
	})
	var spaceResult []*dws.DwsSpaceWitem1h
	for id, v := range spaceAggResultMap {
		dwsR := &dws.DwsSpaceWitem1h{
			SpaceId:   id,
			StartDate: startTime,
			EndDate:   endTime,
		}

		for _, version := range v {
			dwsR.Num += version.Num
			dwsR.ExpireNum += version.ExpireNum
			dwsR.TodoNum += version.TodoNum
			dwsR.CloseNum += version.CloseNum
			dwsR.AbortNum += version.AbortNum
			dwsR.CompleteNum += version.CompleteNum
		}
		spaceResult = append(spaceResult, dwsR)
	}

	err = t.data.Db().Transaction(func(tx *gorm.DB) error {
		//保存dws报表
		for _, v := range result {
			res := tx.Model(v).
				Where("space_id=? AND version_id = ? AND end_date = ?", v.SpaceId, v.VersionId, v.StartDate.Format("2006-01-02 15:04:05")).
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

		for _, v := range spaceResult {
			res := tx.Model(v).
				Where("space_id=? AND end_date = ?", v.SpaceId, v.StartDate.Format("2006-01-02 15:04:05")).
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

	return err
}

func (t *DwdToDwsVersionWitem1hTask) latestVersionRecord(spaceId int64) ([]*dws.DwsVersWitem1h, error) {
	const latestVersionRecordSql = `
SELECT*FROM (
	SELECT*,ROW_NUMBER() over (PARTITION BY version_id ORDER BY _id DESC) AS ranking FROM dws_vers_witem_1h WHERE space_id=?
) AS t 
WHERE t.ranking=1
`

	var latestVersionRecord []*dws.DwsVersWitem1h
	err := t.data.Db().Raw(latestVersionRecordSql, spaceId).Scan(&latestVersionRecord).Error
	if err != nil {
		return nil, err
	}

	return latestVersionRecord, nil
}

func (t *DwdToDwsVersionWitem1hTask) SpaceVersionMap(endTime time.Time) (map[int64][]*versionSearchResult, error) {
	const bizVersionSearchSql = `
SELECT * FROM (
	SELECT 
		*,
		ROW_NUMBER() over ( PARTITION BY version_id ORDER BY _id DESC ) AS ranking 
	FROM dim_version 
	WHERE start_date <= ? AND end_date >= ? 
) AS t 
WHERE t.ranking = 1
		`

	end := endTime.Format("2006-01-02 15:04:05")

	var versions []*versionSearchResult
	err := t.data.Db().Raw(bizVersionSearchSql, end, end).Scan(&versions).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	spaceVersionMap := stream.GroupBy(versions, func(v *versionSearchResult) int64 {
		return v.SpaceId
	})

	return spaceVersionMap, nil
}

func (t *DwdToDwsVersionWitem1hTask) VersionAggResultMap(endTime time.Time) (map[int64]*versionAggSearchResult, error) {
	end := endTime.Format("2006-01-02 15:04:05")

	//符合时间段的最后状态数据
	bizSql := fmt.Sprintf(`
SELECT
    d.*,
    IF(stu.status_type IS NULL, 2, stu.status_type) AS status_type,
    IF(stu.status_key IS NULL, 'progressing', stu.status_key) AS status_key
FROM (
    SELECT 
        *,
        ROW_NUMBER() OVER (PARTITION BY work_item_id ORDER BY end_date DESC) AS rn
    FROM dwd_witem
    WHERE start_date < '%s' 
      AND end_date >= '%s'
) d
LEFT JOIN dim_witem_status stu 
    ON d.status_id = stu.status_id
WHERE d.rn = 1
		`, end, end)

	//汇总 每小时 纬度：空间|版本|状态 指标:任务总数|过期任务总数
	aggSql := fmt.Sprintf(`
			SELECT
			 space_id, 
			 version_id,
			 COUNT(work_item_id) num,
			 SUM(IF( status_type <> 3 AND ( %v > plan_complete_at ) ,1,0)) expire_num,
			 SUM( IF(status_type = 2,1,0) ) todo_num,
			 SUM( IF(status_type = 3 AND status_key != 'close' AND status_key != 'terminated' ,1,0) ) complete_num,
			 SUM( IF(status_key = 'close',1,0) ) close_num,
			 SUM( IF(status_key = 'terminated',1,0) ) abort_num
			FROM
			 (%v) dwd_witem
			GROUP BY 
			 space_id, version_id
		`, time.Now().Unix(), bizSql)

	var aggResult []*versionAggSearchResult
	err := t.data.Db().Raw(aggSql).Scan(&aggResult).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	aggResultMap := stream.ToMap(aggResult, func(_ int, v *versionAggSearchResult) (int64, *versionAggSearchResult) {
		return v.VersionId, v
	})

	return aggResultMap, nil
}

func (t *DwdToDwsVersionWitem1hTask) Stop() {
}
