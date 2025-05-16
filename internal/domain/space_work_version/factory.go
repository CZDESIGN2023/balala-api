package space_work_version

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

func NewSpaceWorkVersion(
	id int64,
	spaceId int64,
	versionKey string,
	versionName string,
	ranking int64,
	oper shared.Oper,
) *SpaceWorkVersion {
	ins := &SpaceWorkVersion{
		Id:          id,
		SpaceId:     spaceId,
		VersionKey:  versionKey,
		VersionName: versionName,
		Ranking:     ranking,
	}

	ins.AddMessage(oper, &domain_message.CreateWorkVersion{
		SpaceId:         spaceId,
		WorkVersionId:   id,
		WorkVersionName: versionName,
	})

	return ins
}
