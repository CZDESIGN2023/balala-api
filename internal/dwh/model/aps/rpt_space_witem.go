package aps

import "time"

type RptSpaceWitemQuery struct {
	//空间Id
	SpaceIds []int64 `json:"space_id"`
	//版本id
	UserId int64 `json:"user_id"`
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

type RptSpaceWitem struct {
	//空间Id
	SpaceId int64 `json:"space_id"`
	//总数量
	Num int64 `json:"num"`
	//逾期数量
	ExpireNum int64 `json:"expire_num"`
	//待办数量
	TodoNum int64 `json:"todo_num"`
	//完成任务数量
	CompleteNum int64 `json:"complete_num"`
	//关闭数量
	CloseNum int64 `json:"close_num"`
	//终止数量
	AbortNum int64 `json:"abort_num"`
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

func (r *RptSpaceWitem) GetStartDate() time.Time {
	return r.StartDate
}

func (r *RptSpaceWitem) GetEndDate() time.Time {
	return r.EndDate
}

func (r *RptSpaceWitem) SetStartDate(t time.Time) {
	r.StartDate = t
}

func (r *RptSpaceWitem) SetEndDate(t time.Time) {
	r.EndDate = t
}

func (r *RptSpaceWitem) Clone() *RptSpaceWitem {
	cpy := *r
	return &cpy
}
