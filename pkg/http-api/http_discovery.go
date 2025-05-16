package http_api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"math/rand"
)

// 替代类型*etcd.Registry
type RegistryInterface interface {
	GetService(ctx context.Context, service string) ([]*registry.ServiceInstance, error)
	// include all other methods that you need
}

// HttpDiscovery 支持取得对象实例的方法
type HttpDiscoveryInterface interface {
	GetName() string
	GetServiceUrl(context.Context, string) (string, error)
}

type HttpDiscovery struct {
	discovery RegistryInterface
	name      string
}

// NewRegistryInterface 为了wire生成, 加个构造, 我可真难
func NewRegistryInterface(e *etcd.Registry) RegistryInterface {
	return e
}

func NewHttpDiscovery(e RegistryInterface, name string) HttpDiscoveryInterface {
	return &HttpDiscovery{
		discovery: e,
		name:      name,
	}
}

func (h *HttpDiscovery) GetName() string {
	return h.name
}

// GetServiceUrl 传入路径生成完整服务路径
func (h *HttpDiscovery) GetServiceUrl(ctx context.Context, path string) (string, error) {

	//str := fmt.Sprintf("%s://%s:%d%s", .Schema, instance.Addr, instance.Port, path)
	ep, err := h.getInstances(ctx)
	if err != nil {
		return "", err
	}
	return ep + path, nil
}

// getInstances 取得连接实现, 负载均衡策略
func (h *HttpDiscovery) getInstances(ctx context.Context) (string, error) {
	instances, err := h.discovery.GetService(ctx, h.name)
	if err != nil {
		return "", err
	}

	if len(instances) == 0 {
		return "", fmt.Errorf("no %s service instances available", h.name)
	}

	//  随机服务器
	instance := instances[rand.Intn(len(instances))]
	ep := instance.Endpoints[rand.Intn(len(instance.Endpoints))]
	return ep, err
}
