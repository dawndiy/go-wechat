package pay

import (
	"context"
	"fmt"
	"time"
)

// MarketingBusifavorService 商家券服务
type MarketingBusifavorService service

// 商家券批次类型
const (
	BusifavorStockTypeNormal   = "NORMAL"   // 固定面额满减券批次
	BusifavorStockTypeDiscount = "DISCOUNT" // 折扣券批次
	BusifavorStockTypeExchange = "EXCHANGE" // 换购券批次
)

// 商家券 code 模式
const (
	CouponCodeModeWechatpayMode  = "WECHATPAY_MODE"  // 系统分配券code。（固定22位纯数字）
	CouponCodeModeMerchantAPI    = "MERCHANT_API"    // 商户发放时接口指定券code。
	CouponCodeModeMerchantUpload = "MERCHANT_UPLOAD" // 商户上传自定义code，发券时系统随机选取上传的券code。
)

// BusifavorStockCreateRequest 创建商家券批次请求
type BusifavorStockCreateRequest struct {
	// 商家券批次名称 必填
	// 批次名称，字数上限为21个字节长度（中文按UTF8编码算字节数）。
	// 示例值：8月1日活动券
	StockName string `json:"stock_name"`

	// 批次归属商户号 必填
	// 批次归属于哪个商户。
	// 示例值：10000022
	BelongMerchant string `json:"belong_merchant"`

	// 批次备注 可选
	// 仅配置商户可见，用于自定义信息。
	// 示例值：活动使用
	Comment string `json:"comment,omitempty"`

	// 适用商品范围 可选
	// 用来描述批次在哪些商品可用，会显示在微信卡包中。
	// 示例值：xxx商品使用
	GoodsName string `json:"goods_name,omitempty"`

	// 批次类型 必填
	// NORMAL：固定面额满减券批次
	// DISCOUNT：折扣券批次
	// EXCHANGE：换购券批次
	StockType string `json:"stock_type"`

	// 核销规则 必填
	CouponUseRule BusifavorCouponUseRule `json:"coupon_use_rule"`
	// 发放规则 必填
	StockSendRule BusifavorStockSendRule `json:"stock_send_rule"`

	// 商户请求单号 必填
	// 商户创建批次凭据号（格式：商户id+日期+流水号），商户侧需保持唯一性。
	// 示例值：100002322019090134234sfdf
	OutRequestNo string `json:"out_request_no"`

	// 自定义入口 可选
	// 卡详情页面，可选择多种入口引导用户。
	CustomEntrance *BusifavorCustomEntrance `json:"custom_entrance,omitempty"`

	// 样式信息 可选
	DisplayPatternInfo *BusifavorDisplayPatternInfo `json:"display_pattern_info,omitempty"`

	// 券code模式 必填
	// 枚举值：
	// WECHATPAY_MODE：系统分配券code。（固定22位纯数字）
	// MERCHANT_API：商户发放时接口指定券code。
	// MERCHANT_UPLOAD：商户上传自定义code，发券时系统随机选取上传的券code。
	CouponCodeMode string `json:"coupon_code_mode"`

	// 事件通知配置 可选
	NotifyConfig *BusifavorNotifyConfig `json:"notify_config,omitempty"`
}

// 商家券核销方式
const (
	BusifavorUseMethodOffline      = "OFF_LINE"      // 线下滴码核销
	BusifavorUseMethodMiniPrograms = "MINI_PROGRAMS" // 线上小程序核销
	BusifavorUseMethodPaymentCode  = "PAYMENT_CODE"  // 微信支付付款码核销
	BusifavorUseMethodSelfConsume  = "SELF_CONSUME"  // 用户自助核销
)

