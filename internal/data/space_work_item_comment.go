package data

import (
	"context"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	domain "go-cs/internal/domain/space_work_item_comment"
	"go-cs/internal/domain/space_work_item_comment/repo"
)

type spaceWorkItemCommentRepo struct {
	baseRepo
}

func NewSpaceWorkItemCommentRepo(data *Data, logger log.Logger) repo.SpaceWorkItemCommentRepo {
	moduleName := "SpaceWorkItemCommentRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &spaceWorkItemCommentRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *spaceWorkItemCommentRepo) SetUserReadTime(ctx context.Context, userId int64, workItemId int64, t time.Time) error {
	r.RemoveUnreadNumForUser(ctx, userId, []int64{workItemId})
	key := fmt.Sprintf("balala:workItemCommentUserReadTime:%v:%v", workItemId, userId)
	return r.data.rdb.Set(ctx, key, t.UnixMilli(), time.Hour*24*30*6).Err()
}

func (r *spaceWorkItemCommentRepo) GetUserReadTime(ctx context.Context, userId int64, workItemId int64) (time.Time, error) {
	key := fmt.Sprintf("balala:workItemCommentUserReadTime:%v:%v", workItemId, userId)
	i, err := r.data.rdb.Get(ctx, key).Int64()
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(i), nil
}

func (r *spaceWorkItemCommentRepo) CreateComment(ctx context.Context, comment *domain.SpaceWorkItemComment) error {

	po := convert.SpaceWorkItemCommentEntityToPo(comment)
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemComment{}).Create(po).Error
	if err != nil {
		r.log.Error(err)
		return err
	}

	comment.Id = po.Id

	return nil
}

func (r *spaceWorkItemCommentRepo) SaveComment(ctx context.Context, comment *domain.SpaceWorkItemComment) error {

	diffs := comment.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkItemComment{}
	mColumns := m.Cloumns()

	updateColumns := make(map[string]interface{})

	for _, diff := range diffs {
		switch diff {
		case domain.Diff_Content:
			updateColumns[mColumns.Content] = comment.Content
			updateColumns[mColumns.UpdatedAt] = time.Now().Unix()
			comment.UpdatedAt = time.Now().Unix()
		case domain.Diff_ReferUserIds:
			updateColumns[mColumns.ReferUserIds] = comment.ReferUserIds.ToJsonString()
		case domain.Diff_Emojis:
			updateColumns[mColumns.Emojis] = comment.Emojis.ToJson()
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	err := r.data.DB(ctx).Model(m).Where("id=?", comment.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *spaceWorkItemCommentRepo) GetComment(ctx context.Context, id int64) (*domain.SpaceWorkItemComment, error) {
	var v *db.SpaceWorkItemComment
	err := r.data.DB(ctx).Model(&db.SpaceWorkItemComment{}).
		Where("id = ?", id).
		Take(&v).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceWorkItemCommentPoToEntity(v), nil
}

func (r *spaceWorkItemCommentRepo) DelComment(ctx context.Context, id int64) (int64, error) {
	res := r.data.DB(ctx).
		Where("id = ?", id).
		Delete(&db.SpaceWorkItemComment{})

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemCommentRepo) GetCommentByIds(ctx context.Context, ids []int64) (domain.SpaceWorkItemComments, error) {
	var rows []*db.SpaceWorkItemComment
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemComment{}).
		Where("id in ?", ids).
		Find(&rows).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return convert.SpaceWorkItemCommentPoToEntities(rows), nil
}

func (r *spaceWorkItemCommentRepo) DelCommentByIds(ctx context.Context, ids []int64) (int64, error) {
	res := r.data.DB(ctx).
		Where("id in ?", ids).
		Delete(&db.SpaceWorkItemComment{})

	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemCommentRepo) CommentMap(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkItemComment, error) {
	byIds, err := r.GetCommentByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(byIds, func(i int, v *domain.SpaceWorkItemComment) (int64, *domain.SpaceWorkItemComment) {
		return v.Id, v
	})

	return m, nil
}

func (r *spaceWorkItemCommentRepo) GetCommentByWorkItemId(ctx context.Context, workItemId int64) (domain.SpaceWorkItemComments, error) {
	var rows []*db.SpaceWorkItemComment
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemComment{}).
		Where("work_item_id = ?", workItemId).
		Find(&rows).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return convert.SpaceWorkItemCommentPoToEntities(rows), nil
}

func (r *spaceWorkItemCommentRepo) QCommentPagination(ctx context.Context, workItemId int64, pos, size int, order string) (domain.SpaceWorkItemComments, error) {
	var exp any
	switch order {
	case "DESC":
		exp = gorm.Expr("id < ?", pos)
	case "ASC":
		exp = gorm.Expr("id > ?", pos)
	}

	var rows []*db.SpaceWorkItemComment
	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemComment{}).
		Where("work_item_id = ?", workItemId).
		Where(exp).
		Order("id " + order).
		Limit(size).
		Find(&rows).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	return convert.SpaceWorkItemCommentPoToEntities(rows), nil
}

