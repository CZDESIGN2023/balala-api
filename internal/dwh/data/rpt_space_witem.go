package data

import (
	"context"
	"errors"
	"fmt"
	aps_model "go-cs/internal/dwh/model/aps"
	"go-cs/internal/utils"
	"strings"

	"gorm.io/gorm"
)

func (r *RptRepo) DashboardRptSpaceWitem1h(ctx context.Context, query *aps_model.RptSpaceWitemQuery) ([]*aps_model.RptSpaceWitem, error) {

	//时间范围是必须的
	querySql := `
SELECT 
    dws.space_id, 
    dws.start_date,
    dws.end_date,
    dws.num, 
    dws.expire_num, 
    dws.todo_num, 
    dws.complete_num, 
    dws.close_num, 
    dws.abort_num
FROM
	dws_space_witem_1h dws
WHERE {{WhereCondition}}
ORDER BY dws.space_id, dws.start_date ASC;
	`

	var rawWhereSql []string
	rawWhereSql = append(rawWhereSql, fmt.Sprintf(" ( dws.start_date < '%s' AND dws.end_date >= '%s' )",
		query.EndDate.Format("2006-01-02 15:04:05"), query.StartDate.Format("2006-01-02 15:04:05")))

	if len(query.SpaceIds) != 0 {
		rawWhereSql = append(rawWhereSql, fmt.Sprintf("AND dws.space_id IN (%v)", strings.Join(utils.ToStrArray(query.SpaceIds), ",")))
	}

	querySql = strings.Replace(querySql, "{{WhereCondition}}", strings.Join(rawWhereSql, " "), 1)

	var rows []*aps_model.RptSpaceWitem
	err := r.data.Db().Raw(querySql).Scan(&rows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return rows, nil
}

func (r *RptRepo) DashboardRptSpaceWitem1day(ctx context.Context, query *aps_model.RptSpaceWitemQuery) ([]*aps_model.RptSpaceWitem, error) {

	//时间范围是必须的
	querySql := `
SELECT 
    dws.space_id, 
    dws.start_date, 
    dws.end_date,
    dws.num, 
    dws.expire_num, 
    dws.todo_num, 
    dws.complete_num, 
    dws.close_num, 
    dws.abort_num
FROM (
    SELECT 
        space_id, 
        start_date, 
        end_date,
        num,
        expire_num,
        todo_num,
        complete_num,
        close_num,
        abort_num,
        ROW_NUMBER() OVER (PARTITION BY space_id, DATE(start_date) ORDER BY start_date DESC) AS rn
    FROM dws_space_witem_1h dws
    WHERE {{WhereCondition}}
) dws
WHERE dws.rn = 1
ORDER BY dws.start_date, dws.space_id ASC;
	`

	var rawWhereSql []string
	rawWhereSql = append(rawWhereSql, fmt.Sprintf(" ( dws.start_date < '%s' AND dws.end_date >= '%s' )",
		query.EndDate.Format("2006-01-02 15:04:05"), query.StartDate.Format("2006-01-02 15:04:05")))

	if len(query.SpaceIds) != 0 {
		rawWhereSql = append(rawWhereSql, fmt.Sprintf("AND dws.space_id IN (%v)", strings.Join(utils.ToStrArray(query.SpaceIds), ",")))
	}

	querySql = strings.Replace(querySql, "{{WhereCondition}}", strings.Join(rawWhereSql, " "), 1)
	var rows []*aps_model.RptSpaceWitem
	err := r.data.Db().Raw(querySql).Scan(&rows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return rows, nil
}

func (r *RptRepo) DashboardRptSpaceWitem1month(ctx context.Context, query *aps_model.RptSpaceWitemQuery) ([]*aps_model.RptSpaceWitem, error) {
	//时间范围是必须的
	querySql := `
SELECT 
    dws.space_id, 
    dws.start_date, 
    dws.end_date,
    dws.num, 
    dws.expire_num, 
    dws.todo_num, 
    dws.complete_num, 
    dws.close_num, 
    dws.abort_num
FROM (
    SELECT 
        space_id, 
        start_date, 
        end_date,
        num,
        expire_num,
        todo_num,
        complete_num,
        close_num,
        abort_num,
        ROW_NUMBER() OVER (PARTITION BY space_id, DATE_FORMAT(start_date, "%Y-%M") ORDER BY start_date DESC) AS rn
    FROM dws_space_witem_1h dws
    WHERE {{WhereCondition}}
) dws
WHERE dws.rn = 1
ORDER BY dws.start_date, dws.space_id ASC;
	`
	var rawWhereSql []string
	rawWhereSql = append(rawWhereSql, fmt.Sprintf(" ( dws.start_date < '%s' AND dws.end_date >= '%s' )",
		query.EndDate.Format("2006-01-02 15:04:05"), query.StartDate.Format("2006-01-02 15:04:05")))

	if len(query.SpaceIds) != 0 {
		rawWhereSql = append(rawWhereSql, fmt.Sprintf("AND dws.space_id IN (%v)", strings.Join(utils.ToStrArray(query.SpaceIds), ",")))
	}

	querySql = strings.Replace(querySql, "{{WhereCondition}}", strings.Join(rawWhereSql, " "), 1)

	var rows []*aps_model.RptSpaceWitem
	err := r.data.Db().Raw(querySql).Scan(&rows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return rows, nil
}
