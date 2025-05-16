package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/space"
	repo "go-cs/internal/domain/space/repo"
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/log"
)

type spaceRepo struct {
	baseRepo
}

func NewSpaceRepo(data *Data, logger log.Logger) repo.SpaceRepo {
	moduleName := "SpaceRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_Space, "", repo.clearCache2)

	return repo
}

func (c *spaceRepo) CreateSpace(ctx context.Context, space *domain.Space) error {

	po := convert.SpaceEntityToPo(space)
	err := c.data.DB(ctx).Model(&db.Space{}).Create(po).Error
	return err
}

func (c *spaceRepo) IsExistBySpaceName(ctx context.Context, userId int64, spaceName string) (bool, error) {
	var rowNum int64
	err := c.data.DB(ctx).Model(&db.Space{}).Where("user_id=? and BINARY space_name=?", userId, spaceName).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *spaceRepo) GetUserSpaceList(ctx context.Context, userId int64) ([]*domain.Space, error) {
	var rows []*db.Space

	memeberTb := (&db.SpaceMember{}).TableName()
	spaceTb := (&db.Space{}).TableName()

	err := c.data.RoDB(ctx).Table(memeberTb+" m").Select("s.*").Joins(
		"INNER JOIN "+spaceTb+" s ON m.space_id = s.id where m.user_id=? ", userId,
	).Find(&rows).Error
	if err != nil {
		c.log.Errorf("error in GetListByUser: err = %v", err)
		return nil, err
	}

	var list []*domain.Space
	for _, v := range rows {
		list = append(list, convert.SpacePoToEntity(v))
	}

	return list, nil
}

func (c *spaceRepo) GetUserSpaceIds(ctx context.Context, userId int64) ([]int64, error) {
	//查询用户加入的所有空间
	var list []int64

	c.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Where("user_id=?", userId).
		Pluck("space_id", &list)

	return list, nil
}

func (c *spaceRepo) GetSpaceByCreator(ctx context.Context, userId int64, spaceId int64) (*domain.Space, error) {
	var space db.Space
	err := c.data.RoDB(ctx).Model(&db.Space{}).Where("id=? and user_id=?", spaceId, userId).Take(&space).Error
	if err != nil {
		return nil, err
	}
	return convert.SpacePoToEntity(&space), nil
}

