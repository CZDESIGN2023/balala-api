package utils

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdCache struct {
	ctx     context.Context
	log     *log.Helper
	etcdCli *clientv3.Client
	name    string
}

func NewEtcdCache(c context.Context, logger log.Logger, cli *clientv3.Client, name string) *EtcdCache {
	return &EtcdCache{
		ctx:     c,
		log:     log.NewHelper(log.With(logger, "module", "data/etcd-cache")),
		etcdCli: cli,
		name:    name,
	}
}

// 上锁
func (e *EtcdCache) Lock(key string) (*concurrency.Mutex, error) {
	s1, err := concurrency.NewSession(e.etcdCli)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	defer s1.Close()

	m1 := concurrency.NewMutex(s1, "/mutex/"+e.name+key)
	// 会话s1获取锁
	if err := m1.Lock(e.ctx); err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return m1, nil
}

// 解锁
func (e *EtcdCache) UnLock(m *concurrency.Mutex) error {
	if err := m.Unlock(e.ctx); err != nil {
		// log.Fatal(err)
		return err
	}
	return nil
}
