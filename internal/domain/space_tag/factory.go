package space_tag

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
)

func NewSpaceTag(
	id int64,
	spaceId int64,
	tagName string,
	oper shared.Oper,
) *SpaceTag {
	ins := &SpaceTag{
		Id:        id,
		TagGuid:   uuid.NewString(),
		SpaceId:   spaceId,
		TagName:   tagName,
		TagStatus: TagStatus_Enable,
		CreatedAt: time.Now().Unix(),
	}

	ins.AddMessage(oper, &domain_message.CreateSpaceTag{
		SpaceId:      spaceId,
		SpaceTagName: tagName,
		SpaceTagId:   id,
	})

	return ins
}
