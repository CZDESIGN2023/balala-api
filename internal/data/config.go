package data

import (
	"context"
	"github.com/tidwall/gjson"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/local_cache"
	"go-cs/pkg/stream"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type configRepo struct {
	baseRepo
	cacheConfig *local_cache.Cache[string, *db.Config] // 配置信息缓存
}

func NewConfigRepo(data *Data, logger log.Logger) biz.ConfigRepo {
	moduleName := "ConfigRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &configRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		cacheConfig: local_cache.NewCache[string, *db.Config](-1),
	}
	return repo
}

// GetById 依照id取得配置(from db_ro)
func (repo *configRepo) GetById(ctx context.Context, id int64) (*db.Config, error) {
	dbModel := repo.data.RoDB(ctx).Where("id = ?", id)
	record := &db.Config{}
	err := dbModel.Take(record).Error

	return record, err
}

// GetByKey 依照key取得配置(from db_ro)
func (repo *configRepo) GetByKey(ctx context.Context, ConfigKey string) (*db.Config, error) {
	dbModel := repo.data.RoDB(ctx).Where("config_key = ?", ConfigKey)
	record := &db.Config{}
	err := dbModel.Take(record).Error

	return record, err
}

func (repo *configRepo) List(ctx context.Context) ([]*db.Config, error) {
	var list []*db.Config
	err := repo.data.RoDB(ctx).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, err
}

func (repo *configRepo) defaultMap() map[string]*db.Config {
	return map[string]*db.Config{
		consts.CONFIG_BALALA_REGISTER_ENTRY: {
			ConfigKey:   consts.CONFIG_BALALA_REGISTER_ENTRY,
			ConfigValue: "0",
		},
		consts.CONFIG_BALALA_ATTACH: {
			ConfigKey: consts.CONFIG_BALALA_ATTACH,
			ConfigValue: utils.ToJSON(map[string]any{
				"value": "100",
				"unit":  "MB",
			}),
		},
	}
}

func (repo *configRepo) Map(ctx context.Context) (map[string]*db.Config, error) {
	list, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *db.Config) (string, *db.Config) {
		return t.ConfigKey, t
	})

	// 填充默认值
	for _, v := range repo.defaultMap() {
		if _, ok := m[v.ConfigKey]; !ok {
			m[v.ConfigKey] = v
		}
	}

	return m, nil
}

func (repo *configRepo) UpdateByKey(ctx context.Context, key string, value string) error {

	err := repo.data.DB(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.
			Model(&db.Config{}).
			Where("config_key = ?", key).
			Update("config_value", value)

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			tx.Create(&db.Config{
				ConfigName:   key,
				ConfigKey:    key,
				ConfigValue:  value,
				ConfigStatus: "1",
			})
		}

		return nil
	})

	return err
}

func (repo *configRepo) CanRegister(ctx context.Context) bool {
	config, _ := repo.GetByKey(ctx, consts.CONFIG_BALALA_REGISTER_ENTRY)
	return config != nil && config.ConfigValue == "1"
}

func (repo *configRepo) AttachSize(ctx context.Context) int64 {
	const DEFAULT_ATTACH_SIZE = 100 * 1024 * 1024

	key := consts.CONFIG_BALALA_ATTACH

	conf, err := repo.GetByKey(ctx, key)
	if err != nil {
		repo.log.Error(err)
	}

	if conf == nil {
		conf = repo.defaultMap()[key]
	}

	if conf == nil {
		return DEFAULT_ATTACH_SIZE
	}

	value := gjson.Get(conf.ConfigValue, "value").Int()
	unit := gjson.Get(conf.ConfigValue, "unit").String()

	if value <= 0 || unit == "" {
		return DEFAULT_ATTACH_SIZE
	}

	switch unit {
	default:
		return DEFAULT_ATTACH_SIZE
	case "KB":
		return value * 1024
	case "MB":
		return value * 1024 * 1024
	case "GB":
		return value * 1024 * 1024 * 1024
	}
}
