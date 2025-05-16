package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"strings"
)

type Client struct {
	accessKeyId     string
	accessKeySecret string
	SignName        string
}

func NewClient(accessKeyId, accessKeySecret, signName string) *Client {
	return &Client{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		SignName:        signName,
	}
}

func (c *Client) Send(phoneNumber string, code, templateCode string) error {
	client, err := c.createClient(tea.String(c.accessKeyId), tea.String(c.accessKeySecret))
	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(c.SignName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
	}

	runtime := &util.RuntimeOptions{}

	tryErr := func() (err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		var _t *tea.SDKError
		if errors.As(tryErr, &_t) {
			err = _t
		}
		// 错误 message
		fmt.Println(tea.StringValue(err.Message))

		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(err.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, err2 := util.AssertAsString(err.Message)
		if err2 != nil {
			return err2
		}
	}
	return err
}

func (c *Client) createClient(accessKeyId *string, accessKeySecret *string) (client *dysmsapi.Client, err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        tea.String(Endpoint),
	}

	client, err = dysmsapi.NewClient(config)
	return
}
