package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo"
	"go-cs/internal/biz"
	"go-cs/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type LogRepo struct {
	baseRepo
}

func NewLogRepo(data *Data, logger log.Logger) biz.LogRepo {
	moduleName := "LogRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &LogRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *LogRepo) UserLoginLogPagination(ctx context.Context, userId int64, pos, size int) ([]*db.UserLoginLog, error) {
	var list []*db.UserLoginLog
	err := r.data.RoDB(ctx).Model(&db.UserLoginLog{}).
		Where("login_user_id = ?", userId).
		Where("id <= ?", pos).
		Order("id DESC").
		Limit(size).
		Find(&list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}

func (r *LogRepo) OpLogPagination(ctx context.Context, searchVo *vo.OpLogPaginationSearchVo) ([]*db.OperLog, error) {
	var list []*db.OperLog

	tx := r.data.RoDB(ctx).Model(&db.OperLog{})
	if len(searchVo.SpaceIds) > 0 {
		tx = tx.Where("space_id in ?", searchVo.SpaceIds)
	}

	if len(searchVo.IncludeModuleType) > 0 {
		tx = tx.Where("module_type in ?", searchVo.IncludeModuleType)
	}

	if searchVo.ModuleType > 0 {
		tx = tx.Where("module_type = ?", searchVo.ModuleType)
	}

	if searchVo.ModuleId > 0 {
		tx = tx.Where("module_id = ?", searchVo.ModuleId)
	}

	if searchVo.OperId > 0 {
		tx = tx.Where("oper_id = ?", searchVo.OperId)
	}

	if searchVo.Pos > 0 {
		tx = tx.Where("id <= ?", searchVo.Pos)
	}

	if searchVo.OperatorType > 0 {
		tx = tx.Where("operator_type = ?", searchVo.OperatorType)
	}

	err := tx.Order("id DESC").
		Limit(searchVo.Size).
		Find(&list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return list, nil
}
