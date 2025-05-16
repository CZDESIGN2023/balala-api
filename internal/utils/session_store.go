package utils

import (
	"errors"
	"go-cs/internal/conf"
	"go-cs/internal/utils/cache"
	"go-cs/internal/utils/sessions"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/securecookie"
)

func NewSessionStore(data *conf.Data, logger log.Logger) (sessions.Store, func(), error) {

	switch data.Sessions.Driver {
	case "cookie":

		csOpt := &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		}

		cookieConf := data.Sessions.GetCookie()
		if cookieConf != nil {
			if cookieConf.Path != "" {
				csOpt.Path = cookieConf.Path
			}

			if cookieConf.MaxAge != 0 {
				csOpt.MaxAge = int(cookieConf.MaxAge)
			}

			if cookieConf.Domain != "" {
				csOpt.Domain = cookieConf.Domain
			}

			if cookieConf.CorssMode != 0 {
				csOpt.SameSite = http.SameSite(cookieConf.CorssMode)
			}

			if cookieConf.HttpOnly {
				csOpt.HttpOnly = true
			}

			if cookieConf.Secure {
				csOpt.Secure = true
			}
		}

		cs := &sessions.CookieStore{
			Options: csOpt,
			Codecs:  securecookie.CodecsFromPairs([]byte(data.Sessions.SessionSecure), []byte(data.Sessions.SessionSecure)),
		}

		cleanup := func() {
		}

		return cs, cleanup, nil

	case "redis":
		var cacheConfig *cache.Config

		//使用默认配置的缓存链接
		redisConf := data.Sessions.GetRedis()
		if redisConf.Addr == "" {
			cacheConfig = cache.NewConfig(data)
		} else {
			cacheConfig = NewRedisConfig(data)
		}

		cacheCli, cleanup, err := cache.NewRedis(cacheConfig, log.DefaultLogger)
		if err != nil {
			return nil, nil, err
		}

		redisStroe, err := sessions.NewRedisStore(cacheCli, []byte(data.Sessions.SessionSecure))
		if err != nil {
			return nil, nil, err
		}

		cleanup2 := func() {
			cleanup()
		}

		return redisStroe, cleanup2, nil
	}

	return nil, nil, errors.New("unsupported sessions driver")
}

func NewRedisConfig(data *conf.Data) *cache.Config {
	return &cache.Config{
		Addr:     data.Sessions.Redis.Addr,
		Password: data.Sessions.Redis.Password,
		DB:       int(data.Sessions.Redis.DbIndex),
		PoolSize: int(data.Sessions.Redis.PoolSize),

		// 超时时间引发的血案, 由于前面单位是秒,后面也是秒,太大,超时后导致溢出值太小导致超时
		ReadTimeout:  time.Duration(data.Sessions.Redis.ReadTimeout),
		WriteTimeout: time.Duration(data.Sessions.Redis.WriteTimeout),
		IdleTimeout:  600,
		MinIdleConns: 10,
	}
}
