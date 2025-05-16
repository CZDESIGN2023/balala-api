package data

import (
	"context"
	"encoding/json"
	"errors"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "go-cs/internal/domain/space_work_object"
	repo "go-cs/internal/domain/space_work_object/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type spaceWorkObjectRepo struct {
	baseRepo
}

func NewSpaceWorkObjectRepo(data *Data, logger log.Logger) repo.SpaceWorkObjectRepo {
	moduleName := "SpaceWorkObjectRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceWorkObjectRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_SpaceWorkObject, "", repo.clearCache2)

	return repo
}

func (c *spaceWorkObjectRepo) CreateSpaceWorkObject(ctx context.Context, workObject *domain.SpaceWorkObject) error {

	po := convert.SpaceWorkObjectEntityToPo(workObject)
	err := c.data.DB(ctx).Model(&db.SpaceWorkObject{}).Create(po).Error
	if err != nil {
		return err
	}

	return err
}

func (c *spaceWorkObjectRepo) SaveSpaceWorkObject(ctx context.Context, workObject *domain.SpaceWorkObject) error {

	diffs := workObject.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkObject{}
	cloumns := m.Cloumns()

	cloums := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			cloums[cloumns.Ranking] = workObject.Ranking
		case domain.Diff_Name:
			cloums[cloumns.WorkObjectName] = workObject.WorkObjectName
		case domain.Diff_Status:
			cloums[cloumns.WorkObjectStatus] = workObject.WorkObjectStatus
		}
	}

	if len(cloums) == 0 {
		return nil
	}

	cloums[cloumns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", workObject.Id).UpdateColumns(cloums).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceWorkObjectRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max any
	row := c.data.DB(ctx).Model(&db.SpaceWorkObject{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return cast.ToInt64(max), err
	}
	return cast.ToInt64(max), nil
}

func (c *spaceWorkObjectRepo) CheckSpaceWorkObjectName(ctx context.Context, space_id int64, workObjectName string) (bool, error) {
	var count int64
	err := c.data.DB(ctx).Model(&db.SpaceWorkObject{}).
		Where("space_id=? and BINARY work_object_name=?", space_id, workObjectName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (c *spaceWorkObjectRepo) QSpaceWorkObjectList(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkObject, error) {

	var list []*db.SpaceWorkObject
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkObject{}).
		Where("space_id=?", spaceId).
		Order("ranking desc").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceWorkObjectPoToEntitys(list), err
}

func (c *spaceWorkObjectRepo) QSpaceWorkObjectById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.SpaceWorkObject, error) {

	var list []*db.SpaceWorkObject
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkObject{}).
		Where("id in ? and space_id=?", ids, spaceId).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceWorkObjectPoToEntitys(list), err
}

func (c *spaceWorkObjectRepo) SpaceWorkObjectMap(ctx context.Context, spaceId int64) (map[int64]*domain.SpaceWorkObject, error) {
	list, err := c.QSpaceWorkObjectList(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(_ int, v *domain.SpaceWorkObject) (int64, *domain.SpaceWorkObject) {
		return v.Id, v
	})

	return m, nil
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObject(ctx context.Context, spaceId int64, workObjectId int64) (*domain.SpaceWorkObject, error) {
	var row *db.SpaceWorkObject
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkObject{}).
		Where("id=? and space_id=?", workObjectId, spaceId).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return convert.SpaceWorkObjectPoToEntity(row), err
}

func (c *spaceWorkObjectRepo) DelWorkObject(ctx context.Context, workObjectId int64) (int64, error) {
	res := c.data.DB(ctx).
		Where("id = ?", workObjectId).
		Unscoped().
		Delete(&db.SpaceWorkObject{})

	return res.RowsAffected, res.Error
}

func (c *spaceWorkObjectRepo) DelWorkObjectBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := c.data.DB(ctx).
		Where("space_id = ?", spaceId).
		Unscoped().
		Delete(&db.SpaceWorkObject{})

	return res.RowsAffected, res.Error
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObjectCount(ctx context.Context, spaceId int64) (int64, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkObject{}).
		Where("space_id=?", spaceId).
		Count(&rowNum).Error
	return rowNum, err
}

