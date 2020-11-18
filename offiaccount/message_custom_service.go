package offiaccount

import (
	"context"

	"github.com/dawndiy/go-wechat/pkg/message"
)

// MessageCustomService 消息管理服务-客服消息
type MessageCustomService service

// Send 客服接口-发消息
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#6
func (s *MessageCustomService) Send(
	ctx context.Context, toUser string, msg message.Message, kfAccount ...string) error {

	body := map[string]interface{}{
		"touser":   toUser,
		"msgtype":  msg.Type(),
		msg.Type(): msg,
	}
	if len(kfAccount) > 0 && kfAccount[0] != "" {
		body["customservice"] = map[string]string{
			"kf_account": kfAccount[0],
		}
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

// SetTyping 客服输入状态
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#8
func (s *MessageCustomService) SetTyping(ctx context.Context, toUser string, command message.StatusCommand) error {
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
