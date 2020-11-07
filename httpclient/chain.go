package httpclient

import "net/http"

// RequestHandler 请求处理器
type RequestHandler interface {
	// Do 处理 *http.Request 发起 HTTP 请求并得到 *http.Response
	Do(*http.Request) (*http.Response, error)
}

// RequestHandlerFunc 请求处理器方法
//
// 它是一个 RequestHandler 的包装方法，直接将 func(*http.Request) (*http.Response, error) 包装成 RequestHandler
type RequestHandlerFunc func(*http.Request) (*http.Response, error)

// Do 处理 *http.Request 发起 HTTP 请求并得到 *http.Response
func (h RequestHandlerFunc) Do(r *http.Request) (*http.Response, error) { return h(r) }

// Chain 将 []func(RequestHandler) RequestHandler 包装成 Middlewares
func Chain(middlwares ...func(RequestHandler) RequestHandler) Middlewares {
	return Middlewares(middlwares)
}

// Middlewares 请求处理中间件列表
type Middlewares []func(RequestHandler) RequestHandler

// RequestHandler 将 h 附加中间件
func (mws Middlewares) RequestHandler(h RequestHandler) RequestHandler {
	return &ChainHandler{mws, h, chain(mws, h)}
}

// RequestHandlerFunc 将 h 附加中间件
func (mws Middlewares) RequestHandlerFunc(h RequestHandlerFunc) RequestHandler {
	return &ChainHandler{mws, h, chain(mws, h)}
}

// ChainHandler 请求处理链
type ChainHandler struct {
	Middlewares Middlewares
	Endpoint    RequestHandler
	chain       RequestHandler
}

// Do 处理 *http.Request 发起 HTTP 请求并得到 *http.Response
func (c *ChainHandler) Do(req *http.Request) (*http.Response, error) {
	return c.chain.Do(req)
}

func chain(middlwares []func(RequestHandler) RequestHandler, endpoint RequestHandler) RequestHandler {
	size := len(middlwares)
	if size == 0 {
		return endpoint
	}

	h := middlwares[size-1](endpoint)
	for i := size - 2; i >= 0; i-- {
		h = middlwares[i](h)
	}
	return h
}
