package data

import (
	"context"
	"fmt"
	aps_model "go-cs/internal/dwh/model/aps"
)

// 公共 SQL 模板
const spaceIncrBaseSQL = `
SELECT
    s.space_id, -- 项目空间 ID
    CAST(DATE_FORMAT(s.gmt_create, '%s') AS DATETIME) AS start_date, -- 根据时间格式化规则转换时间戳
    COUNT(DISTINCT work_item_id) AS num -- 统计完成数量
FROM
    dwd_witem s
    JOIN dim_witem_status w ON s.status_id = w.status_id -- 关联状态表
WHERE
    s.space_id IN ?
  	AND s.end_date >= ? AND s.start_date < ? 
    AND s.gmt_create BETWEEN ? AND ?
GROUP BY
    s.space_id, -- 按项目空间分组
    start_date -- 按日期分组
ORDER BY
    s.space_id ASC, -- 按项目空间升序排序
    start_date ASC; -- 按日期升序排序
`

// querySpaceNewIncrData 通用查询方法
func (r *RptRepo) querySpaceNewIncrData(ctx context.Context, dateFormat string, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	var list []*aps_model.RptSpaceIncrWitem

	// 动态生成 SQL 查询
	sqlQuery := fmt.Sprintf(spaceIncrBaseSQL, dateFormat)

	// 执行 SQL 查询
	err := r.data.Db().Raw(
		sqlQuery,
		query.SpaceIds,
		query.StartDate, query.EndDate,
		query.StartDate, query.EndDate,
	).Scan(&list).Error

	if err != nil {
		// 记录错误日志或提供更多上下文信息（根据实际需求）
		return nil, fmt.Errorf("failed to query space increment data: %w", err)
	}

	return list, nil
}

// DashboardSpaceNewIncrDay 查询每日新增增量数据
func (r *RptRepo) DashboardSpaceNewIncrDay(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceNewIncrData(ctx, "%Y-%m-%d 00:00:00", query)
}

// DashboardSpaceNewIncrHour 查询每小时新增增量数据
func (r *RptRepo) DashboardSpaceNewIncrHour(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceNewIncrData(ctx, "%Y-%m-%d %H:00:00", query)
}

// Month
func (r *RptRepo) DashboardSpaceNewIncrMonth(ctx context.Context, query *aps_model.RptSpaceIncrWitemQuery) ([]*aps_model.RptSpaceIncrWitem, error) {
	return r.querySpaceNewIncrData(ctx, "%Y-%m-01 00:00:00", query)
}
