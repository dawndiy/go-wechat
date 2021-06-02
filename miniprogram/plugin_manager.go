package miniprogram

import "context"

// PluginManagerService 插件管理服务
type PluginManagerService service

// PluginInfo 小程序插件
type PluginInfo struct {
	// 插件 appId
	APPID string `json:"appid"`
	// 插件状态
	// 1 申请中
	// 2 申请通过
	// 3 已拒绝
	// 4 已超时
	Status int `json:"status"`
	// 插件昵称
	Nickname string `json:"nickname"`
	// 插件头像
	HeadImgURL string `json:"headimgurl"`
}

// GetPluginList 查询已添加的插件
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/plugin-management/pluginManager.getPluginList.html
func (s *PluginManagerService) GetPluginList(ctx context.Context) ([]PluginInfo, error) {
	u, err := s.client.apiURL(ctx, "wxa/plugin", nil)
	if err != nil {
		return nil, err
	}

	body := map[string]string{
		"action": "list",
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}

	var data struct {
		// 申请或使用中的插件列表
		PluginList []PluginInfo `json:"plugin_list"`
	}
	_, err = s.client.Do(req, &data)
	return data.PluginList, nil
}

// ApplyPlugin 向插件开发者发起使用插件的申请
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/plugin-management/pluginManager.applyPlugin.html
func (s *PluginManagerService) ApplyPlugin(ctx context.Context, pluginAPPID, reason string) error {
	u, err := s.client.apiURL(ctx, "wxa/plugin", nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"action":       "apply",
		"plugin_appid": pluginAPPID,
		"reason":       reason,
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return nil
}

// UnbindPlugin 删除已添加的插件
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/plugin-management/pluginManager.unbindPlugin.html
func (s *PluginManagerService) UnbindPlugin(ctx context.Context, pluginAPPID string) error {
	u, err := s.client.apiURL(ctx, "wxa/plugin", nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"action":       "unbind",
		"plugin_appid": pluginAPPID,
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return nil
}
