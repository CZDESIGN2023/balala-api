package dwd

import (
	"go-cs/internal/dwh/pkg/model"
	"reflect"
)

type DwdMember struct {
	model.DimModel
	model.ChainModel

	MemberId int64 `json:"member_id" gorm:"column:member_id"`
	SpaceId  int64 `json:"plan_start_at" gorm:"column:space_id"`
	UserId   int64 `json:"user_id" gorm:"column:user_id"`
	RoleId   int64 `json:"role_id" gorm:"column:role_id"`
}

func (w *DwdMember) DeepEqual(y *DwdMember) bool {
	return reflect.DeepEqual(w.KeyValue(), y.KeyValue())
}

func (w *DwdMember) KeyValue() map[string]interface{} {

	return map[string]interface{}{
		"member_id": w.MemberId,
		"space_id":  w.SpaceId,
		"user_id":   w.UserId,
		"role_id":   w.RoleId,
	}
}
