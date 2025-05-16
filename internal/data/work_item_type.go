package data

import (
	"context"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/internal/utils/local_cache"
	"go-cs/pkg/stream"

	domain "go-cs/internal/domain/work_item_type"
	repo "go-cs/internal/domain/work_item_type/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type workItemTypeRepo struct {
	baseRepo
	cache *local_cache.Cache[string, any] // 缓存
}

func NewWorkItemTypeRepo(data *Data, logger log.Logger) repo.WorkItemTypeRepo {
	moduleName := "workItemTypeRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)
	repo := &workItemTypeRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cache: local_cache.NewCache[string, any](-1),
	}
	return repo
}

func (c *workItemTypeRepo) CreateWorkItemType(ctx context.Context, workItemType *domain.WorkItemType) error {
	po := convert.WorkItemTypeEntityToPo(workItemType)
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Create(po).Error
	return err
}

func (c *workItemTypeRepo) GetWorkItemType(ctx context.Context, id int64) (*domain.WorkItemType, error) {
	var row *db.WorkItemType
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Where("id=? ", id).Take(&row).Error
	if err != nil {
		return nil, err
	}

	ent := convert.WorkItemTypePoToEntity(row)
	return ent, err
}

func (c *workItemTypeRepo) DelBySpaceId(ctx context.Context, spaceId int64) error {

	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *workItemTypeRepo) GetWorkItemTypeBySpaceId(ctx context.Context, spaceId int64) ([]*domain.WorkItemType, error) {
	var rows []*db.WorkItemType
	err := c.data.RoDB(ctx).Model(&db.WorkItemType{}).Where("space_id=?", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.WorkItemTypePoToEntities(rows)
	return list, nil
}

func (c *workItemTypeRepo) IsExistByName(ctx context.Context, spaceId int64, name string) (bool, error) {

	var count int64
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Where("space_id=? and BINARY name=?", spaceId, name).Count(&count).Error
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

func (c *workItemTypeRepo) IsExistByKey(ctx context.Context, spaceId int64, key string) (bool, error) {

	var count int64
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Where("space_id=? and `key`=?", spaceId, key).Count(&count).Error
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

func (c *workItemTypeRepo) QWorkItemTypeInfo(ctx context.Context, req query.WorkItemTypeInfoQuery) (*query.WorkItemTypeInfoQueryResult, error) {

	cacheKey := fmt.Sprintf("work_item_type:list:%d", req.SpaceId)

	if v, isok := c.cache.Get(cacheKey); isok {
		return query.BuildWorkItemTypeInfoQueryResult(req.SpaceId, v.([]*domain.WorkItemType)), nil
	}

	list, err := c.GetWorkItemTypeBySpaceId(ctx, req.SpaceId)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, list)
	return query.BuildWorkItemTypeInfoQueryResult(req.SpaceId, list), err
}

func (c *workItemTypeRepo) GetWorkItemTypeByIds(ctx context.Context, ids []int64) ([]*domain.WorkItemType, error) {
	var list []*db.WorkItemType
	err := c.data.DB(ctx).Model(&db.WorkItemType{}).Where("id in ? ", ids).Find(&list).Error
	if err != nil {
		return nil, err
	}

	ret := convert.WorkItemTypePoToEntities(list)
	return ret, err
}

func (c *workItemTypeRepo) WorkItemTypeMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkItemType, error) {
	list, err := c.GetWorkItemTypeByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	return stream.ToMap(list, func(i int, t *domain.WorkItemType) (int64, *domain.WorkItemType) {
		return t.Id, t
	}), nil
}
