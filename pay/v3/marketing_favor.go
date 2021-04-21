package pay

import (
	"context"
	"fmt"
	"time"
)

// MarketingFavorService 代金券服务
type MarketingFavorService service

// CouponStocksRequest 创建代金券批次请求
type CouponStocksRequest struct {
	// 批次名称 必填
	// 示例值：微信支付代金券批次
	StockName string `json:"stock_name"`

	// 批次备注 可选
	// 仅配置商户可见，用于自定义信息
	// 示例值：零售批次
	Comment string `json:"comment,omitempty"`

	// 归属商户号 必填
	// 示例值：98568865
	BelongMerchant string `json:"belong_merchant"`

	// 批次开始时间 必填
	AvailableBeginTime time.Time
	// 批次结束时间 必填
	AvailableEndTime time.Time

	// 发放规则 必填
	// 批次使用规则
	StockUseRule CouponStockUseRule `json:"stock_use_rule"`

	// 样式设置 可选
	// 代金券详情页
	PatternInfo *CouponPatternInfo `json:"pattern_info"`

	// 核销规则 必填
	CouponUseRule CouponUseRule `json:"coupon_use_rule"`

	// 是否无资金流 必填
	NoCash bool `json:"no_cash"`

	// 批次类型 必填
	// 批次类型，枚举值：
	// NORMAL：固定面额满减券批次
	StockType string `json:"stock_type"`

	// 商户单据号 必填
	// 商户创建批次凭据号（格式：商户id+日期+流水号），
	// 可包含英文字母，数字，|，_，*，-等内容，不允许出现其他不合法符号，
	// 商户侧需保持唯一性
	OutRequestNo string `json:"out_request_no"`

	// 扩展属性 可选
	// 扩展属性字段，按json格式，暂时无需填写。
	// 示例值：{'exinfo1':'1234','exinfo2':'3456'}
	ExtInfo interface{} `json:"ext_info,omitempty"`
}

// CouponStockUseRule 使用规则
type CouponStockUseRule struct {
	// 发放总上限 必填
	// 最大发券数
	// 示例值：100
	MaxCoupons uint64 `json:"max_coupons"`

	// 总预算 必填
	// 总消耗金额，单位：分。
	// max_amount需要等于coupon_amount（面额） * max_coupons（发放总上限）
	MaxAmount uint64 `json:"max_amount"`

	// 单天发放上限金额 可选
	// 单天最高消耗金额，单位：分。
	// 示例值：400
	MaxAmountByDay uint64 `json:"max_amount_by_day,omitempty"`

	// 单个用户可领个数 必填
	// 单个用户可领个数，每个用户最多60张券
	// 示例值：3
	MaxCouponsPerUser int `json:"max_coupons_per_user"`

	// 是否开启自然人限制 必填
	NaturalPersonLimit bool `json:"natural_person_limit"`

	// api发券防刷 必填
	PerventAPIAbuse bool `json:"pervent_api_abuse"`

	// 以下字段查询结果可能含有

	// 固定面额批次特定信息 可选
	// 请求无需设置，只有查询结果里可能含有
	FixedNormalCoupon *CouponFixedNormalRule `json:"fixed_normal_coupon"`

	// 券类型 可选
	// NORMAL：满减券
	// CUT_TO：减至券
	// 示例值：NORMAL
	CouponType string `json:"coupon_type,omitempty"`

	// 订单优惠标记 可选
	// 订单优惠标记
	// 特殊规则：单个优惠标记的字符长度为【1，128】,条目个数限制为【1，50】。
	// 示例值：{'123456','23456'}
	GoodsTag []string `json:"goods_tag"`

	// 支付方式 必填
	// MICROAPP：小程序支付
	// APPPAY：APP支付
	// PPAY：免密支付
	// CARD：刷卡支付
	// FACE：人脸支付
	// OTHER：其他支付
	TradeType []string `json:"trade_type"`

	// 是否可叠加其他优惠 可选
	CombineUse bool `json:"combine_use,omitempty"`
}

