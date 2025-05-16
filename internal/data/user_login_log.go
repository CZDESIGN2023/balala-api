package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/user_login_log"
	repo "go-cs/internal/domain/user_login_log/repo"
	"go-cs/internal/utils"
)

type userLoginLogRepo struct {
	baseRepo
}

func NewUserLoginLogRepo(data *Data, logger log.Logger) repo.UserLoginLogRepo {
	moduleName := "UserLoginLogRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &userLoginLogRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *userLoginLogRepo) CreateLoginInfo(ctx context.Context, log *domain.UserLoginLog) error {
	po := convert.UserLoginLogEntityToPo(log)
	err := r.data.DB(ctx).Model(&db.UserLoginLog{}).Save(po).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userLoginLogRepo) GetLatestLoginTime(ctx context.Context) int64 {
	var createdAt int64
	err := r.data.DB(ctx).Model(&db.UserLoginLog{}).Order("id asc").Limit(1).Pluck("created_at", &createdAt).Error
	if err != nil {
		r.log.Error(err)
		return 0
	}

	return createdAt
}
