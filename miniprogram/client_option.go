package miniprogram

import (
	"net/http"
	"net/url"

	"github.com/dawndiy/go-wechat/util/token"
)

// ClientOption 客户端设置选项
type ClientOption func(*Client)

// WithHTTPClient 设置 HTTP 客户端
func WithHTTPClient(httpclient *http.Client) ClientOption {
	return func(c *Client) {
		c.client.Client = httpclient
	}
}

// WithUserAgent 设置 User-Agent
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithAccessTokenStore 设置 Token 存储
func WithAccessTokenStore(store token.Store) ClientOption {
	return func(c *Client) {
		c.accessTokenStore = store
	}
}

// WithAPPIDSecret 设置小程序 APPID 和 Secret
func WithAPPIDSecret(appid, secret string) ClientOption {
	return func(c *Client) {
		c.appid, c.secret = appid, secret
	}
}

// WithBaseURL 设置请求基础 URL
// 默认 https://api.wexin.qq.com
func WithBaseURL(u *url.URL) ClientOption {
	return func(c *Client) {
		c.baseURL = u
	}
}
