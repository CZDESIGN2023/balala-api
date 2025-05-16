package http_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HTTPClientInterface interface {
	DoGet(ctx context.Context, path string) ([]byte, error)
	Get(ctx context.Context, path string, params *map[string]string) ([]byte, error)
	Post(ctx context.Context, path string, body interface{}) ([]byte, error)
}

type HTTPClient struct {
	name      string
	discovery HttpDiscoveryInterface
}

func NewHTTPClient(d HttpDiscoveryInterface) HTTPClientInterface {
	return &HTTPClient{
		name:      d.GetName(),
		discovery: d,
	}
}

func (c *HTTPClient) DoGet(ctx context.Context, path string) ([]byte, error) {

	str, err := c.discovery.GetServiceUrl(ctx, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, str, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s service returned status code %d", c.name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Get sends a GET request to the specified URL with the specified headers.
func (c *HTTPClient) Get(ctx context.Context, path string, params *map[string]string) ([]byte, error) {
	page, err := c.discovery.GetServiceUrl(ctx, path)
	if err != nil {
		return nil, err
	}

	if params != nil && len(*params) > 0 {
		values := url.Values{}
		for k, v := range *params {
			// 添加键值对。
			values.Set(k, v)
		}
		// 使用 Encode 方法将键值对编码成 URL 参数的字符串形式。
		p := values.Encode()
		page += "?" + p
	}

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d, response body: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}

// Post sends a POST request to the specified URL with the specified headers and body.
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	page, err := c.discovery.GetServiceUrl(ctx, path)
	if err != nil {
		return nil, err
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, page, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// for k, v := range headers {
	// 	req.Header.Set(k, v)
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// panic(err)
		return nil, err
	}

	defer resp.Body.Close()

	// 读取 HTTP 响应的内容。
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		// panic(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d, response body: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}
