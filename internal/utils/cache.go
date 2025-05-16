package utils

type Cache interface {
	// ------------------- string -------------------

	Get(key string) (*string, error)
	Set(key string, val string, ttl int32) error
	SetObj(key string, obj interface{}, ttl int32) error
	GetObj(key string, obj interface{}) (interface{}, error)
	Del(key string) error

	// ------------------- int64 -------------------

	GetI(key string) (int64, error)
	SetI(key string, val int64, expire int32) error
	IncrBy(key string, val int64) (int64, error)

	// ------------------- hash -------------------
	HGet(key string, field string) (*string, error)
	HSet(key string, field string, val string) (int64, error)
	HDel(key string, field ...string) (int64, error)
	HKeys(key string) (*[]string, error)
	HGetAll(key string) (*map[string]string, error)

	// ------------------- list -------------------
	LPush(key string, values ...interface{}) (int64, error)
	RPush(key string, values ...interface{}) (int64, error)
	LPop(key string) (string, error)
	RPop(key string) (string, error)
	LLen(key string) (int64, error)
	LRange(key string, start int64, stop int64) (*[]string, error)

	// ------------------- set -------------------
	SAdd(key string, values ...interface{}) (int64, error)
	SRemove(key string, values ...interface{}) (int64, error)
	SMembers(key string) (*[]string, error)
	SCount(key string) (int64, error)
	SIsMember(key string, values ...interface{}) (bool, error)

	// ------------------- 其他 -------------------
	Expire(key string, t int32) (bool, error)
	TTL(key string) (float64, error)
	// Lock(key string) error
}