func (r *spaceWorkItemCommentRepo) DelCommentByWorkItemIds(ctx context.Context, workItemIds []int64) (int64, error) {
	res := r.data.DB(ctx).
		Where("work_item_id in ?", workItemIds).
		Delete(&db.SpaceWorkItemComment{})
	err := res.Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return res.RowsAffected, nil
}

func (r *spaceWorkItemCommentRepo) DelCommentByWorkItemId(ctx context.Context, workItemId int64) (int64, error) {
	return r.DelCommentByWorkItemIds(ctx, []int64{workItemId})
}

func (r *spaceWorkItemCommentRepo) CountCommentByWorkItemId(ctx context.Context, workItemId int64) (int64, error) {
	m, err := r.CountCommentByWorkItemIds(ctx, []int64{workItemId})
	if err != nil {
		return 0, err
	}

	return m[workItemId], nil
}

func (r *spaceWorkItemCommentRepo) CountWorkItemCommentNumByTime(ctx context.Context, workItemId int64, t time.Time) (int64, error) {
	var count int64

	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemComment{}).
		Where("work_item_id = ?", workItemId).
		Where("created_at > ?", t.Unix()).
		Count(&count).Error
	if err != nil {
		r.log.Error(err)
		return 0, err
	}

	return count, nil
}

func (r *spaceWorkItemCommentRepo) CountCommentByWorkItemIds(ctx context.Context, workItemIds []int64) (map[int64]int64, error) {
	type s struct {
		Id    int64
		Count int64
	}
	var list []s

	err := r.data.RoDB(ctx).Model(&db.SpaceWorkItemComment{}).
		Select("work_item_id id, SUM(1) count").
		Where("work_item_id in ?", workItemIds).
		Group("work_item_id").
		Find(&list).Error
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	m := stream.ToMap(list, func(i int, s s) (int64, int64) {
		return s.Id, s.Count
	})

	return m, nil
}

func (r *spaceWorkItemCommentRepo) IncrUnreadNumForUser(ctx context.Context, workItemId int64, userIds []int64) error {
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

func (r *spaceWorkItemCommentRepo) RemoveUnreadNumForUser(ctx context.Context, userId int64, workItemIds []int64) error {
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

func (r *spaceWorkItemCommentRepo) GetUserUnreadNum(ctx context.Context, userId int64, workItemId int64) (int64, error) {
	ids, err := r.UserUnreadNumMapByWorkItemIds(ctx, userId, []int64{workItemId})
	if err != nil {
		return 0, err
	}

	return ids[workItemId], nil
}

func (r *spaceWorkItemCommentRepo) UserUnreadNumMapByWorkItemIds(ctx context.Context, userId int64, workItemIds []int64) (map[int64]int64, error) {
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

func (r *spaceWorkItemCommentRepo) SetUserUnreadNum(ctx context.Context, userId int64, numMap map[int64]int64) error {
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
