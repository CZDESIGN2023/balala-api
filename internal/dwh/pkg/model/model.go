package model

import "time"

type OdsModel struct {
	OdsId   int64 `protobuf:"zigzag64,1,opt,name=_id,proto3" json:"_id" gorm:"column:_id"`
	OdsOpTs int64 `protobuf:"zigzag64,1,opt,name=_op_ts,proto3" json:"_op_ts" gorm:"column:_op_ts"`
}

type DimModel struct {
	GmtCreate   time.Time `json:"gmt_create"`
	GmtModified time.Time `json:"gmt_modified"`
}

type DwdModel struct {
	GmtCreate   time.Time `json:"gmt_create"`
	GmtModified time.Time `json:"gmt_modified"`
}

// 拉链表
type ChainModel struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