func (c *spaceRepo) SaveSpace(ctx context.Context, space *domain.Space) error {

	diffs := space.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.Space{}
	mClounms := m.Cloumns()

	updateColumns := make(map[string]interface{})

	for _, v := range diffs {
		switch v {
		case domain.Diff_SpaceName:
			updateColumns[mClounms.SpaceName] = space.SpaceName
		case domain.Diff_Describe:
			updateColumns[mClounms.Describe] = space.Describe
		case domain.Diff_UserId:
			updateColumns[mClounms.UserId] = space.UserId
		case domain.Diff_Notify:
			updateColumns[mClounms.Notify] = space.Notify
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[mClounms.UpdatedAt] = time.Now().Unix()
	err := c.data.RoDB(ctx).Model(&db.Space{}).Where("id=?", space.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	c.clearCache2([]int64{space.Id})

	return nil
}

func (c *spaceRepo) DelSpace(ctx context.Context, spaceId int64) error {

	var opValue map[string]interface{} = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.Space{}).Unscoped().Where("id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}

	c.clearCache2([]int64{spaceId})

	return nil
}

func (c *spaceRepo) SpaceMap(ctx context.Context, ids []int64) (map[int64]*domain.Space, error) {
	list, err := c.GetSpaceByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.Space) (int64, *domain.Space) {
		return t.Id, t
	})

	return m, nil
}

func (c *spaceRepo) GetSpace(ctx context.Context, id int64) (*domain.Space, error) {
	spaces, err := c.GetSpaceByIds(ctx, []int64{id})
	if err != nil {
		return nil, err
	}

	if len(spaces) > 0 {
		return spaces[0], nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (c *spaceRepo) GetSpaceFromDB(ctx context.Context, spaceId int64) (*domain.Space, error) {

	m := &db.Space{}
	mColumns := m.Cloumns()

	var row *db.Space

	err := c.data.DB(ctx).Model(m).
		Select(mColumns.SelectEx(mColumns.Describe, mColumns.Remark)).
		Where("id=?", spaceId).
		Take(&row).Error

	if err != nil {
		return nil, err
	}

	return convert.SpacePoToEntity(row), nil
}

func (c *spaceRepo) GetSpaceDetail(ctx context.Context, spaceId int64) (*domain.Space, error) {

	m := &db.Space{}
	var row *db.Space

	err := c.data.DB(ctx).Model(m).
		Where("id=?", spaceId).
		Take(&row).Error

	if err != nil {
		return nil, err
	}

	return convert.SpacePoToEntity(row), nil
}

func (c *spaceRepo) GetSpaceFromRedis(ctx context.Context, spaceId int64) (*domain.Space, error) {
	key := NewSpaceKey(spaceId).Key()
	val, err := c.data.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	if len(val) == 0 {
		return nil, nil
	}

	var v db.Space
	_ = json.Unmarshal(val, &v)

	return convert.SpacePoToEntity(&v), nil
}

func (c *spaceRepo) GetSpaceByIdsFromRedis(ctx context.Context, ids []int64) ([]*domain.Space, error) {

	keys := stream.Map(ids, func(v int64) string {
		return NewSpaceKey(v).Key()
	})

	result, err := c.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	list := stream.Map(result, func(v any) *domain.Space {
		var u domain.Space
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, err
}

func (c *spaceRepo) GetSpaceByIdsFromDB(ctx context.Context, ids []int64) ([]*domain.Space, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var rows []*db.Space
	err := c.data.RoDB(ctx).Where("id in (?)", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.Space
	for _, v := range rows {
		list = append(list, convert.SpacePoToEntity(v))
	}

	return list, err
}

func (c *spaceRepo) GetSpaceByIds(ctx context.Context, ids []int64) ([]*domain.Space, error) {
	ids = stream.Unique(ids)
	ids = stream.Filter(ids, func(id int64) bool {
		return id != 0
	})

	if len(ids) == 0 {
		return nil, nil
	}

	fromRedis, err := c.GetSpaceByIdsFromRedis(ctx, ids)
	if err != nil {
		c.log.Error(err)
	}

	if len(ids) == len(fromRedis) {
		return fromRedis, nil
	}

	noCachedIds := stream.Filter(ids, func(id int64) bool {
		return !stream.ContainsFunc(fromRedis, func(space *domain.Space) bool {
			return id == space.Id
		})
	})

	fromDB, err := c.GetSpaceByIdsFromDB(ctx, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewSpaceKey(v.Id).Key()
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

func (c *spaceRepo) clearCache2(ids []int64) {
	keys := stream.Map(ids, func(v int64) string {
		return NewSpaceKey(v).Key()
	})

	c.data.rdb.Del(context.Background(), keys...)
}

func (c *spaceRepo) CreateSpaceConfig(ctx context.Context, spaceConfig *domain.SpaceConfig) error {
	po := convert.SpaceConfigEntityToPo(spaceConfig)
	err := c.data.DB(ctx).Create(po).Error
	return err
}

func (c *spaceRepo) SaveSpaceConfig(ctx context.Context, spaceConfig *domain.SpaceConfig) error {

	diffs := spaceConfig.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceConfig{}
	mColumns := m.Cloumns()

	updateColumns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_SpaceConfig_WorkingDay:
			updateColumns[mColumns.WorkingDay] = spaceConfig.WorkingDay
		case domain.Diff_SpaceConfig_CommentDeletable:
			updateColumns[mColumns.CommentDeletable] = spaceConfig.CommentDeletable
		case domain.Diff_SpaceConfig_CommentDeletableWhenArchived:
			updateColumns[mColumns.CommentDeletableWhenArchived] = spaceConfig.CommentDeletableWhenArchived
		case domain.Diff_SpaceConfig_CommentShowPos:
			updateColumns[mColumns.CommentShowPos] = spaceConfig.CommentShowPos
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[mColumns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", spaceConfig.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceRepo) getSpaceConfig(ctx context.Context, spaceId int64) (*domain.SpaceConfig, error) {
	fromDB, err := c.GetSpaceConfigListFromDB(ctx, []int64{spaceId})
	if err != nil {
		return nil, err
	}

	if len(fromDB) > 0 {
		return convert.SpaceConfigPoToEntity(fromDB[0]), nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (c *spaceRepo) SetSpaceCommentDeletable(ctx context.Context, spaceId int64, val int64) error {
	err := c.updateSpaceConfigFields(ctx, spaceId, map[string]any{
		"comment_deletable": val,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceRepo) updateSpaceConfigFields(ctx context.Context, spaceId int64, kvs map[string]any) error {

	err := c.data.DB(ctx).Model(&db.SpaceConfig{}).Where("space_id = ?", spaceId).UpdateColumns(kvs).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceRepo) HasSpaceConfig(ctx context.Context, spaceId int64) (bool, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.SpaceConfig{}).Where("space_id = ?", spaceId).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	return rowNum > 0, nil
}

func (c *spaceRepo) GetSpaceConfig(ctx context.Context, spaceId int64) (*domain.SpaceConfig, error) {
	config, err := c.getSpaceConfig(ctx, spaceId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config = c.defaultConfig(spaceId)
		_ = c.CreateSpaceConfig(ctx, config)
	}

	return config, nil
}

func (c *spaceRepo) defaultConfig(spaceId int64) *domain.SpaceConfig {
	return &domain.SpaceConfig{
		SpaceId:          spaceId,
		WorkingDay:       `[1,2,3,4,5]`,
		CommentDeletable: 0,
	}
}

func (c *spaceRepo) GetSpaceConfigListFromDB(ctx context.Context, spaceIds []int64) ([]*db.SpaceConfig, error) {
	var list []*db.SpaceConfig
	err := c.data.DB(ctx).Model(&db.SpaceConfig{}).Where("space_id in ?", spaceIds).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *spaceRepo) SpaceConfigMap(ctx context.Context, spaceIds []int64) (map[int64]*domain.SpaceConfig, error) {
	list, err := c.GetSpaceConfigListFromDB(ctx, spaceIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(_ int, v *db.SpaceConfig) (int64, *domain.SpaceConfig) {
		return v.SpaceId, convert.SpaceConfigPoToEntity(v)
	})

	for _, spaceId := range spaceIds {
		if _, ok := m[spaceId]; !ok {
			m[spaceId] = c.defaultConfig(spaceId)
		}
	}

	return m, nil
}

func (c *spaceRepo) DelSpaceConfig(ctx context.Context, spaceId int64) error {

	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.SpaceConfig{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil

}

func (c *spaceRepo) SetTempConfig(ctx context.Context, spaceId int64, confMap map[string]string) error {
	for key, val := range confMap {
		key = fmt.Sprintf("balala:space_temp_config:%d:%s", spaceId, key)
		c.data.rdb.Set(ctx, key, val, time.Hour*24*90)
	}

	return nil
}

func (c *spaceRepo) GetTempConfig(ctx context.Context, userId int64, keys ...string) map[string]string {
	finalKeys := stream.Map(keys, func(v string) string {
		return fmt.Sprintf("balala:space_temp_config:%d:%s", userId, v)
	})

	ret := c.data.rdb.MGet(ctx, finalKeys...)

	values := stream.Map(ret.Val(), func(v any) string {
		return cast.ToString(v)
	})

	for _, key := range finalKeys {
		c.data.rdb.Expire(ctx, key, time.Hour*24*90)
	}

	return stream.Zip(keys, values)
}

func (c *spaceRepo) DelTempConfig(ctx context.Context, userId int64, keys ...string) error {
	finalKeys := stream.Map(keys, func(v string) string {
		return fmt.Sprintf("balala:space_temp_config:%d:%s", userId, v)
	})

	return c.data.rdb.Del(ctx, finalKeys...).Err()
}

func (c *spaceRepo) GetAllSpaceIds() ([]int64, error) {
	var ids []int64
	err := c.data.RoDB(context.Background()).Model(&db.Space{}).Pluck("id", &ids).Error
	return ids, err
}

func (c *spaceRepo) GetAllSpace(ctx context.Context) ([]*domain.Space, error) {
	var list []*db.Space
	err := c.data.DB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}

	ret := stream.Map(list, func(v *db.Space) *domain.Space {
		return convert.SpacePoToEntity(v)
	})

	return ret, nil
}
