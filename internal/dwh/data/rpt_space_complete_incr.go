package data

import (
	"context"
	aps_model "go-cs/internal/dwh/model/aps"
)

// 定义通用 SQL 模板
const spaceCompleteIncrSqlTemplate = `
SELECT
    s.space_id, -- 项目空间 ID
    CAST(FROM_UNIXTIME(s.last_status_at, ?) AS DATETIME) AS start_date, -- 动态时间格式化
    COUNT(DISTINCT work_item_id) AS num -- 统计完成数量
    
FROM
    dwd_witem s
    JOIN dim_witem_status w ON s.status_id = w.status_id -- 关联状态表
    
WHERE
    s.space_id IN ?
    AND s.end_date >= ? AND s.start_date < ?
  	AND s.last_status_at Between ? AND ?
    AND (
        (
            w.status_type = 3 
            AND w.flow_scope = 'state_flow' 
        ) 
        OR w.status_key = 'completed'
    )
GROUP BY
    s.space_id, -- 按项目空间分组
    start_date -- 按日期分组

ORDER BY
    s.space_id ASC, -- 按项目空间升序排序
    start_date ASC; -- 按日期升序排序
`

// 通用查询方法
func (r *RptRepo) querySpaceCompleteIncr(ctx context.Context, dateFormat string, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	var list []*aps_model.RptSpaceIncrWitem

	// 使用动态 SQL 模板执行查询
	err := r.data.Db().Raw(spaceCompleteIncrSqlTemplate,
		dateFormat,
		query.SpaceIds,
		query.StartDate, query.EndDate,
		query.StartDate.Unix(), query.EndDate.Unix(),
	).Scan(&list).Error
	return list, err
}

// 按天统计完成增量
func (r *RptRepo) DashboardSpaceCompleteIncrDay(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceCompleteIncr(ctx, "%Y-%m-%d", query)
}

// 按小时统计完成增量
func (r *RptRepo) DashboardSpaceCompleteIncrHour(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceCompleteIncr(ctx, "%Y-%m-%d %H:00:00", query)
}

// 月
func (r *RptRepo) DashboardSpaceCompleteIncrMonth(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceCompleteIncr(ctx, "%Y-%m-01 00:00:00", query)
}