func (c *spaceWorkObjectRepo) SpaceWorkObjectMapByObjectIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkObject, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	list, err := c.GetSpaceWorkObjectByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.SpaceWorkObject) (int64, *domain.SpaceWorkObject) {
		return t.Id, t
	})

	return m, nil
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObjectByIds(ctx context.Context, ids []int64) ([]*domain.SpaceWorkObject, error) {
	ids = stream.Unique(ids)
	ids = stream.Filter(ids, func(id int64) bool {
		return id != 0
	})

	if len(ids) == 0 {
		return nil, nil
	}

	fromRedis, err := c.GetSpaceWorkObjectByIdsFromRedis(ctx, ids)
	if err != nil {
		c.log.Error(err)
	}

	if len(ids) == len(fromRedis) {
		return fromRedis, nil
	}

	redisUserIds := stream.Map(fromRedis, func(v *domain.SpaceWorkObject) int64 {
		return v.Id
	})

	var noCachedIds []int64
	for _, v := range ids {
		if !stream.Contains(redisUserIds, v) {
			noCachedIds = append(noCachedIds, v)
		}
	}

	fromDB, err := c.GetSpaceWorkObjectByIdsFromDB(ctx, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewWorkObjectKey(v.Id).Key()
		v := utils.ToJSON(v)
		kv[k] = v
	}

	_, _ = c.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range kv {
			pipeline.Set(ctx, k, v, time.Hour*24*3)
		}
		return nil
	})

	return append(fromRedis, fromDB...), nil
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObjectByIdsFromDB(ctx context.Context, ids []int64) ([]*domain.SpaceWorkObject, error) {
	var list []*db.SpaceWorkObject
	err := c.data.RoDB(ctx).Where("id in ?", ids).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceWorkObjectPoToEntitys(list), nil
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObjectByIdsFromRedis(ctx context.Context, ids []int64) ([]*domain.SpaceWorkObject, error) {

	keys := stream.Map(ids, func(v int64) string {
		return NewWorkObjectKey(v).Key()
	})

	result, err := c.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	list := stream.Map(result, func(v any) *domain.SpaceWorkObject {
		var u domain.SpaceWorkObject
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, err
}

//func (c *spaceWorkObjectRepo) clearCache(ctx context.Context, workObjectIds ...int64) {
//	keys := stream.Map(workObjectIds, func(v int64) string {
//		return NewWorkObjectKey(v).Key()
//	})
//	c.data.rdb.Del(ctx, keys...)
//}

func (c *spaceWorkObjectRepo) clearCache2(workObjectIds []int64) {
	keys := stream.Map(workObjectIds, func(v int64) string {
		return NewWorkObjectKey(v).Key()
	})
	c.data.rdb.Del(context.Background(), keys...)
}

func (c *spaceWorkObjectRepo) IsEmpty(ctx context.Context, id int64) (bool, error) {
	var itemId int
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id = ?", id).
		Limit(1).
		Pluck("id", &itemId).Error
	if err != nil {
		return false, err
	}

	return itemId == 0, nil
}

func (c *spaceWorkObjectRepo) SetOrder(ctx context.Context, spaceId int64, fromIdx, toIdx int64) error {
	var srcId int64

	member, err := c.QSpaceWorkObjectList(ctx, spaceId)
	if err != nil {
		return err
	}

	if fromIdx < 0 || fromIdx >= int64(len(member)) {
		return errors.New("wrong index")
	}

	srcRank, dstRank := member[fromIdx].Ranking, member[toIdx].Ranking
	srcId = member[fromIdx].Id

	var updateExp clause.Expr
	var condExp clause.Expr
	if srcRank > dstRank {
		updateExp = gorm.Expr("ranking + 1")
		condExp = gorm.Expr("ranking >= ? AND ranking < ?", dstRank, srcRank)
	} else {
		updateExp = gorm.Expr("ranking - 1")
		condExp = gorm.Expr("ranking > ? AND ranking <= ?", srcRank, dstRank)
	}

	err = c.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&db.SpaceWorkObject{}).
			Where("space_id = ?", spaceId).
			Where(condExp).
			UpdateColumn("ranking", updateExp).Error
		if err != nil {
			return err
		}

		err = tx.Model(&db.SpaceWorkObject{}).
			Where("id = ?", srcId).
			UpdateColumn("ranking", dstRank).Error
		return err
	})

	//objectIds := stream.Map(member, func(t *db.SpaceWorkObject) int64 {
	//	return t.SpaceId
	//})

	//c.clearCache(ctx, objectIds...)
	return err
}

func (c *spaceWorkObjectRepo) GetSpaceWorkObjectBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.SpaceWorkObject, error) {
	var list []*db.SpaceWorkObject
	err := c.data.RoDB(ctx).Where("space_id in ?", spaceIds).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceWorkObjectPoToEntitys(list), nil
}
