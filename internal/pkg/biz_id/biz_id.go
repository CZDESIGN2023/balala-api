package biz_id

import (
	"context"
	"go-cs/internal/utils"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type BusinessId struct {
	Id   int64
	Type string
}

type BusinessIdService struct {
	idAllocMap sync.Map
	db         *gorm.DB
	log        *log.Helper
	once       sync.Once
}

func NewBusinessIdService(db *gorm.DB, logger log.Logger) *BusinessIdService {
	moduleName := "BusinessIdService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	s := &BusinessIdService{
		db:         db,
		log:        hlog,
		idAllocMap: sync.Map{},
	}

	return s
}

func (r *BusinessIdService) init() error {

	type ent struct {
		Id        int64  `gorm:"column:id"`
		TableName string `gorm:"column:table_name"`
	}

	ents := make([]*ent, 0)

	tx := r.db.Begin()
	defer tx.Rollback()

	err := tx.Raw("SELECT id, table_name  FROM business_id ").Find(&ents).Error
	if err != nil {
		return err
	}

	for _, ent := range ents {
		var maxId int64
		err = tx.Raw("select if(max(id), max(id), 0) from " + ent.TableName).Row().Scan(&maxId)
		if err != nil {
			return err
		}

		err = tx.Exec("update business_id set max_id = if(max_id>?, max_id, ?)  where id = ?  ", maxId, maxId, ent.Id).Error
		if err != nil {
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil

}

func (r *BusinessIdService) NewId(ctx context.Context, bizId string) *BusinessId {
	r.once.Do(func() {
		err := r.init()
		if err != nil {
			// 初始化失败，需要记录日志
			r.log.Errorf("init BusinessIdService failed: %v", err)
			panic("init BusinessIdService failed")
		}
	})

	ida, isOk := r.idAllocMap.Load(bizId)
	if !isOk {
		ida = newIdAlloc(bizId, r.db)
		r.idAllocMap.Store(bizId, ida)
	}

	id, err := ida.(*idAlloc).Get()
	if err != nil {
		return nil
	}
	return &BusinessId{Type: bizId, Id: id}
}
