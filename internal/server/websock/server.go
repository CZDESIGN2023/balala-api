package websock

import (
	"errors"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis/v8"
	ws "github.com/gorilla/websocket"
	"go-cs/internal/conf"
	"go-cs/internal/server/auth/server3auth/server3authfunc"
	"go-cs/internal/utils"
	"go-cs/pkg/sync1"
	"google.golang.org/protobuf/proto"
	"sync"
	"sync/atomic"
)

type Server struct {
	upgrader       *ws.Upgrader
	logger         *log.Helper
	httpServer     *http.Server
	sessionMap     sync1.Map[SessionID, *Session]
	userSessionMap sync1.Map[int64, *sync1.Map[SessionID, *Session]]
	register       chan *Session
	unregister     chan *Session
	confJwt        *conf.Jwt
	rdb            *redis.Client
	num            atomic.Int64
}

type UserSessions struct {
	data sync.Map
}

// NewServer 连接管理器使用单例
func NewServer(rdb *redis.Client, confJwt *conf.Jwt, logger log.Logger) *Server {
	moduleName := "WebsocketServer"
	_, helper := utils.InitModuleLogger(logger, moduleName)

	s := &Server{
		logger: helper,
		upgrader: &ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		register:   make(chan *Session),
		unregister: make(chan *Session),
		confJwt:    confJwt,
		rdb:        rdb,
	}

	s.start()

	return s
}

func (s *Server) SetHttpServer(httpServer *http.Server) {
	s.httpServer = httpServer
}

func (s *Server) UpgradeHandler(httpCtx http.Context) error {
	http.SetOperation(httpCtx, "ws.chat")

	ctx, err := server3authfunc.AuthToken(httpCtx, s.confJwt, s.rdb)
	if err != nil {
		s.logger.Debugf("[websocket] auth token failed, err: %v", err)
		err := httpCtx.Result(401, nil)
		if err != nil {
			return err
		}
		return err
	}

	loginUserInfo := utils.GetLoginUser(ctx)
	if loginUserInfo.UserId == 0 {
		s.logger.Error("[websocket] no such user")
		return errors.New("no such user")
	}

	rsp := httpCtx.Response()
	req := httpCtx.Request()

	conn, err := s.upgrader.Upgrade(rsp, req, nil)
	if err != nil {
		s.logger.Error("[websocket] upgrade error: ", err)
		return err
	}

	codec := getCodec(req)

	session := NewSession(conn, s, codec, loginUserInfo.UserId, loginUserInfo.JwtTokenId)
	s.register <- session

	go session.WritePump()
	go session.ReadPump()

	return httpCtx.Result(200, nil)
}

func (s *Server) start() {
	go func() {
		if err := recover(); err != nil {
			s.logger.Errorf("[websocket] register panic: %v", err)
		}

		for {
			select {
			case c := <-s.register:
				s.sessionMap.Store(c.SessionID(), c)

				userSessions, loaded := s.userSessionMap.Load(c.userId)
				if !loaded {
					userSessions = &sync1.Map[SessionID, *Session]{}
					s.userSessionMap.Store(c.userId, userSessions)
				}
				oldSession, loaded := userSessions.Swap(c.SessionID(), c)
				if loaded {
					s.logger.Debugf("[websocket] %s old session close", oldSession)
					oldSession.Close()
				}

				c.OnConnected()
			}
		}
	}()

	go func() {
		if err := recover(); err != nil {
			s.logger.Errorf("[websocket] unregister panic: %v", err)
		}

		for {
			select {
			case c := <-s.unregister:
				s.sessionMap.CompareAndDelete(c.SessionID(), c)

				ret, loaded := s.userSessionMap.Load(c.userId)
				if loaded {
					ret.CompareAndDelete(c.SessionID(), c)
				}
			}
		}
	}()
}

// SendObject 单播数据包
func (s *Server) SendObjectBySessionID(sessionId SessionID, message proto.Message) {
	c, ok := s.sessionMap.Load(sessionId)
	if !ok {
		s.logger.Errorf("[websocket] session not found: %v", sessionId)
	}

	c.SendObject(message)
}

// BroadcastObject 广播对象
func (s *Server) BroadcastObject(message proto.Message) {
	s.sessionMap.Range(func(k SessionID, v *Session) {
		v.SendObject(message)
	})
}

// BroadcastData 广播数据包
func (s *Server) BroadcastData(data []byte) {
	s.sessionMap.Range(func(k SessionID, v *Session) {
		v.SendMessage(data)
	})
}

// SendObject 单播数据包
func (s *Server) SendObject(message proto.Message, userIds ...int64) {
	for _, userId := range userIds {
		userSession, ok := s.userSessionMap.Load(userId)
		if !ok {
			s.logger.Errorf("[websocket] %v session not found", userId)
			continue
		}

		userSession.Range(func(k SessionID, v *Session) {
			v.SendObject(message)
		})
	}
}

func (s *Server) SendData(data []byte, userIds ...int64) {
	for _, userId := range userIds {
		s.SendData2User(data, userId)
	}
}

func (s *Server) SendData2User(data []byte, userId int64) bool {
	userSession, ok := s.userSessionMap.Load(userId)
	if !ok {
		s.logger.Debugf("[websocket] %v session not found", userId)
		return false
	}

	ok = false
	userSession.Range(func(k SessionID, v *Session) {
		v.SendMessage(data)
		ok = true
	})

	return ok
}

func getCodec(req *http.Request) encoding.Codec {
	contentType := req.URL.Query().Get("ct")
	codec := encoding.GetCodec(contentType)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}

	return codec
}
