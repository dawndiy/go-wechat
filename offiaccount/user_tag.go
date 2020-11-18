package offiaccount

import (
	"context"
)

// TagCreate 创建标签
//
// 一个公众号，最多可以创建100个标签。
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagCreate(ctx context.Context, name string) (int, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/create", nil)
	if err != nil {
		return 0, err
	}
	body := map[string]interface{}{
		"tag": map[string]string{
			"name": name,
		},
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return 0, err
	}
	var data struct {
		Tag struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	_, err = s.client.Do(req, &data)
	return data.Tag.ID, err
}

// UserTag 用户标签
type UserTag struct {
	// 标签ID
	ID int `json:"id"`
	// 标签名称
	Name string `json:"name"`
	// 此标签下粉丝数
	Count int64 `json:"count"`
}

// TagGet 获取公众号已创建的标签
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagGet(ctx context.Context) ([]UserTag, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/get", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	var data struct {
		Tags []UserTag `json:"tags"`
	}
	_, err = s.client.Do(req, &data)
	return data.Tags, err
}

// TagUpdate 编辑标签
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagUpdate(ctx context.Context, id int, name string) error {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/update", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"tag": map[string]interface{}{
			"id":   id,
			"name": name,
		},
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// TagDelete 删除标签
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagDelete(ctx context.Context, id int) error {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/delete", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"tag": map[string]interface{}{
			"id": id,
		},
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// TagUsers 获取标签下粉丝列表
// 返回单次返回粉丝数，粉丝 openid 列表，下一个 openid
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagUsers(ctx context.Context, tagID int, nextOpenID string) (int64, []string, string, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/user/tag/get", nil)
	if err != nil {
		return 0, nil, "", err
	}
	body := map[string]interface{}{
		"tagid":       tagID,
		"next_openid": nextOpenID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return 0, nil, "", err
	}
	var data struct {
		Count int64 `json:"count"`
		Data  struct {
			OpenID []string `json:"openid"`
		} `json:"data"`
		NextOpenID string `json:"next_openid"`
	}
	_, err = s.client.Do(req, &data)
	return data.Count, data.Data.OpenID, data.NextOpenID, err
}

// TagBatchTagging 批量为用户打标签
//
// 标签功能目前支持公众号为用户打上最多20个标签。
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagBatchTagging(ctx context.Context, tagID int, openIDList []string) error {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/members/batchtagging", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"tagid":       tagID,
		"openid_list": openIDList,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// TagBatchUntagging 批量为用户取消标签
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagBatchUntagging(ctx context.Context, tagID int, openIDList []string) error {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/members/batchuntagging", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"tagid":       tagID,
		"openid_list": openIDList,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// TagUserTags 获取用户身上的标签列表
//
// 文档: https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (s *UserService) TagUserTags(ctx context.Context, openID string) ([]int, error) {
	u, err := s.client.apiURL(ctx, "cgi-bin/tags/getidlist", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]string{
		"openid": openID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data struct {
		TagIDList []int `json:"tagid_list"`
	}
	_, err = s.client.Do(req, &data)

	return data.TagIDList, err
}