// CouponFixedNormalRule 固定面额批次特定信息
type CouponFixedNormalRule struct {
	// 面额 必填
	// 面额，单位：分
	CouponAmount uint64 `json:"coupon_amount"`
	// 门槛 必填
	// 使用券金额门槛，单位：分
	TransactionMinimum uint64 `json:"transaction_minimum"`
}

// CouponPatternInfo 代金券详情页样式设置
type CouponPatternInfo struct {
	// 使用说明 可选
	// 用于说明详细的活动规则，会展示在代金券详情页。
	// 示例值：微信支付营销代金券
	Description string `json:"description,omitempty"`
	// 商户logo 可选
	// 商户logo ，仅支持通过图片上传API接口获取的图片URL地址。
	// 1、商户logo大小需为120像素*120像素。
	// 2、支持JPG/JPEG/PNG格式，且图片小于1M。
	// 示例值：https://qpic.cn/xxx
	MerchantLogo string `json:"merchant_logo,omitempty"`
	// 商户名称 可选
	// 示例值：微信支付
	MerchantName string `json:"merchant_name,omitempty"`
	// 背景颜色 可选
	// 券的背景颜色，可设置10种颜色，色值请参考卡券背景颜色图。颜色取值为颜色图中的颜色名称
	// 示例值：Color020
	BackgroundColor string `json:"background_color,omitempty"`
	// 券详情图片 可选
	// 券详情图片， 850像素*350像素，且图片大小不超过2M，支持JPG/PNG格式，
	// 仅支持通过图片上传API接口获取的图片URL地址
	CouponImage string `json:"coupon_image,omitempty"`
}

// 支付方式
const (
	TradeTypeMicroAPP = "MICROAPP" // 小程序支付
	TradeTypeAPPPay   = "APPPAY"   // APP支付
	TradeTypePPay     = "PPAY"     // 免密支付
	TradeTypeCard     = "CARD"     // 刷卡支付
	TradeTypeFace     = "FACE"     // 人脸支付
	TradeTypeOther    = "OTHER"    // 其他支付
)

// CouponUseRule 核销规则
type CouponUseRule struct {
	// 券生效时间 可选
	// 需要指定领取后延时生效否填
	CouponAvailableTime *CouponAvailableTime `json:"coupon_available_time,omitempty"`

	// 固定面额满减券使用规则 可选
	// stock_type为NORMAL时必填
	FixedNormalCoupon *CouponFixedNormalRule `json:"fixed_normal_coupon,omitempty"`

	// 折扣券使用规则 可选
	// stock_type为DISCOUNT时必填
	DiscountCoupon *CouponDiscountRule `json:"discount_coupon,omitempty"`

	// 换购券使用规则 可选
	// stock_type为EXCHANGE时必填
	ExchangeCoupon *CouponExchangeRule `json:"exchange_coupon,omitempty"`

	// 订单优惠标记 可选
	GoodsTag []string `json:"goods_tag,omitempty"`

	// 指定支付方式 可选
	// 限定该批次核销的指定支付方式，如零钱、指定银行卡等，需填入支付方式编码，
	// 不在此列表中的银行卡，即暂不支持营销能力。
	// 特殊规则：条目个数限制为【1，1】。
	// 示例值：ICBC_CREDIT
	LimitPay string `json:"limit_pay,omitempty"`

	// 支付方式 可选
	// MICROAPP：小程序支付
	// APPPAY：APP支付
	// PPAY：免密支付
	// CARD：刷卡支付
	// FACE：人脸支付
	// OTHER：其他支付
	// 示例值：MICROAPP
	TradeType string `json:"trade_type,omitempty"`

	// 是否可叠加其他优惠 可选
	CombineUse bool `json:"combine_use,omitempty"`

	// 可核销商品编码 可选
	// 可核销商品编码，按json格式。
	// 特殊规则：单个商品编码的字符长度为【1，128】,条目个数限制为【1，50】。
	// 示例值：['123321','456654']
	AvailableItems []string `json:"available_items,omitempty"`

	// 不参与优惠商品编码 可选
	// 不参与优惠商品编码，按json格式。
	// 特殊规则：单个商品编码的字符长度为【1，128】,条目个数限制为【1，50】。
	// 示例值：['789987','56765']
	UnavailableItems []string `json:"unavailable_items,omitempty"`

	// 可核销商户号 必填
	// 可核销商户号，按json格式。
	// 特殊规则：单个商品号的字符长度为【1，20】,条目个数限制为【1，50】。
	// 示例值：['9856000','9856111']
	AvailableMerchants []string `json:"available_merchants"`
}

