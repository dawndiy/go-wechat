package pay

import (
	"io/ioutil"
	"net/http"
)

// ClientOption 客户端选项
type ClientOption func(*Client)

// WithMCHID 设置商户号
func WithMCHID(mchID string) ClientOption {
	return func(c *Client) {
		c.mchid = mchID
	}
}

// WithAPIKey 设置API密钥
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithAPIClientKeyPair 设置 apiclient_cert.pem , apiclient_key.pem 文件内容
func WithAPIClientKeyPair(cert, key []byte) ClientOption {
	return func(c *Client) {
		c.apiClientCert = cert
		c.apiClientKey = key
	}
}

// WithAPIClientKeyPairFile 设置 apiclient_cert.pem , apiclient_key.pem 文件
func WithAPIClientKeyPairFile(certFile, keyFile string) (ClientOption, error) {
	var cert, key []byte
	var err error
	if cert, err = ioutil.ReadFile(certFile); err != nil {
		return nil, err
	}
	if key, err = ioutil.ReadFile(keyFile); err != nil {
		return nil, err
	}
	return func(c *Client) {
		c.apiClientCert = cert
		c.apiClientKey = key
	}, nil
}

// WithHTTPClient 设置自定义 *http.Client
func WithHTTPClient(httpclient *http.Client) ClientOption {
	return func(c *Client) {
		c.client.Client = httpclient
	}
}
