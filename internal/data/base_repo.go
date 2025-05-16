package data

import (
	"context"
	"errors"
	"go-cs/internal/utils"
	"regexp"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// 存着本机id是否从数据库load取
var initMaxIdMap sync.Map

// MaxStruct 取得最大ID用的结构体
type _maxStruct struct {
	Max int64 `json:"max"`
}

type baseRepo struct {
	data *Data
	log  *log.Helper
}

//////////////////////////////////////////////////
// 以下是自增id支持

func (c *baseRepo) loadMaxIdFromMysql(ctx context.Context, tb string) (int64, error) {
	if !isValidSqlStr(tb) {
		errMsg := "不合法的傳入值, 可能有sql injection風險"
		c.log.Error(errMsg)
		return 0, errors.New(errMsg)
	}

	var result _maxStruct
	err := c.data.DB(ctx).Raw("SELECT MAX( id ) AS max FROM " + tb).WithContext(ctx).Scan(&result).Error
	if err != nil {
		c.log.Error("出错 %v", err)
		return 0, err
	}
	return result.Max, err
}

func (c *baseRepo) loadMaxIdFromMongoDb(ctx context.Context, tb string) (int64, error) {

	var result struct {
		Max int64 `bson:"max"`
	}

	filter := bson.D{}
	option := options.FindOne().SetSort(bson.D{{"id", -1}}) // 按照id字段降序排列
	collection := c.data.mongo.Collection(tb)

	err := collection.FindOne(ctx, filter, option).Decode(&result)
	if err != nil {
		// 有表但没资料
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		c.log.Error("出错 %v", err)
		return 0, err
	}
	return result.Max, err
}

func (c *baseRepo) waitForInitId(ctx context.Context, name string, key string, dbType string) error {
	if _, ok := initMaxIdMap.Load(name); ok {
		return nil
	}

	rdbLock := utils.NewRedisLock(c.data.rdb)
	err := rdbLock.Lock(ctx, key, 10000)
	defer func(rdbLock utils.Mutex, ctx context.Context) {
		err := rdbLock.UnLock(ctx)
		if err != nil {
			c.log.Errorf("解锁错误", key, err)
		}
	}(rdbLock, ctx)

	if err != nil {
		c.log.Errorf("上锁失败", key, err)
		return err
	}

	// 初始化id
	maxId, err := c.data.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		var maxId int64
		var err error
		// 未初始化
		if dbType == "mysql" {
			maxId, err = c.loadMaxIdFromMysql(ctx, name)
		}
		if dbType == "mongoDB" {
			maxId, err = c.loadMaxIdFromMongoDb(ctx, name)
		}
		if err != nil {
			return err
		}

		// 将数据库id写入redis
		err = c.data.rdb.Set(ctx, key, maxId, 0).Err()
		if err != nil {
			return err
		}

	} else if err != nil {
		c.log.Errorf("getNewId ", key, err)
		return err
	}

	// 记录到本机初始化状态
	initMaxIdMap.Store(name, maxId)
	return nil
}

// 从redis里面获得自增ID
func (c *baseRepo) getNewId(ctx context.Context, name string, dbType string) (int64, error) {
	key := "newId:" + name

	err := c.waitForInitId(ctx, name, key, dbType)
	if err != nil {
		return 0, err
	}

	ret, err := c.data.rdb.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return 0, err
	}
	return ret, nil
}

// 从redis里面获得指定字段自增值
func (c *baseRepo) getNewFieldIdx(ctx context.Context, tableName string, dbType string, fieldName string) (int64, error) {
	return c.getNewFieldIdxWithCondition(ctx, tableName, fieldName, dbType, "", "")
}

// 依照where條件, 从redis里面获得指定字段自增值 (redis暫存的key會多一層whereKey)
func (c *baseRepo) getNewFieldIdxWithCondition(ctx context.Context, tableName string, dbType string, fieldName string, whereKey string, whereVal interface{}) (int64, error) {
	// id字段走原本的方式
	if fieldName == "id" {
		return c.getNewId(ctx, tableName, dbType)
	}

	key := "newFieldIdx:" + tableName + ":" + fieldName
	if whereKey != "" {
		key = key + ":" + whereKey
	}
	err := c.waitForInitFieldIdx(ctx, tableName, key, fieldName, whereKey, whereVal)
	if err != nil {
		return 0, err
	}

	ret, err := c.data.rdb.IncrBy(ctx, key, 1).Result()
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (c *baseRepo) waitForInitFieldIdx(ctx context.Context, tableName string, key string, fieldName string, whereKey string, whereVal interface{}) error {
	rdbLock := utils.NewRedisLock(c.data.rdb)
	err := rdbLock.Lock(ctx, key, 10000)
	defer rdbLock.UnLock(ctx)

	if err != nil {
		c.log.Errorf("上锁失败", key, err)
		return err
	}

	// 初始化id
	_, err = c.data.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		// 未初始化
		maxId, err := c.loadFieldMaxIdxFromDb(ctx, tableName, fieldName, whereKey, whereVal)
		if err != nil {
			return err
		}

		// 将数据库id写入redis
		err = c.data.rdb.Set(ctx, key, maxId, 0).Err()
		if err != nil {
			return err
		}

	} else if err != nil {
		c.log.Errorf("getNewFieldIdx ", key, err)
		return err
	}

	return nil
}

func (c *baseRepo) loadFieldMaxIdxFromDb(ctx context.Context, tableName string, fieldName string, whereKey string, whereVal interface{}) (int64, error) {
	if !isValidSqlStr(tableName) || !isValidSqlStr(fieldName) || (whereKey != "" && !isValidSqlStr(whereKey)) {
		errMsg := "不合法的傳入值, 可能有sql injection風險"
		c.log.Error(errMsg)
		return 0, errors.New(errMsg)
	}

	var err error
	var result _maxStruct
	if whereKey == "" {
		err = c.data.DB(ctx).Table(tableName).Select("MAX(" + fieldName + ") AS max").Scan(&result).Error
	} else {
		err = c.data.DB(ctx).Table(tableName).Select("MAX("+fieldName+") AS max").Where(whereKey+" = ?", whereVal).Scan(&result).Error
	}
	if err != nil {
		c.log.Error("出错 %v", err)
		return 0, err
	}
	return result.Max, err
}

// 串接sql字串時只允許: 字母、數字、下劃線_、連字符- , 避免sql injection風險
func isValidSqlStr(sqlStr string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(sqlStr)
}
