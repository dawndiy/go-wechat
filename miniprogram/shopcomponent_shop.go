package miniprogram

import (
	"context"
)

// ShopComponentShopService 自定义交易组件服务
type ShopComponentShopService service

// RegisterApply 接入申请
//
// 通过此接口开通自定义版交易组件，将同步返回接入结果，不再有异步事件回调。
// 如果账户已接入标准版组件，则无法开通，请到微信公众平台取消标准组件的开通。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/enter/enter_apply.html
func (s *ShopComponentShopService) RegisterApply(ctx context.Context) error {
	u, err := s.client.apiURL(ctx, "shop/register/apply", nil)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// ShopRegisterAccessInfo 自定义交易组件接入相关信息
type ShopRegisterAccessInfo struct {
	SPUAuditSuccess      int `json:"spu_audit_success"`
	PayOrderSuccess      int `json:"pay_order_success"`
	SendDeliverySuccess  int `json:"send_delivery_success"`
	AddAftersaleSuccess  int `json:"add_aftersale_success"`
	SPUAuditFinished     int `json:"spu_audit_finished"`
	PayOrderFinished     int `json:"pay_order_finished"`
	SendDeliveryFinished int `json:"send_delivery_finished"`
	AddAftersaleFinished int `json:"add_aftersale_finished"`
	TestAPIFinished      int `json:"test_api_finished"`
	DeployWXAFinished    int `json:"deploy_wxa_finished"`
}

// ShopRegisterSeneGroup 自定义交易组件场景接入相关
type ShopRegisterSeneGroup struct {
	GroupID           int    `json:"group_id"`
	Reason            string `json:"reason"`
	Name              string `json:"name"`
	Status            int    `json:"status"`
	SceneGroupExtList []struct {
		ExtID  int    `json:"ext_id"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	} `json:"scene_group_ext_list"`
}

// RegisterCheck 获取接入状态
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/enter/enter_check.html
func (s *ShopComponentShopService) RegisterCheck(ctx context.Context) (int, *ShopRegisterAccessInfo, []ShopRegisterSeneGroup, error) {

	u, err := s.client.apiURL(ctx, "shop/register/check", nil)
	if err != nil {
		return 0, nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return 0, nil, nil, err
	}

	var data struct {
		Data struct {
			Status         int                     `json:"status"`
			AccessInfo     ShopRegisterAccessInfo  `json:"access_info"`
			SceneGroupList []ShopRegisterSeneGroup `json:"scene_group_list"`
		} `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data.Status, &data.Data.AccessInfo, data.Data.SceneGroupList, err
}

// RegisterFinish 完成接入任务
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/enter/finish_access_info.html
func (s *ShopComponentShopService) RegisterFinish(ctx context.Context, accessInfoItem int) error {
	u, err := s.client.apiURL(ctx, "shop/register/finish_access_info", nil)
	if err != nil {
		return err
	}
	body := map[string]int{
		"access_info_item": accessInfoItem,
	}

	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// RegisterApplyScene 场景接入申请
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/framework/ministore/minishopopencomponent2/API/enter/scene_apply.html
func (s *ShopComponentShopService) RegisterApplyScene(ctx context.Context, sceneGroupID int) error {
	u, err := s.client.apiURL(ctx, "shop/register/apply_scene", nil)
	if err != nil {
		return err
	}
	body := map[string]int{
		"scene_group_id": sceneGroupID,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
