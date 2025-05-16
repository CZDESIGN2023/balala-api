package mocks

import "github.com/gorilla/websocket"

type WebConn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

var _ WebConn = &websocket.Conn{}
