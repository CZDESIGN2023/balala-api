package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	domain "go-cs/internal/domain/space_view"
	repo "go-cs/internal/domain/space_view/repo"
)

type spaceViewRepo struct {
	baseRepo
}

func NewSpaceViewRepo(data *Data, logger log.Logger) repo.SpaceViewRepo {
	moduleName := "SpaceViewRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceViewRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	return repo
}

func (c *spaceViewRepo) CreateUserView(ctx context.Context, item *domain.SpaceUserView) error {

	po := convert.SpaceUserViewEntityToPo(item)
	err := c.data.DB(ctx).Model(&db.SpaceUserView{}).Create(po).Error
	if err != nil {
		c.log.Error(err)
		return err
	}

	item.Id = po.Id
	return err
}

func (c *spaceViewRepo) CreateUserViews(ctx context.Context, items []*domain.SpaceUserView) error {
	po := convert.SpaceUserViewEntityToPos(items)

	err := c.data.DB(ctx).Model(&db.SpaceUserView{}).CreateInBatches(po, 100).Error
	if err != nil {
		c.log.Error(err)
		return err
	}

	return err
}

func (c *spaceViewRepo) CreateGlobalView(ctx context.Context, item *domain.SpaceGlobalView) error {

	po := convert.SpaceGlobalViewEntityToPo(item)
	err := c.data.DB(ctx).Model(&db.SpaceGlobalView{}).Create(po).Error
	if err != nil {
		c.log.Error(err)
		return err
	}

	item.Id = po.Id
	return err
}

func (c *spaceViewRepo) CreateGlobalViews(ctx context.Context, items []*domain.SpaceGlobalView) error {
	po := convert.SpaceGlobalViewEntityToPos(items)

	err := c.data.DB(ctx).Model(&db.SpaceGlobalView{}).CreateInBatches(po, 100).Error
	if err != nil {
		c.log.Error(err)
		return err
	}

	return err
}

func (c *spaceViewRepo) UserViewList(ctx context.Context, userId int64, spaceIds []int64, key string) ([]*domain.SpaceUserView, error) {
	tx := c.data.DB(ctx).Model(&db.SpaceUserView{}).Where("user_id = ?", userId)

	if len(spaceIds) > 0 {
		tx = tx.Where("space_id IN ?", spaceIds)
	}

	if len(key) > 0 {
		tx = tx.Where("`key` = ?", key)
	}

	var rows []*db.SpaceUserView
	err := tx.
		Order("ranking desc").
		Order("id asc").
		Find(&rows).Error
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	outerIds := stream.Map(rows, func(t *db.SpaceUserView) int64 {
		return t.OuterId
	})

	outerIds = stream.Of(outerIds).Unique().List()
	globalViewMap, err := c.GetGlobalViewMapByIds(ctx, outerIds)
	if err != nil {
		c.log.Error(err)
	}

	rets := convert.SpaceUserViewPoToEntities(rows)
	for _, v := range rets {
		v.SetGlobalView(globalViewMap[v.OuterId])
	}

	return rets, err
}

