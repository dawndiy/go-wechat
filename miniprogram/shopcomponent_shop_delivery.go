package miniprogram

import "context"

// DeliveryCompany 快递公司
type DeliveryCompany struct {
	// 快递公司id
	DeliveryID string `json:"delivery_id"`
	// 快递公司名称
	DeliveryName string `json:"delivery_name"`
}

// DeliveryGetCompoanyList 获取快递公司列表
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/delivery/get_company_list.html
func (s *ShopComponentShopService) DeliveryGetCompoanyList(ctx context.Context) ([]DeliveryCompany, error) {
	u, err := s.client.apiURL(ctx, "shop/delivery/get_company_list", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), struct{}{})
	if err != nil {
		return nil, err
	}
	var data struct {
		CompanyList []DeliveryCompany `json:"company_list"`
	}
	_, err = s.client.Do(req, &data)
	return data.CompanyList, err
}

// DeliverySendRequest 交易组件订单发货请求
type DeliverySendRequest struct {
	// 订单ID
	OrderID int64 `json:"order_id,omitempty"`
	// 商家自定义订单ID，与 order_id 二选一
	OutOrderID string `json:"out_order_id,omitempty"`
	// 用户的openid
	OpenID string `json:"openid"`
	// 发货完成标志位, 0: 未发完, 1:已发完
	FinishAllDelivery int `json:"finish_all_delivery"`
	// 快递信息，delivery_type=1时必填
	DeliveryList []DeliverySendItem `json:"delivery_list,omitempty"`
	// 完成发货时间， finish_all_delivery = 1 必传
	ShipDoneTime string `json:"ship_done_time"`
}

// DeliverySendItem 交易组件订单发货快递信息
type DeliverySendItem struct {
	// 快递公司ID
	DeliveryID string `json:"delivery_id"`
	// 快递单号
	WaybillID string `json:"waybill_id"`
	// 物流单对应的商品信息
	ProductInfoList []DeliverySendProductInfo `json:"product_info_list"`
}

// DeliverySendProductInfo 物流单对应的商品信息
type DeliverySendProductInfo struct {
	OutProductID string `json:"out_product_id"`
	OutSKUID     string `json:"out_sku_id"`
	ProductCnt   int64  `json:"product_cnt"`
}

// DeliverySend 订单发货
//
// 新增订单时请指定delivery_type：1: 正常快递, 2: 无需快递, 3: 线下配送, 4: 用户自提。
// 对于没有物流的订单，可以不传delivery_list。
// delivery_id请参考获取快递公司列表接口的数据结果，不能自定义，如果对应不上，请传"delivery_id":"OTHERS"。
// 当finish_all_delivery=1时，把订单状态从20（待发货）流转到30（待收货）。
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/delivery/send.html
func (s *ShopComponentShopService) DeliverySend(ctx context.Context, r DeliverySendRequest) error {
	u, err := s.client.apiURL(ctx, "shop/delivery/send", nil)
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

// DeliveryRecieve 订单确认收货
//
// 把订单状态从30（待收货）流转到100（完成）
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/delivery/recieve.html
func (s *ShopComponentShopService) DeliveryRecieve(ctx context.Context, orderID int64, outOrderID, openid string) error {
	u, err := s.client.apiURL(ctx, "shop/delivery/recieve", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"openid": openid,
	}
	if orderID != 0 {
		body["order_id"] = orderID
	} else {
		body["out_order_id"] = outOrderID
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
