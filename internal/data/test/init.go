package test

import (
	"github.com/elastic/go-elasticsearch/v8"
	kconfig "github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/subosito/gotenv"
	"go-cs/internal/conf"
	v8 "go-cs/internal/utils/es/v8"
	"go-cs/internal/utils/orm/mysql"
	"gorm.io/gorm"
)

var (
	db       *gorm.DB
	redisCli *redis.Client
	es       *elasticsearch.Client //es
	bs       *conf.Bootstrap
)

func initConf(confPath, envPath string) *conf.Bootstrap {
	err := gotenv.Load(envPath)
	if err != nil {
		panic(err)
	}

	c := kconfig.New(
		kconfig.WithSource(
			file.NewSource(confPath),
			env.NewSource(),
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}

func initGorm(conf *conf.Data_Database) *gorm.DB {
	if conf == nil {
		return nil
	}

	gdb, err := mysql.New(mysql.Config{
		Driver: "mysql",
		DSN:    conf.Dsn,
	})
	if err != nil {
		panic(err)
	}

	if conf.Debug {
		gdb = gdb.Debug()
	}

	return gdb
}

func initRedis(conf *conf.Data_Redis) *redis.Client {
	if conf == nil {
		return nil
	}

	return redis.NewClient(&redis.Options{
		Addr: conf.Addr,
	})
}

func Init(configPath, envPath string, debug ...bool) {
	bs = initConf(configPath, envPath)

	if len(debug) > 0 && debug[0] == true {
		bs.Data.Database.Debug = true
		bs.Data.DatabaseRo.Debug = true
	}

	db = initGorm(bs.Data.Database)
	redisCli = initRedis(bs.Data.Redis)
	es, _, _ = v8.NewEsClient(v8.NewConfig(bs.Data), klog.DefaultLogger)
}

func GetDB() *gorm.DB {
	return db
}

func GetRedis() *redis.Client {
	return redisCli
}

func GetEs() *elasticsearch.Client {
	return es
}

func GetConf() *conf.Bootstrap {
	return bs
}
