package river

import (
	"context"
	"errors"
	"fmt"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-mysql-org/go-mysql/canal"
	mysql2 "github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-sql-driver/mysql"
)

type River struct {
	ctx    context.Context
	cancel context.CancelCauseFunc

	Schema string
	Tables []string

	master            *masterInfo
	canal             atomic.Pointer[canal.Canal]
	canalConfig       *canal.Config
	canalEventHandler canal.EventHandler

	syncCh    chan posSaver
	closeOnce sync.Once

	log *log.Helper
}

func (r *River) WaitUntilInitialized() {
	time.Sleep(time.Second * 5) // 等待 river 初始化完成
	return
}

func New(conf *conf.Data, logger log.Logger) *River {
	var err error

	var r = &River{
		syncCh: make(chan posSaver, 4096),
	}

	dsn := conf.Database.Dsn
	_, helper := utils.InitModuleLogger(logger, "River")

	parseDSN, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	dataDir := "./configs/var"
	if r.master, err = loadMasterInfo(dataDir); err != nil {
		panic(err)
	}

	cfg := canal.NewDefaultConfig()
	cfg.Addr = parseDSN.Addr
	cfg.User = parseDSN.User
	cfg.Password = parseDSN.Passwd
	cfg.Dump.ExecutionPath = ""

	r.canalConfig = cfg
	r.canalEventHandler = &MyEventHandler{r: r}
	r.ctx, r.cancel = context.WithCancelCause(context.Background())
	r.log = helper
	r.Schema = parseDSN.DBName
	r.Tables = []string{
		"space_work_item_v2",
		"space_work_version",
		"space_work_object",
		"space_member",
		"space_tag",
		"space",
		"user",
	}

	return r
}

func isValidPos(pos mysql2.Position) bool {
	return pos.Name != ""
}

func (r *River) matchTable(schema, table string) bool {
	return schema == r.Schema && slices.Contains(r.Tables, table)
}

func (r *River) initCanal() {
	c, err := canal.NewCanal(r.canalConfig)
	if err != nil {
		panic(err)
	}

	c.SetEventHandler(r.canalEventHandler)

	r.canal.Store(c)
}

func (r *River) closeAndResetCanal() {
	r.Canal().Close()
	r.initCanal()
}

func (r *River) Canal() *canal.Canal {
	return r.canal.Load()
}

func (r *River) Start() {
	r.initCanal()

	if !isValidPos(r.master.Position()) {
		pos, err := r.Canal().GetMasterPos()
		if err != nil {
			panic(err)
		}

		err = r.master.Save(pos)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		if err := recover(); err != nil {
			r.log.Error(err, string(debug.Stack()))
		}

		for {
			select {
			case <-r.ctx.Done():
				return
			default:
			}

			err := r.Canal().RunFrom(r.master.Position())
			if err != nil {
				r.log.Error(err, string(debug.Stack()))

				// 如果是binlog位置错误，就从最新位置开始
				if IsBinlogPosWrongError(err) {
					oldPos := r.master.Position()

					pos, err := r.Canal().GetMasterPos()
					if err != nil {
						panic(err)
					}

					err = r.master.Save(pos)
					if err != nil {
						panic(err)
					}

					r.log.Infof("binlog position wrong, reset from %s to %s", oldPos, pos)
				}
			}

			time.Sleep(5 * time.Second)

			// 重新创建canal实例
			r.closeAndResetCanal()
		}
	}()

	r.Canal().WaitDumpDone()

	go r.SyncLoop()
}

func (r *River) Stop() {
	r.closeOnce.Do(func() {
		r.cancel(errors.New("river stopped"))
		r.Canal().Close()
		r.master.Close()
	})
}

func (r *River) SyncLoop() {
	defer func() {
		if err := recover(); err != nil {
			r.log.Error(err, string(debug.Stack()))
		}
	}()

	lastSavedTime := time.Now()
	var pos mysql2.Position

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		needSavePos := false

		select {
		case <-r.ctx.Done():
			return
		case v := <-r.syncCh:
			now := time.Now()
			pos = v.pos
			if v.force || now.Sub(lastSavedTime) > 3*time.Second {
				lastSavedTime = now
				needSavePos = true
			}
		case <-ticker.C:
			if isValidPos(pos) && r.master.Position().Compare(pos) != 0 {
				needSavePos = true
			}
		}

		if needSavePos {
			if err := r.master.Save(pos); err != nil {
				log.Errorf("save sync position %s err %v, close sync", pos, err)
				r.cancel(fmt.Errorf("save sync position err %w", err))
				return
			}
		}
	}
}

// IsBinlogPosWrongError 判断是否是binlog位置错误
func IsBinlogPosWrongError(err error) bool {
	return strings.Contains(err.Error(), "ERROR 1236 (HY000)")
}
