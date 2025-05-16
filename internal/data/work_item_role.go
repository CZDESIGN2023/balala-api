package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/internal/utils/local_cache"
	"go-cs/pkg/stream"

	domain "go-cs/internal/domain/work_item_role"
	repo "go-cs/internal/domain/work_item_role/repo"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type workItemRoleRepo struct {
	baseRepo
	cacheConfig *local_cache.Cache[string, *db.Config] // 配置信息缓存
}

func NewWorkItemRoleRepo(data *Data, logger log.Logger) repo.WorkItemRoleRepo {
	moduleName := "workItemRoleRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)
	repo := &workItemRoleRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cacheConfig: local_cache.NewCache[string, *db.Config](-1),
	}
	return repo
}

func (r *workItemRoleRepo) CreateWorkItemRole(ctx context.Context, role *domain.WorkItemRole) error {
	po := convert.WorkItemRoleEntityToPo(role)
	err := r.data.DB(ctx).Model(&db.WorkItemRole{}).Create(po).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *workItemRoleRepo) CreateWorkItemRoles(ctx context.Context, roles domain.WorkItemRoles) error {

	err := r.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, role := range roles {
			po := convert.WorkItemRoleEntityToPo(role)
			err := r.data.DB(ctx).Model(&db.WorkItemRole{}).Create(po).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *workItemRoleRepo) GetWorkItemRole(ctx context.Context, id int64) (*domain.WorkItemRole, error) {
	var row db.WorkItemRole
	err := c.data.DB(ctx).Where("id=? ", id).Take(&row).Error
	if err != nil {
		return nil, err
	}

	return convert.WorkItemRolePoToEntity(&row), err
}

func (c *workItemRoleRepo) GetWorkItemRoles(ctx context.Context, spaceId int64) (domain.WorkItemRoles, error) {
	var rows []*db.WorkItemRole
	err := c.data.DB(ctx).Where("space_id=? ", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.WorkItemRole
	for _, row := range rows {
		list = append(list, convert.WorkItemRolePoToEntity(row))
	}

	return list, err
}

func (c *workItemRoleRepo) QWorkItemRoleList(ctx context.Context, spaceId int64) ([]*domain.WorkItemRole, error) {
	var rows []*db.WorkItemRole
	err := c.data.DB(ctx).Model(db.WorkItemRole{}).Where("space_id = ? ", spaceId).Order("ranking desc").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.WorkItemRole
	for _, row := range rows {
		list = append(list, convert.WorkItemRolePoToEntity(row))
	}

	return list, err
}

func (c *workItemRoleRepo) DelWorkItemRoleBySpaceId(ctx context.Context, spaceId int64) error {

	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.WorkItemRole{}).Unscoped().Where("space_id=?", spaceId).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *workItemRoleRepo) DelWorkItemRole(ctx context.Context, id int64) error {

	var opValue = make(map[string]interface{})
	err := c.data.DB(ctx).Model(&db.WorkItemRole{}).Unscoped().Where("id=?", id).Delete(&opValue).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *workItemRoleRepo) QSpaceWorkItemRoleList(ctx context.Context, req *query.SpaceWorkItemRoleQuery) (result *query.SpaceWorkItemRoleQueryResult, err error) {
	tx := c.data.DB(ctx).Model(&db.WorkItemRole{})

	if req.FlowScope != consts.FlowScope_All {
		tx = tx.Where("flow_scope=?", req.FlowScope)
	}

	var rows []*db.WorkItemRole
	err = tx.Where("space_id=?", req.SpaceId).Order("ranking desc").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result = &query.SpaceWorkItemRoleQueryResult{}
	result.Total = int64(len(rows))
	for _, row := range rows {
		result.List = append(result.List, &query.SpaceWorkItemRoleQueryResult_ListItem{
			Id:             row.Id,
			Name:           row.Name,
			Ranking:        row.Ranking,
			SpaceId:        row.SpaceId,
			WorkItemTypeId: row.WorkItemTypeId,
			Key:            row.Key,
			Status:         row.Status,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
			IsSys:          int64(row.IsSys),
			FlowScope:      row.FlowScope,
		})
	}

	return result, nil
}

func (c *workItemRoleRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max int64

	row := c.data.DB(ctx).Model(&db.WorkItemRole{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return max, err
	}

	return max, nil
}

func (c *workItemRoleRepo) SaveWorkItemRole(ctx context.Context, role *domain.WorkItemRole) error {

	diffs := role.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.WorkItemRole{}
	mColumns := m.Cloumns()

	columns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			columns[mColumns.Ranking] = role.Ranking
		case domain.Diff_Name:
			columns[mColumns.Name] = role.Name
		}
	}

	if len(columns) == 0 {
		return nil
	}

	err := c.data.DB(ctx).Model(m).Where("id=?", role.Id).Updates(columns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *workItemRoleRepo) IsExistByRoleName(ctx context.Context, spaceId int64, roleName string, flowScope consts.FlowScope) (bool, error) {

	var count int64
	err := c.data.DB(ctx).Model(&db.WorkItemRole{}).Where("space_id=? AND flow_scope=? AND BINARY name=?", spaceId, flowScope, roleName).Count(&count).Error
	if err != nil {
		return true, err
	}

	return count > 0, nil
}

func (c *workItemRoleRepo) WorkItemRoleMap(ctx context.Context, spaceId int64) (map[int64]*domain.WorkItemRole, error) {
	var rows []*db.WorkItemRole
	err := c.data.DB(ctx).Where("space_id=? ", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.WorkItemRole
	for _, row := range rows {
		list = append(list, convert.WorkItemRolePoToEntity(row))
	}

	return stream.ToMap(list, func(i int, t *domain.WorkItemRole) (int64, *domain.WorkItemRole) {
		return t.Id, t
	}), nil
}
