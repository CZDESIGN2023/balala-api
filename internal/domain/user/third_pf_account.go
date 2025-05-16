package user

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
)

type ThirdPfInfo struct {
	PfCode        int32  ` bson:"pf_code" json:"pf_code"`
	PfName        string ` bson:"pf_name" json:"pf_name"`
	PfUserKey     string ` bson:"pf_user_key" json:"pf_user_key"`
	PfUserName    string ` bson:"pf_user_name" json:"pf_user_name"`
	PfUserId      int64  ` bson:"pf_user_id" json:"pf_user_id"`
	PfUserAccount string ` bson:"pf_user_account" json:"pf_user_account"`
}

// 数据库表名: third_pf_account
type ThirdPfAccount struct {
	shared.AggregateRoot

	Id        int64       ` bson:"_id" json:"id"`
	UserId    int64       ` bson:"user_id" json:"user_id"`
	PfInfo    ThirdPfInfo `bson:"pf_info" json:"pf_info"`
	Notify    int32       ` bson:"notify" json:"notify"`
	CreatedAt int64       ` bson:"created_at" json:"created_at"`
	UpdatedAt int64       ` bson:"updated_at" json:"updated_at"`
	DeletedAt int64       ` bson:"deleted_at" json:"deleted_at"`
}

func (t *ThirdPfAccount) SetNotify(notify int32, oper shared.Oper) {
	t.Notify = notify

	t.AddDiff(Diff_ThirdPlatform_Notify)

	t.AddMessage(oper, &domain_message.PersonalSetThirdPlatformNotify{
		PlatformCode: t.PfInfo.PfCode,
		Notify:       notify,
	})
}
