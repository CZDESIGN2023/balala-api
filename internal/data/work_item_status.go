package data

import (
	"context"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/work_item_status"
	repo "go-cs/internal/domain/work_item_status/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/local_cache"
	"go-cs/pkg/stream"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type workItemStatusRepo struct {
	baseRepo
	cache *local_cache.Cache[string, any] // 缓存
}

func NewWorkItemStatusRepo(data *Data, logger log.Logger) repo.WorkItemStatusRepo {
	moduleName := "workItemStatusRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &workItemStatusRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cache: local_cache.NewCache[string, any](-1),
	}
	return repo
}

func (c *workItemStatusRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max int64
	row := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return max, err
	}
	return max, nil
}

func (c *workItemStatusRepo) GetWorkItemStatusItem(ctx context.Context, id int64) (*domain.WorkItemStatusItem, error) {
	var po db.WorkItemStatus
	err := c.data.DB(ctx).Model(db.WorkItemStatus{}).Where("id=?", id).First(&po).Error
	if err != nil {
		return nil, err
	}

	item := convert.WorkItemStatusPoToEntity(&po)
	return item, nil
}

func (c *workItemStatusRepo) DelWorkItemStatusBySpaceId(ctx context.Context, spaceId int64) error {
	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}

	c.cache.Delete(fmt.Sprintf("WorkItemStatus:%d:0", spaceId))
	return nil
}

func (c *workItemStatusRepo) CreateWorkItemStatusItems(ctx context.Context, spaceId int64, items []*domain.WorkItemStatusItem) error {

	if len(items) == 0 {
		return nil
	}

	txErr := c.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range items {
			po := convert.WorkItemStatusEntityToPo(v)
			err := tx.Model(&db.WorkItemStatus{}).Create(po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return txErr
	}

	c.cache.Delete(fmt.Sprintf("WorkItemStatus:%d:0", spaceId))
	return nil
}

func (c *workItemStatusRepo) CreateWorkItemStatusItem(ctx context.Context, item *domain.WorkItemStatusItem) error {

	po := convert.WorkItemStatusEntityToPo(item)
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Create(po).Error
	if err != nil {
		return err
	}

	c.cache.Delete(fmt.Sprintf("WorkItemStatus:%d:0", item.SpaceId))
	return nil
}

func (c *workItemStatusRepo) SaveWorkItemStatusItem(ctx context.Context, item *domain.WorkItemStatusItem) error {

	diffs := item.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.WorkItemStatus{}
	mColumns := m.Cloumns()

	columns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			columns[mColumns.Ranking] = item.Ranking
		case domain.Diff_Name:
			columns[mColumns.Name] = item.Name
		case domain.Diff_Status:
			columns[mColumns.Status] = item.Status
		}
	}

	if len(columns) == 0 {
		return nil
	}

	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Where("id=?", item.Id).Updates(columns).Error
	if err != nil {
		return err
	}

	c.cache.Delete(fmt.Sprintf("WorkItemStatus:%d:0", item.SpaceId))

	return nil
}

func (c *workItemStatusRepo) GetWorkItemStatusItemsBySpace(ctx context.Context, spaceId int64) (domain.WorkItemStatusItems, error) {
	var list []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(db.WorkItemStatus{}).Where("space_id=? and work_item_type_id = 0", spaceId).Order("ranking desc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WorkItemStatusItem, 0)
	for _, v := range list {
		items = append(items, convert.WorkItemStatusPoToEntity(v))
	}
	return items, err
}

func (c *workItemStatusRepo) WorkItemStatusMap(ctx context.Context, spaceId int64) (map[int64]*domain.WorkItemStatusItem, error) {
	var list []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(db.WorkItemStatus{}).Where("space_id=? and work_item_type_id = 0", spaceId).Find(&list).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WorkItemStatusItem, 0)
	for _, v := range list {
		items = append(items, convert.WorkItemStatusPoToEntity(v))
	}

	return stream.ToMap(items, func(i int, v *domain.WorkItemStatusItem) (int64, *domain.WorkItemStatusItem) {
		return v.Id, v
	}), nil
}

func (c *workItemStatusRepo) WorkItemStatusKeyMap(ctx context.Context, spaceId int64, keys ...string) (map[string]*domain.WorkItemStatusItem, error) {
	var list []*db.WorkItemStatus
	tx := c.data.DB(ctx).Model(db.WorkItemStatus{}).Where("space_id=? and work_item_type_id = 0", spaceId)
	if len(keys) > 0 {
		tx = tx.Where("`key` in ?", keys)
	}
	err := tx.Find(&list).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WorkItemStatusItem, 0)
	for _, v := range list {
		items = append(items, convert.WorkItemStatusPoToEntity(v))
	}

	return stream.ToMap(items, func(i int, v *domain.WorkItemStatusItem) (string, *domain.WorkItemStatusItem) {
		return v.Key, v
	}), nil
}

