package orm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"go-cs/internal/conf"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func getCustomTLSConfig(caFile string) (*tls.Config, error) {
	tlsConfig := new(tls.Config)
	certs, err := ioutil.ReadFile(caFile)

	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.RootCAs = x509.NewCertPool()
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

	if !ok {
		return tlsConfig, errors.New("Failed parsing pem file")
	}

	return tlsConfig, nil
}

func getDatabaseNameFromDSN(dsn string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	// 获取路径部分，并移除前面的斜线
	dbName := strings.TrimLeft(parsedURL.Path, "/")

	// 如果路径中含有 "?"，只取问号之前的部分
	if strings.Contains(dbName, "?") {
		dbName = strings.Split(dbName, "?")[0]
	}

	return dbName, nil
}

type MongoDBFactory func() (*mongo.Database, error)

// 创建一个用于缓存已经初始化的数据库实例的map
var dbCache = make(map[string]*mongo.Database)
var mu sync.Mutex

// ProvideMongoDatabase 从连接字符串中取得数据库名
func ProvideMongoDatabase(c *conf.Data, kLogger klog.Logger, zLogger *zap.Logger) MongoDBFactory {

	return func() (*mongo.Database, error) {
		dbName, err := getDatabaseNameFromDSN(c.Mongo.Dsn)
		if err != nil {
			return nil, err
		}

		mu.Lock()
		defer mu.Unlock()

		// 从缓存中检查是否存在已经初始化的数据库实例
		if db, exists := dbCache[dbName]; exists {
			return db, nil
		}

		// 如果不存在，则创建新的实例
		client, err := ProvideMongoDBClient(c, kLogger, zLogger)
		if err != nil {
			return nil, err
		}
		db := client.Database(dbName)

		// 将新的数据库实例存入缓存
		dbCache[dbName] = db

		return db, nil
	}
}

func ProvideMongoDBClient(c *conf.Data, kLogger klog.Logger, zLogger *zap.Logger) (*mongo.Client, error) {
	if c.Mongo == nil {
		return nil, fmt.Errorf("mongo配置为空")
	}

	start := time.Now()
	opt := options.Client().ApplyURI(c.Mongo.Dsn)

	zLogger.Info("NewMongoDBClient connect to url", zap.String("url", c.Mongo.Dsn))

	if c.Mongo.CaFilePath != "" {
		tlsConfig, err := getCustomTLSConfig(c.Mongo.CaFilePath)
		if err != nil {
			zLogger.Error("Failed getting TLS configuration", zap.String("cafile", c.Mongo.CaFilePath), zap.Error(err))
			return nil, err
		}
		opt = opt.SetTLSConfig(tlsConfig)
	}

	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, err
	}

	const connectTimeout = 50
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	zLogger.Info("MongoDB connected!", zap.Duration("elapsed", elapsed))

	return client, nil
}

// MongoStoreFactory 传入集合名进行构造工厂方法
type MongoStoreFactory func(collectName string) (*MongoStore, error)

// ProvideMongoStoreFactory 返回一个函数，该函数接收集合名，并返回一个新的MongoStore实例。
func ProvideMongoStoreFactory(c *conf.Data, kLogger klog.Logger, zLogger *zap.Logger) MongoStoreFactory {
	mdb, err := ProvideMongoDatabase(c, kLogger, zLogger)()
	if err != nil {
		zLogger.Error("ProvideMongoStoreFactory,NewMongoDBClient Error!", zap.Error(err))
		return nil
	}
	return func(collectName string) (*MongoStore, error) {
		store, err := NewMongoStore(mdb, collectName)
		if err != nil {
			zLogger.Error("ProvideMongoStoreFactory,NewMongoStore Error!", zap.Error(err))
			return nil, err
		}
		return store, nil
	}
}

// MongoStore 封装了对MongoDB的基本操作
type MongoStore struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoStore 创建一个MongoStore实例
// 从连接字符串中解析数据库名
func NewMongoStore(db *mongo.Database, collectionName string) (*MongoStore, error) {
	collection := db.Collection(collectionName)
	return &MongoStore{db: db, collection: collection}, nil
}

// NewMongoStoreFromDsn 根据连接字符串创建一个MongoStore实例
func NewMongoStoreFromDsn(connectionString string, collectionName string) (*MongoStore, error) {
	opt := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		return nil, err
	}

	// 使用函数获取数据库名
	dbName, err := getDatabaseNameFromDSN(connectionString)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return NewMongoStore(db, collectionName)
}

