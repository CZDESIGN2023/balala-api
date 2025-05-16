package server3auth

import (
	"context"
	"go-cs/internal/conf"
	"go-cs/internal/server/auth/server3auth/server3authfunc"
	"go-cs/internal/utils"
	"go-cs/internal/utils/local_cache"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-redis/redis/v8"
)

type Void struct{}

type AuthServer struct {
	log        *log.Helper
	rdb        *redis.Client
	jwtConf    *conf.Jwt
	loginCache *local_cache.Cache[int64, Void]
}

var instance *AuthServer

func GetAuthServerInstance() *AuthServer {
	return instance
}

func InitAuthServer(rdb *redis.Client, jwtConf *conf.Jwt, logger log.Logger) *AuthServer {
	moduleName := "Server3AuthServer"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	instance = &AuthServer{
		log:        hlog,
		rdb:        rdb,
		jwtConf:    jwtConf,
		loginCache: local_cache.NewCache[int64, Void](-1),
	}
	return instance
}

func Server(rdb *redis.Client) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			newCtx, err := server3authfunc.AuthToken(ctx, instance.jwtConf, rdb)
			if err != nil {
				return nil, err
			}

			return handler(newCtx, req)
		}
	}

}
