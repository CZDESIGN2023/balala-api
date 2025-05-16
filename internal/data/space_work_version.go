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

	"github.com/go-kratos/kratos/v2/log"

	domain "go-cs/internal/domain/space_work_version"
	repo "go-cs/internal/domain/space_work_version/repo"
)

type spaceWorkVersionRepo struct {
	baseRepo
}

func NewSpaceWorkVersionRepo(data *Data, logger log.Logger) repo.SpaceWorkVersionRepo {
	moduleName := "SpaceWorkVersionRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceWorkVersionRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_SpaceWorkVersion, "", repo.clearCache2)

	return repo
}

func (c *spaceWorkVersionRepo) CreateSpaceWorkVersion(ctx context.Context, workVersion *domain.SpaceWorkVersion) error {

	po := convert.SpaceWorkVersionEntityToPo(workVersion)
	err := c.data.DB(ctx).Model(&db.SpaceWorkVersion{}).Create(po).Error
	if err != nil {
		c.log.Debug(err)
		return err
	}

	return err
}

func (c *spaceWorkVersionRepo) CheckSpaceWorkVersionName(ctx context.Context, space_id int64, workVersionName string) (bool, error) {
	var rowNum int64
	err := c.data.DB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=? and BINARY version_name=?", space_id, workVersionName).
		Count(&rowNum).Error
	if err != nil {
		c.log.Debug(err)
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *spaceWorkVersionRepo) QSpaceWorkVersionList(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkVersion, error) {

	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=?", spaceId).
		Order("ranking desc").
		Find(&rows).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	return convert.SpaceWorkVersionPoToEntitys(rows), err
}

func (c *spaceWorkVersionRepo) QSpaceWorkVersionById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.SpaceWorkVersion, error) {

	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("id in ? and space_id=?", ids, spaceId).
		Find(&rows).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	return convert.SpaceWorkVersionPoToEntitys(rows), err
}

func (c *spaceWorkVersionRepo) SpaceWorkVersionMap(ctx context.Context, spaceId int64) (map[int64]*domain.SpaceWorkVersion, error) {

	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=?", spaceId).
		Order("ranking desc").
		Find(&rows).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	list := convert.SpaceWorkVersionPoToEntitys(rows)
	m := stream.ToMap(list, func(_ int, t *domain.SpaceWorkVersion) (int64, *domain.SpaceWorkVersion) {
		return t.Id, t
	})

	return m, err
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionByKey(ctx context.Context, spaceId int64, workVersionKey string) (*domain.SpaceWorkVersion, error) {
	var row *db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=? and version_key=?", spaceId, workVersionKey).
		Take(&row).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}
	return convert.SpaceWorkVersionPoToEntity(row), err
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersion(ctx context.Context, workVersionId int64) (*domain.SpaceWorkVersion, error) {
	var row *db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("id=?", workVersionId).
		Take(&row).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}
	return convert.SpaceWorkVersionPoToEntity(row), err
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionBySpaceId(ctx context.Context, spaceId int64) ([]*domain.SpaceWorkVersion, error) {
	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=?", spaceId).
		Find(&rows).Error
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	return convert.SpaceWorkVersionPoToEntitys(rows), err
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.SpaceWorkVersion, error) {
	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id IN ?", spaceIds).
		Find(&rows).Error
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	return convert.SpaceWorkVersionPoToEntitys(rows), err
}

func (c *spaceWorkVersionRepo) SaveSpaceWorkVersion(ctx context.Context, workVersion *domain.SpaceWorkVersion) error {

	diffs := workVersion.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceWorkVersion{}
	cloumns := m.Cloumns()

	updateColumns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			updateColumns[cloumns.Ranking] = workVersion.Ranking
		case domain.Diff_VersionName:
			updateColumns[cloumns.VersionName] = workVersion.VersionName
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[cloumns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", workVersion.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceWorkVersionRepo) DelWorkVersion(ctx context.Context, workVersionId int64) (int64, error) {
	res := c.data.DB(ctx).
		Where("id = ?", workVersionId).
		Delete(&db.SpaceWorkVersion{})

	return res.RowsAffected, res.Error
}

func (c *spaceWorkVersionRepo) DelWorkVersionBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := c.data.DB(ctx).
		Where("space_id = ?", spaceId).
		Delete(&db.SpaceWorkVersion{})

	return res.RowsAffected, res.Error
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionCount(ctx context.Context, spaceId int64) (int64, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=?", spaceId).
		Count(&rowNum).Error
	return rowNum, err
}

func (c *spaceWorkVersionRepo) SpaceWorkVersionMapByVersionIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceWorkVersion, error) {
	list, err := c.GetSpaceWorkVersionByIds(ctx, ids)
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.SpaceWorkVersion) (int64, *domain.SpaceWorkVersion) {
		return t.Id, t
	})

	return m, nil
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionByIds(ctx context.Context, ids []int64) ([]*domain.SpaceWorkVersion, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	ids = stream.Unique(ids)
	ids = stream.Filter(ids, func(id int64) bool {
		return id != 0
	})

	if len(ids) == 0 {
		return nil, nil
	}

	fromCache, err := c.GetSpaceWorkVersionByIdsFromRedis(ctx, ids)
	if err != nil {
		c.log.Error(err)
	}

	if len(ids) == len(fromCache) {
		return fromCache, nil
	}

	cacheIds := stream.Map(fromCache, func(v *domain.SpaceWorkVersion) int64 {
		return v.Id
	})

	var noCachedIds []int64
	for _, v := range ids {
		if !stream.Contains(cacheIds, v) {
			noCachedIds = append(noCachedIds, v)
		}
	}

	fromDB, err := c.GetSpaceWorkVersionByIdsFromDB(ctx, noCachedIds)
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewWorkVersionKey(v.Id).Key()
		v := utils.ToJSON(v)
		kv[k] = v
	}

	_, _ = c.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range kv {
			pipeline.Set(ctx, k, v, time.Hour*24*3)
		}
		return nil
	})

	return append(fromCache, fromDB...), nil
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionByIdsFromDB(ctx context.Context, ids []int64) ([]*domain.SpaceWorkVersion, error) {
	var rows []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Where("id in ?", ids).
		Find(&rows).Error
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	return convert.SpaceWorkVersionPoToEntitys(rows), nil
}

