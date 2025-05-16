package service

import (
	domain "go-cs/internal/domain/user_login_log"
	"go-cs/internal/domain/user_login_log/repo"

	"time"
)

type UserLgoinLogService struct {
	repo repo.UserLoginLogRepo
}

func NewUserLoginLogService(
	repo repo.UserLoginLogRepo,
) *UserLgoinLogService {
	return &UserLgoinLogService{
		repo: repo,
	}
}

func (s *UserLgoinLogService) NewLog(user domain.LoginUser, loginStauts domain.LoginStauts, loginIp domain.LoginIp, browser string, os string, message string, loginTime time.Time) *domain.UserLoginLog {
	return &domain.UserLoginLog{
		LoginUser:   user,
		LoginIp:     loginIp,
		LoginStauts: loginStauts,
		Browser:     browser,
		Os:          os,
		Msg:         message,
		LoginAt:     loginTime.Unix(),
		CreatedAt:   time.Now().Unix(),
	}
}
