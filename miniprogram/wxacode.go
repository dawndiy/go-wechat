package miniprogram

import (
	"context"
	"io/ioutil"
	"net/http"
)

// WXACodeService 小程序码
type WXACodeService service

// QRCodeBaseConfig 二维码基础配置
type QRCodeBaseConfig struct {
	// 二维码的宽度，单位 px。最小 280px，最大 1280px
	Width int `json:"width,omitempty"`
	// 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调
	AutoColor bool `json:"auto_color,omitempty"`
	// auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	LineColor *LineColor `json:"line_color,omitempty"`
	// 是否需要透明底色，为 true 时，生成透明底色的小程序码
	IsHyaline bool `json:"is_hyaline,omitempty"`
}

// QRCodeConfig 二维码获取配置
type QRCodeConfig struct {
	// 扫码进入的小程序页面路径，最大长度 128 字节，不能为空；对于小游戏，可以只传入 query 部分，来实现传参效果，
	// 如：传入 "?foo=bar"，即可在 wx.getLaunchOptionsSync 接口中的 query 参数获取到 {foo:"bar"}。}
	Path string `json:"path"`

	QRCodeBaseConfig `json:",inline"`
}

// LineColor QRCode RGB 颜色, 十进制表示
type LineColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// CreateQRCode 获取小程序二维码，适用于需要的码数量较少的业务场景。
//
// 通过该接口生成的小程序码，永久有效，有数量限制，详见获取二维码。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.createQRCode.html
func (s *WXACodeService) CreateQRCode(ctx context.Context, qrConfig *QRCodeConfig) (http.Header, []byte, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/wxaapp/createwxaqrcode", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), qrConfig)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return resp.Header, b, err
}

// Get 获取小程序码，适用于需要的码数量较少的业务场景。
//
// 通过该接口生成的小程序码，永久有效，有数量限制，详见获取二维码。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.get.html
func (s *WXACodeService) Get(ctx context.Context, qrConfig *QRCodeConfig) (http.Header, []byte, error) {
	u, err := s.client.apiURL(ctx, "wxa/getwxacode", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), qrConfig)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	return resp.Header, b, err
}

// WXACodeUnlimitedConfig wxacode.getUnlimited 接口获取二维码配置
type WXACodeUnlimitedConfig struct {

	// 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，
	// 其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，请使用其他编码方式）
	Scene string `json:"scene"`

	// 必须是已经发布的小程序存在的页面（否则报错），
	// 例如 pages/index/index, 根路径前不要填加 /,不能携带参数（参数请放在scene字段里），
	// 如果不填写这个字段，默认跳主页面
	Page string `json:"page,omitempty"`

	// 检查page 是否存在，
	// 为 true 时 page 必须是已经发布的小程序存在的页面（否则报错）；
	// 为 false 时允许小程序未发布或者 page 不存在， 但page 有数量上限（60000个）请勿滥用
	CheckPath bool `json:"check_path,omitempty"`

	// 要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"
	EnvVersion string `json:"env_version,omitempty"`

	QRCodeBaseConfig `json:",inline"`
}

// GetUnlimited 获取小程序码，适用于需要的码数量较少的业务场景。
// 通过该接口生成的小程序码，永久有效，有数量限制，详见获取二维码。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (s *WXACodeService) GetUnlimited(ctx context.Context, conf *WXACodeUnlimitedConfig) (http.Header, []byte, error) {
	u, err := s.client.apiURL(ctx, "wxa/getwxacodeunlimit", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), conf)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return resp.Header, b, err
}
