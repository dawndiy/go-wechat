package pay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NotifyCouponSend 券领券事件回调通知
type NotifyCouponSend struct {
	// 业务细分事件类型
	EventType string `json:"event_type"`
	// 券的唯一标识
	CouponCode string `json:"coupon_code"`
	// 批次号
	StockID string `json:"stock_id"`
	// 发放时间
	SendTime time.Time `json:"send_time"`
	// 微信用户在appid下的唯一标识
	OpenID string `json:"openid"`
	// 用户统一标识
	UnionID string `json:"unionid"`
	// 发放渠道
	// 枚举值：
	// BUSICOUPON_SEND_CHANNEL_MINIAPP：小程序
	// BUSICOUPON_SEND_CHANNEL_API：API
	// BUSICOUPON_SEND_CHANNEL_PAYGIFT：支付有礼
	// BUSICOUPON_SEND_CHANNEL_H5：H5
	// BUSICOUPON_SEND_CHANNEL_FTOF：面对面
	// BUSICOUPON_SEND_CHANNEL_MEMBER_CARD_ACT：会员卡活动
	// BUSICOUPON_SEND_CHANNEL_HALL：扫码领券（营销馆）
	SendChannel string `json:"send_channel"`
	// 发券商户号
	SendMerchant string `json:"send_merchant"`
	// 发券附加信息
	// 仅在支付有礼、扫码领券（营销馆）、会员有礼发放渠道，才有该信息
	AttachInfo string `json:"attach_info"`

	bytes []byte `json:"-"`
}

// Bytes 返回解密后原始字节
func (n NotifyCouponSend) Bytes() []byte {
	return n.bytes
}

// ParseNotify 解析商家券领券事件回调通知
func (s *MarketingBusifavorService) ParseNotify(req *http.Request) (*NotifyCouponSend, error) {
	notify, err := s.client.ParseNotify(req, nil)
	if err != nil {
		return nil, err
	}
	if notify.EventType != EventTypeCouponSend {
		return nil, fmt.Errorf("通知类型错误")
	}

	rc := notify.Resource
	data, err := DecodeCiphertext(rc.Algorithm, rc.Ciphertext, rc.Nonce, rc.AssociatedData, s.client.apiKey)
	if err != nil {
		return nil, err
	}

	couponSend := new(NotifyCouponSend)
	err = json.Unmarshal(data, couponSend)
	couponSend.bytes = data
	return couponSend, err
}
