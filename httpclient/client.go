package httpclient

import "net/http"

// Client 一个带中间件的 HTTP Client
type Client struct {
	*http.Client

	middlwares []func(RequestHandler) RequestHandler
}

// NewClient 新建一个 Client
func NewClient() *Client {
	c := new(Client)
	c.Client = new(http.Client)
	return c
}

// Use 添加一个请求中间件
func (c *Client) Use(middlwares ...func(RequestHandler) RequestHandler) {
	c.middlwares = append(c.middlwares, middlwares...)
}

// Do 执行请求
//
// *http.Request 会经过所有中间件最终得到 *http.Response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	chain := Chain(c.middlwares...)
	return chain.RequestHandler(c.Client).Do(req)
}
