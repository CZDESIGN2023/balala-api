package dws

import "time"

type DwsMbrWitem1h struct {
	SpaceId     int64     `gorm:"column:space_id"`
	UserId      int64     `gorm:"column:user_id"`
	Num         int64     `gorm:"column:num"`
	ExpireNum   int64     `gorm:"column:expire_num"`
	TodoNum     int64     `gorm:"column:todo_num"`
	CompleteNum int64     `gorm:"column:complete_num"`
	CloseNum    int64     `gorm:"column:close_num"`
	AbortNum    int64     `gorm:"column:abort_num"`
	StartDate   time.Time `gorm:"column:start_date"`
	EndDate     time.Time `gorm:"column:end_date"`
}

func (DwsMbrWitem1h) TableName() string {
	return "dws_mbr_witem_1h"
}