// CouponAvailableTime 券生效时间
type CouponAvailableTime struct {
	// 固定时间段可用 可选
	FixAvailableTime *CouponFixAvailableTime `json:"fix_available_time,omitempty"`
	// 领取第二天生效 可选
	SecondDayAvailable bool `json:"second_day_available,omitempty"`
	// 领取后有效时间 可选
	// 单位：分钟
	AvailableTimeAfterReceive int `json:"available_time_after_receive,omitempty"`
}

// CouponFixAvailableTime 固定时间段可用
type CouponFixAvailableTime struct {
	// 可用星期数 可选
	// 0代表周日，1代表周一，以此类推
	AvailableWeekDay []int `json:"available_week_day,omitempty"`
	// 当天开始时间 必填
	// 当天开始时间，单位：秒
	BeginTime int64 `json:"begin_time"`
	// 当天结束时间 可选
	// 当天结束时间，单位：秒，默认为23点59分59秒
	EndTime int64 `json:"end_time,omitempty"`
}

// CouponDiscountRule 折扣券使用规则
type CouponDiscountRule struct {
	// 最高折扣金额 必填
	// 单位：分
	DiscountAmountMax uint64 `json:"discount_amount_max"`
	// 折扣百分比 必填
	// 例如88-八八折
	// 示例值：88
	DiscountPercent int `json:"discount_percent"`
	// 门槛 可选
	// 使用券金额门槛，单位：分
	// 示例值：100
	TransactionMinimum uint64
}

// CouponExchangeRule 换购券使用规则
type CouponExchangeRule struct {
	// 可用优惠的商品最高单价 必填
	// 超过商品最高单价的商品不享受单品优惠，单位：分
	SinglePriceMax uint64 `json:"single_price_max"`
	// 单品换购价 必填
	// 不超过商品最高单价的商品可用换购价购买商品，单位：分
	ExchangePrice uint64 `json:"exchange_price"`
}

// CouponStocks 创建代金券批次, 返回批次号和创建时间
//
// 通过调用此接口可创建代金券批次，包括预充值&免充值类型
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_1.shtml
func (s *MarketingFavorService) CouponStocks(ctx context.Context, r *CouponStocksRequest) (string, time.Time, error) {
	api := "marketing/favor/coupon-stocks"

	req, err := s.client.NewRequest(ctx, "POST", api, r)
	if err != nil {
		return "", time.Time{}, err
	}

	var data struct {
		StockID    string    `json:"stock_id"`
		CreateTime time.Time `json:"create_time"`
	}
	_, err = s.client.Do(req, &data)
	return data.StockID, data.CreateTime, err
}