// buildFilter 构建MongoDB的查询过滤器。
//
// 参数:
// - filter: 过滤条件，可以是多种数据类型，如整数、映射或结构体。
//
// 返回值:
// - 构建完成的查询过滤器
//
// 支持的过滤条件包括：
//
//  1. 整数: 当filter为整数时，认为它是文档的ID。
//     例如：filter = 123 -> {_id: 123}
//
//  2. 映射: 使用bson.M来构建复杂的查询条件。
//     例如：filter = bson.M{"name": "Alice"} -> {name: "Alice"}
//     更复杂的例子：
//     - bson.M{"age": bson.M{"$gt": 20}} -> 查询年龄大于20的文档
//     - bson.M{"name": bson.M{"$in": []string{"Alice", "Bob"}}} -> 查询名字为Alice或Bob的文档
//
//  3. 结构体: 使用结构体的字段来构建查询条件，字段的零值将被忽略。
//     例如：struct User{Name: "Alice"} -> {name: "Alice"}
//
// MongoDB常用的操作符:
// - "$eq": 等于
// - "$ne": 不等于
// - "$gt": 大于
// - "$gte": 大于等于
// - "$lt": 小于
// - "$lte": 小于等于
// - "$in": 包含在指定数组中
// - "$nin": 不包含在指定数组中
// - "$and": 与操作
// - "$or": 或操作
// - "$not": 非操作
// ... 更多操作符请参考MongoDB官方文档
//
// 使用举例:
// filter = 123 -> {_id: 123}
// filter = bson.M{"age": bson.M{"$gt": 20, "$lt": 30}} -> {age: {$gt: 20, $lt: 30}}
// filter = User{Name: "Alice", Age: 0} -> {name: "Alice"}  // Age字段的零值被忽略
func buildFilter(filter interface{}) interface{} {
	switch v := filter.(type) {
	case int, int64, uint64, uint32, primitive.ObjectID: // 如果filter是整数，我们认为它是ID
		return bson.M{"_id": bson.M{"$eq": v}}
	case map[string]interface{}:
		bFilter := bson.M{}
		for key, value := range v {
			if strings.HasPrefix(key, "$") {
				bFilter[key] = value
			} else if value != nil {
				bFilter[key] = bson.M{"$eq": value}
			}
		}
		return bFilter
	case bson.M:
		return v
	default: // 当filter是结构体时
		val := reflect.ValueOf(filter)
		bFilter := bson.M{}
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			value := val.Field(i).Interface()
			if !isEmptyValue(value) {
				bFilter[field.Name] = bson.M{"$eq": value}
			}
		}
		return bFilter
	}
}

func isEmptyValue(v interface{}) bool {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return false
}

// buildUpdate 构建MongoDB的更新操作内容。
//
// 参数:
// - update: 更新的内容，可以是多种数据类型，如映射或结构体。
//
// 返回值:
// - 构建完成的更新操作内容。
//
// 支持的更新内容包括：
// 1. 映射: 使用map[string]interface{}或bson.M来构建更新操作。
//
//   - 当键名以"$"开始（如"$set", "$push"等），直接将其视为MongoDB的更新操作符。
//     例如：update = bson.M{"$set": bson.M{"name": "Alice"}} -> {$set: {name: "Alice"}}
//
//   - 当键名不以"$"开始，认为它是需要更新的字段，此时将会自动为它添加"$set"操作符。
//     例如：update = bson.M{"name": "Alice"} -> {$set: {name: "Alice"}}
//
//     2. 结构体: 使用结构体的非零值字段来构建"$set"更新操作。
//     例如：struct User{Name: "Alice"} -> {$set: {name: "Alice"}}
//
// MongoDB常用的更新操作符:
// - "$set": 设置字段的值
// - "$unset": 移除字段
// - "$inc": 增加字段的值
// - "$push": 向数组字段添加一个值
// - "$pull": 从数组字段移除一个值
// - "$pop": 从数组字段移除第一个或最后一个值
// ... 更多更新操作符请参考MongoDB官方文档
//
// 使用举例:
// update = bson.M{"$set": bson.M{"name": "Alice", "age": 25}} -> {$set: {name: "Alice", age: 25}}
// update = bson.M{"age": bson.M{"$inc": 1}} -> {age: {$inc: 1}}  // 年龄字段加1
// update = User{Name: "Alice", Age: 25} -> {$set: {name: "Alice", age: 25}}
func buildUpdate(update interface{}) interface{} {
	switch v := update.(type) {
	case map[string]interface{}:
		bUpdate := bson.M{}
		for key, value := range v {
			// 如果有$开头的操作符，直接使用
			if strings.HasPrefix(key, "$") {
				bUpdate[key] = value
			} else if value != nil { // 否则，认为它是一个需要更新的字段
				bUpdate[key] = value
			}
		}
		return bson.M{"$set": bUpdate}
	case bson.M:
		// 检查 bson.M 是否已经包含了以 '$' 开头的键
		containsDollarKey := false
		for key := range v {
			if strings.HasPrefix(key, "$") {
				containsDollarKey = true
				break
			}
		}

		// 如果不包含，那么将整个bson.M放入"$set"下
		if !containsDollarKey {
			v = bson.M{"$set": v}
		}
		return v
	default: // 当update是结构体时
		val := reflect.ValueOf(update)
		bUpdate := bson.M{}
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			value := val.Field(i).Interface()
			if !isEmptyValue(value) {
				bUpdate[field.Name] = value
			}
		}
		return bson.M{"$set": bUpdate}
	}
}

