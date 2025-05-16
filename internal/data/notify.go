package data

import (
	"context"
	"fmt"
	"go-cs/api/notify"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/log"

	domain "go-cs/internal/domain/notify_snapshot"
	repo "go-cs/internal/domain/notify_snapshot/repo"
)

type notifyRepo struct {
	baseRepo
}

func NewNotifyRepo(data *Data, logger log.Logger) repo.NotifySnapshotRepo {
	moduleName := "NotifyRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &notifyRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (c *notifyRepo) CreateNotify(ctx context.Context, in *domain.NotifySnapShot) error {
	po := convert.NotifySnapShotEntityToPo(in)
	err := c.data.DB(ctx).Model(db.Notify{}).Create(po).Error
	return err
}

func (c *notifyRepo) GetUserRelatedCommentIds(ctx context.Context, userId int64, pos, size int) (ids []int64, nextPos int64, hasNext bool, err error) {
	err = c.data.RoDB(ctx).Model(&db.Notify{}).
		Where("user_id = ? and typ = ?", userId, notify.Event_AddCommentEvent).
		Where("id < ?", pos).
		Order("id desc").
		Limit(size+1).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, 0, false, err
	}

	if len(ids) > size {
		hasNext = true
		nextPos = ids[size]
		ids = ids[:size]
	}

	return ids, nextPos, hasNext, nil
}

func (c *notifyRepo) GetNotifyByIds(ctx context.Context, userId int64, ids []int64) ([]*domain.NotifySnapShot, error) {
	var list []*db.Notify
	err := c.data.RoDB(ctx).Model(&db.Notify{}).
		Where("user_id = ?", userId).
		Where("id in ?", ids).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return convert.NotifySnapShotPoToEntitys(list), nil
}

func (c *notifyRepo) SaveOfflineNotify(ctx context.Context, userId int64, data []byte) error {
	key := fmt.Sprintf("balala:notify:%d", userId)

	c.data.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.RPush(ctx, key, data)
		p.LTrim(ctx, key, -300, -1)         //保留最近300条
		p.Expire(ctx, key, time.Hour*24*14) //保留14天
		return nil
	})

	return nil
}

func (c *notifyRepo) GetDelOfflineNotify(ctx context.Context, userId int64) ([]string, error) {
	key := fmt.Sprintf("balala:notify:%d", userId)
	list := c.data.rdb.LRange(ctx, key, 0, -1).Val()
	c.data.rdb.Del(ctx, key)

	return list, nil
}
