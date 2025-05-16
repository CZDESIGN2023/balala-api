package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"time"
)

type AdminRepo struct {
	baseRepo
}

func NewAdminRepo(data *Data, logger log.Logger) biz.AdminRepo {
	moduleName := "AdminRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &AdminRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
	return repo
}

func (r *AdminRepo) SearchUser(ctx context.Context, keyword string, pos, size int, sorts []vo.OrderBy) (list []*db.User, nextPos int, hasNext bool, err error) {
	var rows []*db.User

	tx := r.data.RoDB(ctx).Model(&db.User{}).Where("user_status = 1")
	if keyword != "" {
		tx = tx.Where("user_pinyin like ? or user_nickname like ? or user_name like ?", ",%"+keyword+"%,", "%"+keyword+"%", "%"+keyword+"%")
	}

	for _, order := range sorts {
		tx = tx.Order(order.Field + " " + order.Order)
	}

	err = tx.
		Order("id desc").
		Offset(pos * size).
		Limit(size + 1).
		Find(&rows).Error
	if err != nil {
		return nil, 0, false, err
	}

	if len(rows) > size {
		hasNext = true
		nextPos = pos + 1
		rows = rows[:size]
	}

	return rows, nextPos, hasNext, nil
}
func (r *AdminRepo) SetAdminNeedResetPwdStatus(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("balala:admin:needResetPwd:%d", userId)
	r.data.rdb.Set(ctx, key, time.Now().UnixMilli(), 0)
	return nil
}
func (r *AdminRepo) GetAdminNeedResetPwdStatus(ctx context.Context, userId int64) (bool, error) {
	key := fmt.Sprintf("balala:admin:needResetPwd:%d", userId)
	val := r.data.rdb.Exists(ctx, key).Val()
	return val == 1, nil
}
func (r *AdminRepo) DelAdminNeedResetPwdStatus(ctx context.Context, userId int64) error {
	key := fmt.Sprintf("balala:admin:needResetPwd:%d", userId)
	r.data.rdb.Del(ctx, key)
	return nil
}

func (r *AdminRepo) TotalUser(ctx context.Context, keyword string) (int64, error) {
	tx := r.data.RoDB(ctx).Model(&db.User{}).Where("user_status = 1")
	if keyword != "" {
		tx = tx.Where("user_pinyin like ? or user_nickname like ? or user_name like ?", ",%"+keyword+"%,", "%"+keyword+"%", "%"+keyword+"%")
	}

	var count int64
	err := tx.
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
