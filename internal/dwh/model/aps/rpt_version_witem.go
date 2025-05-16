package aps

import "time"

type RptVersionWitemQuery struct {
	//空间Id
	SpaceIds []int64 `json:"space_id"`
	//版本id
	VersionId int64 `json:"version_id"`
	//任务状态id
	StatusId int64 `json:"status_id"`
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

type RptVersionWitem struct {
	//空间Id
	SpaceId int64 `json:"space_id"`
	//版本id
	VersionId int64 `json:"version_id"`
	//版本名称
	VersionName string `json:"version_name"`
	//任务状态id
	StatusId int64 `json:"status_id"`
	//任务状态名称
	StatusName string `json:"status_name"`
	// 总数量
	Num int64 `json:"num"`
	// 逾期数量
	ExpireNum int64 `json:"expire_num"`
	// 待办数量
	TodoNum int64 `json:"todo_num"`
	// 完成任务数量
	CompleteNum int64 `json:"complete_num"`
	// 关闭数量
	CloseNum int64 `json:"close_num"`
	// 终止数量
	AbortNum int64 `json:"abort_num"`
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

func (r *RptVersionWitem) GetStartDate() time.Time {
	return r.StartDate
}

func (r *RptVersionWitem) GetEndDate() time.Time {
	return r.EndDate
}

func (r *RptVersionWitem) SetStartDate(t time.Time) {
	r.StartDate = t
}

func (r *RptVersionWitem) SetEndDate(t time.Time) {
	r.EndDate = t
}

func (r *RptVersionWitem) Clone() *RptVersionWitem {
	cpy := *r
	return &cpy
}