// BusifavorCouponUseRule 商家券核销规则
type BusifavorCouponUseRule struct {
	// 券可核销时间 必填
	CouponAvailableTime BusifavorCouponAvailableTime `json:"coupon_available_time"`

	// 固定面额满减券使用规则 stock_type为NORMAL时必填
	FixedNormalCoupon *BusifavorFixedNormalCoupon `json:"fixed_normal_coupon,omitempty"`
	// 折扣券使用规则 stock_type为DISCOUNT时必填
	DiscountCoupon *BusifavorDiscountCoupon `json:"discount_coupon,omitempty"`
	// 换购券使用规则 stock_type为EXCHANGE时必填
	ExchangeCoupon *BusifavorExchangeCoupon `json:"exchange_coupon,omitempty"`

	// 核销方式 必填
	// 枚举值：
	// OFF_LINE：线下滴码核销，点击券“立即使用”跳转展示券二维码详情。
	// MINI_PROGRAMS：线上小程序核销，点击券“立即使用”跳转至配置的商家小程序（需要添加小程序appid和path）。
	// PAYMENT_CODE：微信支付付款码核销，点击券“立即使用”跳转至微信支付钱包付款码。
	// SELF_CONSUME：用户自助核销，点击券“立即使用”跳转至用户自助操作核销界面（当前暂不支持用户自助核销）。
	UseMethod string `json:"use_method"`

	// 小程序appid 可选
	// 核销方式为线上小程序核销才有效
	MiniProgramsAPPID string `json:"mini_programs_appid,omitempty"`
	// 小程序path 可选
	// 核销方式为线上小程序核销才有效
	MiniProgramsPath string `json:"mini_programs_path,omitempty"`
}

// BusifavorCouponAvailableTime 商家券可核销时间
type BusifavorCouponAvailableTime struct {

	// 批次开始时间 必填
	AvailableBeginTime time.Time `json:"available_begin_time"`

	// 批次结束时间 必填
	AvailableEndTime time.Time `json:"available_end_time"`

	// 生效后N天内有效 可选
	// 日期区间内，券生效后x天内有效。
	// 例如生效当天内有效填1，生效后2天内有效填2，以此类推……注意，
	// 用户在有效期开始前领取商家券，则从有效期第1天开始计算天数，
	// 用户在有效期内领取商家券，则从领取当天开始计算天数，
	// 无论用户何时领取商家券，商家券在活动有效期结束后均不可用。
	// 可配合wait_days_after_receive一同填写，也可单独填写。
	// 单独填写时，有效期内领券后立即生效，生效后x天内有效。
	// 示例值：3
	AvaliableDayAfterReceive int `json:"avaliable_day_after_receive,omitempty"`

	// 固定周期有效时间段
	// 可以设置多个星期下的多个可用时间段，比如每周二10点到18点，用户自定义字段。
	AvaliableWeek           *BusifavorCouponAvaliableWeek            `json:"avaliable_week,omitempty"`
	IrregularyAvaliableTime []BusifavorCouponIrregularyAvaliableTime `json:"irregulary_avaliable_time,omitempty"`

	// 领取后N天开始生效 可选
	// 日期区间内，用户领券后需等待x天开始生效。
	// 例如领券后当天开始生效则无需填写，领券后第2天开始生效填1，以此类推……
	// 用户在有效期开始前领取商家券，则从有效期第1天开始计算天数，
	// 用户在有效期内领取商家券，则从领取当天开始计算天数。
	// 无论用户何时领取商家券，商家券在活动有效期结束后均不可用。
	// 需配合available_day_after_receive一同填写，不可单独填写。
	// 示例值：7
	WaitDayAfterReceive int `json:"wait_day_after_receive,omitempty"`
}

// BusifavorCouponIrregularyAvaliableTime 无规律的有效时间段
type BusifavorCouponIrregularyAvaliableTime struct {
	// 开始时间
	BeginTime time.Time `json:"begin_time,omitempty"`
	// 结束时间
	EndTime time.Time `json:"end_time,omitempty"`
}

// BusifavorCouponAvaliableWeek 固定周期有效时间段
type BusifavorCouponAvaliableWeek struct {
	// 可用星期数
	// 0代表周日，1代表周一，以此类推
	// 当填写available_day_time时，week_day必填
	// 示例值：1, 2
	WeekDay []int `json:"week_day,omitempty"`
	// 当天可用时间段
	// 可以填写多个时间段，最多不超过2个。
	AvaliableDayTime []BusifavorBeginEndSecend `json:"avaliable_day_time,omitempty"`
}

// BusifavorBeginEndSecend 当天可用时间段
type BusifavorBeginEndSecend struct {
	// 当天可用开始时间，单位：秒，1代表当天0点0分1秒。
	// 示例值：3600
	BeginTime int64 `json:"begin_time"`
	// 当天可用结束时间，单位：秒，86399代表当天23点59分59秒。
	// 示例值：86399
	EndTime int64 `json:"end_time"`
}

// BusifavorFixedNormalCoupon 商家券固定面额满减券使用规则
type BusifavorFixedNormalCoupon struct {
	// 优惠金额 可选
	// 单位：分
	DiscountAmount int `json:"discount_amount,omitempty"`
	// 消费门槛 可选
	// 单位：分
	TransactionMinimum int `json:"transaction_minimum,omitempty"`
}

