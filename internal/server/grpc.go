package server

import (
	loginV1 "go-cs/api/login/v1"
	userv1 "go-cs/api/user/v1"
	"go-cs/internal/server/auth/server3auth"

	"github.com/go-redis/redis/v8"

	"go-cs/internal/conf"
	"go-cs/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	j *conf.Jwt,
	login *service.LoginService,
	user *service.UserService,
	rdb *redis.Client,
	logger log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			RemoteAddrMiddleware(),
			selector.Server(server3auth.Server(rdb)).Match(NewWhiteListMatcher()).Build(),
			tracing.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	userv1.RegisterUserServer(srv, user)
	loginV1.RegisterLoginServer(srv, login)

	return srv
}
