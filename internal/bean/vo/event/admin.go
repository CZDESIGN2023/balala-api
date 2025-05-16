package event

import (
	"go-cs/api/notify"
	"go-cs/internal/consts"
)

type AdminChangeUserNickname struct {
	Event    notify.Event
	Operator int64
	UserId   int64
	OldValue string
	NewValue string
}

type AdminChangeUserRole struct {
	Event    notify.Event
	Operator int64
	UserId   int64
	OldValue consts.SystemRole
	NewValue consts.SystemRole
}

type AdminCancelUser struct {
	Event    notify.Event
	Operator int64
	UserId   int64
}
