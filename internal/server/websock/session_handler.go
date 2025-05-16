package websock

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/encoding/proto"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cast"
	"go-cs/internal/bean"
	"io"
	netHttp "net/http"
	"net/http/httptest"
)

const WebsocketSession = "websocket-session"

func (s *Session) MessageHandler(data []byte) ([]byte, error) {
	s.srv.logger.Infof("[websocket] %v, handle data: %v", s, string(data))

	codec := s.codec

	var req bean.Request
	err := codec.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}

	httpReq, err := s.convertToHttpRequest(&req, codec.Name())
	if err != nil {
		return nil, err
	}

	httpRsp := httptest.NewRecorder()

	s.srv.httpServer.ServeHTTP(httpRsp, httpReq)

	beanRsp := &bean.Response{
		Random: req.Random,
		Body:   httpRsp.Body.Bytes(),
	}

	marshal, err := codec.Marshal(beanRsp)
	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (s *Session) convertToHttpRequest(req *bean.Request, codecName string) (*http.Request, error) {
	var reqBody io.Reader
	switch req.Method {
	case
		bean.Request_GET,
		bean.Request_DELETE:
	case
		bean.Request_POST,
		bean.Request_PUT:
		reqBody = bytes.NewBuffer(req.Data)
	default:
		return nil, errors.New("unknown http method")
	}

	httpMethod := bean.Request_Method_name[int32(req.Method)]
	httpReq, err := netHttp.NewRequest(httpMethod, req.Path, reqBody)
	httpReq.RequestURI = httpReq.URL.RequestURI()

	// 带上一个上下文，表明是websocket转发的请求
	httpReq = httpReq.WithContext(context.WithValue(context.Background(), WebsocketSession, []string{cast.ToString(s.userId), string(s.SessionID())}))

	switch codecName {
	case proto.Name:
		httpReq.Header.Set("Content-Type", "application/proto")
		httpReq.Header.Set("Accept", "application/proto")
	default:
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")
	}

	return httpReq, err
}
