package offiaccount

import (
	"context"
	"net/url"
)

// UserService 用户管理服务
type UserService service

// UserInfo 用户信息
type UserInfo struct {
	// 用户是否订阅该公众号标识，值为0时，
	// 代表此用户没有关注该公众号，拉取不到其余信息。
	Subscribe int `json:"subscribe"`
	// 用户的标识，对当前公众号唯一
	OpenID string `json:"openid"`
	// 用户的昵称
	NickName string `json:"nickname"`
	// 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Sex int `json:"sex"`
	// 用户所在城市
	City string `json:"city"`
	// 用户所在国家
	Country string `json:"country"`
	// 用户所在省份
	Province string `json:"province"`
	// 用户的语言，简体中文为zh_CN
	Language string `json:"language"`
	// 用户头像，最后一个数值代表正方形头像大小
	// （有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	HeadImgURL string `json:"headimgurl"`
	// 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	SubscribeTime int64 `json:"subscribe_time"`
	// 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	UnionID string `json:"unionid"`
	// 公众号运营者对粉丝的备注，
	// 公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	Remark string `json:"remark"`
	// 用户所在的分组ID（兼容旧的用户分组接口）
	GroupID string `json:"groupid"`
	// 用户被打上的标签ID列表
	TagIDList []int `json:"tagid_list"`
	// 返回用户关注的渠道来源，
	// ADD_SCENE_SEARCH 公众号搜索，
	// ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，
	// ADD_SCENE_PROFILE_CARD 名片分享，
	// ADD_SCENE_QR_CODE 扫描二维码，
	// ADD_SCENE_PROFILE_LINK 图文页内名称点击，
	// ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，
	// ADD_SCENE_PAID 支付后关注，
	// ADD_SCENE_WECHAT_ADVERTISEMENT 微信广告，
	// ADD_SCENE_OTHERS 其他
	SubscribeScene string `json:"subscribe_scene"`
	// 二维码扫码场景（开发者自定义）
	QRScene int `json:"qr_scene"`
	// 二维码扫码场景描述（开发者自定义）
	QRSceneStr string `json:"qr_scene_str"`
}

// 获取用户基本信息(UnionID机制)
//
// 在关注者与公众号产生消息交互后，公众号可获得关注者的OpenID（加密后的微信号，
// 每个用户对每个公众号的OpenID是唯一的。对于不同公众号，同一用户的openid不同）。
// 公众号可通过本接口来根据OpenID获取用户基本信息，包括昵称、头像、性别、所在城市、语言和关注时间
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
func (s *UserService) Info(ctx context.Context, openID string, lang string) (*UserInfo, error) {
	v := url.Values{}
	v.Set("openid", openID)
	if lang != "" {
		v.Set("lang", lang)
	}
	u, err := s.client.apiURL(ctx, "cgi-bin/user/info", v)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	info := new(UserInfo)
	_, err = s.client.Do(req, info)

	return info, err
}