func (c *workItemStatusRepo) GetWorkItemStatusItemsBySpaceIds(ctx context.Context, spaceId []int64) (domain.WorkItemStatusItems, error) {
	var list []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(db.WorkItemStatus{}).Where("space_id in ? ", spaceId).Order("ranking desc").Find(&list).Error
	if err != nil {
		return nil, err
	}

	items := make([]*domain.WorkItemStatusItem, 0)
	for _, v := range list {
		items = append(items, convert.WorkItemStatusPoToEntity(v))
	}
	return items, err
}

func (c *workItemStatusRepo) GetWorkItemStatusInfo(ctx context.Context, spaceId int64) (*domain.WorkItemStatusInfo, error) {

	cacheKey := fmt.Sprintf("WorkItemStatus:%d:0", spaceId)
	if v, isok := c.cache.Get(cacheKey); isok {
		statusInfo := domain.BuildWorkItemStatusInfo(spaceId, v.(domain.WorkItemStatusItems))
		return statusInfo, nil
	}

	list, err := c.GetWorkItemStatusItemsBySpace(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	c.cache.Set(cacheKey, list)

	statusInfo := domain.BuildWorkItemStatusInfo(spaceId, list)

	return statusInfo, nil
}

func (c *workItemStatusRepo) GetWorkItemStatusInfoBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.WorkItemStatusInfo, error) {
	items, err := c.GetWorkItemStatusItemsBySpaceIds(ctx, spaceIds)
	if err != nil {
		return nil, err
	}

	itemGroup := stream.GroupBy(items, func(v *domain.WorkItemStatusItem) int64 {
		return v.SpaceId
	})

	infos := make([]*domain.WorkItemStatusInfo, 0)
	for _, v := range itemGroup {
		info := domain.BuildWorkItemStatusInfo(v[0].SpaceId, v)
		infos = append(infos, info)
	}
	return infos, nil
}

func (c *workItemStatusRepo) IsExistByWorkItemStatusName(ctx context.Context, spaceId int64, roleName string, flowScope consts.FlowScope) (bool, error) {

	var count int64
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Where("space_id=? AND flow_scope=? AND BINARY name=?", spaceId, flowScope, roleName).Count(&count).Error
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

func (c *workItemStatusRepo) QSpaceWorkItemStatusList(ctx context.Context, spaceIds []int64, scope consts.FlowScope) ([]*domain.WorkItemStatusItem, error) {
	var rows []*db.WorkItemStatus

	tx := c.data.DB(ctx).Model(&db.WorkItemStatus{})
	if scope != consts.FlowScope_All {
		tx = tx.Where("flow_scope IN ?", []consts.FlowScope{consts.FlowScope_All, scope})
	}

	err := tx.Where("space_id IN ?", spaceIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.WorkItemStatusPoToEntitys(rows)
	return list, nil
}

func (c *workItemStatusRepo) QSpaceWorkItemStatusById(ctx context.Context, spaceId int64, ids []int64) ([]*domain.WorkItemStatusItem, error) {
	var rows []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Where("space_id=? and id in ?", spaceId, ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.WorkItemStatusPoToEntitys(rows)
	return list, nil
}

func (c *workItemStatusRepo) DelSpaceWorkItemStatusItem(ctx context.Context, spaceId int64, id int64) error {
	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Unscoped().Where("id=? and space_id=?", id, spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	c.cache.Delete(fmt.Sprintf("WorkItemStatus:%d:0", spaceId))
	return nil
}

func (c *workItemStatusRepo) StatusMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkItemStatusItem, error) {
	var rows []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Where("id in ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.WorkItemStatusPoToEntitys(rows)
	return stream.ToMap(list, func(_ int, v *domain.WorkItemStatusItem) (int64, *domain.WorkItemStatusItem) {
		return v.Id, v
	}), nil
}

func (c *workItemStatusRepo) StatusMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64]*domain.WorkItemStatusItem, error) {
	var rows []*db.WorkItemStatus
	err := c.data.DB(ctx).Model(&db.WorkItemStatus{}).Where("space_id in ?", spaceIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.WorkItemStatusPoToEntitys(rows)
	return stream.ToMap(list, func(_ int, v *domain.WorkItemStatusItem) (int64, *domain.WorkItemStatusItem) {
		return v.Id, v
	}), nil
}
