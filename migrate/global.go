package migrate

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils/errs"
	"go-cs/migrate/migrations"
	"go-cs/migrate/pkg"
	"gorm.io/gorm"
)

type Global struct {
	logger     log.Logger
	log        *log.Helper
	db         *gorm.DB
	migrations migrations.Migrations
}

func NewGlobal(logger log.Logger, gdb *gorm.DB, migrations migrations.Migrations) (*Global, error) {
	return &Global{
		logger:     logger,
		log:        log.NewHelper(logger, log.WithMessageKey("Migrate Global")),
		db:         gdb,
		migrations: migrations,
	}, nil
}

func (g *Global) MigrateTo(dstVerStr string) error {
	if dstVerStr == "" {
		return errors.New("destination version is empty")
	}

	dstVer := pkg.NewVersion(dstVerStr)
	srcVer := pkg.NewVersion(getSrcVersion(g.db))

	g.log.Info("srcVer: ", srcVer, ", dstVer: ", dstVer)

	cmpResult := srcVer.CompareTo(dstVer)

	// srcVer == dstVer，说明数据库版本号等于目标版本号，直接返回
	if cmpResult == 0 {
		g.log.Info("srcVer: ", srcVer, ", dstVer: ", dstVer, ", version equal, skip")
		return nil
	}

	// srcVer > dstVer，说明数据库版本号比目标版本号高，当前应用的版本低于数据库版本，直接终止迁移
	if cmpResult == 1 {
		return fmt.Errorf("source version (%v) is higher than destination version (%v)", srcVer, dstVer)
	}

	m := pkg.NewManager(g.logger, g.db, g.migrations, srcVer, dstVer)

	// 检查是否满足最低版本要求
	if migrations.RequiredMinVersion.CompareTo(srcVer) == 1 {
		return fmt.Errorf("required minimum version: %v, actual version: %v", migrations.RequiredMinVersion, dstVer)
	}

	return m.Migrate()
}

func getSrcVersion(gdb *gorm.DB) string {
	conf := db.Config{}

	var srcVer string
	res := gdb.Model(&db.Config{}).Where("config_key = ?", consts.CONFIG_BALALA_VERSION).Find(&conf)
	if res.Error != nil && !errs.IsDbRecordNotFoundErr(res.Error) {
		panic(res.Error)
	}

	if conf.ConfigValue != "" {
		srcVer = conf.ConfigValue
	}

	return srcVer
}
