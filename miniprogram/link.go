package miniprogram

import "context"

// LinkServie URL 链接服务
type LinkServie service

// GenWXAShortLink 获取小程序 Short Link
//
// 适用于微信内拉起小程序的业务场景。目前只开放给电商类目(具体包含以下一级类目：电商平台、商家自营、跨境电商)。
// 通过该接口，可以选择生成到期失效和永久有效的小程序短链
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/short-link/shortlink.generate.html
func (s *LinkServie) GenWXAShortLink(ctx context.Context, pageURL string, pageTitle string, isPermanent bool) (string, error) {
	u, err := s.client.apiURL(ctx, "wxa/genwxashortlink", nil)
	if err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"page_url":     pageURL,
		"page_title":   pageTitle,
		"is_permanent": isPermanent,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}
	var data struct {
		Link string `json:"link"`
	}
	_, err = s.client.Do(req, &data)
	return data.Link, err
}

// URLGenerateOption URL 生成选项
type URLGenerateOption struct {
	Path           string `json:"path"`
	Query          string `json:"query"`
	IsExpire       bool   `json:"is_expire"`
	ExpireType     int    `json:"expire_type"`
	ExpireTime     int64  `json:"expire_time"`
	ExpireInterval int    `json:"expire_interval"`
}

// GenerateURLLink 获取小程序 URL Link
//
// 适用于短信、邮件、网页、微信内等拉起小程序的业务场景。
// 通过该接口，可以选择生成到期失效和永久有效的小程序链接，有数量限制
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-link/urllink.generate.html
func (s *LinkServie) GenerateURLLink(ctx context.Context, option URLGenerateOption) (string, error) {
	u, err := s.client.apiURL(ctx, "wxa/generate_urllink", nil)
	if err != nil {
		return "", err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), option)
	if err != nil {
		return "", err
	}
	var data struct {
		URLLink string `json:"url_link"`
	}
	_, err = s.client.Do(req, &data)
	return data.URLLink, err
}

// GenerateURLScheme 获取小程序 scheme 码
//
// 用于短信、邮件、外部网页、微信内等拉起小程序的业务场景。
// 通过该接口，可以选择生成到期失效和永久有效的小程序码，有数量限制
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-scheme/urlscheme.generate.html
func (s *LinkServie) GenerateURLScheme(ctx context.Context, option URLGenerateOption) (string, error) {
	u, err := s.client.apiURL(ctx, "wxa/generatescheme", nil)
	if err != nil {
		return "", err
	}
	body := map[string]interface{}{
		"jump_wxa": map[string]string{
			"path":  option.Path,
			"query": option.Query,
		},
		"is_expire":       option.IsExpire,
		"expire_type":     option.ExpireType,
		"expire_time":     option.ExpireTime,
		"expire_interval": option.ExpireInterval,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}
	var data struct {
		OpenLink string `json:"openlink"`
	}
	_, err = s.client.Do(req, &data)
	return data.OpenLink, err
}