// StockStart 激活代金券批次
//
// 制券成功后，通过调用此接口激活批次，如果是预充值代金券，激活时会从商户账户余额中锁定本批次的营销资金
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_3.shtml
func (s *MarketingFavorService) StockStart(ctx context.Context, mchid, stockID string) (time.Time, string, error) {
	api := fmt.Sprintf("marketing/favor/stocks/%s/start", stockID)

	body := map[string]string{
		"stock_creator_mchid": mchid,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return time.Time{}, "", err
	}

	var data struct {
		StartTime time.Time `json:"start_time"`
		StockID   string    `json:"stock_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.StartTime, data.StockID, err
}

// UserCouponSent 发放代金券返回代金券id
//
// 商户平台/API完成制券后，可使用发放代金券接口发券。
// 通过调用此接口可发放指定批次给指定用户，发券场景可以是小程序、H5、APP等。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_2.shtml
func (s *MarketingFavorService) UserCouponSent(ctx context.Context, couponReq *CouponUserCouponsRequest) (string, error) {
	api := fmt.Sprintf("marketing/favor/users/%s/coupons", couponReq.OpenID)

	if couponReq.StockCreatorMCHID == "" {
		couponReq.StockCreatorMCHID = s.client.mchid
	}

	req, err := s.client.NewRequest(ctx, "POST", api, couponReq)
	if err != nil {
		return "", err
	}

	var data struct {
		CouponID string `json:"coupon_id"`
	}

	_, err = s.client.Do(req, &data)
	return data.CouponID, err
}

// StockPause 暂停代金券批次
//
// 通过此接口可暂停指定代金券批次。暂停后，该代金券批次暂停发放，用户无法通过任何渠道再领取该批次的券
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_13.shtml
func (s *MarketingFavorService) StockPause(ctx context.Context, mchid, stockID string) (time.Time, string, error) {
	api := fmt.Sprintf("marketing/favor/stocks/%s/pause", stockID)

	body := map[string]string{
		"stock_creator_mchid": mchid,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return time.Time{}, "", err
	}

	var data struct {
		StartTime time.Time `json:"start_time"`
		StockID   string    `json:"stock_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.StartTime, data.StockID, err
}

// StockRestart 重启代金券批次
//
// 通过此接口可重启指定代金券批次。重启后，该代金券批次可以再次发放
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_14.shtml
func (s *MarketingFavorService) StockRestart(ctx context.Context, mchid, stockID string) (time.Time, string, error) {
	api := fmt.Sprintf("marketing/favor/stocks/%s/restart", stockID)

	body := map[string]string{
		"stock_creator_mchid": mchid,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return time.Time{}, "", err
	}

	var data struct {
		StartTime time.Time `json:"start_time"`
		StockID   string    `json:"stock_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.StartTime, data.StockID, err
}

// CouponUserCouponsRequest 发放代金券请求
type CouponUserCouponsRequest struct {
	// 用户openid 必填
	// openid信息，用户在appid下的唯一标识
	// 示例值：2323dfsdf342342
	OpenID string `json:"-"`

	// 批次号 必填
	// 微信为每个批次分配的唯一id
	// 示例值：9856000
	StockID string `json:"stock_id"`

	// 商户单据号 必填
	// 商户此次发放凭据号（格式：商户id+日期+流水号），
	// 可包含英文字母，数字，|，_，*，-等内容，
	// 不允许出现其他不合法符号，商户侧需保持唯一性。
	// 示例值： 89560002019101000121
	OutRequestNo string `json:"out_request_no"`

	// 公众账号ID , 如不设置默认使用 client.APPID
	// 微信为发券方商户分配的公众账号ID，
	// 接口传入的所有appid应该为公众号的appid（在mp.weixin.qq.com申请的），
	// 不能为APP的appid（在open.weixin.qq.com申请的）
	// 示例值：wx233544546545989
	APPID string `json:"appid"`

	// 创建批次的商户号, 如不设置默认使用 client.mchid
	// 批次创建方商户号
	// 示例值：8956000
	StockCreatorMCHID string `json:"stock_creator_mchid"`

	// 指定面额发券，面额 可选
	// 指定面额发券场景，券面额，其他场景不需要填，单位：分
	// 示例值：100
	CouponValue uint64 `json:"coupon_value,omitempty"`

	// 指定面额发券，券门槛 可选
	// 定面额发券批次门槛，其他场景不需要，单位：分
	// 示例值：100
	CouponMinimum uint64 `json:"coupon_minimum,omitempty"`
}

// 代金券批次状态
const (
	StockStatusUnactivated = "unactivated" // 未激活
	StockStatusAudit       = "audit"       // 审核中
	StockStatusRunning     = "running"     // 运行中
	StockStatusStoped      = "stoped"      // 已停止
	StockStatusPaused      = "paused"      // 暂停发放
)

// CouponStocksOptions 代金券批次查询选项
type CouponStocksOptions struct {
	// 列表选项 必填
	ListOptions

	// 创建批次的商户号 必填
	// 示例值：9856888
	StockCreatorMchID string `url:"stock_creator_mchid"`

	// 起始创建时间 可选
	CreateStartTime time.Time `url:"create_start_time,omitempty"`
	// 终止创建时间 可选
	CreateEndTime time.Time `url:"create_end_time,omitempty"`

	// 批次状态 可选
	// 批次状态，枚举值：
	// StockStatusUnactivated unactivated：未激活
	// StockStatusAudit       audit：审核中
	// StockStatusRunning     running：运行中
	// StockStatusStoped      stoped：已停止
	// StockStatusPaused      paused：暂停发放
	Status string `url:"status,omitempty"`
}

// 批次类型
const (
	StockTypeNormal      = "NORMAL"       // 代金券批次
	StockTypeDiscountCut = "DISCOUNT_CUT" // 立减与折扣
	StockTypeOther       = "OTHER"        // 其他
)

// CouponStock 批次详情
type CouponStock struct {
	// 批次号 必填
	// 微信为每个代金券批次分配的唯一id。
	// 示例值：9836588
	StockID string `json:"stock_id"`

	// 创建批次的商户号 必填
	// 微信为创建方商户分配的商户号。
	// 示例值：123456
	StockCreatorMchID string `json:"stock_creator_mchid"`

	// 批次名称 必填
	// 示例值：微信支付批次
	StockName string `json:"stock_name"`

	// 批次状态 必填
	// StockStatusUnactivated unactivated：未激活
	// StockStatusAudit       audit：审核中
	// StockStatusRunning     running：运行中
	// StockStatusStoped      stoped：已停止
	// StockStatusPaused      paused：暂停发放
	Status string `json:"status"`

	// 创建时间 必填
	CreateTime time.Time `json:"create_time"`

	// 使用说明 必填
	// 批次描述信息
	// 示例值：微信支付营销
	Description string `json:"description"`

	// 满减券批次使用规则 可选
	// 普通发券批次特定信息
	StockUseRule *CouponStockUseRule `json:"stock_use_rule,omitempty"`

	// 可用开始时间 必填
	AvailableBeginTime time.Time `json:"available_begin_time"`
	// 可用结束时间 必填
	AvailableEndTime time.Time `json:"available_end_time"`

	// 已发券数量 必填
	DistributedCoupons int `json:"distributed_coupons"`

	// 是否无资金流 必填
	NoCash bool `json:"no_cash"`

	// 激活批次的时间 可选
	StartTime time.Time `json:"start_time,omitempty"`
	// 终止批次的时间 可选
	StopTime time.Time `json:"stop_time,omitempty"`

	// 减至批次特定信息 可选
	CutToMessage *CouponCutToMessage

	// 是否单品优惠 必填
	SingleItem bool `json:"singleitem"`

	// 批次类型 必填
	StockType string `json:"stock_type"`
}

// CouponCutToMessage 单品优惠特定信息
type CouponCutToMessage struct {
	// 可用优惠的商品最高单价 必填
	// 单位：分
	SinglePriceMax uint64 `json:"single_price_max"`
	// 减至后的优惠单价 必填
	// 单位：分
	CutToPrice int64 `json:"cut_to_price"`
}

// Stocks 条件查询批次列表
//
// 通过此接口可查询多个批次的信息，包括批次的配置信息以及批次概况数据。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_4.shtml
func (s *MarketingFavorService) Stocks(ctx context.Context, stocksReq *CouponStocksOptions) (*ListOptions, []CouponStock, error) {
	api := "marketing/favor/stocks"
	if stocksReq.StockCreatorMchID == "" {
		stocksReq.StockCreatorMchID = s.client.mchid
	}
	u, err := addOptions(api, stocksReq)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var data struct {
		*ListOptions `json:",inline"`
		Data         []CouponStock `json:"data"`
	}
	_, err = s.client.Do(req, &data)

	return data.ListOptions, data.Data, err
}

// Stock 查询批次详情
//
// 通过此接口可查询批次信息，包括批次的配置信息以及批次概况数据
// 如 stockCreateMCHID 则使用 client.MCHID
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_5.shtml
func (s *MarketingFavorService) Stock(ctx context.Context, stockID, stockCreateMCHID string) (*CouponStock, error) {
	api := fmt.Sprintf("marketing/favor/stocks/%s", stockID)
	if stockCreateMCHID == "" {
		stockCreateMCHID = s.client.mchid
	}
	api += "?stock_creator_mchid=" + stockCreateMCHID

	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}

	stock := new(CouponStock)
	_, err = s.client.Do(req, stock)
	return stock, err
}

// 代金券状态
const (
	CouponStatusSended  = "SENDED"  // 可用
	CouponStatusUsed    = "USED"    // 已实扣
	CouponStatusExpired = "EXPIRED" // 已过期
)

// 券类别
const (
	CouponTypeNormal = "NORMAL" // 满减券
	CouponTypeCutTo  = "CUT_TO" // 减至券
)

// Coupon 代金券详情
type Coupon struct {
	// 创建批次的商户号 必填
	StockCreatorMCHID string `json:"stock_creator_mchid"`

	// 批次号 必填
	StockID string `json:"stock_id"`

	// 代金券id 必填
	CouponID string `json:"coupon_id"`

	// 单品优惠特定信息 可选
	CutToMessage *CouponCutToMessage `json:"cut_to_message,omitempty"`

	// 代金券名称 必填
	CouponName string `json:"coupon_name"`

	// 代金券状态 必填
	// SENDED：可用
	// USED：已实扣
	// EXPIRED：已过期
	Status string `json:"status"`

	// 使用说明 必填
	// 代金券描述说明字段
	Description string `json:"description"`

	// 领券时间 必填
	CreateTime string `json:"create_time"`

	// 券类型 必填
	// NORMAL：满减券
	// CUT_TO：减至券
	CouponType string `json:"coupon_type"`

	// 是否无资金流 必填
	NoCash bool `json:"no_cash"`

	// 可用开始时间 必填
	AvailableBeginTime time.Time `json:"available_begin_time"`
	// 可用结束时间 必填
	AvailableEndTime time.Time `json:"available_end_time"`

	// 是否单品优惠 必填
	SingleItem bool `json:"singleitem"`

	// 满减券信息 可选
	NormalCouponInformation *CouponFixedNormalRule `json:"normal_coupon_information,omitempty"`

	// 已实扣代金券核销信息 可选
	ConsumeInformation *ConsumeInformation `json:"consume_information,omitempty"`
}

// UserCouponQuery 查询代金券详情
//
// 通过此接口可查询代金券信息，包括代金券的基础信息、状态。
// 如代金券已核销，会包括代金券核销的订单信息（订单号、单品信息等）
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_6.shtml
func (s *MarketingFavorService) UserCouponQuery(ctx context.Context, couponID, appID, openID string) (*Coupon, error) {
	api := fmt.Sprintf("marketing/favor/users/%s/coupons/%s?appid=%s", openID, couponID, appID)
	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}
	coupon := new(Coupon)
	_, err = s.client.Do(req, coupon)
	return coupon, err
}

// StockMerchants 查询代金券可用商户
//
// 通过调用此接口可查询批次的可用商户号，判断券是否在某商户号可用，来决定是否展示
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_7.shtml
func (s *MarketingFavorService) StockMerchants(ctx context.Context, stockCreateMCHID, stockID string, opt ListOptions) (*ListOptions, error) {
	api := fmt.Sprintf("marketing/favor/stocks/favor/stocks/%s/merchants", stockID)
	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
	}
	var data struct {
		// TODO
	}
	_, err = s.client.Do(req, &data)

	return nil, err
}

// Callbacks 设置消息通知地址
//
// 用于设置接收营销事件通知的URL，可接收营销相关的事件通知，包括核销、发放、退款等
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/convention/chapter3_12.shtml
func (s *MarketingFavorService) Callbacks(ctx context.Context, notifyURL string, sw bool) (time.Time, string, error) {
	api := "marketing/favor/callbacks"
	body := map[string]interface{}{
		"mchid":      s.client.mchid,
		"notify_url": notifyURL,
		"switch":     sw,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return time.Time{}, "", err
	}

	var data struct {
		UpdateTime time.Time `json:"update_time"`
		NotifyURL  string    `json:"notify_url"`
	}

	_, err = s.client.Do(req, &data)
	return data.UpdateTime, data.NotifyURL, err
}
