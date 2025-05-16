package domain

import (
	shared "go-cs/internal/pkg/domain"
)

type LoginStauts int32

var (
	LoginSuccess LoginStauts = 1
	LoginFailed  LoginStauts = 2
)

type LoginUser struct {
	LoginUserId       int64  ` bson:"login_user_id" json:"login_user_id"`
	LoginUserName     string ` bson:"login_user_name" json:"login_user_name"`
	LoginUserNickname string ` bson:"login_user_nickname" json:"login_user_nickname"`
}

type LoginIp struct {
	IpAddr   string ` bson:"ip_addr" json:"ip_addr"`
	Location string ` bson:"location" json:"location"`
}

type UserLoginLog struct {
	shared.AggregateRoot

	Id          int64       ` bson:"_id" json:"id"`
	LoginUser   LoginUser   ` bson:"login_user" json:"login_user"`
	LoginIp     LoginIp     ` bson:"ipaddr" json:"ipaddr"`
	Browser     string      ` bson:"browser" json:"browser"`
	Os          string      ` bson:"os" json:"os"`
	LoginStauts LoginStauts ` bson:"status" json:"status"`
	Msg         string      ` bson:"msg" json:"msg"`
	LoginAt     int64       ` bson:"login_at" json:"login_at"`
	CreatedAt   int64       ` bson:"created_at" json:"created_at"`
	UpdatedAt   int64       ` bson:"updated_at" json:"updated_at"`
	DeletedAt   int64       ` bson:"deleted_at" json:"deleted_at"`
}
