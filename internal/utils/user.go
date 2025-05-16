package utils

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"go-cs/internal/utils/auth"
	"go-cs/internal/utils/oper"
	"math/rand"
)

type LoginUserInfo struct {
	UserId     int64
	JwtTokenId string
}

func (user *LoginUserInfo) GetType() oper.OperatorType {
	return oper.OperatorTypeUser
}

func (user *LoginUserInfo) GetId() int64 {
	return user.UserId
}

func GetLoginUserInfo(ctx context.Context) (*LoginUserInfo, error) {
	info := GetLoginUser(ctx)
	if info == nil || info.UserId == 0 {
		return nil, errors.New("GetLoginUserInfo: 查无用户信息")
	}

	return info, nil
}

func GetLoginUser(ctx context.Context) (info *LoginUserInfo) {
	tokenInfo, ok := auth.FromContext(ctx)
	if !ok {
		return &LoginUserInfo{}
	}

	return &LoginUserInfo{
		UserId:     tokenInfo.UserId,
		JwtTokenId: tokenInfo.JwtTokenId,
	}
}

func GetLoginUserId(ctx context.Context) int64 {
	return GetLoginUser(ctx).UserId
}

// EncryptPassword 密码加密
func EncryptUserPassword(password, salt string) string {
	md5PasswordStr := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	md5SaltStr := fmt.Sprintf("%x", md5.Sum([]byte(salt)))
	encPwd := fmt.Sprintf("%x", md5.Sum([]byte(md5PasswordStr+md5SaltStr)))
	return encPwd
}

func GenerateUserName(prefix string) string {

	length := rand.Intn(15-len(prefix)) + 5 + len(prefix)

	length -= len(prefix)

	const userNameChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	buf := make([]byte, length)
	for i, _ := range buf {
		randIdx := rand.Intn(len(userNameChars))
		buf[i] = userNameChars[randIdx]
	}

	return prefix + string(buf)
}