// Create
// 插入一个文档并返回该文档的ID。
// 参数:
// - ctx: 上下文对象，用于请求的超时和取消
// - document: 需要插入的文档
// 返回值:
// - 插入文档的ID
// - 如果有错误发生，则返回错误
// 使用举例:
// docID, err := db.Create(ctx, &User{Name: "Alice", Age: 25})
func (store *MongoStore) Create(ctx context.Context, document interface{}) (int64, error) {
	res, err := store.collection.InsertOne(ctx, document)
	if err != nil {
		return 0, err
	}
	return res.InsertedID.(int64), nil
}

// Delete
// 删除与给定过滤条件匹配的文档。过滤条件可以是一个整数ID、一个过滤映射或一个结构体。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件，例如：23 (ID)、bson.M{"name": "Alice"} 或 User{Name: "Alice"}
// 返回值:
// - 如果有错误发生，则返回错误
// 使用举例:
// err := db.Delete(ctx, 23)
// err := db.Delete(ctx, bson.M{"name": "Alice"})
func (store *MongoStore) Delete(ctx context.Context, filter interface{}) error {
	switch v := filter.(type) {
	case int, int32, int64, uint64, uint32, primitive.ObjectID: // 当 filter 是整数时，我们认为它是ID
		_, err := store.collection.DeleteOne(ctx, bson.M{"_id": v})
		return err
	default:
		_, err := store.collection.DeleteMany(ctx, buildFilter(filter))
		return err
	}
}

// Update
// 更新与给定过滤条件匹配的第一个文档的指定字段。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件
// - field: 要更新的字段名
// - value: 要更新的字段值
// 返回值:
// - 如果有错误发生，则返回错误
// 使用举例:
// err := db.Update(ctx, bson.M{"name": "Alice"}, "age", 26)
func (store *MongoStore) Update(ctx context.Context, filter interface{}, updates ...interface{}) error {
	bUpdate := bson.M{}
	for i := 0; i < len(updates); i += 2 {
		key, ok := updates[i].(string)
		if !ok || i+1 >= len(updates) {
			return fmt.Errorf("invalid updates format")
		}
		value := updates[i+1]
		bUpdate[key] = value
	}

	updateData := bson.M{"$set": bUpdate}

	_, err := store.collection.UpdateOne(ctx, buildFilter(filter), updateData)
	return err
}

// Updates
// 更新与给定过滤条件匹配的所有文档。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件
// - update: 更新内容，可以包含多种MongoDB更新操作符如"$set", "$inc"等
// 返回值:
// - 如果有错误发生，则返回错误
// 使用举例:
// err := db.Updates(ctx, bson.M{"name": "Alice"}, bson.M{"$set": bson.M{"age": 26}})
// err := db.Updates(ctx, bson.M{"name": "Alice"}, bson.M{"$inc": bson.M{"age": 1}}) //年龄加1
func (store *MongoStore) Updates(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := store.collection.UpdateMany(ctx, buildFilter(filter), buildUpdate(update))
	return err
}

// Find
// 查询与给定过滤条件匹配的所有文档并返回。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件，例如使用bson.M{"age": bson.M{"$gt": 20}}可以过滤年龄大于20的文档
// - results: 用于存放查询结果的切片指针
// 返回值:
// - 如果有错误发生，则返回错误
// 使用举例:
// var users []User
// err := db.Find(ctx, bson.M{"age": bson.M{"$gt": 20}}, &users)
func (store *MongoStore) Find(ctx context.Context, filter interface{}, results interface{}) error {
	//opts := options.Find().SetLimit(limit)
	cursor, err := store.collection.Find(ctx, buildFilter(filter))
	if err != nil {
		return err
	}
	// 这里我们使用results，它应该是一个切片的指针，例如：*[]User
	if err := cursor.All(ctx, results); err != nil {
		return err
	}
	return nil
}

// Take
// 查询与给定过滤条件匹配的第一个文档并返回。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件
// - result: 用于存放查询结果的结构体指针
// 返回值:
// - 如果有错误发生，则返回错误
// 使用举例:
// var user User
// err := db.Take(ctx, bson.M{"name": "Alice"}, &user)
func (store *MongoStore) Take(ctx context.Context, filter interface{}, result interface{}) error {
	// 使用buildFilter构建过滤器
	singleResult := store.collection.FindOne(ctx, buildFilter(filter))

	// 解码查询结果到提供的结构体
	if err := singleResult.Decode(result); err != nil {
		return err
	}
	return nil
}

// Count
// 统计与给定过滤条件匹配的文档数量。
// 参数:
// - ctx: 上下文对象
// - filter: 过滤条件，例如使用bson.M{"age": bson.M{"$lt": 20}}可以统计年龄小于20的文档数量
// 返回值:
// - 匹配文档的数量
// - 如果有错误发生，则返回错误
// 使用举例:
// count, err := db.Count(ctx, bson.M{"age": bson.M{"$lt": 20}})
func (store *MongoStore) Count(ctx context.Context, filter interface{}) (int64, error) {
	count, err := store.collection.CountDocuments(ctx, buildFilter(filter))
	return count, err
}
