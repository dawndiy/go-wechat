package offiaccount

import (
	"context"
	"net/url"

	"github.com/dawndiy/go-wechat/pkg/token"
)

const tokenGrantType = "client_credential"

// BaseService 基础服务
type BaseService service

// Token 获取 AccessToken
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func (s *BaseService) GetAccessToken(ctx context.Context, appid, secret string) (*token.AccessToken, error) {
	u, _ := url.Parse("cgi-bin/token")
	v := url.Values{}
	v.Set("grant_type", tokenGrantType)
	v.Set("appid", appid)
	v.Set("secret", secret)
	u.RawQuery = v.Encode()

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	accessToken := new(token.AccessToken)
	_, err = s.client.Do(req, accessToken)

	return accessToken, err
}
