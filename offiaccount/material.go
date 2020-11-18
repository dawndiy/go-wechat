package offiaccount

import (
	"context"
	"net/url"

	"github.com/dawndiy/go-wechat/pkg/upload"
)

// MaterialService 素材管理服务
type MaterialService service

// Upload 新增临时素材
//
// 公众号经常有需要用到一些临时性的多媒体素材的场景，
// 例如在使用接口特别是发送消息时，对多媒体文件、多媒体消息的获取和调用等操作，
// 是通过media_id来进行的。素材管理接口对所有认证的订阅号和服务号开放。
// 通过本接口，公众号可以新增临时素材（即上传临时多媒体文件）
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
func (s *MaterialService) UploadTempMedia(
	ctx context.Context, mediaType string, up upload.Uploader) (*upload.Result, error) {

	v := url.Values{}
	v.Set("type", mediaType)
	u, err := s.client.apiURL(ctx, "cgi-bin/media/upload", v)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewUploadRequest(ctx, u.String(), "media", up.Name(), up.Reader())
	if err != nil {
		return nil, err
	}

	result := new(upload.Result)
	_, err = s.client.Do(req, result)

	return result, err
}

// MaterialAddNewsArticle
type MaterialAddNewsArticle struct {
	Title              string `json:"title"`
	ThumbMediaID       string `json:"thumb_media_id"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	ShowCoverPic       int    `json:"show_cover_pic"`
	Content            string `json:"content"`
	ContentSourceURL   string `json:"content_source_url"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}

func (s *MaterialService) AddNews(ctx context.Context) {
}

// UploadImage 新增永久素材-上传图文消息内的图片获取URL
//
// 本接口所上传的图片不占用公众号的素材库中图片数量的100000个的限制。图片仅支持jpg/png格式，大小必须在1MB以下。
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
func (s *MaterialService) UploadImage(ctx context.Context, up upload.Uploader) (string, error) {

	u, err := s.client.apiURL(ctx, "cgi-bin/media/uploadimg", nil)
	if err != nil {
		return "", err
	}

	req, err := s.client.NewUploadRequest(ctx, u.String(), "media", up.Name(), up.Reader())
	if err != nil {
		return "", err
	}

	var data struct {
		URL string `json:"url"`
	}
	_, err = s.client.Do(req, &data)

	return data.URL, err
}
