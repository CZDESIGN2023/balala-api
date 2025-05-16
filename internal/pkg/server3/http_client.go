package server3

import (
	"context"
	"encoding/json"
	db "go-cs/internal/bean/biz"
	. "go-cs/pkg/http-api"
	"io"
	"net/http"
)

type Server3Interface interface {
	SendGift(ctx context.Context, fromUserId int64, toUserId int64, itemId int64, count int32) (*SendGiftReply, error)
}

type Server3Api struct {
	api HTTPClientInterface
}

func NewServer3Api(e RegistryInterface) *Server3Api {
	d := NewHttpDiscovery(e, "server3")
	return &Server3Api{
		api: NewHTTPClient(d),
	}
}

func GetTest(ctx context.Context, path string) (*[]byte, error) {
	page := "http://10.5.20.235:20010" + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, page, nil)
	if err != nil {
		return nil, err
	}

	// for k, v := range headers {
	// 	req.Header.Set(k, v)
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// panic(err)
		return nil, err
	}
	ret := respBody
	return &ret, nil
}

/////////////////////////////////////////

// SendGift 送礼
func (c *Server3Api) SendGift(ctx context.Context, fromUserId int64, toUserId int64, itemId int64, count int32) (*SendGiftReply, error) {
	info := make(map[string]interface{})
	info["from_user_id"] = fromUserId
	info["to_user_id"] = toUserId
	info["item_id"] = itemId
	info["count"] = count
	params, err := MakeSign(info, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.api.DoGet(ctx, "/open_api_go_im/chat_send_gift?"+params)
	if err != nil {
		return nil, err
	}
	data := &SendGiftReply{}
	if err := json.Unmarshal(resp, data); err != nil {
		return nil, err
	}

	return data, nil
}

// FetchUserInfo 获取用户信息回来
func (c *Server3Api) FetchUserInfo(ctx context.Context, ChatUserId int64) (*db.User, error) {
	info := make(map[string]interface{})
	info["user_id"] = ChatUserId

	params, err := MakeSign(info, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.api.DoGet(ctx, "/open_api_go_im/get_user_base_info?"+params)
	if err != nil {
		return nil, err
	}
	data := &db.User{}
	if err := json.Unmarshal(resp, data); err != nil {
		return nil, err
	}

	return data, nil
}
