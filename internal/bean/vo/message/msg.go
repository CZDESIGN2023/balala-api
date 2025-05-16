package message

import "go-cs/api/notify"

type Msg struct {
	Type notify.MsgType `json:"type"` //消息类型
	Data any            `json:"data"` // 消息内容
}
