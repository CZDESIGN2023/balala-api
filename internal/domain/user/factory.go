package user

import (
	"go-cs/internal/consts"
	shared "go-cs/internal/pkg/domain"
	"time"
)

func NewUser(id int64, name string, nickName string, nickNamePy string, pwd string, pwdSalt string, avatar string, role consts.SystemRole, oper shared.Oper) *User {
	user := &User{
		Id:           id,
		UserName:     name,
		UserNickname: nickName,
		UserPinyin:   nickNamePy,
		UserPassword: pwd,
		UserSalt:     pwdSalt,
		UserStatus:   1,
		Avatar:       avatar,
		Role:         role,
		CreatedAt:    time.Now().Unix(),
	}

	return user
}
