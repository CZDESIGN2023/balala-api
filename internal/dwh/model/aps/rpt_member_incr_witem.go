package aps

import "time"

type RptMemberIncrWitemQuery struct {
	//空间Id
	SpaceIds []int64
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

type RptMemberIncrWitem struct {
	SpaceId   int64
	UserId    int64
	StartDate time.Time
	Num       int64
}
