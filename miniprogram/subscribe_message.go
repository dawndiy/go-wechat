package miniprogram

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// SubscribeMessageService 订阅消息
type SubscribeMessageService service

// SubscribeMessageCategory 小程序账号的类目
type SubscribeMessageCategory struct {
	// 类目id，查询公共库模版时需要
	ID int `json:"id"`
	// 类目的中文名
	Name string `json:"name"`
}

// GetCategory 获取小程序账号的类目
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getCategory.html
func (s *SubscribeMessageService) GetCategory(ctx context.Context) ([]SubscribeMessageCategory, error) {
	u, err := s.client.apiURL(ctx, "wxaapi/newtmpl/getcategory", nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data []SubscribeMessageCategory `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

// 订阅模版类型
const (
	SubscribeMessageTypeOnce = 2 // 一次性订阅
	SubscribeMessageTypeLong = 3 // 长期订阅
)

// SubscribeMessageTemplate 个人模板
type SubscribeMessageTemplate struct {
	// 添加至帐号下的模板 id，发送小程序订阅消息时所需
	PriTmplID string `json:"priTmplId"`
	// 模版标题
	Title string `json:"title"`
	// 模版内容
	Content string `json:"content"`
	// 模板内容示例
	Example string `json:"example"`
	// 模版类型，2 为一次性订阅，3 为长期订阅
	Type int `json:"type"`
}

// GetTemplateList 获取当前帐号下的个人模板列表
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getTemplateList.html
func (s *SubscribeMessageService) GetTemplateList(ctx context.Context) ([]SubscribeMessageTemplate, error) {
	u, err := s.client.apiURL(ctx, "wxaapi/newtmpl/gettemplate", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data []SubscribeMessageTemplate `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data, err
}

// SubscribeMessagePubTemplateTitle 模板标题
type SubscribeMessagePubTemplateTitle struct {
	// 模版标题 id
	TID int `json:"tid"`
	// 模版标题
	Title string `json:"title"`
	// 模版类型，2 为一次性订阅，3 为长期订阅
	Type int `json:"type"`
	// 模版所属类目 id
	CategoryID string `json:"categoryId"`
}

// GetPubTemplateTitleList 获取帐号所属类目下的公共模板标题
//
// ids: 类目 id 列表.
// start: 用于分页，表示从 start 开始。从 0 开始计数.
// limit: 用于分页，表示拉取 limit 条记录。最大为 30
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getPubTemplateTitleList.html
func (s *SubscribeMessageService) GetPubTemplateTitleList(
	ctx context.Context, ids []string, start, limit int) (int64, []SubscribeMessagePubTemplateTitle, error) {

	v := url.Values{}
	v.Set("ids", strings.Join(ids, ","))
	v.Set("start", fmt.Sprint(start))
	v.Set("limit", fmt.Sprint(limit))

	u, err := s.client.apiURL(ctx, "wxaapi/newtmpl/getpubtemplatetitles", v)
	if err != nil {
		return 0, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, nil, err
	}

	var data struct {
		Count int64                              `json:"count"`
		Data  []SubscribeMessagePubTemplateTitle `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Count, data.Data, err
}

// SubscribeMessagePubTemplateKeyWords 模板标题下的关键词
type SubscribeMessagePubTemplateKeyWords struct {
	// 关键词 id，选用模板时需要
	KID int `json:"kid"`
	// 关键词内容
	Name string `json:"name"`
	// 关键词内容对应的示例
	Example string `json:"example"`
	// 参数类型
	Rule string `json:"rule"`
}

// GetPubTemplateKeyWordsByID 获取模板标题下的关键词列表
//
// tid: 模板标题 id，可通过接口获取.
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.getPubTemplateKeyWordsById.html
func (s *SubscribeMessageService) GetPubTemplateKeyWordsByID(
	ctx context.Context, tid int) (int64, []SubscribeMessagePubTemplateKeyWords, error) {

	v := url.Values{}
	v.Set("tid", fmt.Sprint(tid))

	u, err := s.client.apiURL(ctx, "wxaapi/newtmpl/getpubtemplatekeywords", v)
	if err != nil {
		return 0, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, nil, err
	}

	var data struct {
		Count int64                                 `json:"count"`
		Data  []SubscribeMessagePubTemplateKeyWords `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Count, data.Data, err
}

// SubscribeMessage 订阅消息
type SubscribeMessage struct {
	// 接收者（用户）的 openid
	ToUser string `json:"touser"`
	// 所需下发的订阅模板id
	TemplateID string `json:"template_id"`
	// 点击模板卡片后的跳转页面，仅限本小程序内的页面。
	// 支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Page string `json:"page,omitempty"`
	// 模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }
	Data interface{} `json:"data"`
	// 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	MiniProgramState string `json:"miniprogram_state,omitempty"`
	// 进入小程序查看”的语言类型，
	// 支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
	Lang string `json:"lang,omitempty"`
}

// Send 发送订阅消息
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
func (s *SubscribeMessageService) Send(ctx context.Context, msg *SubscribeMessage) error {
	u, err := s.client.apiURL(ctx, "cgi-bin/message/subscribe/send", nil)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), msg)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
