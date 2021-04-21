package pay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NotifyCouponUse 代金券核销事件回调通知
type NotifyCouponUse struct {
	// 创建批次的商户号 必填
	StockCreatorMCHID string `json:"stock_creator_mchid"`

	// 批次号 必填
	// 微信为每个代金券批次分配的唯一id
	StockID string `json:"stock_id"`

	// 代金券id 必填
	// 微信为代金券唯一分配的id
	CouponID string `json:"coupon_id"`

	// 单品优惠特定信息 可选
	SingleItemDicountOff *SingleItemDicountOff `json:"single_item_dicount_off,omitempty"`

	// 减至优惠特定信息 可选
	// 减至优惠限定字段，仅减至优惠场景有返回
	DiscountTo *DiscountTo `json:"discount_to,omitempty"`

	// 代金券名称 必填
	CouponName string `json:"coupon_name"`

	// 代金券状态 必填
	// SENDED：可用
	// USED：已实扣
	// EXPIRED：已过期
	Status string `json:"status"`

	// 代金券描述 必填
	Description string `json:"description"`

	// 领券时间 必填
	CreateTime time.Time `json:"create_time"`

	// 券类型 必填
	// NORMAL：满减券
	// CUT_TO：减至券
	CouponType string `json:"coupon_type"`

	// 是否无资金流 必填
	NoCash bool `json:"no_cash"`

	// 可用开始时间 必填
	AvailableBeginTime time.Time `json:"available_begin_time"`
	// 可用结束时间 必填
	AvailableEndTime time.Time `json:"available_end_time"`

	// 是否单品优惠
	SingleItem bool `json:"singleitem"`

	// 普通满减券信息 可选
	// 普通满减券面额、门槛信息
	NormalCouponInformation *CouponFixedNormalRule `json:"normal_coupon_information,omitempty"`

	// 实扣代金券信息 可选
	ConsumeInformation *ConsumeInformation `json:"consume_information,omitempty"`

	// 单品信息 可选
	// 商户下单接口传的单品信息
	GoodsDetail *ConsumeGoodsDetail `json:"goods_detail,omitempty"`
}

// SingleItemDicountOff 单品优惠特定信息
type SingleItemDicountOff struct {
	// 单品最高优惠价格 必填
	SinglePriceMax int64 `json:"single_price_max"`
}

// DiscountTo 减至优惠特定信息
type DiscountTo struct {
	// 减至后优惠单价 可选
	// 减至后优惠单价，单位：分
	CutToPrice int64 `json:"cut_to_price,omitempty"`

	// 最高价格 可选
	// 可享受优惠的最高价格，单位：分
	MaxPrice int64 `json:"max_price,omitempty"`
}

// ConsumeInformation 实扣代金券信息
type ConsumeInformation struct {
	// 核销时间 必填
	ConsumeTime time.Time `json:"consume_time"`

	// 核销商户号 必填
	ConsumeMCHID string `json:"consume_mchid"`

	// 核销订单号 必填
	TransactionID string `json:"transaction_id"`

	// 单品信息 可选
	// 商户下单接口传的单品信息
	GoodsDetail []ConsumeGoodsDetail `json:"goods_detail,omitempty"`
}

// ConsumeGoodsDetail 单品信息
type ConsumeGoodsDetail struct {
	// 单品编码 必填
	// 单品券创建时录入的单品编码
	GoodsID string `json:"goods_id"`

	// 单品数量 必填
	Quantity int `json:"quantity"`

	// 单品单价 必填
	Price int `json:"price"`

	// 优惠金额 必填
	DiscountAmount int `json:"discount_amount"`
}

// ParseNotify 解析核销事件回调通知
//
func (s *MarketingFavorService) ParseNotify(req *http.Request) (*NotifyCouponUse, error) {
	notify, err := s.client.ParseNotify(req, nil)
	if err != nil {
		return nil, err
	}

	if notify.EventType != EventTypeCouponUse {
		return nil, fmt.Errorf("通知类型错误")
	}

	rc := notify.Resource
	data, err := DecodeCiphertext(rc.Algorithm, rc.Ciphertext, rc.Nonce, rc.AssociatedData, s.client.apiKey)
	if err != nil {
		return nil, err
	}

	couponUse := new(NotifyCouponUse)
	err = json.Unmarshal(data, couponUse)
	return couponUse, err
}
