package data

import (
	"context"
	aps_model "go-cs/internal/dwh/model/aps"
	"gorm.io/gorm"
	"log"
)

// 定义差集查询结果的结构体
type TaskDiffResult struct {
	SpaceID     int64  `json:"space_id"`
	UserID      int64  `json:"user_id"`
	Num         int64  `json:"num"`
	DiffTaskIDs string `json:"diff_task_ids"`
}

const memberNewIncrSQL = `
WITH DateFilteredItems AS (
    SELECT *,
           ROW_NUMBER() OVER (PARTITION BY work_item_id ORDER BY _id DESC) AS rn,
           ROW_NUMBER() OVER (PARTITION BY work_item_id ORDER BY _id ASC) AS rn2
    FROM dwd_witem
    WHERE space_id IN ?
      AND end_date >= ?
      AND start_date < ?
),
RecentTasks AS (
    SELECT space_id, participator, work_item_id
    FROM DateFilteredItems
    CROSS JOIN JSON_TABLE(participators, "$[*]" 
              COLUMNS (participator VARCHAR(50) PATH '$')) AS jt
    WHERE rn = 1 
      AND end_date >= ?  -- 保证任务在这个时间点前没有被删除
),
OldTasks AS (
    SELECT space_id, participator, work_item_id
    FROM DateFilteredItems
    CROSS JOIN JSON_TABLE(participators, "$[*]" 
              COLUMNS (participator VARCHAR(50) PATH '$')) AS jt
    WHERE rn2 = 1 
      AND start_date < ?
)
SELECT 
    r.space_id, 
    r.participator AS user_id, 
    COUNT(*) AS num
-- ,JSON_ARRAYAGG(r.work_item_id) AS task_ids                   -- 聚合任务 ID（JSON 数组）
FROM RecentTasks r
LEFT JOIN OldTasks o
    ON r.space_id = o.space_id
    AND r.participator = o.participator
    AND r.work_item_id = o.work_item_id
WHERE o.work_item_id IS NULL
GROUP BY r.space_id, r.participator;
`

// 定义模板参数结构体
type SubQueryParams struct {
	StartDate string // 开始日期
	EndDate   string // 结束日期
}

// 通用查询方法
func (r *RptRepo) queryMemberNewIncr(ctx context.Context, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	startDate := query.StartDate.Format("2006-01-02 15:04:05")
	endDate := query.EndDate.Format("2006-01-02 15:04:05")

	finalSQL := r.data.Db().ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(memberNewIncrSQL,
			query.SpaceIds, startDate, endDate,
			endDate,
			startDate,
		).Find(&[]TaskDiffResult{})
	})

	//fmt.Print(finalSQL)

	var results []TaskDiffResult
	err := r.data.Db().Raw(finalSQL).Scan(&results).Error
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	var list []*aps_model.RptMemberIncrWitem

	for _, result := range results {
		list = append(list, &aps_model.RptMemberIncrWitem{
			SpaceId: result.SpaceID,
			UserId:  result.UserID,
			Num:     result.Num,
		})
	}

	return list, err
}

// 月
func (r *RptRepo) DashboardMemberNewIncr(ctx context.Context, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	return r.queryMemberNewIncr(ctx, query)
}
