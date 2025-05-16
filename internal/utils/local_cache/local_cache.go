package local_cache

import (
	"sync"
	"sync/atomic"
	"time"
)

// Cache 泛型cache
// 初始化方式:
// local_cache.NewCache[key型別, value型別](ttl)
// 例如:
// local_cache.NewCache[string, *db.Config](-1)
// local_cache.NewCache[string, *db.Config](time.Hour * 24)
//
// 注意事項:
// 應盡量使用Reload或SetData方式初始化map, 避免一個一個Set
// 如果一個一個Set, 而且又有ttl需求, 記得最後要SetUpdateTime
//
// ====================================================================
// 如果需要每個key有不同的ttl
// 可以使用github.com/Code-Hex/go-generics-cache
// 用法可參考data/channel_config.go
// ====================================================================
type Cache[K comparable, V any] struct {
	data       map[K]V
	rwMutex    sync.RWMutex  // 讀寫鎖
	reloadFlag atomic.Bool   // reload標記
	ttl        time.Duration // cache存活時間長度, <=0時永不過期
	updateTime int64         // 最後更新时间
	expireTime int64         // 過期時間 = now + ttl, (此設定值在ttl>0時才進行判斷)
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		data:       map[K]V{},
		ttl:        ttl,
		updateTime: -1,
	}
}

// Reload 重新載入並更新全部緩存資料 (適合每次Get不到就嘗試讀取並更新緩存的情境使用, 避免緩存擊穿)
func (c *Cache[K, V]) Reload(reloadFunc func() (data map[K]V)) {
	if reloadFunc == nil {
		return
	}
	// 檢查reload標記, 確保一次只有一個線程可以進入
	// 避免cache過期後大量請求同時reload造成db負擔
	if c.reloadFlag.CompareAndSwap(false, true) {
		defer c.reloadFlag.Store(false)
		newData := reloadFunc()
		if newData != nil {
			c.rwMutex.Lock()
			defer c.rwMutex.Unlock()
			c.data = newData
			c.updateTime = time.Now().UnixNano()
			c.expireTime = time.Now().Add(c.ttl).UnixNano()
		}
	}
}

// SetData 一次性更新全部緩存資料 (適合只有一個thread負責更新的時候使用)
func (c *Cache[K, V]) SetData(data map[K]V) {
	if data == nil {
		return
	}
	c.rwMutex.Lock()
	c.data = data
	c.updateTime = time.Now().UnixNano()
	c.expireTime = time.Now().Add(c.ttl).UnixNano()
	c.rwMutex.Unlock()
}

// Get 取得緩存值, val不一定會是nil(要看當時建立的類型), 必要時應該搭配found判斷取值是否成功
func (c *Cache[K, V]) Get(k K) (val V, found bool) {
	c.rwMutex.RLock()
	v, found := c.data[k]
	c.rwMutex.RUnlock()
	if !found {
		return
	}
	if c.ttl > 0 && c.expireTime < time.Now().UnixNano() {
		return
	}
	return v, true
}

func (c *Cache[K, V]) Set(k K, v V) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.data[k] = v
}

func (c *Cache[K, V]) Delete(key K) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	delete(c.data, key)
}

func (c *Cache[K, V]) GetKeys() []K {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	keys := make([]K, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

// GetUpdateTime 取得緩存過期時間, 供外部邏輯判斷要不要reload
func (c *Cache[K, V]) GetUpdateTime() int64 {
	return c.updateTime
}

// SetUpdateTime 應該不需要用到, 除非map的初始化是一個一個Set才需要最後在這更新過期時間
func (c *Cache[K, V]) SetUpdateTime(updateTime int64) {
	c.updateTime = updateTime
	c.expireTime = time.Now().Add(c.ttl).UnixNano()
}
