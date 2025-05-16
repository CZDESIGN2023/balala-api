package server

import (
	"go-cs/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdClient(c *conf.Etcd) *etcd.Registry {
	if len(c.Endpoints) == 0 {
		c.Endpoints = append(c.Endpoints, "127.0.0.1:2379")
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: c.Timeout.AsDuration(),
	})
	if err != nil {
		panic(err)
	}
	// new reg with etcd client
	reg := etcd.New(client)

	return reg
}
