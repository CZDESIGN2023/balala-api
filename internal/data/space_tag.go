package data

import (
	"context"
	"encoding/json"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/space_tag"
	"go-cs/internal/domain/space_tag/repo"
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/log"
)

type spaceTagRepo struct {
	baseRepo
}

func NewSpaceTagRepo(data *Data, logger log.Logger) repo.SpaceTagRepo {
	moduleName := "SpaceTagRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceTagRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_SpaceTag, "", repo.clearCache2)

	return repo
}

func (c *spaceTagRepo) CreateTag(ctx context.Context, tag *domain.SpaceTag) error {
	po := convert.SpaceTagEntityToPo(tag)
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Create(po).Error
	return err
}

func (s *spaceTagRepo) SaveTag(ctx context.Context, tag *domain.SpaceTag) error {

	diffs := tag.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceTag{}
	mColumns := m.Cloumns()

	cloums := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_TagName:
			cloums[mColumns.TagName] = tag.TagName
		case domain.Diff_TagStatus:
			cloums[mColumns.TagStatus] = tag.TagStatus
		}
	}

	if len(cloums) == 0 {
		return nil
	}

	cloums[mColumns.UpdatedAt] = time.Now().Unix()
	err := s.data.DB(ctx).Model(&db.SpaceTag{}).Where("id=?", tag.Id).UpdateColumns(cloums).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *spaceTagRepo) FilterExistSpaceTagIds(ctx context.Context, spaceId int64, tagIds []int64) ([]int64, error) {
	var list []*db.SpaceTag
	err := s.data.RoDB(ctx).Model(&db.SpaceTag{}).Select("id").Where("space_id=? and id in ?", spaceId, tagIds).Find(&list).Error
	if err != nil {
		return nil, err
	}

	var existTagIds []int64
	for _, v := range list {
		existTagIds = append(existTagIds, v.Id)
	}

	return existTagIds, nil
}

func (c *spaceTagRepo) CheckTagNameIsExist(ctx context.Context, spaceId int64, tagName string) (bool, error) {
	var rowNum int64
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Where("space_id=? and BINARY tag_name=?", spaceId, tagName).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *spaceTagRepo) GetSpaceTag(ctx context.Context, spaceId int64, tagId int64) (*domain.SpaceTag, error) {
	var row *db.SpaceTag
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Where("id=? and space_id=?", tagId, spaceId).Take(&row).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceTagPoToEntity(row), err
}

func (c *spaceTagRepo) GetSpaceTags(ctx context.Context, spaceId int64, tagIds []int64) ([]*domain.SpaceTag, error) {
	var rows []*db.SpaceTag
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Where("space_id=? and id in ? ", spaceId, tagIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceTagPoToEntitys(rows), err
}

func (c *spaceTagRepo) DelSpaceTag(ctx context.Context, tagId int64) error {

	err := c.data.DB(ctx).Unscoped().Where("id=?", tagId).Delete(&db.SpaceTag{}).Error
	if err != nil {
		return err
	}

	return err
}

func (c *spaceTagRepo) DelSpaceTagBySpaceId(ctx context.Context, spaceId int64) error {
	err := c.data.DB(ctx).Unscoped().Where("space_id=?", spaceId).Delete(&db.SpaceTag{}).Error
	if err != nil {
		return err
	}
	return err
}

func (c *spaceTagRepo) TagMap(ctx context.Context, ids []int64) (map[int64]*domain.SpaceTag, error) {
	list, err := c.GetTagByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.SpaceTag) (int64, *domain.SpaceTag) {
		return t.Id, t
	})

	return m, nil
}

func (c *spaceTagRepo) GetTagByIds(ctx context.Context, tagIds []int64) ([]*domain.SpaceTag, error) {
	tagIds = stream.Unique(tagIds)
	tagIds = stream.Filter(tagIds, func(id int64) bool {
		return id != 0
	})

	if len(tagIds) == 0 {
		return nil, nil
	}

	fromRedis, err := c.GetTagByIdsFromRedis(ctx, tagIds)
	if err != nil {
		c.log.Error(err)
	}

	if len(tagIds) == len(fromRedis) {
		return fromRedis, nil
	}

	noCachedIds := stream.Filter(tagIds, func(id int64) bool {
		return !stream.ContainsFunc(fromRedis, func(tag *domain.SpaceTag) bool {
			return id == tag.Id
		})
	})

	fromDB, err := c.GetTagByIdsFromDB(ctx, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewTagKey(v.Id).Key()
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

func (p *spaceTagRepo) GetTagByIdsFromDB(ctx context.Context, tagIds []int64) ([]*domain.SpaceTag, error) {
	if len(tagIds) == 0 {
		return nil, nil
	}

	var rows []*db.SpaceTag
	err := p.data.DB(ctx).Where("id in (?)", tagIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceTagPoToEntitys(rows), err
}

func (c *spaceTagRepo) GetTagByIdsFromRedis(ctx context.Context, tagIds []int64) ([]*domain.SpaceTag, error) {

	keys := stream.Map(tagIds, func(v int64) string {
		return NewTagKey(v).Key()
	})

	result, err := c.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	list := stream.Map(result, func(v any) *domain.SpaceTag {
		var u domain.SpaceTag
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, err
}

//func (p *spaceTagRepo) clearCache(ctx context.Context, ids ...int64) {
//	keys := stream.Map(ids, func(v int64) string {
//		return NewTagKey(v).Key()
//	})
//
//	p.data.rdb.Del(ctx, keys...)
//}

func (p *spaceTagRepo) clearCache2(ids []int64) {
	keys := stream.Map(ids, func(v int64) string {
		return NewTagKey(v).Key()
	})

	p.data.rdb.Del(context.Background(), keys...)
}

func (c *spaceTagRepo) QSpaceTagList(ctx context.Context, spaceId int64) ([]*domain.SpaceTag, error) {
	var rows []*db.SpaceTag
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Where("space_id=? ", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return convert.SpaceTagPoToEntitys(rows), err
}

func (c *spaceTagRepo) GetSpaceTagCount(ctx context.Context, spaceId int64) (int64, error) {
	var count int64
	err := c.data.DB(ctx).Model(&db.SpaceTag{}).Where("space_id=? ", spaceId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}
