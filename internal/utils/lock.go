package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-cs/pkg/disgo"
	"time"
)

type Mutex interface {
	// Lock 传入锁名, 及超时时间
	Lock(ctx context.Context, name string, ms int) error
	UnLock(ctx context.Context) error
}

type redisLock struct {
	redisClient *redis.Client
	lock        *disgo.DistributedLock
}

var defaultInst redisLock

// DefaultInstance 默认锁实例
func DefaultInstance() Mutex {
	return defaultInst
}

// NewRedisLock 创建新的redis实例, 线程安全
func NewRedisLock(rc *redis.Client) Mutex {
	return &redisLock{redisClient: rc}
}

func (r redisLock) Lock(ctx context.Context, name string, ms int) error {
	lock, err := disgo.GetLock(r.redisClient, "lock:"+name)
	// 获得锁失败
	if err != nil {
		return err
	}
	// 将锁实例保留下来
	r.lock = lock
	_, err = lock.Lock(ctx, time.Millisecond*time.Duration(ms))
	if err != nil {
		// TODO:如果何清理战场呢
		r.lock = nil
		return err
	}
	return nil
}

func (r redisLock) UnLock(ctx context.Context) error {
	if r.lock != nil {
		_, err := r.lock.Release(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
