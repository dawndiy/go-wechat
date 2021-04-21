package pay

import (
	"context"
	"fmt"
	"time"
)

// EcommerceRefundService 电商收付通-退款
type EcommerceRefundService service

// RefundApply 退款申请
type RefundApply struct {
	// 二级商户号 必填
	// 微信支付分配二级商户的商户号
	SubMCHID string `json:"sub_mchid"`

	// 电商平台APPID 必填
	// 电商平台在微信公众平台申请服务号对应的APPID，
	// 申请商户功能的时候微信支付会配置绑定关系。
	SpAPPID string `json:"sp_appid"`

	// 二级商户APPID 可选
	// 二级商户在微信申请公众号成功后分配的帐号ID，
	// 需要电商平台侧配置绑定关系才能传参。
	SubAPPID string `json:"sub_appid,omitempty"`

	// 微信订单号 可选
	// 原支付交易对应的微信订单号
	TransactionID string `json:"transaction_id,omitempty"`

	// 商户订单号 可选
	// 原支付交易对应的商户订单号
	OutTradeNo string `json:"out_trade_no,omitempty"`

	// 商户退款单号 必填
	// 商户系统内部的退款单号，商户系统内部唯一，
	// 只能是数字、大小写字母_-|*@，同一退款单号多次请求只退一笔
	OutRefundNo string `json:"out_refund_no"`

	// 订单金额信息 必填
	Amount RefundAmount `json:"amount"`

	// 退款结果回调url 可选
	// 异步接收微信支付退款结果通知的回调地址，
	// 通知url必须为外网可访问的url，不能携带参数。
	// 如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效，优先回调当前传的地址。
	NotifyURL string `json:"notify_url,omitempty"`
}

// RefundAmount 订单金额信息
type RefundAmount struct {
	// 退款金额 必填
	// 退款金额，币种的最小单位，只能为整数，不能超过原订单支付金额
	Refund int64 `json:"refund"`

	// 原订单金额 必填
	// 原支付交易的订单总金额，币种的最小单位，只能为整数
	Total int64 `json:"total"`

	// 退款币种 可选
	// 符合ISO 4217标准的三位字母代码，目前只支持人民币：CNY
	Currency string `json:"currency,omitempty"`

	/////////// 查询结果字段 ///////////

	// 用户退款金额 查询必填
	// 退款给用户的金额，不包含所有优惠券金额
	PayerRefund int64 `json:"payer_refund,omitempty"`

	// 优惠退款金额 查询可选
	// 优惠券的退款金额，原支付单的优惠按比例退款
	DiscountRefound int64 `json:"discount_refund,omitempty"`
}

// RefundApplyResult 退款申请结果
type RefundApplyResult struct {
	// 微信退款单号 必填
	// 微信支付退款订单号
	RefundID string `json:"refund_id"`

	// 商户退款单号 必填
	// 商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔
	OutRefundNo string `json:"out_refund_no"`

	// 退款创建时间 必填
	CreateTime time.Time `json:"create_time"`

	// 订单金额信息 必填
	Amount RefundAmount `json:"amount"`

	// 优惠退款功能信息 可选
	PromotionDetail []RefundPromotionDetail `json:"promotion_detail"`
}

// RefundPromotionDetail 优惠退款详情
type RefundPromotionDetail struct {
	// 券ID 必填
	// 券或者立减优惠id
	PromotionID string `json:"promotion_id"`

	// 优惠范围 必填
	// GLOBAL：全场代金券
	// SINGLE：单品优惠
	Scope string `json:"scope"`

	// 优惠类型 必填
	// COUPON：充值型代金券，商户需要预先充值营销经费
	// DISCOUNT：免充值型优惠券，商户不需要预先充值营销经费
	Type string `json:"type"`

	// 优惠券面额 必填
	// 用户享受优惠的金额（优惠券面额=微信出资金额+商家出资金额+其他出资方金额 ）
	Amount int64 `json:"amount"`

	// 优惠退款金额 必填
	// 代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，
	// 说明详见《代金券或立减优惠》
	RefundAmount int64 `json:"refund_amount"`
}