// BusifavorDiscountCoupon 商家券折扣券使用规则
type BusifavorDiscountCoupon struct {
	// 折扣比例 可选
	// 例如：88为八八折。
	DiscountPercent int `json:"discount_percent,omitempty"`
	// 消费门槛 可选
	// 单位：分
	TransactionMinimum int `json:"transaction_minimum,omitempty"`
}

// BusifavorExchangeCoupon 商家券换购券使用规则
type BusifavorExchangeCoupon struct {
	// 单品换购价 可选
	// 单位：分
	ExchangePrice int `json:"exchange_price,omitempty"`
	// 消费门槛 可选
	// 单位：分
	TransactionMinimum int `json:"transaction_minimum,omitempty"`
}

// BusifavorStockSendRule 商家券发放规则
type BusifavorStockSendRule struct {
	// 批次最大发放个数 必填
	MaxCoupons int `json:"max_coupons"`
	// 用户最大可领个数 必填
	// 用户可领个数，每个用户最多100张券 。
	MaxCouponsPerUser int `json:"max_coupons_per_user"`
	// 单天发放上限个数 可选
	// 单天发放上限个数（stock_type为DISCOUNT或EXCHANGE时可传入此字段控制单天发放上限）
	MaxCouponsByDay int `json:"max_coupons_by_day,omitempty"`
	// 是否开启自然人限制 可选
	// 不填默认否
	NaturalPersonLimit bool `json:"natural_person_limit,omitempty"`
	// 可疑账号拦截 可选
	// 不填默认否
	PreventAPIAbuse bool `json:"prevent_api_abuse,omitempty"`
	// 是否允许转赠 可选
	// 不填默认否
	Transferable bool `json:"transferable,omitempty"`
	// 是否允许分享链接 可选
	// 不填默认否
	Shareable bool `json:"shareable,omitempty"`
}

// 商家券自定义入口 Code 展示模式
const (
	BusifavorCodeDisplayModeNotShow = "NOT_SHOW" // 不展示code
	BusifavorCodeDisplayModeBarCode = "BARCODE"  // 一维码
	BusifavorCodeDisplayModeQRCode  = "QRCODE"   // 二维码
)

// BusifavorCustomEntrance 商家券自定义入口
type BusifavorCustomEntrance struct {
	// 小程序入口 可选
	MiniProgramsInfo *BusifavorCustomEntranceMiniInfo `json:"mini_programs_info,omitempty"`
	// 商户公众号appid 可选
	// 可配置商户公众号，从券详情可跳转至公众号，用户自定义字段。
	APPID string `json:"appid,omitempty"`
	// 营销馆id 可选
	// 填写微信支付营销馆的馆id，用户自定义字段。 营销馆需在商户平台 创建。
	HallID string `json:"hall_id,omitempty"`
	// 可用门店id 可选
	StoreID string `json:"store_id,omitempty"`
	// code展示模式
	// 枚举值：
	// NOT_SHOW：不展示code
	// BARCODE：一维码
	// QRCODE：二维码
	CodeDisplayMode string `json:"code_display_mode,omitempty"`
}

// BusifavorCustomEntranceMiniInfo 小程序入口信息
type BusifavorCustomEntranceMiniInfo struct {
	// 商家小程序appid  必填
	MiniProgramsAPPID string `json:"mini_programs_appid"`
	// 商家小程序path 必填
	MiniProgramsPath string `json:"mini_programs_path"`
	// 入口文案 必填
	EntranceWords string `json:"entrance_words"`
	// 引导文案 可选
	GuidingWords string `json:"guiding_words,omitempty"`
}

// 卡券背景颜色图
const (
	BusifavorBackgroundColor010 = "Color010" // #63B359
	BusifavorBackgroundColor020 = "Color020" // #2C9F67
	BusifavorBackgroundColor030 = "Color030" // #509Fc9
	BusifavorBackgroundColor040 = "Color040" // #5885CF
	BusifavorBackgroundColor050 = "Color050" // #9062C0
	BusifavorBackgroundColor060 = "Color060" // #D09A45
	BusifavorBackgroundColor070 = "Color070" // #E4B138
	BusifavorBackgroundColor080 = "Color080" // #EE903C
	BusifavorBackgroundColor090 = "Color090" // #DD6549
	BusifavorBackgroundColor100 = "Color100" // #CC463D
)

