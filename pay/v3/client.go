package pay

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dawndiy/go-wechat/httpclient"
	"github.com/dawndiy/go-wechat/internal"
)

const (
	// URLv3 接口 URL
	URLv3 = "https://api.mch.weixin.qq.com/v3/"
	// UserAgent 默认 UserAgent
	UserAgent = "go-wechat/pay"
)

type service struct {
	client *Client
}

// Client 微信支付接口客户端
type Client struct {
	client    *httpclient.Client
	baseURL   *url.URL
	userAgent string

	mchid  string // 商户号
	apiKey string // 商户APIv3密钥

	apiClientCert []byte // apiclient_cert.pem
	apiClientKey  []byte // apiclient_key.pem

	wechatPayCerts []WeChatPayCert // 平台证书(非商户)

	common service

	Base *BaseService

	// 营销
	Marketing struct {
		// 代金券
		Favor *MarketingFavorService
		// 商家券
		Busifavor *MarketingBusifavorService
		// 媒体上传
		Media *MarketingMediaService
	}

	// 基础支付
	Pay struct {
		// 合单支付
		Combine *PayCombineService
		// 账单
		//Bill *PaymentBillService
	}
}

// NewClient 新建一个微信支付接口客户端
func NewClient(opts ...ClientOption) *Client {
	c := new(Client)
	c.client = httpclient.NewClient()
	c.baseURL, _ = url.Parse(URLv3)
	c.userAgent = UserAgent

	c.common.client = c

	c.Base = (*BaseService)(&c.common)

	c.Marketing.Favor = (*MarketingFavorService)(&c.common)
	c.Marketing.Busifavor = (*MarketingBusifavorService)(&c.common)
	c.Marketing.Media = (*MarketingMediaService)(&c.common)

	c.Pay.Combine = (*PayCombineService)(&c.common)

	for _, opt := range opts {
		opt(c)
	}

	return c
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
	req.Header.Set("Accept", "application/json")

	if err := c.signRequest(req); err != nil {
		return nil, err
	}

	return req, nil
}

// Do 处理一个 HTTP 请求并返回响应对象和错误
//
// 如果服务端响应业务错误，返回 *http.Response 的同时，error 也会返回解析过后的 ErrorResponse
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	return internal.DoJSONRequest(c.client, req, v, c.checkResponse)
}

// Use 使用请求中间件
func (c *Client) Use(middlwares ...func(httpclient.RequestHandler) httpclient.RequestHandler) {
	c.client.Use(middlwares...)
}

// signRequest 生成请求签名
//
// 文档: https://wechatpay-api.gitbook.io/wechatpay-api-v3/qian-ming-zhi-nan-1/qian-ming-sheng-cheng
func (c *Client) signRequest(req *http.Request) error {

	if c.apiClientCert == nil || c.apiClientKey == nil {
		return fmt.Errorf("cert or key not set")
	}

	if c.mchid == "" {
		return fmt.Errorf("mchid not set")
	}

	httpMethod := req.Method
	urlStr := req.URL.RequestURI()
	timestamp := time.Now().Unix()
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	nonceStr := hex.EncodeToString(randBytes)
	var body string
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body.Close()
		body = string(b)
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}

	// HTTP请求方法\n
	// URL\n
	// 请求时间戳\n
	// 请求随机串\n
	// 请求报文主体\n
	elems := []string{
		httpMethod,
		urlStr,
		fmt.Sprint(timestamp),
		nonceStr,
		body,
	}
	sign, err := CalcSign(elems, c.apiClientKey)
	if err != nil {
		return err
	}

	serialNum, err := GetCertSerialNumber(c.apiClientCert)
	if err != nil {
		return err
	}

	wechatAuthorization := fmt.Sprintf(
		`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%d",serial_no="%s"`,
		c.mchid,
		nonceStr,
		sign,
		timestamp,
		serialNum,
	)

	req.Header.Set("Authorization", wechatAuthorization)
	return nil
}

func (c *Client) getWeChatPayCert(serialNo string) (*WeChatPayCert, error) {
	if serialNo == "" {
		return nil, fmt.Errorf("serialNo is empty")
	}
	for _, v := range c.wechatPayCerts {
		if v.SerialNo == serialNo {
			return &v, nil
		}
	}

	certs, err := c.Base.Certificates(context.Background())
	if err != nil {
		return nil, err
	}
	c.wechatPayCerts = certs

	for _, v := range c.wechatPayCerts {
		if v.SerialNo == serialNo {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("serialNo '%s' not found", serialNo)
}

func (c *Client) checkResponse(r *http.Response) error {
	// Certificates 接口设置 disableSignCheck 不校验应答签名
	// 429 - Too Many Requests 没有应答签名
	if !disableSignCheckFromContext(r.Request.Context()) || r.StatusCode != 429 {
		if err := c.checkResponseSign(r); err != nil {
			return err
		}
	}
	return CheckResponse(r)
}

// checkResponseSign 响应签名验证
//
// 文档: https://wechatpay-api.gitbook.io/wechatpay-api-v3/qian-ming-zhi-nan-1/qian-ming-yan-zheng
func (c *Client) checkResponseSign(r *http.Response) error {
	timestamp := r.Header.Get("Wechatpay-Timestamp")
	nonce := r.Header.Get("Wechatpay-Nonce")
	serial := r.Header.Get("Wechatpay-Serial")
	respSignB64 := r.Header.Get("Wechatpay-Signature")

	respSignBytes, err := base64.StdEncoding.DecodeString(respSignB64)
	if err != nil {
		return err
	}

	cert, err := c.getWeChatPayCert(serial)
	if err != nil {
		return err
	}

	body, err := readResponse(r)
	if err != nil {
		return err
	}

	elems := []string{timestamp, nonce, string(body)}
	raw := strings.Join(elems, "\n") + "\n"
	h := sha256.Sum256([]byte(raw))

	err = rsa.VerifyPKCS1v15(cert.Certficate.PublicKey.(*rsa.PublicKey), crypto.SHA256, h[:], respSignBytes)
	return err
}
