package space_work_version

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

type SpaceWorkVersions []*SpaceWorkVersion

type SpaceWorkVersion struct {
	shared.AggregateRoot

	Id            int64  `json:"id"`
	SpaceId       int64  `json:"space_id"`
	VersionKey    string `json:"version_key"`
	VersionName   string `json:"version_name"`
	VersionStatus int64  `json:"version_status"`
	Remark        string `json:"remark"`
	Ranking       int64  `json:"ranking"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	DeletedAt     int64  `json:"deleted_at"`
}

func (s *SpaceWorkVersions) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)
	for _, role := range *s {
		logs = append(logs, role.GetMessages()...)
	}
	return logs
}

func (s *SpaceWorkVersion) UpdateRanking(ranking int64, oper shared.Oper) {
	if s.Ranking == ranking {
		return
	}

	oldValue := s.Ranking

	// 更新排序字段
	s.Ranking = ranking
	s.AddDiff(Diff_Ranking)

	s.AddMessage(oper, &domain_message.ModifyWorkVersion{
		SpaceId:         s.SpaceId,
		WorkVersionId:   s.Id,
		WorkVersionName: s.VersionName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "ranking",
				OldValue: oldValue,
				NewValue: s.Ranking,
			},
		},
	})
}

func (s *SpaceWorkVersion) UpdateName(name string, oper shared.Oper) {
	if s.VersionName == name {
		return
	}

	oldValue := s.VersionName

	s.VersionName = name
	s.AddDiff(Diff_VersionName)

	s.AddMessage(oper, &domain_message.ModifyWorkVersion{
		SpaceId:         s.SpaceId,
		WorkVersionId:   s.Id,
		WorkVersionName: s.VersionName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: name,
			},
		},
	})
}

func (s *SpaceWorkVersion) IsSameSpace(spaceId int64) bool {
	return s.SpaceId == spaceId
}

func (s *SpaceWorkVersion) OnDelete(oper shared.Oper) {

	if s.DeletedAt > 0 {
		return
	}

	s.DeletedAt = time.Now().Unix()

	s.AddMessage(oper, &domain_message.DeleteWorkVersion{
		SpaceId:         s.SpaceId,
		WorkVersionId:   s.Id,
		WorkVersionName: s.VersionName,
	})
}

func (s *SpaceWorkVersion) IsDefaultVersion() bool {
	return s.VersionKey == consts.DefaultWorkItemVersionKey
}
