package offiaccount

import "context"

// AccountService 帐号管理服务
type AccountService service

// 二维码类型
const (
	QRCodeCreateActionQRScene         = "QR_SCENE"           // 临时的整型参数值
	QRCodeCreateActionQRStrScene      = "QR_STR_SCENE"       // 临时的字符串参数值
	QRCodeCreateActionQRLimitScene    = "QR_LIMIT_SCENE"     // 永久的整型参数值
	QRCodeCreateActionQRLimitStrScene = "QR_LIMIT_STR_SCENE" // 永久的字符串参数值
)

// QRCodeActionInfo 二维码详细信息
type QRCodeActionInfo struct {
	// 场景值ID，临时二维码时为32位非0整型，
	// 永久二维码时最大值为100000（目前参数只支持1--100000）
	SceneID int32 `json:"scene_id,omitempty"`
	// 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
	SceneStr string `json:"scene_str,omitempty"`
}

// QRCodeCreate 生成带参数的二维码
//
// 为了满足用户渠道推广分析和用户帐号绑定等场景的需要，
// 公众平台提供了生成带参数二维码的接口。使用该接口可以获得多个带不同场景值的二维码，
// 用户扫描后，公众号可以接收到事件推送。
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
func (s *AccountService) QRCodeCreate(ctx context.Context,
	expiresIn int64, actionName string, actionInfo *QRCodeActionInfo) (string, int64, string, error) {

	u, err := s.client.apiURL(ctx, "cgi-bin/qrcode/create", nil)
	if err != nil {
		return "", 0, "", err
	}
	body := map[string]interface{}{
		"expire_seconds": expiresIn,
		"action_name":    actionName,
		"action_info": map[string]interface{}{
			"scene": actionInfo,
		},
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)

	var data struct {
		Ticket        string `json:"ticket"`
		ExpireSeconds int64  `json:"expire_seconds"`
		URL           string `json:"url"`
	}
	_, err = s.client.Do(req, &data)

	return data.Ticket, data.ExpireSeconds, data.URL, err
}

const (
	ShortURLActionLong2Short = "long2short" // 代表长链接转短链接
)

// ShortURL 长链接转短链接接口
//
// 将一条长链接转成短链接。 主要使用场景：
// 开发者用于生成二维码的原链接（商品、支付二维码等）
// 太长导致扫码速度和成功率下降，将原长链接通过此接口转成短链接再生成二维码将大大提升扫码速度和成功率。
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Account_Management/URL_Shortener.html
func (s *AccountService) ShortURL(ctx context.Context, action, longURL string) (string, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/shorturl", nil)
	if err != nil {
		return "", err
	}
	body := map[string]string{
		"action":   action,
		"long_url": longURL,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return "", err
	}

	var data struct {
		ShortURL string `json:"short_url"`
	}

	_, err = s.client.Do(req, &data)

	return data.ShortURL, err
}
