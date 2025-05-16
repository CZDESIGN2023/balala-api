package data

import (
	"context"
	"errors"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/data/convert"
	"go-cs/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	domain "go-cs/internal/domain/file_info"
	repo "go-cs/internal/domain/file_info/repo"
)

type FileInfoRepo struct {
	baseRepo
}

func NewFileInfoRepo(data *Data, logger log.Logger) repo.FileInfoRepo {
	moduleName := "FileInfoRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &FileInfoRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}
}

func (r *FileInfoRepo) CreateFileInfo(ctx context.Context, info *domain.FileInfo) error {

	po := convert.FileInfoEntityToPo(info)
	err := r.data.DB(ctx).Model(&db.FileInfo{}).Create(po).Error
	return err
}

func (r *FileInfoRepo) GetFileInfo(ctx context.Context, id int64) (*domain.FileInfo, error) {
	var row *db.FileInfo
	err := r.data.RoDB(ctx).Model(&db.FileInfo{}).Where("id=?", id).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return convert.FileInfoPoToEntity(row), nil
}

func (r *FileInfoRepo) GetFileInfoByOwner(ctx context.Context, id int64, ownerId int64) (*domain.FileInfo, error) {
	var row *db.FileInfo
	err := r.data.RoDB(ctx).Model(&db.FileInfo{}).Where("id=? and owner=?", id, ownerId).Take(&row).Error
	if err != nil {
		return nil, err
	}
	return convert.FileInfoPoToEntity(row), nil
}

func (r *FileInfoRepo) GetFileInfoByIds(ctx context.Context, ids []int64) ([]*domain.FileInfo, error) {
	var list []*db.FileInfo
	err := r.data.RoDB(ctx).Model(&db.FileInfo{}).Where("id in ?", ids).Find(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return convert.FileInfoPoToEntities(list), nil
}
