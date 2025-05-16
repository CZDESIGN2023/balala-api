package local_cache

import (
	"go-cs/internal/utils"
	"strconv"
)

// NestedCache 兩層嵌套的泛型cache
type NestedCache[K1 comparable, K2 comparable, V any] struct {
	data           map[K1]map[K2]V
	shardedLocks   *utils.ShardedLocks  // 分片鎖
	shardedKeyFunc func(key1 K1) string // 需要把key1轉為string, 後續計算hash決定分片
}

func createNestedCache[K1 comparable, K2 comparable, V any](numShards uint64, shardedKeyFunc func(key K1) string) *NestedCache[K1, K2, V] {
	return &NestedCache[K1, K2, V]{
		data:           map[K1]map[K2]V{},
		shardedLocks:   utils.NewShardedLocks(numShards),
		shardedKeyFunc: shardedKeyFunc,
	}
}

func NewStrKey1NestedCache[K2 comparable, V any](numShards uint64) *NestedCache[string, K2, V] {
	return createNestedCache[string, K2, V](numShards, strShardedFunc)
}

func NewInt32Key1NestedCache[K2 comparable, V any](numShards uint64) *NestedCache[int32, K2, V] {
	return createNestedCache[int32, K2, V](numShards, int32ShardedFunc)
}

func NewInt64Key1NestedCache[K2 comparable, V any](numShards uint64) *NestedCache[int64, K2, V] {
	return createNestedCache[int64, K2, V](numShards, int64ShardedFunc)
}

func strShardedFunc(key string) string {
	return key
}

func int32ShardedFunc(key int32) string {
	return strconv.FormatInt(int64(key), 10)
}

func int64ShardedFunc(key int64) string {
	return strconv.FormatInt(key, 10)
}

func (c *NestedCache[K1, K2, V]) Get(k1 K1, k2 K2) (val V, found bool) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.RLock()
	defer lock.RUnlock()

	innerMap, found := c.getInnerMap(k1, false)
	if !found {
		return
	}

	val, found = innerMap[k2]
	if !found {
		return
	}
	return val, true
}

func (c *NestedCache[K1, K2, V]) Set(k1 K1, k2 K2, v V) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.Lock()
	defer lock.Unlock()

	innerMap, _ := c.getInnerMap(k1, true)
	innerMap[k2] = v
}

// MSet 指定key1, 對整個innerMap更新
func (c *NestedCache[K1, K2, V]) MSet(k1 K1, innerMap map[K2]V) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.Lock()
	defer lock.Unlock()

	c.data[k1] = innerMap
}

func (c *NestedCache[K1, K2, V]) Delete(k1 K1, k2 K2) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.Lock()
	defer lock.Unlock()

	innerMap, found := c.getInnerMap(k1, false)
	if found {
		delete(innerMap, k2)
	}
}

// GetKeys 取得內層map所有的key
// 返回的第二個參數可以用來判斷是否從未初始化過, 如果found == false, 外部就可以撈取db資料做第一次Set的動作
func (c *NestedCache[K1, K2, V]) GetKeys(k1 K1) ([]K2, bool) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.Lock()
	defer lock.Unlock()

	innerMap, found := c.getInnerMap(k1, false)
	if found {
		keys := make([]K2, 0, len(innerMap))
		for k := range innerMap {
			keys = append(keys, k)
		}
		return keys, found
	}
	return make([]K2, 0), found
}

func (c *NestedCache[K1, K2, V]) DeleteByKey1(k1 K1) {
	lockKey := c.shardedKeyFunc(k1)
	lock := c.shardedLocks.GetLockForKey(lockKey)
	lock.Lock()
	defer lock.Unlock()

	delete(c.data, k1)
}

// 取得內層map (不可獨立使用, 應透過其他加鎖func操作)
func (c *NestedCache[K1, K2, V]) getInnerMap(k1 K1, autoCreate bool) (map[K2]V, bool) {
	innerMap, found := c.data[k1]
	if innerMap == nil && autoCreate {
		cache := map[K2]V{}
		c.data[k1] = cache
		innerMap = cache
	}
	return innerMap, found
}
