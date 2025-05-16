package data

import (
	"context"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	domain "go-cs/internal/domain/space_file_info"
	repo "go-cs/internal/domain/space_file_info/repo"
)

type spaceFileInfoRepo struct {
	baseRepo
}

func NewSpaceFileInfoRepo(data *Data, logger log.Logger) repo.SpaceFileInfoRepo {
	moduleName := "SpaceFileInfoRelationRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &spaceFileInfoRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (c *spaceFileInfoRepo) CreateSpaceFileInfo(ctx context.Context, in *domain.SpaceFileInfo) error {
	po := convert.SpaceFileInfoEntityToPo(in)
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).Create(po).Error
	return err
}

func (c *spaceFileInfoRepo) GetSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) (*domain.SpaceFileInfo, error) {
	var row *db.SpaceFileInfo
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).Where("space_id=? and file_info_id=? and source_id=? and source_type=1", spaceId, fileInfoId, workItemId).Take(&row).Error
	return convert.SpaceFileInfoPoToEntity(row), err
}

func (c *spaceFileInfoRepo) GetSpaceWorkItemFileInfoByWorkItemIds(ctx context.Context, spaceId int64, workItemIds []int64) (domain.SpaceFileInfos, error) {
	var rows []*db.SpaceFileInfo
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).Where("space_id=? and source_id in ? and source_type=1", spaceId, workItemIds).Find(&rows).Error
	return convert.SpaceFileInfoPoSliceToEntitys(rows), err
}

func (c *spaceFileInfoRepo) GetSpaceWorkItemFileInfoById(ctx context.Context, id int64) (*domain.SpaceFileInfo, error) {
	var row *db.SpaceFileInfo
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).Where("id=?", id).Take(&row).Error
	return convert.SpaceFileInfoPoToEntity(row), err
}

func (c *spaceFileInfoRepo) SoftDelSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) error {
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).
		Where("space_id=? and file_info_id=? and source_id=? and source_type=1", spaceId, fileInfoId, workItemId).
		UpdateColumns(map[string]any{
			"status":     2,
			"updated_at": time.Now().Unix(),
			"deleted_at": time.Now().Unix(),
		}).Error
	return err
}

func (c *spaceFileInfoRepo) HardDelSpaceWorkItemFileInfo(ctx context.Context, fileInfoId int64, spaceId int64, workItemId int64) error {
	err := c.data.DB(ctx).Unscoped().
		Where("space_id=? and file_info_id=? and source_id=? and source_type=1", spaceId, fileInfoId, workItemId).
		Unscoped().
		Delete(&db.SpaceFileInfo{}).Error
	return err
}

func (c *spaceFileInfoRepo) SoftDelSpaceWorkItemsAllFile(ctx context.Context, workItemIds []int64) error {
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).
		Where("source_id in ? and source_type=1", workItemIds).
		UpdateColumns(map[string]any{
			"status":     2,
			"updated_at": time.Now().Unix(),
			"deleted_at": time.Now().Unix(),
		}).Error
	return err
}

func (c *spaceFileInfoRepo) SoftDelFileBySpaceId(ctx context.Context, spaceId int64) error {
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).
		Where("space_id=?", spaceId).
		UpdateColumns(map[string]any{
			"status":     2,
			"updated_at": time.Now().Unix(),
			"deleted_at": time.Now().Unix(),
		}).Error
	return err
}

func (c *spaceFileInfoRepo) CountWorkItemFileNum(ctx context.Context, workItemId int64) (int64, error) {
	var count int64
	err := c.data.DB(ctx).Model(&db.SpaceFileInfo{}).
		Where("source_id = ? and source_type=1 and status = 1 and deleted_at = 0", workItemId).
		Count(&count).Error
	return count, err
}

func (c *spaceFileInfoRepo) SaveFileDownToken(ctx context.Context, token string, fileInfoId int64, expSecond int64) error {
	_, err := c.data.rdb.SetEX(ctx, "file_download_token:file_id:"+token, fileInfoId, time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *spaceFileInfoRepo) GetFileDownToken(ctx context.Context, token string) (string, error) {
	tokenInfo, err := c.data.rdb.Get(ctx, "file_download_token:file_id:"+token).Result()
	if err != nil {
		return "", err
	}

	return tokenInfo, nil
}

func (c *spaceFileInfoRepo) QFileInfoList(ctx context.Context, spaceId int64, workItemId int64) ([]*domain.SpaceFileInfo, error) {
	var rows []*db.SpaceFileInfo

	err := c.data.RoDB(ctx).
		Model(&db.SpaceFileInfo{}).
		Where("space_id=? and source_type=1 and source_id = ? and status=1 ", spaceId, workItemId).
		Order("id desc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return convert.SpaceFileInfoPoSliceToEntitys(rows), err
}
