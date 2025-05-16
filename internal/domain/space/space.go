package space

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"time"
)

type Space struct {
	shared.AggregateRoot

	Id          int64  `bson:"_id" json:"id"`
	UserId      int64  `bson:"user_id" json:"user_id"`
	SpaceGuid   string `bson:"space_guid" json:"space_guid"`
	SpaceName   string `bson:"space_name" json:"space_name"`
	SpaceStatus int32  `bson:"space_status" json:"space_status"`
	Remark      string `bson:"remark" json:"remark"`
	Describe    string `bson:"describe" json:"describe"`
	Notify      int64  `bson:"notify" json:"notify"`
	CreatedAt   int64  `bson:"created_at" json:"created_at"`
	UpdatedAt   int64  `bson:"updated_at" json:"updated_at"`
	DeletedAt   int64  `bson:"deleted_at" json:"deleted_at"`

	SpaceConfig *SpaceConfig
}

// 初始化默认配置
func (s *Space) InitDefaultConfig() {
	if s.SpaceConfig != nil {
		return
	}

	s.SpaceConfig = &SpaceConfig{
		SpaceId:    s.Id,
		WorkingDay: "[1,2,3,4,5]",
		CreatedAt:  time.Now().Unix(),
	}
}

func (s *Space) UpdateName(newName string, oper shared.Oper) {
	if s.SpaceName == newName {
		return
	}

	oldName := s.SpaceName
	s.SpaceName = newName

	s.AddDiff(Diff_SpaceName)

	s.AddMessage(oper, &domain_message.ModifySpace{
		SpaceId:   s.Id,
		SpaceName: oldName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "spaceName",
				OldValue: oldName,
				NewValue: newName,
			},
		},
	})
}

func (s *Space) UpdateDescribe(describe string, oper shared.Oper) {

	oldValue := s.Describe
	s.Describe = describe

	s.AddDiff(Diff_Describe)

	s.AddMessage(oper, &domain_message.ModifySpace{
		SpaceId:   s.Id,
		SpaceName: s.SpaceName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "describe",
				OldValue: oldValue,
				NewValue: describe,
			},
		},
	})

}

func (s *Space) TransferSpace(ownerId int64, oper shared.Oper) {

	if s.UserId == ownerId {
		return
	}

	oldValue := s.UserId

	s.UserId = ownerId
	s.AddDiff(Diff_UserId)

	s.AddMessage(oper, &domain_message.TransferSpace{
		SpaceId:      s.Id,
		SpaceName:    s.SpaceName,
		UserId:       oldValue,
		TargetUserId: ownerId,
	})
}

func (s *Space) UpdateNotify(notify int64, oper shared.Oper) {
	// if s.Notify == notify {
	// 	return
	// }

	s.Notify = notify

	s.AddDiff(Diff_Notify)

	s.AddMessage(oper, &domain_message.SetSpaceNotify{
		SpaceId:   s.Id,
		SpaceName: s.SpaceName,
		Notify:    int(notify),
	})
}

func (s *Space) IsCreator(uid int64) bool {
	return s.UserId == uid
}

func (s *Space) OnDelete(oper shared.Oper) {

	if s.DeletedAt > 0 {
		return
	}

	s.DeletedAt = time.Now().Unix()

	s.AddMessage(oper, &domain_message.DelSpace{
		SpaceId:   s.Id,
		SpaceName: s.SpaceName,
	})
}
