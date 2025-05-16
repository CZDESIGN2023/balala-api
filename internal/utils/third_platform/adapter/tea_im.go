package adapter

import (
	"go-cs/internal/utils/third_platform"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
)

type IMClient struct {
	raw *tea_im.Client
}

func NewIMClient(platformCode, privateKey, domain string) third_platform.IClient {
	if platformCode == "" || privateKey == "" || domain == "" {
		return nil
	}

	return &IMClient{
		raw: tea_im.NewClient(platformCode, privateKey, domain),
	}
}

func (c *IMClient) Bind(chatToken string) error {
	_, err := c.raw.Bind(chatToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *IMClient) Unbind(chatToken string) error {
	_, err := c.raw.UnBind(chatToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *IMClient) Push(msg any, chatTokens []string) error {
	if msg == nil || len(chatTokens) == 0 {
		return nil
	}

	_, err := c.raw.PushRobotMessage(msg.(*tea_im.RobotMessage), chatTokens)
	if err != nil {
		return err
	}

	return nil
}

func (c *IMClient) GetUserInfo(key string) (*third_platform.UserInfo, error) {
	userInfo, err := c.raw.GetUserInfo(key)
	if err != nil {
		return nil, err
	}

	return &third_platform.UserInfo{
		Id:       int64(userInfo.Id),
		UserName: userInfo.UserName,
		NickName: userInfo.NickName,
		HeadAddr: userInfo.HeadAddr,
	}, nil
}
