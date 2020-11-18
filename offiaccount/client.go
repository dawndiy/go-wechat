package offiaccount

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dawndiy/go-wechat/httpclient"
	"github.com/dawndiy/go-wechat/internal"
	"github.com/dawndiy/go-wechat/pkg/token"
)

const (
	// URL 微信接口UR
	URL = "https://api.weixin.qq.com"
)

type service struct {
	client *Client
}

// Client 微信公众号接口客户端
type Client struct {
	client *httpclient.Client
	common service

	accessTokenStore token.Store

	userAgent string
	baseURL   *url.URL

	appid  string // 公众号 appid
	secret string // 公众号 appSecret

	Base    *BaseService
	Account *AccountService
	Message struct {
		CustomService   *MessageCustomService
		TemplateService *MessageTemplateService
	}
}

// NewClient 新建一个公众号接口客户端
func NewClient(opts ...ClientOption) *Client {
	c := new(Client)
	c.client = httpclient.NewClient()
	c.baseURL, _ = url.Parse(URL)

	for _, opt := range opts {
		opt(c)
	}

	c.common.client = c
	c.Base = (*BaseService)(&c.common)
	c.Account = (*AccountService)(&c.common)
	c.Message.CustomService = (*MessageCustomService)(&c.common)
	c.Message.TemplateService = (*MessageTemplateService)(&c.common)
	return c
}

// Use 使用请求中间件
func (c *Client) Use(middlwares ...func(httpclient.RequestHandler) httpclient.RequestHandler) {
	c.client.Use(middlwares...)
}

// NewRequest 新建一个接口请求
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req, err := internal.NewJSONRequest(ctx, method, u, body)
	if err != nil {
		return nil, err
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

// NewUploadRequest 新建上传请求
func (c *Client) NewUploadRequest(
	ctx context.Context, urlStr, fieldname, filename string, reader io.Reader) (*http.Request, error) {

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return internal.NewUploadRequest(ctx, u, fieldname, filename, reader)
}

// Do 执行一个接口请求
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {

	return internal.DoJSONRequest(c.client, req, v, internal.CheckJSONResponse)
}

// GetAccessToken 获取接口调用凭证
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	if token, err := c.accessTokenStore.Get(); err == nil {
		return token, nil
	}

	if c.appid == "" || c.secret == "" {
		return "", fmt.Errorf("appid or secret not set")
	}

	accessToken, err := c.Base.GetAccessToken(ctx, c.appid, c.secret)
	if err != nil {
		return "", err
	}

	err = c.accessTokenStore.Save(accessToken.Value, accessToken.ExpiresIn)
	return accessToken.Value, err
}

// apiURL 传入 API PATH 和 url.Values ，自动附加上 access_token 返回 *url.URL
func (c *Client) apiURL(ctx context.Context, path string, value url.Values) (*url.URL, error) {
	var u *url.URL
	var err error
	if u, err = url.Parse(path); err != nil {
		return nil, err
	}
	var accessToken string
	if accessToken, err = c.GetAccessToken(ctx); err != nil {
		return nil, err
	}
	if value == nil {
		value = url.Values{}
	}
	value.Set("access_token", accessToken)

	u.RawQuery = value.Encode()
	return u, err
}