func (c *spaceViewRepo) SaveSpaceUserView(ctx context.Context, workVersion *domain.SpaceUserView) error {

	diffs := workVersion.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceUserView{}
	columns := m.Cloumns()

	updateColumns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Ranking:
			updateColumns[columns.Ranking] = workVersion.Ranking
		case domain.Diff_Name:
			updateColumns[columns.Name] = workVersion.Name
		case domain.Diff_Status:
			updateColumns[columns.Status] = workVersion.Status
		case domain.Diff_QueryConfig:
			updateColumns[columns.QueryConfig] = workVersion.QueryConfig
		case domain.Diff_TableConfig:
			updateColumns[columns.TableConfig] = workVersion.TableConfig
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[columns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", workVersion.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceViewRepo) SaveSpaceGlobalView(ctx context.Context, workVersion *domain.SpaceGlobalView) error {

	diffs := workVersion.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceGlobalView{}
	columns := m.Cloumns()

	updateColumns := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_Name:
			updateColumns[columns.Name] = workVersion.Name
		case domain.Diff_QueryConfig:
			updateColumns[columns.QueryConfig] = workVersion.QueryConfig
		case domain.Diff_TableConfig:
			updateColumns[columns.TableConfig] = workVersion.TableConfig
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[columns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", workVersion.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *spaceViewRepo) DeleteUserViewById(ctx context.Context, id int64) error {
	res := c.data.DB(ctx).
		Where("id = ?", id).
		Delete(&db.SpaceUserView{})

	return res.Error
}

func (c *spaceViewRepo) DeleteUserViewByUserId(ctx context.Context, userId, spaceId int64) error {
	res := c.data.DB(ctx).
		Where("user_id = ? AND space_id = ?", userId, spaceId).
		Delete(&db.SpaceUserView{})

	return res.Error
}

func (c *spaceViewRepo) DeleteUserViewBySpaceId(ctx context.Context, spaceId int64) error {
	res := c.data.DB(ctx).
		Where("space_id = ?", spaceId).
		Delete(&db.SpaceUserView{})

	return res.Error
}

func (c *spaceViewRepo) DeleteUserViewByOuterId(ctx context.Context, outerId int64) error {
	res := c.data.DB(ctx).
		Where("outer_id = ?", outerId).
		Delete(&db.SpaceUserView{})

	return res.Error
}

func (c *spaceViewRepo) DeleteGlobalViewById(ctx context.Context, id int64) error {
	res := c.data.DB(ctx).
		Where("id = ?", id).
		Delete(&db.SpaceGlobalView{})

	return res.Error
}

func (c *spaceViewRepo) DeleteGlobalViewBySpaceId(ctx context.Context, spaceId int64) error {
	res := c.data.DB(ctx).
		Where("space_id = ?", spaceId).
		Delete(&db.SpaceGlobalView{})

	return res.Error
}

func (c *spaceViewRepo) GetMaxRanking(ctx context.Context, spaceId int64) (int64, error) {
	var max any
	row := c.data.DB(ctx).Model(&db.SpaceUserView{}).Select("MAX(ranking)").Where("space_id=?", spaceId).Row()
	err := row.Scan(&max)
	if err != nil {
		return cast.ToInt64(max), err
	}
	return cast.ToInt64(max), nil
}

func (c *spaceViewRepo) GetUserViewById(ctx context.Context, id int64) (*domain.SpaceUserView, error) {
	var row *db.SpaceUserView
	err := c.data.RoDB(ctx).Model(&db.SpaceUserView{}).Where("id=?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}

	ret := convert.SpaceUserViewPoToEntity(row)

	if row.OuterId != 0 {
		globalView, _ := c.GetGlobalViewById(ctx, row.OuterId)
		ret.SetGlobalView(globalView)
	}

	return ret, nil
}

func (c *spaceViewRepo) GetGlobalViewById(ctx context.Context, id int64) (*domain.SpaceGlobalView, error) {
	var row *db.SpaceGlobalView
	err := c.data.RoDB(ctx).Model(&db.SpaceGlobalView{}).Where("id = ?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceGlobalViewPoToEntity(row), nil
}

func (c *spaceViewRepo) GetGlobalViewByIds(ctx context.Context, ids []int64) ([]*domain.SpaceGlobalView, error) {
	var row []*db.SpaceGlobalView
	err := c.data.RoDB(ctx).Model(&db.SpaceGlobalView{}).Where("id IN ?", ids).Find(&row).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceGlobalViewPoToEntities(row), nil
}

func (c *spaceViewRepo) GetGlobalViewMapByIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceGlobalView, error) {
	var row []*db.SpaceGlobalView
	err := c.data.RoDB(ctx).Model(&db.SpaceGlobalView{}).Where("id IN ?", ids).Find(&row).Error
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(convert.SpaceGlobalViewPoToEntities(row), func(_ int, t *domain.SpaceGlobalView) (int64, *domain.SpaceGlobalView) {
		return t.Id, t
	})

	return m, nil
}

func (c *spaceViewRepo) GetGlobalViewBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.SpaceGlobalView, error) {
	var rows []*db.SpaceGlobalView
	err := c.data.RoDB(ctx).Model(&db.SpaceGlobalView{}).Where("space_id IN ?", spaceIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceGlobalViewPoToEntities(rows), nil
}

func (c *spaceViewRepo) GetGlobalViewMap(ctx context.Context, spaceId int64) (map[int64]*domain.SpaceGlobalView, error) {
	list, err := c.GetGlobalViewList(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(_ int, t *domain.SpaceGlobalView) (int64, *domain.SpaceGlobalView) {
		return t.Id, t
	})

	return m, err
}

func (c *spaceViewRepo) GetGlobalViewList(ctx context.Context, spaceId int64) ([]*domain.SpaceGlobalView, error) {
	var rows []*db.SpaceGlobalView
	err := c.data.RoDB(ctx).Model(&db.SpaceGlobalView{}).Where("space_id=?", spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	list := convert.SpaceGlobalViewPoToEntities(rows)

	return list, err
}

func (c *spaceViewRepo) GetUserViewMap(ctx context.Context, userId, spaceId int64) (map[int64]*domain.SpaceUserView, error) {
	var rows []*db.SpaceUserView
	err := c.data.RoDB(ctx).Model(&db.SpaceUserView{}).Where("space_id=? AND user_id = ?", spaceId, userId).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	list := convert.SpaceUserViewPoToEntities(rows)

	globalViewMap, err := c.GetGlobalViewMap(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(_ int, t *domain.SpaceUserView) (int64, *domain.SpaceUserView) {
		t.SetGlobalView(globalViewMap[t.OuterId])
		return t.Id, t
	})

	return m, err
}

func (c *spaceViewRepo) GetUserViewMapByIds(ctx context.Context, ids []int64) (map[int64]*domain.SpaceUserView, error) {
	var rows []*db.SpaceUserView
	err := c.data.RoDB(ctx).Model(&db.SpaceUserView{}).Where("id IN ?", ids).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	globalViewIds := stream.Map(rows, func(t *db.SpaceUserView) int64 {
		return t.OuterId
	})

	globalViewIds = stream.Filter(globalViewIds, func(t int64) bool {
		return t != 0
	})

	globalView, er := c.GetGlobalViewMapByIds(ctx, globalViewIds)
	if er != nil {
		c.log.Error(er)
	}

	list := convert.SpaceUserViewPoToEntities(rows)
	m := stream.ToMap(list, func(_ int, t *domain.SpaceUserView) (int64, *domain.SpaceUserView) {
		t.SetGlobalView(globalView[t.OuterId])
		return t.Id, t
	})

	return m, err
}
