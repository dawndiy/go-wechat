package pay

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/dawndiy/go-wechat/pkg/upload"
)

// MarketingMediaService 营销媒体上传服务
type MarketingMediaService service

// NewUploadRequest 特定的媒体上传请求
func (s *MarketingMediaService) NewUploadRequest(
	ctx context.Context, urlStr, filename string, reader io.Reader) (*http.Request, error) {

	fileBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write(fileBytes)
	fileSHA256 := hash.Sum(nil)
	meta := map[string]string{
		"filename": filename,
		"sha256":   hex.EncodeToString(fileSHA256),
	}
	metaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	u, err := s.client.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if s.client.userAgent != "" {
		req.Header.Set("User-Agent", s.client.userAgent)
	}
	req.Header.Set("Accept", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewReader(metaBytes)) // 设置 body 用于 signRequest 方法技术签名
	if err = s.client.signRequest(req); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	mpWriter := multipart.NewWriter(buf)
	metaPart := make(textproto.MIMEHeader)
	metaPart.Set("Content-Disposition", `form-data; name="meta";`)
	metaPart.Set("Content-Type", "application/json")
	w, err := mpWriter.CreatePart(metaPart)
	if err != nil {
		return nil, err
	}
	if _, err = w.Write(metaBytes); err != nil {
		return nil, err
	}

	filePart := make(textproto.MIMEHeader)
	filePart.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s";`, filename))
	filePart.Set("Content-Type", http.DetectContentType(fileBytes))
	w, err = mpWriter.CreatePart(filePart)
	if err != nil {
		return nil, err
	}
	if _, err = w.Write(fileBytes); err != nil {
		return nil, err
	}

	if err = mpWriter.Close(); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(buf)

	return req, nil
}

// ImageUpload 图片上传
//
// 通过本接口上传图片后可获得图片url地址。图片url可在微信支付营销相关的API使用，包括商家券、代金券、支付有礼等。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/chapter3_1.shtml
func (s *MarketingMediaService) ImageUpload(ctx context.Context, up upload.Uploader) (string, error) {
	req, err := s.NewUploadRequest(ctx, "marketing/favor/media/image-upload", up.Name(), up.Reader())
	if err != nil {
		return "", err
	}
	var data struct {
		MediaURL string `json:"media_url"`
	}
	_, err = s.client.Do(req, &data)
	return data.MediaURL, err
}
