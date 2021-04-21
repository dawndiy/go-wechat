package pay

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

// 加密方式
const (
	// AEAD_AES_256_GCM 加密方式
	AEAD_AES_256_GCM = "AEAD_AES_256_GCM"
)

// CalcSign 计算签名
//
// elems 需要签名的字符串列表，不需要添加 "\n"
// keyPEMBytes 私钥文件内容 pem
func CalcSign(elems []string, keyPEMBytes []byte) (string, error) {
	// 构造签名串
	raw := strings.Join(elems, "\n") + "\n"

	// SHA256 with RSA
	h := sha256.Sum256([]byte(raw))
	keyBlock, _ := pem.Decode(keyPEMBytes)
	if keyBlock == nil {
		return "", fmt.Errorf("key decode error")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", err
	}
	signBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}

	// Sign Base64
	sign := base64.StdEncoding.EncodeToString(signBytes)

	return sign, nil
}

// GetCertSerialNumber 获取证书序号
func GetCertSerialNumber(certPEMBytes []byte) (string, error) {
	certBlock, _ := pem.Decode(certPEMBytes)
	if certBlock == nil {
		return "", fmt.Errorf("cert decode error")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return "", err
	}
	serialNum := strings.ToUpper(hex.EncodeToString(cert.SerialNumber.Bytes()))
	return serialNum, nil
}

// DecodeCiphertext 证书和回调报文解密
//
// 文档: https://wechatpay-api.gitbook.io/wechatpay-api-v3/qian-ming-zhi-nan-1/zheng-shu-he-hui-tiao-bao-wen-jie-mi
func DecodeCiphertext(algorithm, ciphertext, nonce, associatedData, apiKey string) ([]byte, error) {

	switch algorithm {
	case AEAD_AES_256_GCM:
	default:
		return nil, fmt.Errorf("algorithm '%s' not supported", algorithm)
	}

	var (
		data   []byte
		block  cipher.Block
		aesgcm cipher.AEAD
		err    error
	)

	if data, err = base64.StdEncoding.DecodeString(ciphertext); err != nil {
		return nil, err
	}
	if block, err = aes.NewCipher([]byte(apiKey)); err != nil {
		return nil, err
	}

	if aesgcm, err = cipher.NewGCM(block); err != nil {
		return nil, err
	}

	var plaintextBytes []byte
	if plaintextBytes, err = aesgcm.Open(nil, []byte(nonce), data, []byte(associatedData)); err != nil {
		return nil, err
	}
	return plaintextBytes, nil
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func readResponse(r *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	return body, err
}
