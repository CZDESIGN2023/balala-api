package dws

import "time"

type DwsVersWitem1h struct {
	Id          int64     `gorm:"column:_id" json:"_id"`
	SpaceId     int64     `gorm:"column:space_id" json:"space_id"`
	VersionId   int64     `gorm:"column:version_id" json:"version_id"`
	Num         int64     `gorm:"column:num" json:"num"`
	ExpireNum   int64     `gorm:"column:expire_num" json:"expire_num"`
	TodoNum     int64     `gorm:"column:todo_num" json:"todo_num"`
	CompleteNum int64     `gorm:"column:complete_num" json:"complete_num"`
	CloseNum    int64     `gorm:"column:close_num" json:"close_num"`
	AbortNum    int64     `gorm:"column:abort_num" json:"abort_num"`
	StartDate   time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate     time.Time `gorm:"column:end_date" json:"end_date"`
}

func (receiver *DwsVersWitem1h) TableName() string {
	return "dws_vers_witem_1h"
}
