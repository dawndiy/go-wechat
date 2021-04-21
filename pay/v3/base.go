package pay

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

// BaseService 基础服务
type BaseService service

// WeChatPayCert 微信支付平台证书信息
type WeChatPayCert struct {
	SerialNo          string    `json:"serial_no"`
	EffectiveTime     time.Time `json:"effective_time"`
	ExpireTime        time.Time `json:"expire_time"`
	EncryptCertficate struct {
		Algorithm      string `json:"algorithm"`
		Nonce          string `json:"nonce"`
		AssociatedData string `json:"associated_data"`
		Ciphertext     string `json:"ciphertext"`
	} `json:"encrypt_certificate"`

	Certficate *x509.Certificate `json:"-"`
}

// Certificates 获取平台证书列表
//
// 获取商户当前可用的平台证书列表。微信支付提供该接口，帮助商户后台系统实现平台证书的平滑更换
//
// 文档: https://wechatpay-api.gitbook.io/wechatpay-api-v3/jie-kou-wen-dang/ping-tai-zheng-shu
func (s *BaseService) Certificates(ctx context.Context) ([]WeChatPayCert, error) {
	api := "certificates"
	ctx = contextWithDisableSignCheck(ctx, true)
	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data []WeChatPayCert `json:"data"`
	}

	if _, err = s.client.Do(req, &data); err == nil {
		for i, v := range data.Data {
			b, e := DecodeCiphertext(
				v.EncryptCertficate.Algorithm, v.EncryptCertficate.Ciphertext,
				v.EncryptCertficate.Nonce, v.EncryptCertficate.AssociatedData, s.client.apiKey)

			if e != nil {
				return data.Data, e
			}

			certBlock, _ := pem.Decode(b)
			if certBlock == nil {
				return data.Data, fmt.Errorf("cert decode error")
			}
			cert, e := x509.ParseCertificate(certBlock.Bytes)
			if e != nil {
				return data.Data, e
			}
			data.Data[i].Certficate = cert
		}
	}

	return data.Data, err
}

type ctxKey string

const ctxKeyDisableSignCheck ctxKey = "disable_sign_check"

func contextWithDisableSignCheck(ctx context.Context, disableSignCheck bool) context.Context {
	return context.WithValue(ctx, ctxKeyDisableSignCheck, disableSignCheck)
}

func disableSignCheckFromContext(ctx context.Context) bool {
	v := ctx.Value(ctxKeyDisableSignCheck)
	if v == nil {
		return false
	}
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}

// ListOptions 列表选项
type ListOptions struct {
	// 分页页码
	// 页码从0开始，默认第0页
	Offset int `json:"offset" url:"offset"`
	// 分页大小
	Limit int `json:"limit" url:"limit"`
	// 总数量
	TotalCount int `json:"total_count" url:"-"`
}
