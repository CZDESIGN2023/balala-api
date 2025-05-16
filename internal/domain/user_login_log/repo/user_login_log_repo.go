package repo

import (
	"context"
	domain "go-cs/internal/domain/user_login_log"
)

type UserLoginLogRepo interface {
	CreateLoginInfo(ctx context.Context, in *domain.UserLoginLog) error
	GetLatestLoginTime(ctx context.Context) int64
}
