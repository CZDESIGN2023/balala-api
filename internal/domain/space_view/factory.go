package space_view

import (
	shared "go-cs/internal/pkg/domain"
)

func NewSpaceUserView(
	spaceId int64,
	key string,
	name string,
	ranking int64,
	typ int64,
	outerId int64,
	queryConfig string,
	tableConfig string,
	oper shared.Oper,
) *SpaceUserView {
	ins := &SpaceUserView{
		UserId:      oper.GetId(),
		SpaceId:     spaceId,
		Key:         key,
		Name:        name,
		Ranking:     ranking,
		Type:        typ,
		OuterId:     outerId,
		QueryConfig: queryConfig,
		TableConfig: tableConfig,
		Status:      1,
	}

	//ins.AddMessage(oper, &domain_message.CreateWorkVersion{
	//	SpaceId:         spaceId,
	//	WorkVersionId:   id,
	//	WorkVersionName: versionName,
	//})

	return ins
}

func NewSpaceGlobalView(
	spaceId int64,
	key string,
	name string,
	ranking int64,
	typ int64,
	queryConfig string,
	tableConfig string,
	oper shared.Oper,
) *SpaceGlobalView {
	ins := &SpaceGlobalView{
		SpaceId:     spaceId,
		Key:         key,
		Name:        name,
		Ranking:     ranking,
		Type:        typ,
		QueryConfig: queryConfig,
		TableConfig: tableConfig,
		Status:      1,
	}

	return ins
}
