package miniprogram

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/dawndiy/go-wechat/httpclient"
	"github.com/dawndiy/go-wechat/internal"
	"github.com/dawndiy/go-wechat/pkg/token"
	"github.com/google/go-querystring/query"
)

const (
	// URL 微信接口URL
	URL = "https://api.weixin.qq.com"
	// UserAgent 默认 User-Agent
	UserAgent = "go-wechat/miniprogram"
)

type service struct {
	client *Client
}

type Service struct {
	Client *Client
}

// Client 微信小程序服务端接口客户端
type Client struct {
	client *httpclient.Client
	common service

	accessTokenStore token.Store

	userAgent string
	baseURL   *url.URL

	appid  string // 小程序 appid
	secret string // 小程序 appSecret

	Auth                   *AuthService
	Analysis               *AnalysisService
	WXACode                *WXACodeService
	CustomerServiceMessage *CustomerServiceMessageService
	SubscribeMessage       *SubscribeMessageService
	// 插件管理
	PluginManager *PluginManagerService
	// 链接
	Link *LinkServie
	// 自定义交易组件
	Shop *ShopComponentShopService
}

// NewClient 新建一个新的小程序服务端接口客户端
func NewClient(opts ...ClientOption) *Client {

	c := new(Client)
	c.client = httpclient.NewClient()
	c.accessTokenStore = &token.MemoryStore{}
	c.userAgent = UserAgent
	c.baseURL, _ = url.Parse(URL)

	for _, opt := range opts {
		opt(c)
	}

	c.common.client = c
	c.Auth = (*AuthService)(&c.common)
	c.Analysis = (*AnalysisService)(&c.common)
	c.WXACode = (*WXACodeService)(&c.common)
	c.CustomerServiceMessage = (*CustomerServiceMessageService)(&c.common)
	c.SubscribeMessage = (*SubscribeMessageService)(&c.common)
	c.PluginManager = (*PluginManagerService)(&c.common)
	c.Link = (*LinkServie)(&c.common)
	c.Shop = (*ShopComponentShopService)(&c.common)

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
	ctx context.Context, urlStr, fieldname, filename string, reader io.Reader, fields map[string]string) (*http.Request, error) {

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return internal.NewUploadRequest(ctx, u, fieldname, filename, reader, fields)
}

// Do 执行一个接口请求
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {

	return internal.DoJSONRequest(c.client, req, v, CheckResponse)
}

// GetAccessToken 获取接口调用凭证
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	token, err := c.accessTokenStore.Get()
	if err == nil {
		return token, nil
	}

	if c.appid == "" || c.secret == "" {
		return "", err
	}

	accessToken, err := c.Auth.GetAccessToken(ctx, c.appid, c.secret)
	if err != nil {
		return "", err
	}

	err = c.accessTokenStore.Save(accessToken.Value, accessToken.ExpiresIn)
	return accessToken.Value, err
}

// URLWithToken 传入 API 路径(不需要开头 '/' ) 得到带 TOKEN 的 URL
func (c *Client) URLWithToken(ctx context.Context, path string, value url.Values) (*url.URL, error) {
	return c.apiURL(ctx, path, value)
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

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}
	values := u.Query()
	for k, v := range qs {
		for _, val := range v {
			values.Add(k, val)
		}
	}

	u.RawQuery = values.Encode()
	return u.String(), nil
}
