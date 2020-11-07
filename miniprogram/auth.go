package miniprogram

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

// 授权类型 grant_type
const (
	tokenGrantType  = "client_credential"
	jsCodeGrantType = "authorization_code"
)

const (
	apiToken             = "cgi-bin/token"
	apiSNSJSCode2Session = "sns/jscode2session"
	apiWXAGetPaidUnionID = "wxa/getpaidunionid"
)

// AuthService 认证服务
type AuthService service

// CodeSession 登录凭证
type CodeSession struct {
	// 用户唯一标识
	OpenID string `json:"openid"`
	// 会话密钥
	SessionKey string `json:"session_key"`
	// 用户在开放平台的唯一标识符，
	// 在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
	UnionID string `json:"unionid"`
}

// Code2Session 登录凭证校验。
// 通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程。更多使用方法详见 小程序登录。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
func (s *AuthService) Code2Session(ctx context.Context, jsCode, appid, secret string) (*CodeSession, error) {
	u, _ := url.Parse(apiSNSJSCode2Session)
	v := url.Values{}
	v.Set("appid", appid)
	v.Set("secret", secret)
	v.Set("js_code", jsCode)
	v.Set("grant_type", jsCodeGrantType)
	u.RawQuery = v.Encode()

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	sess := new(CodeSession)
	_, err = s.client.Do(req, sess)

	return sess, err

}

// AccessToken 接口调用凭据
type AccessToken struct {
	Value     string `json:"access_token"` // 获取到的凭证
	ExpiresIn int64  `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

// GetAccessToken 获取小程序全局唯一后台接口调用凭据（access_token）。
// 调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (s *AuthService) GetAccessToken(ctx context.Context, appid, secret string) (*AccessToken, error) {
	u, _ := url.Parse(apiToken)
	v := url.Values{}
	v.Set("grant_type", tokenGrantType)
	v.Set("appid", appid)
	v.Set("secret", secret)
	u.RawQuery = v.Encode()

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	accessToken := new(AccessToken)
	_, err = s.client.Do(req, accessToken)

	return accessToken, err
}

// PaidUnianIDOptions 获取支付用户 UnionId 选项
type PaidUnianIDOptions struct {
	TransactionID string `url:"transaction_id,omitempty"` // 微信支付订单号
	MCHID         string `url:"mchid,omitempty"`          // 微信支付分配的商户号，和商户订单号配合使用
	OutTradeNo    string `url:"out_trade_no,omitempty"`   // 微信支付商户订单号，和商户号配合使用
}

// Valid 选项是否有效
func (o *PaidUnianIDOptions) Valid() bool {
	if o.TransactionID != "" {
		return true
	}
	if o.MCHID != "" && o.OutTradeNo != "" {
		return true
	}
	return false
}

// GetPaidUnianID 用户支付完成后，获取该用户的 UnionId，无需用户授权。
//
// 注意：调用前需要用户完成支付，且在支付后的五分钟内有效。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/user-info/auth.getPaidUnionId.html
func (s *AuthService) GetPaidUnianID(ctx context.Context, openID string, opt *PaidUnianIDOptions) (string, error) {

	if opt == nil || !opt.Valid() {
		return "", fmt.Errorf("options invalid")
	}

	v, err := query.Values(opt)
	if err != nil {
		return "", err
	}

	u, err := s.client.apiURL(ctx, apiWXAGetPaidUnionID, v)
	if err != nil {
		return "", err
	}

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return "", err
	}

	var data struct {
		UnionID string `json:"unionid"`
	}
	_, err = s.client.Do(req, &data)
	return data.UnionID, err
}
