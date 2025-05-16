package space_view

import (
	"go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

type SpaceGlobalView struct {
	shared.AggregateRoot

	Id          int64
	SpaceId     int64
	Key         string
	Name        string
	Status      int64
	Type        int64
	QueryConfig string
	TableConfig string
	Ranking     int64
	CreatedAt   int64
	UpdatedAt   int64
}

func (s *SpaceGlobalView) SetName(name string, oper shared.Oper) {
	if s.Name == name {
		return
	}

	oldValue := s.Name

	s.Name = name
	s.AddDiff(Diff_Name)

	s.AddMessage(oper, &message.SetSpaceViewName{
		DomainMessageBase: shared.DomainMessageBase{},
		SpaceId:           s.SpaceId,
		ViewId:            s.Id,
		ViewType:          s.Type,
		ViewOldName:       oldValue,
		ViewNewName:       s.Name,
	})

	s.AddMessage(oper, &message.UpdateSpaceView{
		Field:    "name",
		SpaceId:  s.SpaceId,
		ViewId:   s.Id,
		ViewType: s.Type,
		ViewName: s.Name,
		ViewKey:  s.Key,
	})
}

func (s *SpaceGlobalView) SetQueryConfig(query string, oper shared.Oper) {
	s.AddMessage(oper, &message.UpdateSpaceView{
		SpaceId:  s.SpaceId,
		ViewId:   s.Id,
		ViewType: s.Type,
		ViewName: s.Name,
		ViewKey:  s.Key,
		Field:    "queryConfig",
	})

	if s.QueryConfig == query {
		return
	}

	s.QueryConfig = query
	s.AddDiff(Diff_QueryConfig)
}

func (s *SpaceGlobalView) SetTableConfig(tableConfig string, oper shared.Oper) {
	if s.TableConfig == tableConfig {
		return
	}

	s.TableConfig = tableConfig
	s.AddDiff(Diff_TableConfig)

	s.AddMessage(oper, &message.UpdateSpaceView{
		SpaceId:  s.SpaceId,
		ViewId:   s.Id,
		ViewType: s.Type,
		ViewName: s.Name,
		ViewKey:  s.Key,
		Field:    "tableConfig",
	})
}

func (s *SpaceGlobalView) CreateUserView(userId int64) *SpaceUserView {
	return &SpaceUserView{
		Key:     s.Key,
		Type:    s.Type,
		OuterId: s.Id,
		SpaceId: s.SpaceId,
		UserId:  userId,
		Status:  1,
	}
}

func (s *SpaceGlobalView) OnDelete(oper shared.Oper) {
	s.AddMessage(oper, &message.DeleteSpaceView{
		SpaceId:  s.SpaceId,
		ViewId:   s.Id,
		ViewType: s.Type,
		ViewName: s.Name,
	})
}
