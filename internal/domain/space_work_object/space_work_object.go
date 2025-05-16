package space_tag

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

type WorkObjectStatus int32

var (
	WorkObjectStatus_Enable  = 1
	WorkObjectStatus_Disable = 2
)

type SpaceWorkObjects []*SpaceWorkObject

type SpaceWorkObject struct {
	shared.AggregateRoot

	Id               int64            `json:"id,omitempty"`
	SpaceId          int64            `json:"space_id,omitempty"`           //空间id
	UserId           int64            `json:"user_id,omitempty"`            //创建用户id
	WorkObjectGuid   string           `json:"work_object_guid,omitempty"`   //空间Guid
	WorkObjectName   string           `json:"work_object_name,omitempty"`   //空间名称
	WorkObjectStatus WorkObjectStatus `json:"work_object_status,omitempty"` //空间状态;0:禁用,1:正常,2:未验证
	Remark           string           `json:"remark,omitempty"`             //备注
	Describe         string           `json:"describe,omitempty"`           //描述信息
	Ranking          int64            `json:"ranking,omitempty"`            //排序值
	CreatedAt        int64            `json:"created_at,omitempty"`         //创建时间
	UpdatedAt        int64            `json:"updated_at,omitempty"`         //更新时间
	DeletedAt        int64            `json:"deleted_at,omitempty"`         //删除时间
}

func (s *SpaceWorkObjects) GetMessages() shared.DomainMessages {
	logs := make(shared.DomainMessages, 0)
	for _, role := range *s {
		logs = append(logs, role.GetMessages()...)
	}
	return logs
}

func (s *SpaceWorkObject) UpdateName(name string, oper shared.Oper) {
	if s.WorkObjectName == name {
		return
	}

	oldValue := s.WorkObjectName

	s.WorkObjectName = name
	s.AddDiff(Diff_Name)

	s.AddMessage(oper, &domain_message.ModifyWorkObject{
		SpaceId:        s.SpaceId,
		WorkObjectId:   s.Id,
		WorkObjectName: s.WorkObjectName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "name",
				OldValue: oldValue,
				NewValue: name,
			},
		},
	})

}

func (s *SpaceWorkObject) UpdateStatus(status WorkObjectStatus) {
	if s.WorkObjectStatus == status {
		return
	}

	s.WorkObjectStatus = status
	s.AddDiff(Diff_Status)
}

func (s *SpaceWorkObject) UpdateRanking(ranking int64, oper shared.Oper) {
	if s.Ranking == ranking {
		return
	}

	oldValue := s.Ranking

	s.Ranking = ranking
	s.AddDiff(Diff_Ranking)

	s.AddMessage(oper, &domain_message.ModifyWorkObject{
		SpaceId:        s.SpaceId,
		WorkObjectId:   s.Id,
		WorkObjectName: s.WorkObjectName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "ranking",
				OldValue: oldValue,
				NewValue: s.Ranking,
			},
		},
	})
}

func (s *SpaceWorkObject) OnDelete(oper shared.Oper) {
	s.DeletedAt = time.Now().Unix()

	s.AddMessage(oper, &domain_message.DeleteWorkObject{
		SpaceId:        s.SpaceId,
		WorkObjectId:   s.Id,
		WorkObjectName: s.WorkObjectName,
	})
}
