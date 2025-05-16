package third_platform

import "go-cs/api/comm"

type UserInfo struct {
	Id       int64  `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nick_name"`
	HeadAddr string `json:"head_addr"`
}

type IClient interface {
	Bind(chatToken string) error
	Unbind(chatToken string) error
	Push(msg any, chatTokens []string) error
	GetUserInfo(key string) (*UserInfo, error)
}

type Client struct {
	clientMap map[comm.ThirdPlatformCode]IClient
}

func NewClient() *Client {
	return &Client{
		clientMap: make(map[comm.ThirdPlatformCode]IClient),
	}
}

func (c *Client) Add(pfCode comm.ThirdPlatformCode, client IClient) {
	c.clientMap[pfCode] = client
}

func (c *Client) ByPfCode(pfCode comm.ThirdPlatformCode) IClient {
	return c.clientMap[pfCode]
}
