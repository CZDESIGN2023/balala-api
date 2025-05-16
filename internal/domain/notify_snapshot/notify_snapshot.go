package notify_snapshot

import shared "go-cs/internal/pkg/domain"

type NotifyEventType int64

type NotifySnapShot struct {
	shared.AggregateRoot

	Id        int64           `json:"id"`
	SpaceId   int64           `json:"spaceId"`
	UserId    int64           `json:"userId"`
	Typ       NotifyEventType `json:"typ"`
	Doc       string          `json:"doc"`
	CreatedAt int64           `json:"createdAt"`
	UpdatedAt int64           `json:"updatedAt"`
	DeletedAt int64           `json:"deletedAt"`
}
