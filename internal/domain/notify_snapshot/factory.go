package notify_snapshot

import "time"

func NewNotifySnapShot(spaceId int64, uid int64, typ NotifyEventType, doc string) *NotifySnapShot {
	return &NotifySnapShot{
		SpaceId:   spaceId,
		UserId:    uid,
		Typ:       typ,
		Doc:       doc,
		CreatedAt: time.Now().Unix(),
	}
}