// BusifavorDisplayPatternInfo 商家券样式信息
type BusifavorDisplayPatternInfo struct {
	// 使用须知 可选
	// 用于说明详细的活动规则，会展示在代金券详情页。
	// 示例值：xxx门店可用
	Description string `json:"description,omitempty"`

	// 商户logo 可选
	// 商户logo的URL地址，仅支持通过《图片上传API》接口获取的图片URL地址。
	// 1、商户logo大小需为120像素*120像素。
	// 2、支持JPG/JPEG/PNG格式，且图片小于1M。
	// 示例值：https://qpic.cn/xxx
	MerchantLogoURL string `json:"merchant_logo_url,omitempty"`

	// 商户名称 可选
	// 商户名称，字数上限为16个字节长度（中文按UTF8编码算字节数）
	// 示例值：微信支付
	MerchantName string `json:"merchant_name,omitempty"`

	// 背景颜色 可选
	// 券的背景颜色，可设置10种颜色，色值请参考卡券背景颜色图。
	// 颜色取值为颜色图中的颜色名称。
	// 示例值：Color020
	BackgroundColor string `json:"background_color,omitempty"`

	// 券详情图片
	// 券详情图片，850像素*350像素，且图片大小不超过2M，支持JPG/PNG格式，
	// 仅支持通过《图片上传API》接口获取的图片URL地址。
	// 示例值：https://qpic.cn/xxx
	CouponImageURL string `json:"coupon_image_url,omitempty"`
}

// BusifavorNotifyConfig 商家券事件通知配置
type BusifavorNotifyConfig struct {
	// 事件通知appid
	// 用于回调通知时，计算返回操作用户的openid（诸如领券用户），支持小程序or公众号的APPID；
	// 如该字段不填写，则回调通知中涉及到用户身份信息的openid与unionid都将为空。
	// 示例值：wx23232232323
	NotifyAPPID string `json:"notify_appid,omitempty"`
}

