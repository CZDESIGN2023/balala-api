package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go-cs/api/comm"
	loginV1 "go-cs/api/login/v1"
	userv1 "go-cs/api/user/v1"
	"time"

	//auth "go-cs/internal/server/auth"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
)

type Reply struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Code    int32       `json:"code"`
}

func HttpAccessLog(logger log.Logger) middleware.Middleware {
	l := log.NewHelper(logger, log.WithMessageKey("http_access"))

	ignoreMap := map[string]bool{
		loginV1.OperationLoginLogin:     true,
		userv1.OperationUserChangeMyPwd: true,
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			start := time.Now()

			reply, err = handler(ctx, req)

			dur := time.Now().Sub(start)

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return
			}

			ht, ok := tr.(*http.Transport)
			if !ok {
				return
			}

			request := ht.Request()
			if request == nil || ignoreMap[tr.Operation()] {
				return
			}

			uid := utils.GetLoginUser(ctx).UserId
			if v, ok := reply.(interface{ GetError() *comm.ErrorInfo }); ok && v.GetError() != nil {
				err := v.GetError()
				l.Infof("%v uid:%v, %v %v, req=%v, code:%v, msg:%s", dur, uid, request.Method, request.RequestURI, utils.ToJSON(req), err.Code, err.Message)
			} else {
				l.Infof("%v uid:%v, %v %v, req=%v", dur, uid, request.Method, request.RequestURI, utils.ToJSON(req))
			}

			return reply, err
		}
	}
}

func ErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	se := errs.Cast(err)
	if se == nil {
		http.DefaultErrorEncoder(w, r, err)
		return
	}

	codec, _ := http.CodecForRequest(r, "Accept")

	reply := &Reply{
		Code:    se.Code,
		Message: se.Message,
	}

	marshal, err := codec.Marshal(reply)
	if err != nil {
		w.WriteHeader(500)
	}

	w.Header().Set("Content-Type", contentType(codec.Name()))
	w.WriteHeader(200)
	_, _ = w.Write(marshal)
}

func ResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	var m = map[string]bool{
		"/my/info": true,
		"/login":   true,
	}
	if !m[r.RequestURI] {
		return http.DefaultResponseEncoder(w, r, v)
	}

	reply := &Reply{
		Code:    200,
		Data:    v,
		Message: "success",
	}

	codec, _ := http.CodecForRequest(r, "Accept")

	bytes2, _ := codec.Marshal(v)
	reply.Data = json.RawMessage(bytes2)

	marshal, err := codec.Marshal(reply)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", contentType(codec.Name()))
	w.WriteHeader(200)
	w.Write(marshal)
	return nil
}

func contentType(name string) string {
	switch name {
	case "json":
		return "application/json"
	case "protobuf":
		return "application/proto"
	default:
		return ""
	}
}
