package miniprogram

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dawndiy/go-wechat/pkg/message"
	"github.com/dawndiy/go-wechat/pkg/upload"
)

// CustomerServiceMessageService 客服消息服务
type CustomerServiceMessageService service

// GetTempMedia 获取客服消息内的临时素材。
//
// 即下载临时的多媒体文件。目前小程序仅支持下载图片文件。
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.getTempMedia.html
func (s *CustomerServiceMessageService) GetTempMedia(
	ctx context.Context, mediaID string) (http.Header, []byte, error) {

	v := url.Values{}
	v.Set("media_id", mediaID)
	u, err := s.client.apiURL(ctx, "cgi-bin/media/get", v)

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
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

// Send 发送客服消息给用户
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.send.html
func (s *CustomerServiceMessageService) Send(ctx context.Context, toUser string, msg message.Message) error {
	body := map[string]interface{}{
		"touser":   toUser,
		"msgtype":  msg.Type(),
		msg.Type(): msg,
	}
	u, err := s.client.apiURL(ctx, "cgi-bin/message/custom/send", nil)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// SetTyping 下发客服当前输入状态给用户
//
// command, 命令, Typing 对用户下发"正在输入"状态, CancelTyping 取消对用户的"正在输入"状态
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.setTyping.html
func (s *CustomerServiceMessageService) SetTyping(
	ctx context.Context, toUser string, command message.StatusCommand) error {

	u, err := s.client.apiURL(ctx, "cgi-bin/message/custom/send", nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"touser":  toUser,
		"command": string(command),
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// UploadTempMedia 把媒体文件上传到微信服务器。目前仅支持图片。用于发送客服消息或被动回复用户消息
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.uploadTempMedia.html
func (s *CustomerServiceMessageService) UploadTempMedia(
	ctx context.Context, mediaType string, up upload.Uploader) (*upload.Result, error) {

	v := url.Values{}
	v.Set("type", mediaType)
	u, err := s.client.apiURL(ctx, "cgi-bin/media/upload", v)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewUploadRequest(ctx, u.String(), "media", up.Name(), up.Reader(), nil)
	if err != nil {
		return nil, err
	}

	result := new(upload.Result)
	_, err = s.client.Do(req, result)

	return result, err
}