func (c *spaceWorkVersionRepo) GetSpaceWorkVersionByIdsFromRedis(ctx context.Context, ids []int64) ([]*domain.SpaceWorkVersion, error) {

	keys := stream.Map(ids, func(v int64) string {
		return NewWorkVersionKey(v).Key()
	})

	result, err := c.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	list := stream.Map(result, func(v any) *domain.SpaceWorkVersion {
		var u domain.SpaceWorkVersion
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, err
}

//func (c *spaceWorkVersionRepo) clearCache(ctx context.Context, workVersionIds ...int64) {
//	keys := stream.Map(workVersionIds, func(v int64) string {
//		return NewWorkVersionKey(v).Key()
//	})
//	c.data.rdb.Del(ctx, keys...)
//}

func (c *spaceWorkVersionRepo) clearCache2(workVersionIds []int64) {
	keys := stream.Map(workVersionIds, func(v int64) string {
		return NewWorkVersionKey(v).Key()
	})
	c.data.rdb.Del(context.Background(), keys...)
}

func (c *spaceWorkVersionRepo) IsEmpty(ctx context.Context, id int64) (bool, error) {
	var itemId int
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("id = ?", id).
		Limit(1).
		Pluck("id", &itemId).Error

	if err != nil {
		c.log.Debug(err)
		return false, err
	}

	return itemId == 0, nil
}

func (c *spaceWorkVersionRepo) SetOrder(ctx context.Context, spaceId, fromIdx, toIdx int64) error {
	var srcId int64

	var versions []*db.SpaceWorkVersion
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkVersion{}).
		Where("space_id=?", spaceId).
		Order("ranking desc").
		Find(&versions).Error
	if err != nil {
		c.log.Debug(err)
		return err
	}

	if fromIdx < 0 || fromIdx >= int64(len(versions)) {
		return errors.New("wrong index")
	}

	srcRank, dstRank := versions[fromIdx].Ranking, versions[toIdx].Ranking
	srcId = versions[fromIdx].Id

	var updateExp clause.Expr
	var condExp clause.Expr
	if srcRank > dstRank {
		updateExp = gorm.Expr("ranking + 1")
		condExp = gorm.Expr("ranking >= ? AND ranking <= ?", dstRank, srcRank)
	} else {
		updateExp = gorm.Expr("ranking - 1")
		condExp = gorm.Expr("ranking >= ? AND ranking <= ?", srcRank, dstRank)
	}

	err = c.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&db.SpaceWorkVersion{}).
			Where("space_id = ?", spaceId).
			Where(condExp).
			UpdateColumn("ranking", updateExp).Error
		if err != nil {
			return err
		}

		err = tx.Model(&db.SpaceWorkVersion{}).
			Where("id = ?", srcId).
			UpdateColumn("ranking", dstRank).Error
		return err
	})

	//objectIds := stream.Map(member, func(t *db.SpaceWorkVersion) int64 {
	//	return t.SpaceId
	//})

	//c.clearCache(ctx, objectIds...)
	return err
}

func (c *spaceWorkVersionRepo) VersionMap(ctx context.Context, workVersionIds []int64) (map[int64]string, error) {

	fromDB, err := c.GetSpaceWorkVersionByIdsFromDB(ctx, workVersionIds)
	if err != nil {
		c.log.Debug(err)
		return nil, err
	}

	kv := map[int64]string{}
	for _, v := range fromDB {
		kv[v.Id] = v.VersionName
	}

	return kv, nil
}

func (c *spaceWorkVersionRepo) GetVersionRelationCount(ctx context.Context, spaceId int64, workVersionId int64) (int64, error) {
	var count int64
	err := c.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ? and version_id=?", spaceId, workVersionId).
		Count(&count).Error

	if err != nil {
		c.log.Debug(err)
		return 0, err
	}
	return count, nil
}

func (c *spaceWorkVersionRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max any
	row := c.data.DB(ctx).Model(&db.SpaceWorkVersion{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return cast.ToInt64(max), err
	}
	return cast.ToInt64(max), nil
}
