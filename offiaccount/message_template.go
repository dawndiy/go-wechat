package offiaccount

import (
	"context"
)

// MessageTemplateService 消息管理服务-模板消息
type MessageTemplateService service

// MessageTemplate 模板信息
type MessageTemplate struct {
	// 模板ID
	TemplateID string `json:"template_id"`
	// 模板标题
	Title string `json:"title"`
	// 模板所属行业的一级行业
	PrimaryIndustry string `json:"primary_industry"`
	// 模板所属行业的二级行业
	DeputyIndustry string `json:"deputy_industry"`
	// 模板内容
	Content string `json:"content"`
	// 模板示例
	Example string `json:"example"`
}

// GetAllPrivateTemplate 获取模板列表
//
// 获取已添加至帐号下所有模板列表，可在微信公众平台后台中查看模板列表信息。
// 为方便第三方开发者，提供通过接口调用的方式来获取帐号下所有模板信息
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html#3
func (s *MessageTemplateService) GetAllPrivateTemplate(ctx context.Context) ([]MessageTemplate, error) {

	u, err := s.client.apiURL(ctx, "cgi-bin/template/get_all_private_template", nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		TemplateList []MessageTemplate `json:"template_list"`
	}

	_, err = s.client.Do(req, &data)
	return data.TemplateList, err
}

// DelPrivateTemplate 删除模板
//
// 删除模板可在微信公众平台后台完成，为方便第三方开发者，提供通过接口调用的方式来删除某帐号下的模板
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html#4
func (s *MessageTemplateService) DelPrivateTemplate(ctx context.Context, templateID string) error {

	u, err := s.client.apiURL(ctx, "cgi-bin/template/del_private_template", nil)
	if err != nil {
		return err
	}
	body := map[string]string{
		"template_id": templateID,
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}

// MessageTemplateMsg 发送模板消息
type MessageTemplateMsg struct {
	// 接收者openid
	// 必填
	ToUser string `json:"touser"`
	// 模板ID
	// 必填
	TemplateID string `json:"template_id"`
	// 模板跳转链接（海外帐号没有跳转能力）
	// 选填
	URL string `json:"url,omitempty"`
	// 跳小程序所需数据，不需跳小程序可不用传该数据
	// 选填
	MiniProgram *MessageTemplateMsgMiniProgram `json:"miniprogram,omitempty"`
	// 模板数据
	// 必填
	Data map[string]MessageTemplateMsgColorValue `json:"data"`
}

// MessageTemplateMsgColorValue 模板数据
type MessageTemplateMsgColorValue struct {
	// 模板内容字体颜色，不填默认为黑色
	Color string `json:"color,omitempty"`
	// 内容值
	Value string `json:"value"`
}

// MessageTemplateMsgMiniProgram 跳小程序所需数据
type MessageTemplateMsgMiniProgram struct {
	// 所需跳转到的小程序appid
	// （该小程序appid必须与发模板消息的公众号是绑定关联关系，暂不支持小游戏）
	APPID string `json:"appid"`
	// 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），
	// 要求该小程序已发布，暂不支持小游戏
	PagePath string `json:"pagepath,omitempty"`
}

// Send 发送模板消息
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html#5
func (s *MessageTemplateService) Send(ctx context.Context, msg *MessageTemplateMsg) (int64, error) {

	u, err := s.client.apiURL(ctx, "cgi-bin/message/template/send", nil)
	if err != nil {
		return 0, err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), msg)
	if err != nil {
		return 0, err
	}

	var data struct {
		MsgID int64 `json:"msgid"`
	}

	_, err = s.client.Do(req, &data)
	return data.MsgID, err
}
