package tea_im

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go-cs/pkg/stream"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	platformCode string
	privateKey   string
	httpClient   *http.Client
	domain       string
	debug        bool
}

func NewClient(platformCode, privateKey, domain string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		platformCode: platformCode,
		privateKey:   privateKey,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   time.Second * 10,
		},
		domain: domain,
		debug:  true,
	}
}

func (c *Client) Bind(chatToken string) (*Response, error) {
	// 请求参数
	argsMap := map[string]string{
		"pf_code":    c.platformCode,
		"chat_token": chatToken,
	}

	// 签名
	argsMap["sign"] = c.sign(argsMap)

	marshal, _ := json.Marshal(argsMap)
	req, err := http.NewRequest(http.MethodPost, c.domain+API_BIND_URL, bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	all, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) UnBind(chatToken string) (*Response, error) {
	// 请求参数
	argsMap := map[string]string{
		"pf_code":    c.platformCode,
		"chat_token": chatToken,
	}

	// 签名
	argsMap["sign"] = c.sign(argsMap)

	marshal, _ := json.Marshal(argsMap)
	req, err := http.NewRequest(http.MethodPost, c.domain+API_UNBIND_URL, bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	all, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) Push(content string, userIds []int64) (*Response, error) {
	userIdsStr := stream.Map(userIds, func(v int64) string {
		return strconv.FormatInt(v, 10)
	})

	// 请求参数
	argsMap := map[string]string{
		"pf_code":  c.platformCode,
		"content":  content,
		"user_ids": strings.Join(userIdsStr, ","),
	}

	// 签名
	argsMap["sign"] = c.sign(argsMap)

	marshal, _ := json.Marshal(argsMap)
	req, err := http.NewRequest(http.MethodPost, c.domain+API_PUSH_URL, bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	all, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}

func (c *Client) PushRobotMessage(content *RobotMessage, chatTokens []string) (*Response, error) {
	// 请求参数
	argsMap := map[string]string{
		"pf_code":     c.platformCode,
		"content":     string(content.Marshal()),
		"chat_tokens": strings.Join(chatTokens, ","),
	}

	// 签名
	argsMap["sign"] = c.sign(argsMap)

	marshal, _ := json.Marshal(argsMap)
	req, err := http.NewRequest(http.MethodPost, c.domain+API_PUSH_URL, bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	all, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(response.Msg)
	}

	return &response, nil
}

func (c *Client) GetUserInfo(chatToken string) (*UserInfo, error) {
	// 请求参数
	argsMap := map[string]string{
		"pf_code":    c.platformCode,
		"chat_token": chatToken,
	}

	// 签名
	argsMap["sign"] = c.sign(argsMap)

	marshal, _ := json.Marshal(argsMap)
	req, err := http.NewRequest(http.MethodPost, c.domain+API_GET_USER_INFO, bytes.NewReader(marshal))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	all, err := c.doRequest(req)

	var response Response
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if response.Code != 0 {
		return nil, errors.New(response.Msg)
	}

	var userInfo UserInfo
	err = json.Unmarshal(response.Data, &userInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshal userInfo: %w", err)
	}

	return &userInfo, nil
}

func (c *Client) doRequestWithRetry(req *http.Request) (body []byte, err error) {
	retryLimit := 3

	for i := 0; i < retryLimit; i++ {
		body, err = c.doRequest(req)
		if err == nil {
			return
		}
	}

	return
}

func (c *Client) doRequest(req *http.Request) (body []byte, err error) {
	if c.debug {
		dumpReq(req)
	}

	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readAll: %w", err)
	}

	return
}

func dumpReq(req *http.Request) {
	dumpRequest, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(fmt.Errorf("dump: %w", err))
	}
	fmt.Println(string(dumpRequest))
}
