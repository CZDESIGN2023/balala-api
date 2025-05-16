package orm

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

const connectionString = "mongodb://localhost:27017/testDB"
const testCollection = "testUsers"

func init() {
	// 由于我们使用了固定的测试数据库和集合名称，我们希望在测试开始前清空集合
	store, err := NewMongoStoreFromDsn(connectionString, testCollection)
	if err != nil {
		panic(err)
	}
	//defer store.Close(context.Background())

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	store.db.Collection(testCollection).Drop(timeout)
}

func TestCount(t *testing.T) {
	// 创建一个MongoStore实例
	store, err := NewMongoStoreFromDsn(connectionString, testCollection)
	assert.NoError(t, err)
	defer func(store *MongoStore, ctx context.Context) {
		//err := store.Close(ctx)
		if err != nil {
			log.Debugf("%v", err)
		}
	}(store, context.Background())

	// 计算总用户数量
	totalCount, err := store.Count(context.Background(), bson.M{})
	assert.NoError(t, err)

	// 插入两个用户
	user1 := User{
		Name: "Bob",
		Age:  25,
	}
	_, err = store.Create(context.Background(), user1)
	assert.NoError(t, err)

	user2 := User{
		Name: "Charlie",
		Age:  29,
	}
	_, err = store.Create(context.Background(), user2)
	assert.NoError(t, err)

	// 计算年龄大于等于28的用户数量
	count, err := store.Count(context.Background(), bson.M{"age": bson.M{"$gte": 28}})
	assert.NoError(t, err)
	assert.Equal(t, totalCount+1, count)
}

func TestCRUD(t *testing.T) {
	// 创建一个MongoStore实例
	store, err := NewMongoStoreFromDsn(connectionString, testCollection)
	assert.NoError(t, err)
	//defer store.Close(context.Background())

	// C: 创建一个新用户
	user := User{
		Name: "Alice",
		Age:  30,
	}
	id, err := store.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, id)

	// R: 读取用户
	var fetchedUser User
	err = store.Take(context.Background(), bson.M{"_id": id}, &fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", fetchedUser.Name)
	assert.Equal(t, 30, fetchedUser.Age)

	// U: 更新用户
	updateData := bson.M{
		"name": "Alice Updated",
		"age":  31,
	}
	err = store.Updates(context.Background(), bson.M{"_id": id}, updateData)
	assert.NoError(t, err)

	err = store.Take(context.Background(), bson.M{"_id": id}, &fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, "Alice Updated", fetchedUser.Name)
	assert.Equal(t, 31, fetchedUser.Age)

	// D: 删除用户
	err = store.Delete(context.Background(), id)
	assert.NoError(t, err)

	ret, err := store.Count(context.Background(), id)
	assert.Equal(t, ret, int64(0)) // 查询应该会出错，因为文档已被删除
}

func TestBuildUpdate(t *testing.T) {
	tests := []struct {
		input  interface{}
		output interface{}
	}{
		{ // 测试map[string]interface{}的情况
			map[string]interface{}{"name": "John", "age": 25},
			bson.M{"$set": bson.M{"name": "John", "age": 25}},
		},
		{ // 测试bson.M的情况
			bson.M{"$inc": bson.M{"age": 1}},
			bson.M{"$inc": bson.M{"age": 1}},
		},
		{ // 测试结构体的情况
			//struct {
			//	Name string
			//	Age  int
			//}{"John", 25},
			//bson.M{"$set": bson.M{"Name": "John", "Age": 25}},
			struct {
				love_value   int64
				hide_self    int32
				block_anchor int32
			}{20, 1, 1},
			bson.M{"$set": bson.M{"love_value": 20, "hide_self": 1, "block_anchor": 1}},
		},
	}

	for _, tt := range tests {
		result := buildUpdate(tt.input)
		if !reflect.DeepEqual(result, tt.output) {
			t.Errorf("expected %v, but got %v", tt.output, result)
		}
	}
}

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		output interface{}
	}{
		{
			"Integer as ID",
			123,
			bson.M{"_id": bson.M{"$eq": 123}},
		},
		{
			"Map without $ prefix",
			map[string]interface{}{"name": "John", "age": 25},
			bson.M{"name": bson.M{"$eq": "John"}, "age": bson.M{"$eq": 25}},
		},
		//{
		//	"Map with $ prefix",
		//	map[string]interface{}{"$or": []bson.M{{"name": "John"}, {"age": 25}}},
		//	map[string]interface{}{"$or": []bson.M{{"name": "John"}, {"age": 25}}},
		//},
		{
			"Bson.M type",
			bson.M{"name": "John"},
			bson.M{"name": "John"},
		},
		{
			"Struct",
			struct {
				Name string
				Age  int
			}{"John", 25},
			bson.M{"Name": bson.M{"$eq": "John"}, "Age": bson.M{"$eq": 25}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildFilter(tt.input)
			//if !reflect.DeepEqual(result, tt.output) {
			//	t.Errorf("expected %v, but got %v", tt.output, result)
			//}
			if diff := cmp.Diff(tt.output, result); diff != "" {
				t.Errorf("expected %v, but got %v", tt.output, result)
			}
		})
	}
}
