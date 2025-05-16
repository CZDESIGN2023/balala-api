package notify

import (
	"context"
	"encoding/json"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/bean/vo/message"
	"go-cs/internal/utils"
)

func (s *Notify) WsConnected(e *event.WsConnected) {

	ctx := context.Background()

	list, err := s.notifySnapShotRepo.GetDelOfflineNotify(ctx, e.UserId)
	if err != nil {
		return
	}

	// 发送离线消息
	for _, data := range list {
		m := utils.ToJSONBytes(message.Msg{
			Type: notify.MsgType_mt_Notify,
			Data: json.RawMessage(data),
		})
		s.ws.SendData2User(m, e.UserId)
	}
}
