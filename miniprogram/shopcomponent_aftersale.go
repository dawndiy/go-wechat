package miniprogram

import "context"

// ShopAftersaleInfo 创建售后请求
type ShopAftersaleInfo struct {
	// 商家自定义订单ID 必填
	OutOrderID string `json:"out_order_id"`
	// 商家自定义售后ID 必填
	OutAftersaleID string `json:"out_aftersale_id"`
	// 商家小程序该售后单的页面path，不存在则使用订单path
	Path string `json:"path"`
	// 退款金额,单位：分
	Refund int64 `json:"refund"`
	// 用户的openid
	OpenID string `json:"openid"`
	// 售后类型，1:退款,2:退款退货,3:换货 必填
	Type int `json:"type"`
	// 发起申请时间，yyyy-MM-dd HH:mm:ss 必填
	CreateTime string `json:"create_time"`
	// 0:未受理,
	// 1:用户取消,
	// 2:商家受理中,
	// 3:商家逾期未处理,
	// 4:商家拒绝退款,
	// 5:商家拒绝退货退款,
	// 6:待买家退货,
	// 7:退货退款关闭,
	// 8:待商家收货,
	// 11:商家退款中,
	// 12:商家逾期未退款,
	// 13:退款完成,
	// 14:退货退款完成,
	// 15:换货完成,
	// 16:待商家发货,
	// 17:待用户确认收货,
	// 18:商家拒绝换货,
	// 19:商家已收到货
	Status int `json:"status"`
	// 0:订单可继续售后, 1:订单无继续售后
	FinishAllAftersale int `json:"finish_all_aftersale"`
	// 退货相关商品列表
	ProductInfos []ShopAftersaleProductInfo `json:"product_infos"`
}

type ShopAftersaleProductInfo struct {
	// 商家自定义商品ID
	OutProductID string `json:"out_product_id"`
	// 商家自定义sku ID, 如果没有则不填
	OutSKUID string `json:"out_sku_id"`
	// 参与售后的商品数量
	ProductCnt int64 `json:"product_cnt"`
	// 退款原因
	RefundReason string `json:"refund_reason,omitempty"`
	// 买家收货地址
	RefundAddress string `json:"refund_address,omitempty"`
	// 退款金额
	Orderamt int64 `json:"orderamt,omitempty"`
}

// AftersaleAdd  创建售后
//
// 订单原始状态为10, 200, 250时会返回错误码100000
// finish_all_aftersale = 1时订单状态会流转到200（全部售后结束，不可继续售后）
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/add.html
func (s *ShopComponentShopService) AftersaleAdd(ctx context.Context, r ShopAftersaleInfo) error {
	u, err := s.client.apiURL(ctx, "shop/aftersale/add", nil)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), r)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// AftersaleGet 获取订单下售后单
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/get.html
func (s *ShopComponentShopService) AftersaleGet(ctx context.Context, orderID, outOrderID, openid string) (*ShopAftersaleInfo, error) {
	u, err := s.client.apiURL(ctx, "shop/aftersale/get", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]string{
		"openid": openid,
	}
	if orderID != "" {
		body["order_id"] = orderID
	} else {
		body["out_order_id"] = outOrderID
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	info := new(ShopAftersaleInfo)
	_, err = s.client.Do(req, info)
	return info, err
}

type ShopAftersaleUpdateRequest struct {
	OutOrderID     string `json:"out_order_id"`
	OpenID         string `json:"openid"`
	OutAftersaleID string `json:"out_aftersale_id"`
	// 0:未受理,
	// 1:用户取消,
	// 2:商家受理中,
	// 3:商家逾期未处理,
	// 4:商家拒绝退款,
	// 5:商家拒绝退货退款,
	// 6:待买家退货,
	// 7:退货退款关闭,
	// 8:待商家收货,
	// 11:商家退款中,
	// 12:商家逾期未退款,
	// 13:退款完成,
	// 14:退货退款完成,
	// 15:换货完成,
	// 16:待商家发货,
	// 17:待用户确认收货,
	// 18:商家拒绝换货,
	// 19:商家已收到货
	Status int `json:"status"`
	// 0:售后未结束,
	// 1:售后结束且订单状态流转
	FinishAllAftersale int `json:"finish_all_aftersale"`
}

// AftersaleUpdate 更新售后
//
// 只能更新售后状态
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/update.html
func (s *ShopComponentShopService) AftersaleUpdate(ctx context.Context, r ShopAftersaleUpdateRequest) error {
	u, err := s.client.apiURL(ctx, "shop/aftersale/get", nil)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), r)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
