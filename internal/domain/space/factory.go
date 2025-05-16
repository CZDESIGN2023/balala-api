package space

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

func NewSpace(id int64, uid int64, name string, describe string, oper shared.Oper) *Space {
	space := &Space{
		Id:          id,
		UserId:      uid,
		SpaceGuid:   uuid.NewString(),
		SpaceName:   name,
		SpaceStatus: 1,
		Describe:    describe,
		Notify:      1,
		CreatedAt:   time.Now().Unix(),
	}
	space.InitDefaultConfig()

	space.AddMessage(oper, &domain_message.CreateSpace{
		SpaceId:   id,
		SpaceName: name,
	})

	return space
}
