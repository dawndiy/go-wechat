package miniprogram

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

// 客服发送消息类别
const (
	CSMSendTypeText            = "text"
	CSMSendTypeImage           = "image"
	CSMSendTypeLink            = "link"
	CSMSendTypeMiniProgramPage = "miniprogrampage"
)

// MessageText 文本消息
type MessageText struct {
	Content string `json:"content"` // 文本消息内容
}

// MessageImage 图片消息
type MessageImage struct {
	MediaID string `json:"media_id"` // 发送的图片的媒体ID，通过 新增素材接口 上传图片文件获得。
}

// MessageLink 图文链接
type MessageLink struct {
	// 消息标题
	Title string `json:"title"`
	// 图文链接消息
	Description string `json:"description"`
	// 图文链接消息被点击后跳转的链接
	URL string `json:"url"`
	// 图文链接消息的图片链接，支持 JPG、PNG 格式，较好的效果为大图 640 X 320，小图 80 X 80
	ThumbURL string `json:"thumb_url"`
}

// MessageMiniProgram 小程序卡片
type MessageMiniProgram struct {
	// 消息标题
	Title string `json:"title"`
	// 小程序的页面路径，跟app.json对齐，支持参数，比如pages/index/index?foo=bar
	PagePath string `json:"page_path"`
	// 小程序消息卡片的封面， image 类型的 media_id，通过 新增素材接口 上传图片文件获得，建议大小为 520*416
	ThumbMediaID string `json:"thumb_media_id"`
}

// Send 发送客服消息给用户
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.send.html
func (s *CustomerServiceMessageService) Send(ctx context.Context, toUser, msgType string, message interface{}) error {
	switch msgType {
	case CSMSendTypeText,
		CSMSendTypeImage,
		CSMSendTypeLink,
		CSMSendTypeMiniProgramPage:
	default:
		return fmt.Errorf("unknown msgType '%s'", msgType)
	}
	body := map[string]interface{}{
		"touser":  toUser,
		"msgtype": msgType,
	}
	msgKey := ""
	switch message.(type) {
	case MessageText:
		msgKey = "text"
	case MessageImage:
		msgKey = "image"
	case MessageLink:
		msgKey = "link"
	case MessageMiniProgram:
		msgKey = "miniprogrampage"
	default:
		return fmt.Errorf("unknown message object")
	}
	body[msgKey] = message

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

// 客服输入状态
const (
	CSMTyping       = "Typing"
	CSMCancelTyping = "CancelTyping"
)

// SetTyping 下发客服当前输入状态给用户
//
// command, 命令, Typing 对用户下发"正在输入"状态, CancelTyping 取消对用户的"正在输入"状态
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/customer-message/customerServiceMessage.setTyping.html
func (s *CustomerServiceMessageService) SetTyping(ctx context.Context, toUser, command string) error {

	switch command {
	case CSMTyping, CSMCancelTyping:
	default:
		return fmt.Errorf("unknown command '%s'", command)
	}

	u, err := s.client.apiURL(ctx, "cgi-bin/message/custom/send", nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"touser":  toUser,
		"command": command,
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

	req, err := s.client.NewUploadRequest(ctx, u.String(), "media", up.Name(), up.Reader())
	if err != nil {
		return nil, err
	}

	result := new(upload.Result)
	_, err = s.client.Do(req, result)

	return result, err
}
