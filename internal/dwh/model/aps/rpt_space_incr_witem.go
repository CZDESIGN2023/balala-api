package aps

import "time"

type RptSpaceIncrWitemQuery struct {
	//空间Id
	SpaceIds []int64 `json:"space_id"`
	//统计开始时间: 2024-08-20 00:00:00
	StartDate time.Time
	//统计结束时间: 2024-08-20 01:00:00
	EndDate time.Time
}

type RptSpaceIncrWitem struct {
	SpaceId   int64 `json:"space_id"`
	StartDate time.Time
	Num       int64
}
