package pay

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// 通知事件类型
const (
	EventTypeCouponUse  = "COUPON.USE"  // 代金券用券回调通知
	EventTypeCouponSend = "COUPON.SEND" // 领券事件回调通知

	EventTypeTransactionSuccess = "TRANSACTION.SUCCESS" // 支付成功通知

	EventTypeRefundSuccess  = "REFUND.SUCCESS"  // 退款成功通知
	EventTypeRefundAbnormal = "REFUND.ABNORMAL" // 退款异常通知
	EventTypeRefundClosed   = "REFUND.CLOSED"   // 退款关闭通知
)

// Notify 通知
type Notify struct {
	// 通知ID 必填
	ID string `json:"id"`

	// 通知创建时间 必填
	// 通知创建的时间，格式为yyyyMMddHHmmss。
	// 示例值：20180225112233
	CreateTime string `json:"create_time"`

	// 通知类型 必填
	EventType string `json:"event_type"`

	// 通知数据类型 必填
	ResourceType string `json:"resource_type"`

	// 回调摘要 必填
	Summary string `json:"summary"`

	// 通知数据 必填
	Resource NotifyResource `json:"resource"`
}

// NotifyResource 通知数据
type NotifyResource struct {
	// 加密算法类型 必填
	// 对开启结果数据进行加密的加密算法，目前只支持AEAD_AES_256_GCM。
	Algorithm string `json:"algorithm"`

	// 数据密文 必填
	// Base64编码后的开启/停用结果数据密文
	Ciphertext string `json:"ciphertext"`

	// 附加数据 可选
	AssociatedData string `json:"associated_data,omitempty"`

	// 随机串 必填
	// 加密使用的随机串
	Nonce string `json:"nonce"`

	// 原始回调类型 必填
	OriginalType string `json:"original_type"`
}

// ParseNotify 解析回调通知
//
// 调用方需要自己处理 *http.Request 的关闭操作
func (c *Client) ParseNotify(req *http.Request, v interface{}) (*Notify, error) {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	notify := new(Notify)
	if err = json.Unmarshal(body, notify); err != nil {
		return notify, err
	}

	if v != nil {
		rc := notify.Resource
		data, err := DecodeCiphertext(rc.Algorithm, rc.Ciphertext, rc.Nonce, rc.AssociatedData, c.apiKey)
		if err != nil {
			return notify, err
		}
		if err = json.Unmarshal(data, v); err != nil {
			return notify, err
		}
	}

	return notify, nil
}
