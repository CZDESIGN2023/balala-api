package data

import (
	"context"
	"fmt"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"

	"go-cs/internal/domain/space_work_item_comment/repo"
)

type spaceWorkItemCommentCacheRepo struct {
	baseRepo
}

func NewSpaceWorkItemCommentCacheRepo(data *Data, logger log.Logger) repo.SpaceWorkItemCommentCacheRepo {
	moduleName := "SpaceWorkItemCommentCacheRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &spaceWorkItemCommentCacheRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *spaceWorkItemCommentCacheRepo) SetUserReadTime(ctx context.Context, userId int64, workItemId int64, t time.Time) error {
	r.RemoveUnreadNumForUser(ctx, userId, []int64{workItemId})
	key := fmt.Sprintf("balala:workItemCommentUserReadTime:%v:%v", workItemId, userId)
	return r.data.rdb.Set(ctx, key, t.UnixMilli(), time.Hour*24*30*6).Err()
}

func (r *spaceWorkItemCommentCacheRepo) GetUserReadTime(ctx context.Context, userId int64, workItemId int64) (time.Time, error) {
	key := fmt.Sprintf("balala:workItemCommentUserReadTime:%v:%v", workItemId, userId)
	i, err := r.data.rdb.Get(ctx, key).Int64()
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(i), nil
}

func (r *spaceWorkItemCommentCacheRepo) IncrUnreadNumForUser(ctx context.Context, workItemId int64, userIds []int64) error {
	keys := stream.Map(userIds, func(userId int64) string {
		return fmt.Sprintf("balala:unreadNum:%v:%v", userId, workItemId)
	})

	_, err := r.data.rdb.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, key := range keys {
			pipeliner.Incr(ctx, key)
			pipeliner.Expire(ctx, key, time.Hour*24*180)
		}
		return nil
	})
	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *spaceWorkItemCommentCacheRepo) RemoveUnreadNumForUser(ctx context.Context, userId int64, workItemIds []int64) error {
	if len(workItemIds) == 0 {
		return nil
	}

	keys := stream.Map(workItemIds, func(workItemId int64) string {
		return fmt.Sprintf("balala:unreadNum:%v:%v", userId, workItemId)
	})

	err := r.data.rdb.Del(ctx, keys...).Err()
	if err != nil {
		r.log.Error(err)
	}

	return nil
}

func (r *spaceWorkItemCommentCacheRepo) GetUserUnreadNum(ctx context.Context, userId int64, workItemId int64) (int64, error) {
	ids, err := r.UserUnreadNumMapByWorkItemIds(ctx, userId, []int64{workItemId})
	if err != nil {
		return 0, err
	}

	return ids[workItemId], nil
}

func (r *spaceWorkItemCommentCacheRepo) UserUnreadNumMapByWorkItemIds(ctx context.Context, userId int64, workItemIds []int64) (map[int64]int64, error) {
	var m = make(map[int64]int64, len(workItemIds))
	if len(workItemIds) == 0 {
		return m, nil
	}

	keys := stream.Map(workItemIds, func(workItemId int64) string {
		return fmt.Sprintf("balala:unreadNum:%v:%v", userId, workItemId)
	})

	res, err := r.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		r.log.Error(err)
		return m, err
	}

	for i, v := range res {
		if v == nil {
			continue
		}

		m[workItemIds[i]] = cast.ToInt64(v)
	}

	return m, nil
}

func (r *spaceWorkItemCommentCacheRepo) SetUserUnreadNum(ctx context.Context, userId int64, numMap map[int64]int64) error {
	if len(numMap) == 0 {
		return nil
	}

	kvMap := stream.MapKV(numMap, func(k, v int64) (string, string) {
		key := fmt.Sprintf("balala:unreadNum:%v:%v", userId, k)
		val := cast.ToString(v)
		return key, val
	})

	err := r.data.rdb.MSet(ctx, kvMap).Err()
	if err != nil {
		r.log.Error(err)
	}

	return nil
}