// Apply 退款申请API
//
// 当交易发生之后一段时间内，由于买家或者卖家的原因需要退款时，
// 卖家可以通过退款接口将支付款退还给买家，微信支付将在收到退款请求并且验证成功之后，
// 按照退款规则将支付款按原路退到买家帐号上。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/ecommerce/refunds/chapter3_1.shtml
func (s *EcommerceRefundService) Apply(ctx context.Context, apply *RefundApply) (*RefundApplyResult, error) {
	api := "ecommerce/refunds/apply"

	req, err := s.client.NewRequest(ctx, "POST", api, apply)
	if err != nil {
		return nil, err
	}
	result := new(RefundApplyResult)

	_, err = s.client.Do(req, result)
	return result, err
}

// RefundQueryResult 查询退款结果
type RefundQueryResult struct {
	// 微信退款单号 必填
	// 微信支付退款订单号
	RefundID string `json:"refund_id"`

	// 商户退款单号 必填
	// 商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔
	OutRefundNo string `json:"out_refund_no"`

	// 微信订单号 必填
	// 微信支付交易订单号
	TransactionID string `json:"transaction_id"`

	// 商户订单号 必填
	// 返回的原交易订单号
	OutTradeNo string `json:"out_trade_no"`

	// 退款渠道 可选
	// ORIGINAL：原路退款
	// BALANCE：退回到余额
	// OTHER_BALANCE：原账户异常退到其他余额账户
	// OTHER_BANKCARD：原银行卡异常退到其他银行卡
	Channel string `json:"channel"`

	// 退款入账账户 可选
	// 取当前退款单的退款入账方。
	// 退回银行卡：{银行名称}{卡类型}{卡尾号}
	// 退回支付用户零钱：支付用户零钱
	// 退还商户：商户基本账户、商户结算银行账户
	// 退回支付用户零钱通：支付用户零钱通
	UserReceivedAccount string `json:"user_received_account"`

	// 退款成功时间 可选
	SuccessTime time.Time `json:"success_time"`

	// 退款创建时间 必填
	CreateTime time.Time `json:"create_time"`

	// 退款状态 必填
	// SUCCESS：退款成功
	// REFUNDCLOSE：退款关闭
	// PROCESSING：退款处理中
	// ABNORMAL：退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，
	//           可前往【服务商平台—>交易中心】，手动处理此笔退款
	Status string `json:"status"`

	// 订单金额信息 必填
	Amount RefundAmount `json:"amount"`

	// 优惠退款功能信息 可选
	PromotionDetail []RefundPromotionDetail `json:"promotion_detail"`
}

// QueryByRefundID 查询退款API - 通过微信支付退款单号查询退款
//
// 提交退款申请后，通过调用该接口查询退款状态。退款有一定延时，
// 用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/ecommerce/refunds/chapter3_2.shtml
func (s *EcommerceRefundService) QueryByRefundID(ctx context.Context, refundID, subMCHID string) (*RefundQueryResult, error) {
	return s.query(ctx, "id", refundID, subMCHID)
}

// QueryByOutRefundNo 查询退款API - 通过商户退款单号查询退款
//
// 提交退款申请后，通过调用该接口查询退款状态。退款有一定延时，
// 用零钱支付的退款20分钟内到账，银行卡支付的退款3个工作日后重新查询退款状态。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/ecommerce/refunds/chapter3_2.shtml
func (s *EcommerceRefundService) QueryByOutRefundNo(ctx context.Context, outRefundNo, subMCHID string) (*RefundQueryResult, error) {
	return s.query(ctx, "out-refund-no", outRefundNo, subMCHID)
}

func (s *EcommerceRefundService) query(ctx context.Context, key, value, subMCHID string) (*RefundQueryResult, error) {
	api := fmt.Sprintf("ecommerce/refunds/%s/%s?sub_mchid=%s", key, value, subMCHID)

	req, err := s.client.NewRequest(ctx, "POST", api, nil)
	if err != nil {
		return nil, err
	}
	result := new(RefundQueryResult)

	_, err = s.client.Do(req, result)
	return result, err
}
