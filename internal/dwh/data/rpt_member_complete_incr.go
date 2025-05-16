package data

import (
	"context"
	"fmt"
	aps_model "go-cs/internal/dwh/model/aps"
	"gorm.io/gorm"
)

// 定义通用 SQL 模板
const memberCompleteIncrSqlTemplate = `
WITH DistinctData AS (
    SELECT 
        s.space_id,                                               -- 项目空间 ID
        CAST(FROM_UNIXTIME(s.last_status_at, ?) AS DATETIME) AS start_date,  -- 将时间戳转换为日期格式
        jt.participator AS user_id,                               -- 提取参与人
        s.work_item_id                                             -- 任务 ID
    FROM 
        dwd_witem s
    JOIN dim_witem_status w 
        ON s.status_id = w.status_id                               -- 关联状态表
    CROSS JOIN 
        JSON_TABLE(
            s.participators,                                       -- JSON 数组字段
            "$[*]"                                                -- 展开参与人数组
            COLUMNS (
                participator VARCHAR(50) PATH '$'
            )
        ) AS jt
    WHERE 
        s.space_id IN (?)
        AND s.end_date >= ?
		AND s.start_date < ? 
        AND s.last_status_at BETWEEN ? AND ?
        AND (
			w.status_key = 'completed'                          -- 或者直接匹配 key 为 "completed"（条件1）
            OR (w.status_type = 3 AND w.flow_scope = 'state_flow')    -- 筛选完成状态的任务（条件2）
        )
),
DeduplicatedData AS (
	SELECT DISTINCT
			space_id,
			start_date,
			user_id,
			work_item_id
	FROM DistinctData
)
SELECT 
	space_id,
	start_date,
	user_id,
	COUNT(work_item_id) AS num                      -- 统计完成数量
	-- ,JSON_ARRAYAGG(work_item_id) AS task_ids                   -- 聚合任务 ID（JSON 数组）
FROM DeduplicatedData
GROUP BY 
	space_id, 
	start_date, 
	user_id
ORDER BY 
    space_id ASC,                                                 -- 按项目空间升序排序
    user_id ASC,                                                  -- 按参与人升序排序
    start_date ASC;                                               -- 按日期升序排序
`

// 通用查询方法
func (r *RptRepo) queryMemberCompleteIncr(ctx context.Context, dateFormat string, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	var list []*aps_model.RptMemberIncrWitem

	finalSQL := r.data.Db().ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(memberCompleteIncrSqlTemplate,
			dateFormat,
			query.SpaceIds,
			query.EndDate, query.EndDate,
			query.StartDate.Unix(), query.EndDate.Unix(),
		).Find(&[]*aps_model.RptMemberIncrWitem{})
	})

	fmt.Print(finalSQL)

	err := r.data.Db().Raw(finalSQL).Scan(&list).Error
	return list, err
}

// 按天统计完成增量
func (r *RptRepo) DashboardMemberCompleteIncrDay(ctx context.Context, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	return r.queryMemberCompleteIncr(ctx, "%Y-%m-%d", query)
}

// 按小时统计完成增量
func (r *RptRepo) DashboardMemberCompleteIncrHour(ctx context.Context, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	return r.queryMemberCompleteIncr(ctx, "%Y-%m-%d %H:00:00", query)
}

// 月
func (r *RptRepo) DashboardMemberCompleteIncrMonth(ctx context.Context, query *aps_model.RptMemberIncrWitemQuery) ([]*aps_model.RptMemberIncrWitem, error) {
	return r.queryMemberCompleteIncr(ctx, "%Y-%m-01 00:00:00", query)
}
