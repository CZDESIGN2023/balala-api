package space_tag

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

type TagStatus int32

var (
	TagStatus_Enable  TagStatus = 1
	TagStatus_Disable TagStatus = 0
)

type SpaceTag struct {
	shared.AggregateRoot

	Id        int64     `json:"id,omitempty"`
	SpaceId   int64     `json:"space_id,omitempty"`   //空间id
	TagGuid   string    `json:"tag_guid,omitempty"`   //tag唯一标识符
	TagName   string    `json:"tag_name,omitempty"`   //tag名称
	TagStatus TagStatus `json:"tag_status,omitempty"` //Tag状态;0:禁用,1:正常
	CreatedAt int64     `json:"created_at,omitempty"` //创建时间
	UpdatedAt int64     `json:"updated_at,omitempty"` //更新时间
	DeletedAt int64     `json:"deleted_at,omitempty"` //删除时间
}

func (s *SpaceTag) ChangeName(newName string, oper shared.Oper) {
	if s.TagName == newName {
		return
	}

	oldValue := s.TagName
	s.UpdateName(newName)

	s.AddMessage(oper, &domain_message.ModifySpaceTag{
		SpaceId:      s.SpaceId,
		SpaceTagId:   s.Id,
		SpaceTagName: s.TagName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: newName,
			},
		},
	})

}

func (s *SpaceTag) UpdateName(newName string) {
	if s.TagName == newName {
		return
	}

	s.TagName = newName
	s.AddDiff(Diff_TagName)
}

func (s *SpaceTag) UpdateStatus(status TagStatus) {
	if s.TagStatus == status {
		return
	}

	s.TagStatus = status
	s.AddDiff(Diff_TagStatus)
}

func (s *SpaceTag) OnDelete(oper shared.Oper) {

	if s.DeletedAt > 0 {
		return
	}

	s.DeletedAt = time.Now().Unix()

	s.AddMessage(oper, &domain_message.DeleteSpaceTag{
		SpaceId:      s.SpaceId,
		SpaceTagId:   s.Id,
		SpaceTagName: s.TagName,
	})
}
