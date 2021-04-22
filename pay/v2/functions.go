package pay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// 签名方式
const (
	SignTypeMD5        = "MD5"
	SignTypeHMACSHA256 = "HMAC-SHA256"
)

// CalcSign 计算签名
//
// values 为完整的字段集合，如果有字段值为空将不参与计算，如果包含 sign 也不会参与计算
//
// ◆ 参数名ASCII码从小到大排序（字典序）；
// ◆ 如果参数的值为空不参与签名；
// ◆ 参数名区分大小写；
// ◆ 验证调用返回或微信主动通知签名时，传送的sign参数不参与签名，将生成的签名与该sign值作校验。
// ◆ 微信接口可能增加字段，验证签名时必须支持增加的扩展字段
func CalcSign(values url.Values, signType, signKey string) string {

	var keys []string
	for k := range values {
		if k == "sign" {
			continue
		}
		// 非空参数值 才签名
		if values.Get(k) != "" {
			keys = append(keys, k)
		}
	}

	var s string

	sort.Strings(keys)
	for i, k := range keys {
		kv := fmt.Sprintf("%s=%s", k, values.Get(k))
		if i == 0 {
			s += kv
			continue
		}
		s += "&" + kv
	}
	s += "&key=" + signKey

	var sign string
	switch signType {
	case SignTypeHMACSHA256:
		h := hmac.New(sha256.New, []byte(signKey))
		h.Write([]byte(s))
		b := h.Sum(nil)
		sign = strings.ToUpper(hex.EncodeToString(b))
	default:
		// 默认 MD5
		h := md5.New()
		h.Write([]byte(s))
		b := h.Sum(nil)
		sign = strings.ToUpper(hex.EncodeToString(b))
	}
	return sign
}

// CalcPaySign 计算小程序、H5需要的支付签名
func CalcPaySign(appid, timestamp, nonceStr, prepayID, signType, signKey string) string {
	v := url.Values{}
	v.Set("appId", appid)
	v.Set("timeStamp", timestamp)
	v.Set("nonceStr", nonceStr)
	v.Set("package", "prepay_id="+prepayID)
	v.Set("signType", signType)
	paySign := CalcSign(v, signType, signKey)
	return paySign
}
