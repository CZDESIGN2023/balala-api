package websock

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/encoding"
	jsonc "github.com/go-kratos/kratos/v2/encoding/json"
	protoc "github.com/go-kratos/kratos/v2/encoding/proto"
	"github.com/go-kratos/kratos/v2/log"
	ws "github.com/gorilla/websocket"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	"go-cs/pkg/bus"
	"google.golang.org/protobuf/proto"
	"runtime"
	"sync"
	"time"
)

const channelBufSize = 256

type SessionID string

type Session struct {
	Id        SessionID
	userId    int64
	srv       *Server
	conn      *ws.Conn
	send      chan []byte
	codec     encoding.Codec
	logger    *log.Helper
	closeOnce sync.Once
	doneC     chan struct{}
	msgType   int
	No        int64
}

func NewSession(conn *ws.Conn, s *Server, codec encoding.Codec, userId int64, u4 string) *Session {
	c := &Session{
		Id:     SessionID(u4),
		userId: userId,
		conn:   conn,
		send:   make(chan []byte, channelBufSize),
		srv:    s,
		codec:  codec,
		logger: s.logger,
		doneC:  make(chan struct{}),
		No:     s.num.Add(1),
	}

	switch codec.Name() {
	case jsonc.Name:
		c.msgType = ws.TextMessage
	case protoc.Name:
		c.msgType = ws.BinaryMessage
	}

	return c
}

func (s *Session) OnConnected() {
	bus.Emit(notify.InternalEvent_WsConnected, &event.WsConnected{
		UserId: s.userId,
	})
}

func (s *Session) Conn() *ws.Conn {
	return s.conn
}

func (s *Session) SessionID() SessionID {
	return s.Id
}

func (s *Session) SendMessage(message []byte) {
	select {
	case s.send <- message:
	default:
		s.logger.Errorf("[websocket] %v, to many message", s)
	}
}

func (s *Session) SendObject(obj proto.Message) {
	marshal, err := s.codec.Marshal(obj)
	if err != nil {
		s.logger.Errorf("[websocket] encode error: %v", err)
		return
	}

	s.SendMessage(marshal)
}

func (s *Session) printStack(depth int) {
	// 设置调用栈的最大深度

	// 创建一个切片用于存储调用栈信息
	stack := make([]uintptr, depth)

	// 获取调用栈信息
	n := runtime.Callers(5, stack)
	frames := runtime.CallersFrames(stack[:n])

	// 遍历调用栈帧并打印信息
	for {
		frame, more := frames.Next()
		s.logger.Infof("%s stack frame: Function %s, File: %s, Line: %d", s, frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	_, file, line, _ := runtime.Caller(4)
	s.logger.Infof("[websocket] %v, close connection, file: %v, line: %v", s, file, line)
}

func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.logger.Debugf("[websocket] %v, connection closed", s)

		//s.printStack(3)

		close(s.doneC)
		s.srv.unregister <- s

		if err := s.conn.Close(); err != nil {
			s.logger.Errorf("[websocket] %v, disconnect error: ", s, err)
		}
	})
}

func (s *Session) Closed() bool {
	select {
	case <-s.doneC:
		return true
	default:
		return false
	}
}

func (s *Session) WritePump() {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Errorf("[websocket] %v, write pump panic: %v", s, err)
		}

		s.Close()
	}()

	for {
		select {
		case <-s.doneC:
			s.logger.Debugf("[websocket] %v, session closed: %v", s, s.SessionID())
			return
		case msg := <-s.send:
			if err := s.conn.WriteMessage(s.msgType, msg); err != nil {
				s.logger.Errorf("[websocket] %v, write message error: %v", s, err)
				return
			}
		}
	}
}

func (s *Session) ReadPump() {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Debugf("[websocket] %v, read pump panic: %v", s, err)
		}

		s.Close()
	}()

	heartbeatDuration := time.Second * 10
	heartbeat := time.AfterFunc(heartbeatDuration, func() {
		if s.Closed() { //已经关闭就不打印日志了
			return
		}

		s.logger.Debugf("[websocket] %v, heartbeat timeout", s)
		s.Close()
	})

	for {
		msgType, data, err := s.conn.ReadMessage()
		if err != nil {
			s.logger.Debugf("[websocket] %v, read message error: %v", s, err)
			return
		}

		switch msgType {
		case ws.BinaryMessage, ws.TextMessage:
			if string(data) == "ping" {
				if !heartbeat.Reset(heartbeatDuration) {
					s.logger.Errorf("[websocket] %v, reset heartbeat failed", s)
					return //重置心跳时间失败
				}

				err := s.conn.WriteMessage(ws.PongMessage, []byte("pong"))
				if err != nil {
					s.logger.Errorf("[websocket] %v, write pong message error: %v", s, err)
					return
				}
				continue
			}

			// 处理业务消息
			reply, err := s.MessageHandler(data)
			if err != nil {
				s.logger.Errorf("[websocket] %v, message handle error: %v", s, err)
				return
			}
			s.logger.Debug("[websocket] %v, message handle reply: %s", s, reply)

			if reply != nil {
				s.SendMessage(reply)
			}
		case ws.PingMessage:
			err := s.conn.WriteMessage(ws.PongMessage, []byte("pong"))
			if err != nil {
				s.logger.Errorf("[websocket] %v, write pong message error: %v", s, err)
				return
			}
		case ws.CloseMessage:
			return
		case ws.PongMessage:
			return
		}
	}
}

func (s *Session) String() string {
	return fmt.Sprintf("No: %v, userId: %v", s.No, s.userId)
}
