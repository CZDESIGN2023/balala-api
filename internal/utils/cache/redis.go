package cache

import (
	"context"
	"go-cs/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr               string
	Username           string
	Password           string
	DB                 int
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

type Server3Config struct {
	Config
}

func NewConfig(data *conf.Data) *Config {
	return &Config{
		Addr:     data.Redis.Addr,
		Password: data.Redis.Password,
		DB:       int(data.Redis.DbIndex),
		PoolSize: int(data.Redis.PoolSize),

		// 超时时间引发的血案, 由于前面单位是秒,后面也是秒,太大,超时后导致溢出值太小导致超时
		ReadTimeout:  time.Duration(data.Redis.ReadTimeout),
		WriteTimeout: time.Duration(data.Redis.WriteTimeout),
		IdleTimeout:  600,
		MinIdleConns: 10,
	}
}

func NewServer3Config(data *conf.Data) *Server3Config {
	return &Server3Config{
		Config: Config{
			Addr:     data.Server3Redis.Addr,
			Password: data.Server3Redis.Password,
			DB:       int(data.Server3Redis.DbIndex),
			PoolSize: int(data.Server3Redis.PoolSize),

			ReadTimeout:  time.Duration(data.Server3Redis.ReadTimeout),
			WriteTimeout: time.Duration(data.Server3Redis.WriteTimeout),
			IdleTimeout:  600,
			MinIdleConns: 10,
		},
	}
}

// New 创建 cache 客户端
func NewRedis(config *Config, logger log.Logger) (*redis.Client, func(), error) {
	if config == nil {
		return nil, func() {}, nil
	}

	option := &redis.Options{
		Addr: config.Addr,
	}
	if config.Username != "" {
		option.Username = config.Username
	}
	if config.Password != "" {
		option.Password = config.Password
	}
	if config.DB != 0 {
		option.DB = config.DB
	}
	if config.MaxRetries != 0 {
		option.MaxRetries = config.MaxRetries
	}
	if config.MinRetryBackoff != 0 {
		option.MinRetryBackoff = config.MinRetryBackoff * time.Second
	}
	if config.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = config.MaxRetryBackoff * time.Second
	}
	if config.DialTimeout != 0 {
		option.DialTimeout = config.DialTimeout * time.Second
	}
	if config.ReadTimeout != 0 {
		option.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.WriteTimeout != 0 {
		option.WriteTimeout = config.WriteTimeout * time.Second
	}
	if config.PoolSize != 0 {
		option.PoolSize = config.PoolSize
	}
	if config.MinIdleConns != 0 {
		option.MinIdleConns = config.MinIdleConns
	}
	if config.MaxConnAge != 0 {
		option.MaxConnAge = config.MaxConnAge * time.Second
	}
	if config.PoolTimeout != 0 {
		option.PoolTimeout = config.PoolTimeout * time.Second
	}
	if config.IdleTimeout != 0 {
		option.IdleTimeout = config.IdleTimeout * time.Second
	}
	if config.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = config.IdleCheckFrequency * time.Second
	}

	client := redis.NewClient(option)

	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the cache client")

		if err := client.Close(); err != nil {
			log.NewHelper(logger).Error(err.Error())
		}
	}

	return client, cleanup, nil
}

type Server3Redis struct {
	*redis.Client
}

// NewServer3Redis 创建server3 redis
// func NewServer3Redis(config *Config, logger log.Logger) (*redis.Client, func(), error) {
func NewServer3Redis(config *Server3Config, logger log.Logger) (Server3Redis, func(), error) {
	if config == nil {
		return Server3Redis{}, func() {}, nil
	}

	option := &redis.Options{
		Addr: config.Addr,
	}
	if config.Username != "" {
		option.Username = config.Username
	}
	if config.Password != "" {
		option.Password = config.Password
	}
	if config.DB != 0 {
		option.DB = config.DB
	}
	if config.MaxRetries != 0 {
		option.MaxRetries = config.MaxRetries
	}
	if config.MinRetryBackoff != 0 {
		option.MinRetryBackoff = config.MinRetryBackoff * time.Second
	}
	if config.MaxRetryBackoff != 0 {
		option.MaxRetryBackoff = config.MaxRetryBackoff * time.Second
	}
	if config.DialTimeout != 0 {
		option.DialTimeout = config.DialTimeout * time.Second
	}
	if config.ReadTimeout != 0 {
		option.ReadTimeout = config.ReadTimeout * time.Second
	}
	if config.WriteTimeout != 0 {
		option.WriteTimeout = config.WriteTimeout * time.Second
	}
	if config.PoolSize != 0 {
		option.PoolSize = config.PoolSize
	}
	if config.MinIdleConns != 0 {
		option.MinIdleConns = config.MinIdleConns
	}
	if config.MaxConnAge != 0 {
		option.MaxConnAge = config.MaxConnAge * time.Second
	}
	if config.PoolTimeout != 0 {
		option.PoolTimeout = config.PoolTimeout * time.Second
	}
	if config.IdleTimeout != 0 {
		option.IdleTimeout = config.IdleTimeout * time.Second
	}
	if config.IdleCheckFrequency != 0 {
		option.IdleCheckFrequency = config.IdleCheckFrequency * time.Second
	}

	client := redis.NewClient(option)

	ctx := context.Background()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return Server3Redis{}, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the server3 redis")

		if err := client.Close(); err != nil {
			log.NewHelper(logger).Error(err.Error())
		}
	}

	return Server3Redis{Client: client}, cleanup, nil
}
