package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	configV1 "go-cs/api/config/v1"
	loginV1 "go-cs/api/login/v1"
	userv1 "go-cs/api/user/v1"
	"google.golang.org/grpc/peer"
	"net/url"
)

const addrHeader = "X-RemoteAddr"

// 把operation放这里去白名单
var whitelistOps = map[string]bool{
	loginV1.OperationLoginLogin:       true,
	userv1.OperationUserRegUser:       true,
	configV1.OperationConfigList:      true,
	userv1.OperationUserCheckUserName: true,
	// loginV1.OperationLoginLogout:            true,
	loginV1.OperationLoginGetLoginValidCode: true,
	userv1.OperationUserChangeMyPwd:         true,
}

// === HTTP协义 ===
// 基于 operation 去决定用不用跳过验证
func NewWhiteListMatcher() selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		if whitelisted, ok := whitelistOps[operation]; ok && whitelisted {
			log.Info("[白名单] " + operation)
			return false
		}
		return true
	}
}

// === gRPC协义 ===
// water: 暂时没有

// === 共用 ===
// 从请求里拿ip然后放到元信息中
func RemoteAddrMiddleware() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var addr string
			if tp, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := tp.(*http.Transport); ok {
					addr = ht.Request().RemoteAddr
					log.Info("[http] remote ip: " + addr)
				} else if _, ok := tp.(*grpc.Transport); ok {
					if peerInfo, ok := peer.FromContext(ctx); ok {
						addr = peerInfo.Addr.String()
						log.Info("[grpc] peer ip: " + addr)
					}
				}
			}

			remoteUrl, err := url.Parse("http://" + addr)
			if err == nil {
				if md, ok := metadata.FromServerContext(ctx); ok {
					md.Set(addrHeader, remoteUrl.Hostname())
					ctx = metadata.NewServerContext(ctx, md)
				} else {
					log.Fatal("异常 context", ctx)
				}
			}

			return h(ctx, req)
		}
	}
}
