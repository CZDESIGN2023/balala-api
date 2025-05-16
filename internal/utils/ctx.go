package utils

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gogf/gf/text/gregex"
	"google.golang.org/grpc/peer"
)

func GetRequestFromTransport(ctx context.Context) *http.Request {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil
	}

	if ht, ok := tr.(*http.Transport); ok {
		return ht.Request()
	}

	return nil
}

func GetTransPortFormCtx(ctx context.Context) *http.Transport {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil
	}

	if ht, ok := tr.(*http.Transport); ok {
		return ht
	}
	return nil
}

func GetRequestUrlPrefix(ctx context.Context) string {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return ""
	}

	if ht, ok := tr.(*http.Transport); ok {
		host := ht.Request().Host
		if host != "" {
			if ht.Request().TLS != nil {
				return "https://" + host
			}
			return "http://" + host
		}
	}

	return ""
}

func GetIpFrom(ctx context.Context) string {
	if v := GetIpFromHttpRequestHeader(ctx); v != "" {
		return v
	}

	if v := GetIpFromMetadata(ctx); v != "" {
		return v
	}

	if v := GetIpFromTransport(ctx); v != "" {
		return v
	}

	return ""
}

func GetIpFromHttpRequestHeader(ctx context.Context) string {

	var clientIp string
	r, ok := http.RequestFromServerContext(ctx)
	if ok {
		// fmt.Println("r.Header", r.Header)
		realIps := r.Header.Get("X-Forwarded-For")
		if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
			ipArray := strings.Split(realIps, ",")
			clientIp = ipArray[0]
			// fmt.Println("X-Forwarded-For", realIps, clientIp)
		}

		if clientIp == "" {
			clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
			// fmt.Println("TTP_X_FORWARDED_FOR", clientIp)
		}

		if clientIp == "" {
			clientIp = r.Header.Get("X-Real-IP")
			// fmt.Println("X-Real-IP", clientIp)
		}

		if clientIp == "" {
			array, _ := gregex.MatchString(`(.+):(\d+)`, r.RemoteAddr)
			if len(array) > 1 {
				clientIp = strings.Trim(array[1], "[]")
				// fmt.Println("RemoteAddr", r.RemoteAddr)
			}
		}
	}

	return clientIp
}

func GetIpFromMetadata(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get("X-RemoteAddr")
	}
	return ""
}

func GetIpFromTransport(ctx context.Context) string {
	var addr string
	if tp, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := tp.(*http.Transport); ok {
			addr = ht.Request().RemoteAddr
		} else if _, ok := tp.(*grpc.Transport); ok {
			if peerInfo, ok := peer.FromContext(ctx); ok {
				addr = peerInfo.Addr.String()
			}
		}
	}

	remoteUrl, err := url.Parse("http://" + addr)
	if err == nil {
		return remoteUrl.Hostname()
	}
	return ""
}
