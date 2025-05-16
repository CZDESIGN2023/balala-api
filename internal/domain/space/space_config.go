package space

import (
	"encoding/json"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"slices"
	"strings"
	"time"
)

type WorkingDay string

func (w *WorkingDay) ToSlice() []int64 {
	var days []int64
	json.Unmarshal([]byte(*w), &days)
	return days
}

func (w *WorkingDay) IsWorkingDay(day time.Weekday) bool {
	var days []int64
	json.Unmarshal([]byte(*w), &days)
	return slices.Contains(days, int64(day))
}

type SpaceConfig struct {
	shared.DomainEntity

	Id                           int64
	SpaceId                      int64
	WorkingDay                   WorkingDay
	CommentDeletable             int64
	CommentDeletableWhenArchived int64
	CommentShowPos               int64
	CreatedAt                    int64
	UpdatedAt                    int64
	DeletedAt                    int64
}

func (s *SpaceConfig) UpdateWorkingDay(days []string, oper shared.Oper) {
	oldVal := s.WorkingDay
	s.WorkingDay = WorkingDay("[" + strings.Join(days, ",") + "]")
	s.UpdatedAt = time.Now().Unix()

	s.AddDiff(Diff_SpaceConfig_WorkingDay)

	s.AddMessage(oper, &domain_message.SetWorkingDay{
		SpaceId:     s.SpaceId,
		WeekDays:    s.WorkingDay.ToSlice(),
		OldWeekDays: oldVal.ToSlice(),
	})
}

func (s *SpaceConfig) UpdateCommentDeletable(val int64, oper shared.Oper) {
	s.CommentDeletable = val
	s.UpdatedAt = time.Now().Unix()

	s.AddDiff(Diff_SpaceConfig_CommentDeletable)

	s.AddMessage(oper, &domain_message.SetCommentDeletable{
		SpaceId:   s.SpaceId,
		Deletable: s.CommentDeletable,
	})
}

func (s *SpaceConfig) UpdateCommentDeletableWhenArchived(val int64, oper shared.Oper) {
	s.CommentDeletableWhenArchived = val
	s.UpdatedAt = time.Now().Unix()

	s.AddDiff(Diff_SpaceConfig_CommentDeletableWhenArchived)

	s.AddMessage(oper, &domain_message.SetCommentDeletableWhenArchived{
		SpaceId: s.SpaceId,
		Value:   s.CommentDeletableWhenArchived,
	})
}

func (s *SpaceConfig) UpdateCommentShowPos(val int64, oper shared.Oper) {
	s.CommentShowPos = val
	s.UpdatedAt = time.Now().Unix()

	s.AddDiff(Diff_SpaceConfig_CommentShowPos)

	s.AddMessage(oper, &domain_message.SetCommentShowPos{
		SpaceId: s.SpaceId,
		Value:   s.CommentShowPos,
	})
}
