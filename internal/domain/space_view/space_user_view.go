package space_view

import (
	"go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

type SpaceViewUsers []*SpaceUserView

type SpaceUserView struct {
	shared.AggregateRoot

	Id          int64
	UserId      int64
	SpaceId     int64
	Key         string
	Name        string
	Status      int64
	Type        int64
	OuterId     int64
	QueryConfig string
	TableConfig string
	Ranking     int64
	CreatedAt   int64
	UpdatedAt   int64

	globalView *SpaceGlobalView
}

func (s *SpaceUserView) SetGlobalView(globalView *SpaceGlobalView) {
	if globalView == nil {
		return
	}

	s.globalView = globalView
	s.Name = globalView.Name
	s.QueryConfig = globalView.QueryConfig
	s.TableConfig = globalView.TableConfig

	//if s.TableConfig == "" {
	//}
}

func (s *SpaceUserView) GetGlobalView() *SpaceGlobalView {
	return s.globalView
}

func (s *SpaceViewUsers) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)
	for _, role := range *s {
		logs = append(logs, role.GetMessages()...)
	}
	return logs
}

func (s *SpaceUserView) UpdateRanking(ranking int64, oper shared.Oper) {
	if s.Ranking == ranking {
		return
	}

	s.Ranking = ranking
	s.AddDiff(Diff_Ranking)
}

func (s *SpaceUserView) SetName(name string, oper shared.Oper) {
	if s.Name == name {
		return
	}

	oldValue := s.Name

	s.Name = name
	s.AddDiff(Diff_Name)

	s.AddMessage(oper, &message.SetSpaceViewName{
		SpaceId:     s.SpaceId,
		ViewId:      s.Id,
		ViewType:    s.Type,
		ViewOldName: oldValue,
		ViewNewName: s.Name,
	})
}

func (s *SpaceUserView) SetStatus(status int64, oper shared.Oper) {
	if s.Status == status {
		return
	}

	s.Status = status

	s.AddDiff(Diff_Status)

	s.AddMessage(oper, &message.SetSpaceViewStatus{
		SpaceId:  s.SpaceId,
		ViewType: s.Type,
		ViewName: s.Name,
		Status:   status,
	})
}

func (s *SpaceUserView) SetQueryConfig(query string) {
	if s.QueryConfig == query {
		return
	}

	s.QueryConfig = query
	s.AddDiff(Diff_QueryConfig)
}

func (s *SpaceUserView) SetTableConfig(tableConfig string) {
	if s.TableConfig == tableConfig {
		return
	}

	s.TableConfig = tableConfig
	s.AddDiff(Diff_TableConfig)
}

func (s *SpaceUserView) OnDelete(oper shared.Oper) {
	s.AddMessage(oper, &message.DeleteSpaceView{
		SpaceId:  s.SpaceId,
		ViewId:   s.Id,
		ViewType: s.Type,
		ViewName: s.Name,
	})
}
