package space_tag

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

func NewSpaceWorkObject(
	id int64,
	spaceId int64,
	name string,
	ranking int64,
	uid int64,
	oper shared.Oper,
) *SpaceWorkObject {
	ins := &SpaceWorkObject{
		Id:               id,
		WorkObjectGuid:   uuid.NewString(),
		SpaceId:          spaceId,
		WorkObjectName:   name,
		WorkObjectStatus: WorkObjectStatus(WorkObjectStatus_Enable),
		CreatedAt:        time.Now().Unix(),
		Ranking:          ranking,
		UserId:           uid,
	}

	ins.AddMessage(oper, &domain_message.CreateWorkObject{
		SpaceId:        spaceId,
		WorkObjectName: name,
		WorkObjectId:   id,
	})

	return ins
}
