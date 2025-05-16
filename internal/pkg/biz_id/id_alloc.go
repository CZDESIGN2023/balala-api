package biz_id

import (
	"errors"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type idAlloc struct {
	businessId string     // 业务id
	ch         chan int64 // id缓冲池
	min, max   int64      // id段最小值，最大值
	log        *log.Helper
	db         *gorm.DB

	locker sync.Mutex
}

func (u *idAlloc) Get() (int64, error) {
	select {
	case <-time.After(1 * time.Second):
		return 0, errors.New("get uid timeout")
	case uid := <-u.ch:
		return uid, nil
	}
}

func (u *idAlloc) productId() {

	for {
		u.locker.Lock()
		if u.min >= u.max {
			u.reLoad()
		}
		u.min++
		u.locker.Unlock()
		u.ch <- u.min
	}
}

func (u *idAlloc) reLoad() error {

	var err error
	for {
		err = u.reloadFormDb()
		if err == nil {
			return nil
		}
		u.log.Error(err)
		time.Sleep(time.Second)
	}
}

func (u *idAlloc) reloadFormDb() error {

	var maxId int64
	var step int64

	tx := u.db.Begin()
	defer tx.Rollback()

	row := tx.Raw("SELECT max_id,step FROM business_id WHERE business_id = ? FOR UPDATE", u.businessId).Row()
	err := row.Scan(&maxId, &step)
	if err != nil {
		return err
	}

	err = tx.Exec("UPDATE business_id SET max_id = ? WHERE business_id = ?", maxId+step, u.businessId).Error
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	u.min = maxId
	u.max = maxId + step

	// fmt.Println("reload id:", u.businessId, " maxId:", maxId, " min:", u.min, " max:", u.max)

	return nil
}

func newIdAlloc(businessId string, db *gorm.DB) *idAlloc {
	idAlloc := &idAlloc{
		businessId: businessId,
		ch:         make(chan int64, 20),
		db:         db,
	}

	idAlloc.reLoad()
	go idAlloc.productId()
	return idAlloc
}
