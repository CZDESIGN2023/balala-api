package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/utils"

	oper_repo "go-cs/internal/domain/oper_log/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type operLogRepo struct {
	baseRepo
}

func NewOperLogRepo(data *Data, logger log.Logger) oper_repo.OperLogRepo {
	moduleName := "OperLogRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &operLogRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *operLogRepo) AddOperLog(ctx context.Context, in *db.OperLog) error {
	err := r.data.DB(ctx).Model(&db.OperLog{}).Create(in).Error
	return err
}
