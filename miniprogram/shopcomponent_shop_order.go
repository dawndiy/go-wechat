package miniprogram

import "context"

// SceneCheck 检查场景值是否在支付校验范围内
//
// 微信后台会对符合支付校验范围内的场景值下的收银台进行支付（ticket/订单信息）校验
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/platform-capabilities/business-capabilities/ministore/minishopopencomponent2/API/order/check_scene.html
func (s *ShopComponentShopService) SceneCheck(ctx context.Context, scene int) (bool, error) {
	u, err := s.client.apiURL(ctx, "shop/scene/check", nil)
	if err != nil {
		return false, err
	}
	body := map[string]int{
		"scene": scene,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return false, err
	}
	var data struct {
		IsMatched int `json:"is_matched"`
	}
	_, err = s.client.Do(req, &data)
	return data.IsMatched == 1, err
}

type ShopOrderAddRequest struct {
	// 创建时间 必填
	CreateTime string `json:"create_time"`
	// 商家自定义订单ID 必填
	OutOrderID string `json:"out_order_id"`
	// 用户的openid 必填
	OpenID string `json:"openid"`
	// 商家小程序该订单的页面path，用于微信侧订单中心跳转 必填
	Path string `json:"path"`
	// 下单时小程序的场景值，可通getLaunchOptionsSync或onLaunch/onShow拿到
	Scene string `json:"scene"`
	// 订单详情
	OrderDetail struct {
		// 商品列表
		ProductInfos []ShopOrderProductInfo `json:"product_infos"`
		// 支付信息
		PayInfo ShopOrderPayInfo `json:"pay_info"`
	} `json:"order_detail"`
}

type ShopOrderProductInfo struct {
	// 商家自定义商品ID 必填
	OutProductID string `json:"out_product_id"`
	// 商家自定义商品skuID，可填空字符串（如果这个product_id下没有sku）
	OutSKUID string `json:"out_sku_id"`
	// 购买的数量 必填
	ProductCnt int64 `json:"product_cnt"`
	// 生成订单时商品的售卖价（单位：分），可以跟上传商品接口的价格不一致
	SalePrice int64 `json:"sale_price"`
	// 扣除优惠后单件sku的均摊价格（单位：分），如果没优惠则与sale_price一致
	RealPrice int64 `json:"real_price"`
	// 生成订单时商品的头图 必填
	HeadImg string `json:"head_img"`
	// 生成订单时商品的标题 必填
	Title string `json:"title"`
	// 绑定的小程序商品路径 必填
	Path string `json:"path"`
}

type ShopOrderPayInfo struct {
	// 支付方式，0，微信支付，1: 货到付款，2：商家会员储蓄卡（默认0）
	PayMethodType int
	// 预支付ID 必填
	PrepayID string
	// 预付款时间（拿到prepay_id的时间）
	PrepayTime string
}

type ShopOrderAddResult struct {
	OrderID          int64  `json:"order_id"`
	OutOrderID       string `json:"out_order_id"`
	Ticket           string `json:"ticket"`
	TicketExpireTime string `json:"ticket_expire_time"`
	FinalPrice       int64  `json:"final_price"`
}

func (s *ShopComponentShopService) OrderAdd(ctx context.Context, r ShopOrderAddRequest) (*ShopOrderAddResult, error) {
	u, err := s.client.apiURL(ctx, "shop/order/add", nil)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), r)
	if err != nil {
		return nil, err
	}
	result := new(ShopOrderAddResult)
	_, err = s.client.Do(req, result)
	return result, err
}

type ShopOrderPayRequest struct {
	// 订单ID 选填
	OrderID int64 `json:"order_id,omitempty"`
	// 商家自定义订单ID，与 order_id 二选一
	OutOrderID string `json:"out_order_id,omitempty"`
	// 用户的openid 必填
	OpenID string `json:"openid"`
	// 类型，默认1:支付成功,2:支付失败,3:用户取消,4:超时未支付;5:商家取消;10:其他原因取消
	ActionType int `json:"action_type"`
	// 其他具体原因 可选
	ActionRemark string `json:"action_remark,omitempty"`
	// 支付订单号，action_type=1且order/add时传的pay_method_type=0时必填
	TransactionID string `json:"transaction_id,omitempty"`
	// 支付完成时间，action_type=1时必填
	PayTime string `json:"pay_time,omitempty"`
}

func (s *ShopComponentShopService) OrderPay(ctx context.Context, r ShopOrderPayRequest) error {
	u, err := s.client.apiURL(ctx, "shop/order/pay", nil)
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
