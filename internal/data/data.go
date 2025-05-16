package data

import (
	"context"
	"go-cs/internal/conf"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/cache"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/internal/utils/orm"
	"go-cs/internal/utils/sessions"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/elastic/go-elasticsearch/v8"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(

	wire.NewSet(
		NewData,
		NewReadonlyGormDB,
		NewTransaction,
	),
	wire.NewSet(orm.ProvideMongoDatabase, orm.ProvideMongoStoreFactory),
	wire.NewSet(esV8.NewEsClient, esV8.NewConfig),
	utils.NewSessionStore,

	wire.NewSet(orm.NewConfig, orm.NewOrm),
	wire.NewSet(cache.NewConfig, cache.NewRedis),
	wire.NewSet(NewUserRepo),

	wire.NewSet(
		NewConfigRepo,
		NewLoginRepo,
		NewSpaceRepo,
		NewSpaceMemberRepo,
		NewSpaceWorkObjectRepo,
		NewSpaceTagRepo,
		NewSearchRepo,
		NewFileInfoRepo,
		//NewStaticsRepo,
		NewStaticsEsRepo,
		NewSpaceFileInfoRepo,
		NewUserLoginLogRepo,
		NewOperLogRepo,
		// NewSpaceWorkItemFlowRepo,
		NewSpaceWorkItemCommentRepo,
		NewLogRepo,
		NewNotifyRepo,
		NewSpaceWorkVersionRepo,
		NewAdminRepo,
		NewWorkItemTypeRepo,
		NewSpaceViewRepo,
	),

	wire.NewSet(
		NewWorkItemStatusRepo,
		NewWorkFlowRepo,
		NewWorkItemRoleRepo,
		NewSpaceWorkItemRepo,
		NewSpaceWorkItemEsRepo,
	),
)

// Data .
type Data struct {
	db           *gorm.DB        // 数据库
	dbRo         *ReadonlyGormDB // 唯讀數據庫
	rdb          *redis.Client
	mongo        *mongo.Database
	es           *elasticsearch.Client //es
	log          *log.Helper
	sessionStore sessions.Store
	conf         *conf.Data
}

type contextTxKey struct{}

func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {

	//上下文已存在事务对象，则直接传递
	if d.hasTx(ctx) {
		return fn(ctx)
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func (d *Data) RoDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.dbRo.DB
}

func (d *Data) hasTx(ctx context.Context) bool {
	_, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	return ok
}

func NewTransaction(d *Data) trans.Transaction {
	return d
}

// NewData .
func NewData(dbw *gorm.DB, dbRo *ReadonlyGormDB, rdb *redis.Client, es *elasticsearch.Client, sessionStore sessions.Store, conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	moduleName := "im-service/data"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	d := &Data{
		db:           dbw,
		dbRo:         dbRo,
		rdb:          rdb,
		log:          hlog,
		sessionStore: sessionStore,
		es:           es,
		conf:         conf,
	}

	cleanup := func() {
		hlog.Info("closing the data resources")
	}

	return d, cleanup, nil
}

type ReadonlyGormDB struct {
	*gorm.DB
}

func NewReadonlyGormDB(c *conf.Data, logger log.Logger) (*ReadonlyGormDB, func(), error) {
	moduleName := "database_ro"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	dsn := c.DatabaseRo.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return &ReadonlyGormDB{}, nil, err
	}

	if c.DatabaseRo.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return &ReadonlyGormDB{}, nil, err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Second * 25)

	cleanup := func() {
		hlog.Info("closing the database resources")
		if err := sqlDB.Close(); err != nil {
			hlog.Error(err)
		}
	}

	return &ReadonlyGormDB{DB: db}, cleanup, err
}