// StockCreate 创建商家券
// 返回 批次号，创建时间
//
// 商户可以通过该接口创建商家券。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_1.shtml
func (s *MarketingBusifavorService) StockCreate(
	ctx context.Context, r *BusifavorStockCreateRequest) (string, time.Time, error) {

	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/stocks", r)
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

// 商家券批次状态
const (
	BusifavorStockStateUnaudit = "UNAUDIT" // 审核中
	BusifavorStockStateRunning = "RUNNING" // 运行中
	BusifavorStockStateStopped = "STOPED"  // 已停止
	BusifavorStockStatePaused  = "PAUSED"  // 暂停发放
)

// BusifavorStockDetail 商家券批次详情
type BusifavorStockDetail struct {
	// 批次号
	StockID string `json:"stock_id"`
	// 商家券批次名称
	StockName string `json:"stock_name"`
	// 批次归属商户号
	BelongMerchant string `json:"belong_merchant"`
	// 批次备注
	Comment string `json:"comment"`
	// 适用商品范围
	// 用来描述批次在哪些商品可用，会显示在微信卡包中。
	GoodsName string `json:"goods_name"`
	// 批次类型
	// NORMAL：固定面额满减券批次
	// DISCOUNT：折扣券批次
	// EXCHANGE：换购券批次
	StockType string `json:"stock_type"`
	// 核销规则
	CouponUseRule BusifavorCouponUseRule `json:"coupon_use_rule"`
	// 发放规则
	StockSendRule BusifavorStockSendRule `json:"stock_send_rule"`
	// 自定义入口
	CustomEntrance *BusifavorCustomEntrance `json:"custom_entrance"`
	// 样式信息
	DisplayPatternInfo *BusifavorDisplayPatternInfo `json:"display_pattern_info"`
	// 批次状态
	// UNAUDIT：审核中
	// RUNNING：运行中
	// STOPED：已停止
	// PAUSED：暂停发放
	StockState string `json:"stock_state"`
	// 券code模式
	// WECHATPAY_MODE：系统分配券code。
	// MERCHANT_API：商户发放时接口指定券code。
	// MERCHANT_UPLOAD：商户上传自定义code，发券时系统随机选取上传的券code。
	CouponCodeMode string `json:"coupon_code_mode"`

	// 券code数量
	CouponCodeCount *BusifavorCouponCodeCount `json:"coupon_code_count"`
	// 事件通知配置
	NotifyConfig *BusifavorNotifyConfig `json:"notify_config"`
	// 批次发放情况
	SendCountInfomation *BusifavorSendCountInfo `json:"send_count_infomation"`
}

// BusifavorCouponCodeCount 券code数量
type BusifavorCouponCodeCount struct {
	// 该批次总共已上传的code总数
	TotalCount int64 `json:"total_count"`
	// 该批次当前可用的code数
	AvaliableCount int64 `json:"avaliable_count"`
}

// BusifavorSendCountInfo 批次发放情况
type BusifavorSendCountInfo struct {
	// 已发放券张数
	// 批次已发放的券数量，满减、折扣、换购类型会返回该字段
	TotalSendNum int64 `json:"total_send_num"`
	// 已发放券金额
	// 批次已发放的预算金额，满减券类型会返回该字段
	TotalSendAmount int64 `json:"total_send_amount"`
	// 单天已发放券张数
	// 批次当天已发放的券数量，设置了单天发放上限的满减、折扣、换购类型返回该字段
	TodaySendNum int64 `json:"today_send_num"`
	// 单天已发放券金额
	// 批次当天已发放的预算金额，设置了当天发放上限的满减券类型返回该字段
	TodaySendAmount int64 `json:"today_send_amount"`
}

// StockDetail 查询商家券详情
//
// 商户可通过该接口查询已创建的商家券批次详情信息。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_2.shtml
func (s *MarketingBusifavorService) StockDetail(ctx context.Context, stockID string) (*BusifavorStockDetail, error) {
	req, err := s.client.NewRequest(ctx, "GET", fmt.Sprintf("marketing/busifavor/stocks/%s", stockID), nil)
	if err != nil {
	}
	detail := new(BusifavorStockDetail)
	_, err = s.client.Do(req, detail)
	return detail, err
}

// BusifavorCouponUseRequest 核销用户券请求
type BusifavorCouponUseRequest struct {
	// 券的唯一标识 必填
	CouponCode string `json:"coupon_code"`
	// 批次号 可选
	// 商户自定义code的批次，核销时必须填写批次号
	StockID string `json:"stock_id,omitempty"`
	// 支持传入与当前调用接口商户号有绑定关系的appid。
	// 支持小程序appid与公众号appid。核销接口返回的openid会在该传入appid下进行计算获得。
	APPID string `json:"appid"`
	// 请求核销时间
	UseTime time.Time `json:"use_time"`
	// 核销请求单据号
	// 每次核销请求的唯一标识，商户需保证唯一。
	UseRequestNo string `json:"use_request_no"`
	// 用户标识 可选
	// 用户的唯一标识，做安全校验使用，非必填
	OpenID string `json:"open_id,omitempty"`
}

// CouponsUse 核销用户券
// 返回 批次号，用户标识，系统核销券成功的时间
//
// 在用户满足优惠门槛后，服务商可通过该接口核销用户微信卡包中具体某一张商家券。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_3.shtml
func (s *MarketingBusifavorService) CouponsUse(
	ctx context.Context, r *BusifavorCouponUseRequest) (string, string, time.Time, error) {

	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/coupons/use", r)
	if err != nil {
		return "", "", time.Time{}, err
	}

	var data struct {
		StockID          string    `json:"stock_id"`
		OpenID           string    `json:"openid"`
		WeChatPayUseTime time.Time `json:"wechatpay_use_time"`
	}
	_, err = s.client.Do(req, &data)

	return data.StockID, data.OpenID, data.WeChatPayUseTime, err
}

// BusifavorUserCouponsOptions 根据过滤条件查询用户券
type BusifavorUserCouponsOptions struct {
	// 用户标识  必填
	OpenID string `url:"-"`
	// appid 必填
	// 支持传入与当前调用接口商户号有绑定关系的appid。
	// 支持小程序appid与公众号appid。
	APPID string `url:"appid"`
	// 批次号 可选
	StockID string `url:"stock_id,omitempty"`
	// 券状态 可选
	CouponState string `url:"coupon_state,omitempty"`
	// 创建批次的商户号 可选
	CreatorMerchant string `url:"creator_merchant,omitempty"`
	// 批次归属商户号 可选
	BelongMerchant string `url:"belong_merchant,omitempty"`
	// 批次发放商户号 可选
	SenderMerchant string `url:"sender_merchant,omitempty"`
	// 分页页码 可选
	Offset int `url:"offset,omitempty"`
	// 分页大小 可选
	Limit int `url:"limit,omitempty"`
}

// BusifavorUserCoupon 用户券
type BusifavorUserCoupon struct {
	// 批次归属商户号
	BelongMerchant string `json:"belong_merchant"`
	// 商家券批次名称
	StockName string `json:"stock_name"`
	// 批次备注
	Comment string `json:"comment"`
	// 适用商品范围
	GoodsName string `json:"goods_name"`
	// 批次类型
	StockType string `json:"stock_type"`
	// 是否允许转赠
	Transferable bool `json:"transferable"`
	// 是否允许分享领券链接
	Shareable bool `json:"shareable"`
	// 券状态
	CouponState string `json:"coupon_state"`
	// 样式信息
	DisplayPatternInfo *BusifavorDisplayPatternInfo `json:"display_pattern_info"`
	// 券核销规则
	CouponUseRule BusifavorCouponUseRule `json:"coupon_use_rule"`
	// 自定义入口
	CustomEntrance *BusifavorCustomEntrance `json:"custom_entrance"`
	// 券code
	CouponCode string `json:"coupon_code"`
	// 批次号
	StockID string `json:"stock_id"`
	// 券可使用开始时间
	AvaliableStartTime time.Time `json:"avaliable_start_time"`
	// 券过期时间
	ExpireTime time.Time `json:"expire_time"`
	// 券领券时间
	ReceiveTime time.Time `json:"receive_time"`
	// 发券请求单号
	SendRequestNo string `json:"send_request_no"`
	// 核销请求单号
	UseRequestNo string `json:"use_request_no"`
	// 券核销时间
	UseTime time.Time `json:"use_time"`
}

// UserCoupons 根据过滤条件查询用户券
//
// 商户自定义筛选条件（如创建商户号、归属商户号、发放商户号等），查询指定微信用户卡包中满足对应条件的所有商家券信息。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_4.shtml
func (s *MarketingBusifavorService) UserCoupons(
	ctx context.Context, q *BusifavorUserCouponsOptions) (*ListOptions, []BusifavorUserCoupon, error) {

	u, err := addOptions(fmt.Sprintf("marketing/busifavor/users/%s/coupons", q.OpenID), q)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	var data struct {
		Data        []BusifavorUserCoupon `json:"data"`
		ListOptions `json:",inline"`
	}
	_, err = s.client.Do(req, &data)
	return &data.ListOptions, data.Data, err
}

// UserCoupon 查询用户单张券详情
//
// 服务商可通过该接口查询微信用户卡包中某一张商家券的详情信息。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_5.shtml
func (s *MarketingBusifavorService) UserCoupon(
	ctx context.Context, couponCode, appid, openid string) (*BusifavorUserCoupon, error) {

	api := fmt.Sprintf(
		"marketing/busifavor/users/%s/coupons/%s/appids/%s",
		openid, couponCode, appid,
	)
	req, err := s.client.NewRequest(ctx, "POST", api, nil)
	if err != nil {
		return nil, err
	}
	userCoupon := new(BusifavorUserCoupon)
	_, err = s.client.Do(req, userCoupon)
	return userCoupon, err
}

// BusifavorCouponCodesUploadResult 上传预存 code 结果
type BusifavorCouponCodesUploadResult struct {
	// 批次号
	StockID string `json:"stock_id"`
	// 去重后上传code总数
	TotalCount int64 `json:"total_count"`
	// 上传成功code个数
	SuccessCount int64 `json:"success_count"`
	// 上传成功的code列表
	SuccessCodes []string `json:"success_codes"`
	// 上传成功时间
	SuccessTime time.Time `json:"success_time"`
	// 上传失败code个数
	FailCount int64 `json:"fail_count"`
	// 上传失败的code及原因
	FailCodes []struct {
		// 上传失败的券code
		// 商户通过API上传的券code
		CouponCode string `json:"coupon_code"`
		// 上传失败错误码
		Code string `json:"code"`
		// 上传失败错误信息
		Message string `json:"code"`
	} `json:"fail_codes"`
	// 已存在的code列表
	ExistCodes []string `json:"exist_codes"`
	// 本次请求中重复的code列表
	DuplicateCodes []string `json:"duplicate_codes"`
}

// StockCouponCodesUpload 上传预存code
//
// 商家券的Code码可由微信后台随机分配，同时支持商户自定义。
// 如商家已有自己的优惠券系统，可直接使用自定义模式。
// 即商家预先向微信支付上传券Code，当券在发放时，微信支付自动从已导入的Code中随机取值（不能指定），派发给用户。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_6.shtml
func (s *MarketingBusifavorService) StockCouponCodesUpload(
	ctx context.Context, stockID, uploadRequestNo string, codeList []string) (*BusifavorCouponCodesUploadResult, error) {

	api := fmt.Sprintf("marketing/busifavor/stocks/%s/couponcodes", stockID)
	body := map[string]interface{}{
		"coupon_code_list":  codeList,
		"upload_request_no": uploadRequestNo,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return nil, err
	}
	result := new(BusifavorCouponCodesUploadResult)
	_, err = s.client.Do(req, result)

	return result, err
}

// CallbacksSet 设置商家券事件通知地址
//
// 用于设置接收商家券相关事件通知的URL，可接收商家券相关的事件通知、包括发放通知等。
// 需要设置接收通知的URL，并在商户平台开通营销事件推送的能力，即可接收到相关通知。
//
// 注意：
// • 仅可以收到由商户自己创建的批次相关的通知
// • 需要设置apiv3秘钥，否则无法收到回调。
// • 如果需要领券回调中的参数openid。需要创券时候传入 notify_appid 参数。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_7.shtml
func (s *MarketingBusifavorService) CallbacksSet(ctx context.Context, notifyURL, mchid string) error {
	body := map[string]string{
		"notify_url": notifyURL,
	}
	if mchid != "" {
		body["mchid"] = mchid
	}
	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/callbacks", body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// CallbacksDetail 查询商家券事件通知地址
//
// 通过调用此接口可查询设置的通知URL。
//
// 注意：
// • 仅可以查询由请求商户号设置的商家券通知url
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_8.shtml
func (s *MarketingBusifavorService) CallbacksDetail(ctx context.Context, mchid string) (string, error) {
	api := "marketing/busifavor/callbacks?mchid=" + mchid
	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
		return "", err
	}
	var data struct {
		NotifyURL string `json:"notify_url"`
	}
	_, err = s.client.Do(req, &data)
	return data.NotifyURL, err
}

// CouponsAssociate 关联订单信息
//
// 将有效态（未核销）的商家券与订单信息关联，用于后续参与摇奖&返佣激励等操作的统计。
//
// 注意：
// • 仅对有关联订单需求的券进行该操作
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_9.shtml
func (s *MarketingBusifavorService) CouponsAssociate(
	ctx context.Context, stockID, couponCode, outTradeNo, outRequestNo string) (time.Time, error) {

	body := map[string]string{
		"stock_id":       stockID,
		"coupon_code":    couponCode,
		"out_trade_no":   outTradeNo,
		"out_request_no": outRequestNo,
	}
	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/coupons/associate", body)
	if err != nil {
		return time.Time{}, err
	}
	var data struct {
		WeChatPayAssociateTime time.Time `json:"wechatpay_associate_time"`
	}
	_, err = s.client.Do(req, &data)

	return data.WeChatPayAssociateTime, err
}

// CouponsDisassociate 取消关联订单信息
//
// 取消商家券与订单信息的关联关系
//
// 注意：
// • 建议取消前调用查询接口，查到当前关联的商户单号并确认后，再进行取消操作
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_10.shtml
func (s *MarketingBusifavorService) CouponsDisassociate(
	ctx context.Context, stockID, couponCode, outTradeNo, outRequestNo string) (time.Time, error) {

	body := map[string]string{
		"stock_id":       stockID,
		"coupon_code":    couponCode,
		"out_trade_no":   outTradeNo,
		"out_request_no": outRequestNo,
	}
	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/coupons/disassociate", body)
	if err != nil {
		return time.Time{}, err
	}
	var data struct {
		WeChatPayAssociateTime time.Time `json:"wechatpay_associate_time"`
	}
	_, err = s.client.Do(req, &data)

	return data.WeChatPayAssociateTime, err
}

// BusifavorStockBudgetUpdateRequest 修改批次预算请求
type BusifavorStockBudgetUpdateRequest struct {
	// 批次号 必填
	StockID string `json:"-"`
	// 目标批次最大发放个数 可选
	TargetMaxCoupons int `json:"target_max_coupons,omitempty"`
	// 目标单天发放上限个数 可选
	TargetMaxCouponsByDay int `json:"target_max_coupons_by_day,omitempty"`
	// 当前批次最大发放个数 可选
	// 当传入target_max_coupons大于0时，current_max_coupons必传
	CurrentMaxCoupons int `json:"current_max_coupons,omitempty"`
	// 当前单天发放上限个数 可选
	// 当传入target_max_coupons_by_day大于0时，current_max_coupons_by_day必填
	CurrentMaxCouponsByDay int `json:"current_max_coupons_by_day,omitempty"`
	// 修改预算请求单据号 必填
	ModifyBudgetRequestNo string `json:"modify_budget_request_no"`
}

// StockBudgetUpdate 修改批次预算
// 返回 批次当前最大发放个数, 当前单天发放上限个数(可选)
//
// 商户可以通过该接口修改批次单天发放上限数量或者批次最大发放数量
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_11.shtml
func (s *MarketingBusifavorService) StockBudgetUpdate(
	ctx context.Context, r *BusifavorStockBudgetUpdateRequest) (int, int, error) {

	api := fmt.Sprintf("marketing/busifavor/stocks/%s/budget", r.StockID)
	req, err := s.client.NewRequest(ctx, "PATCH", api, r)
	if err != nil {
		return 0, 0, err
	}
	var data struct {
		MaxCoupons      int `json:"max_coupons"`
		MaxCouponsByDay int `json:"max_coupons_by_day"`
	}
	_, err = s.client.Do(req, &data)

	return data.MaxCoupons, data.MaxCouponsByDay, err
}

// BusifavorStockUpdateRequest 商家券基本信息修改请求
type BusifavorStockUpdateRequest struct {
	// 批次号 必填
	StockID string `json:"-"`
	// 自定义入口 可选
	CustomEntrance *BusifavorCustomEntrance `json:"custom_entrance,omitempty"`
	// 商家券批次名称 可选
	StockName string `json:"stock_name,omitempty"`
	// 批次备注 可选
	Comment string `json:"comment,omitempty"`
	// 适用商品范围 可选
	// 用来描述批次在哪些商品可用，会显示在微信卡包中。
	GoodsName string `json:"goods_name,omitempty"`
	// 商户请求单号 必填
	// 商户修改批次凭据号（格式：商户id+日期+流水号），商户侧需保持唯一性
	OutRequestNo string `json:"out_request_no,omitempty"`
	// 样式信息 可选
	DisplayPatternInfo *BusifavorDisplayPatternInfo `json:"display_pattern_info,omitempty"`
	// 核销规则 可选
	CouponUseRule BusifavorCouponUseRule `json:"coupon_use_rule,omitempty"`
	// 发放规则 可选
	StockSendRule BusifavorStockSendRule `json:"stock_send_rule,omitempty"`
	// 事件通知配置 可选
	NotifyConfig *BusifavorNotifyConfig `json:"notify_config,omitempty"`
}

// StockUpdate 修改商家券基本信息
//
// 商户可以通过该接口修改商家券基本信息
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_12.shtml
func (s *MarketingBusifavorService) StockUpdate(ctx context.Context, r *BusifavorStockUpdateRequest) error {
	api := fmt.Sprintf("marketing/busifavor/stocks/%s", r.StockID)
	req, err := s.client.NewRequest(ctx, "PATCH", api, r)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}

// CouponsReturn 申请退券
//
// 商户可以通过该接口为已核销的券申请退券
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_13.shtml
func (s *MarketingBusifavorService) CouponsReturn(
	ctx context.Context, couponCode, stockID, returnRequestNo string) (time.Time, error) {

	body := map[string]string{
		"coupon_code":       couponCode,
		"stock_id":          stockID,
		"return_request_no": returnRequestNo,
	}
	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/coupons/return", body)
	if err != nil {
		return time.Time{}, err
	}

	var data struct {
		WeChatPayReturnTime time.Time `json:"wechatpay_return_time"`
	}
	_, err = s.client.Do(req, &data)

	return data.WeChatPayReturnTime, err
}

// ConponsDeactivate 使券失效
//
// 商户可以通过该接口将可用券进行失效处理，券被失效后无法再被核销
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/marketing/busifavor/chapter3_14.shtml
func (s *MarketingBusifavorService) ConponsDeactivate(
	ctx context.Context, couponCode, stockID, deactivateRequestNo, deactivateReason string) (time.Time, error) {

	body := map[string]string{
		"coupon_code":           couponCode,
		"stock_id":              stockID,
		"deactivate_request_no": deactivateRequestNo,
		"deactivate_reason":     deactivateReason,
	}
	req, err := s.client.NewRequest(ctx, "POST", "marketing/busifavor/coupons/deactivate", body)
	if err != nil {
		return time.Time{}, err
	}

	var data struct {
		WeChatPayDeactivateTime time.Time `json:"wechatpay_deactivate_time"`
	}
	_, err = s.client.Do(req, &data)

	return data.WeChatPayDeactivateTime, err
}
