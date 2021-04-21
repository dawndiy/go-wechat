package pay

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RefundNotify 退款通知
type RefundNotify struct {
	// 电商平台商户号 必填
	// 微信支付分配给电商平台的商户号
	SpMCHID string `json:"sp_mchid"`

	// 二级商户号 必填
	// 微信支付分配给二级商户的商户号
	SubMCHID string `json:"sub_mchid"`

	// 商户订单号 必填
	OutTradeNo string `json:"out_trade_no"`

	// 微信订单号 必填
	TransactionID string `json:"transaction_id"`

	// 商户退款单号 必填
	OutRefundNo string `json:"out_refund_no"`

	// 微信退款单号 必填
	RefundID string `json:"refund_id"`

	// 退款状态 必填
	// SUCCESS：退款成功
	// CLOSE：退款关闭
	// ABNORMAL：退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，
	//           可前往【服务商平台—>交易中心】，手动处理此笔退款
	RefundStatus string `json:"refund_status"`

	// 退款成功时间 可选
	SuccessTime time.Time `json:"success_time"`

	// 退款入账账户
	// 退回银行卡：{银行名称}{卡类型}{卡尾号}
	// 退回支付用户零钱: 支付用户零钱
	// 退还商户: 商户基本账户、商户结算银行账户
	// 退回支付用户零钱通：支付用户零钱通
	UserReceivedAccount string `json:"user_received_account"`

	// 金额信息 必填
	Amount RefundNotifyAmount `json:"amount"`
}

// RefundNotifyAmount 退款通知金额信息
type RefundNotifyAmount struct {
	// 订单金额 必填
	// 订单总金额，单位为分，只能为整数
	Total int64 `json:"total"`

	// 退款金额 必填
	// 退款金额，币种的最小单位，只能为整数，不能超过原订单支付金额，
	// 如果有使用券，后台会按比例退
	Refund int64 `json:"refund"`

	// 用户支付金额 必填
	// 用户实际支付金额，单位为分，只能为整数
	PayerTotal int64 `json:"payer_total"`

	// 用户退款金额 必填
	// 退款给用户的金额，不包含所有优惠券金额
	PayerRefund int64 `json:"payer_refund"`
}

// ParseNotify 解析退款结果通知
func (s *EcommerceRefundService) ParseNotify(req *http.Request) (*RefundNotify, error) {
	notify, err := s.client.ParseNotify(req, nil)
	if err != nil {
		return nil, err
	}

	switch notify.EventType {
	case EventTypeRefundSuccess,
		EventTypeRefundAbnormal,
		EventTypeRefundClosed:
	default:
		return nil, fmt.Errorf("通知类型错误")
	}

	rc := notify.Resource
	data, err := DecodeCiphertext(rc.Algorithm, rc.Ciphertext, rc.Nonce, rc.AssociatedData, s.client.apiKey)
	if err != nil {
		return nil, err
	}

	refundNotify := new(RefundNotify)
	err = json.Unmarshal(data, refundNotify)
	return refundNotify, err
}
