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

// ShopECAftersaleAddRequest 生成售后单请求
type ShopECAftersaleAddRequest struct {
	// 商家自定义订单ID 必填
	OutOrderID string `json:"out_order_id,omitempty"`
	// 和out_order_id二选一
	OrderID int64 `json:"order_id,omitempty"`
	// 商家自定义售后ID 必填
	OutAftersaleID string `json:"out_aftersale_id"`
	// 用户的openid 必填
	OpenID string `json:"openid"`
	// 售后类型，1:退款,2:退款退货 必填
	Type int `json:"type"`
	// 退款金额，单位分
	Orderamt int64 `json:"orderamt"`
	// 退货相关商品列表 必填
	ProductInfo ShopECAftersaleProductInfo `json:"product_info"`
	// 退款原因 必填
	RefundReason string `json:"refund_reason"`
	// 退款原因类型 必填
	// INCORRECT_SELECTION = 1; // 拍错/多拍
	// NO_LONGER_WANT = 2; // 不想要了
	// NO_EXPRESS_INFO = 3; // 无快递信息
	// EMPTY_PACKAGE = 4; // 包裹为空
	// REJECT_RECEIVE_PACKAGE = 5; // 已拒签包裹
	// NOT_DELIVERED_TOO_LONG = 6; // 快递长时间未送达
	// NOT_MATCH_PRODUCT_DESC = 7; // 与商品描述不符
	// QUALITY_ISSUE = 8; // 质量问题
	// SEND_WRONG_GOODS = 9; // 卖家发错货
	// THREE_NO_PRODUCT = 10; // 三无产品
	// FAKE_PRODUCT = 11; // 假冒产品
	// OTHERS = 12; // 其它
	RefundReasonType int `json:"refund_reason_type"`
}

type ShopECAftersaleProductInfo struct {
	// 商家自定义商品ID
	OutProductID string `json:"out_product_id,omitempty"`
	// 微信侧商品ID，和out_product_id二选一
	ProductID int64 `json:"product_id,omitempty"`
	// 商家自定义sku ID, 如果没有则不填
	OutSKUID string `json:"out_sku_id,omitempty"`
	// 微信侧sku_id
	SKUID int64 `json:"sku_id,omitempty"`
	// 参与售后的商品数量 必填
	ProductCnt int64 `json:"product_cnt"`
}

// ECAfterSaleAdd  生成售后单
//
// 创建售后之前，请商家确保同步了默认退款地址
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/add_new.html
func (s *ShopComponentShopService) ECAfterSaleAdd(ctx context.Context, r ShopECAftersaleAddRequest) (int64, error) {
	u, err := s.client.apiURL(ctx, "shop/ecaftersale/add", nil)
	if err != nil {
		return 0, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), r)
	if err != nil {
		return 0, err
	}
	var data struct {
		AftersaleID int64 `json:"aftersale_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.AftersaleID, err

}

// ShopECAftersaleOrderInfo 售后单详情
type ShopECAftersaleOrderInfo struct {
	OutAftersaleID string                     `json:"out_aftersale_id"`
	AftersaleID    int64                      `json:"aftersale_id"`
	OutOrderID     string                     `json:"out_order_id"`
	OrderID        int64                      `json:"order_id"`
	ProductInfo    ShopECAftersaleProductInfo `json:"product_info"`
	Type           int                        `json:"type"`
	ReturnInfo     struct {
		OrderReturnTime int64  `json:"order_return_time"`
		DeliveryID      string `json:"delivery_id"`
		WaybillID       string `json:"waybill_id"`
		DeliveryName    string `json:"delivery_name"`
	} `json:"return_info"`
	Orderamt         int64  `json:"orderamt"`
	RefundReason     string `json:"refund_reason"`
	RefundReasonType int    `json:"refund_reason_type"`
	// AFTERSALESTATUS_INVALID = 0;
	// USER_CANCELD = 1; // 用户取消申请
	// MERCHANT_PROCESSING = 2; // 商家受理中
	// MERCHANT_REJECT_REFUND = 4; // 商家拒绝退款
	// MERCHANT_REJECT_RETURN = 5; // 商家拒绝退货退款
	// USER_WAIT_RETURN = 6; // 待买家退货
	// RETURN_CLOSED = 7; // 退货退款关闭
	// MERCHANT_WAIT_RECEIPT = 8; // 待商家收货
	// MERCHANT_OVERDUE_REFUND = 12; // 商家逾期未退款
	// MERCHANT_REFUND_SUCCESS = 13; // 退款完成
	// MERCHANT_RETURN_SUCCESS = 14; // 退货退款完成
	// PLATFORM_REFUNDING = 15; // 平台退款中
	// PLATFORM_REFUND_FAIL = 16; // 平台退款失败
	// USER_WAIT_CONFIRM = 17; // 待用户确认
	// MERCHANT_REFUND_RETRY_FAIL = 18; // 商家打款失败，客服关闭售后
	// MERCHANT_FAIL = 19; // 售后关闭
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// ECAfterSaleGet 获取售后单详情
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/get_new.html
func (s *ShopComponentShopService) ECAfterSaleGet(ctx context.Context, aftersaleID int64, outAftersaleID string) (*ShopECAftersaleOrderInfo, error) {
	u, err := s.client.apiURL(ctx, "shop/ecaftersale/get", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]interface{}{}
	if aftersaleID != 0 {
		body["aftersale_id"] = aftersaleID
	} else if outAftersaleID != "" {
		body["out_aftersale_id"] = outAftersaleID
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data struct {
		AftersalesOrder ShopECAftersaleOrderInfo `json:"after_sales_order"`
	}
	_, err = s.client.Do(req, &data)
	return &data.AftersalesOrder, err
}

// ECAfterSaleAcceptRefund 同意退款
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/aftersale/acceptrefund.html
func (s *ShopComponentShopService) ECAfterSaleAcceptRefund(ctx context.Context, aftersaleID int64, outAftersaleID string) error {
	u, err := s.client.apiURL(ctx, "shop/ecaftersale/acceptrefund", nil)
	if err != nil {
		return err
	}
	body := map[string]interface{}{}
	if aftersaleID != 0 {
		body["aftersale_id"] = aftersaleID
	} else if outAftersaleID != "" {
		body["out_aftersale_id"] = outAftersaleID
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
