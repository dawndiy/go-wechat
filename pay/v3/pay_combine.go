package pay

import (
	"context"
	"fmt"
	"time"
)

// 交易类型
const (
	TradeTypeNATIVE = "NATIVE" // 扫码支付
	TradeTypeJSAPI  = "JSAPI"  // 公众号支付
	TradeTypeAPP    = "APP"    // APP支付
	TradeTypeMWEB   = "MWEB"   // H5支付
)

// 交易状态
const (
	TradeStateSUCCESS    = "SUCCESS"    // 支付成功
	TradeStateREFUND     = "REFUND"     // 转入退款
	TradeStateNOTPAY     = "NOTPAY"     // 未支付
	TradeStateCLOSED     = "CLOSED"     // 已关闭
	TradeStateUSERPAYING = "USERPAYING" // 用户支付中
	TradeStatePAYERROR   = "PAYERROR"   // 支付失败(其他原因，如银行返回失败)
)

// PayCombineService 基础支付-合单支付
type PayCombineService service

// CombineOrder 支付订单
type CombineOrder struct {
	// 合单商户appid 必填
	// 合单发起方的appid
	// 示例值：wxd678efh567hg6787
	CombineAPPID string `json:"combine_appid"`

	// 合单商户号 必填
	// 合单发起方商户号
	// 示例值：1900000109
	CombineMCHID string `json:"combine_mchid"`

	// 合单商户订单号 必填
	// 合单支付总订单号，要求32个字符内，
	// 只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一
	// 示例值：P20150806125346
	CombineOutTradeNo string `json:"combine_out_trade_no"`

	// 场景信息 可选, H5必填
	// 支付场景信息描述
	SceneInfo *CombineSceneInfo `json:"scene_info,omitempty"`

	// 子单信息 必填
	// 最多支持子单条数：50
	SubOrders []CombineSubOrder `json:"sub_orders"`

	// 支付者信息 可选, JSAPI必填
	CombinePayerInfo *CombinePayerInfo `json:"combine_payer_info,omitempty"`

	// 交易起始时间 可选
	// 订单生成时间
	TimeStart time.Time `json:"time_start,omitempty"`

	// 交易结束时间 可选
	// 订单失效时间
	TimeExpire time.Time `json:"time_expire,omitempty"`

	// 通知地址 必填
	// 接收微信支付异步通知回调地址，
	// 通知url必须为直接可访问的URL，不能携带参数
	NotifyURL string `json:"notify_url"`

	// 指定支付方式 可选
	// no_credit：指定不能使用信用卡支付
	// 特殊规则：长度最大限制32个字节
	// 示例值：no_credit
	LimitPay []string `json:"limit_pay,omitempty"`
}

// CombineSceneInfo 合单支付-场景信息
type CombineSceneInfo struct {
	// 商户端设备号 可选, H5必填
	// 终端设备号（门店号或收银设备ID）。
	// 特殊规则：长度最小7个字节
	// 示例值：POS1:1
	DeviceID string `json:"device_id,omitempty"`

	// 用户终端IP 必填
	// 用户端实际ip
	// 格式: ip(ipv4+ipv6)
	// 示例值：14.17.22.32
	PayerClientIP string `json:"payer_client_ip"`

	// H5场景信息 可选
	H5Info *CombineSceneInfoH5 `json:"h5_info,omitempty"`
}

// CombineSceneInfoH5 H5场景信息
type CombineSceneInfoH5 struct {
	// 场景类型 必填
	// iOS：IOS移动应用；
	// Android：安卓移动应用；
	// Wap：WAP网站应用；
	// 示例值：iOS
	Type string `json:"type"`

	// 应用名称 可选
	// 示例值：王者荣耀
	AppName string `json:"app_name,omitempty"`

	// 网站URL 可选
	// 示例值：https://pay.qq.com
	AppURL string `json:"app_url,omitempty"`

	// iOS平台BundleID 可选
	// 示例值：com.tencent.wzryiOS
	BundleID string `json:"bundle_id,omitempty"`

	// Android平台PackageName 可选
	// 示例值：com.tencent.tmgp.sgame
	PackageName string `json:"package_name,omitempty"`
}

// CombineSubOrder 合单支付-子单信息
type CombineSubOrder struct {
	// 子单商户号 必填
	// 子单发起方商户号，必须与发起方appid有绑定关系。
	// 示例值：1900000109
	MCHID string `json:"mchid"`

	// 附加信息 必填
	// 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	// 示例值：深圳分店
	Attach string `json:"attach"`

	// 订单金额 必填
	Amount CombineOrderAmount `json:"amount"`

	// 子单商户订单号 必填
	// 商户系统内部订单号，要求32个字符内，
	// 只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	// 特殊规则：最小字符长度为6
	// 示例值：20150806125346
	OutTradeNo string `json:"out_trade_no"`

	// 二级商户号 必填
	// 二级商户商户号，由微信支付生成并下发。
	// 注意：仅适用于电商平台 服务商
	// 示例值：1900000109
	SubMCHID string `json:"sub_mchid"`

	// 商品描述 必填
	// 商品简单描述。需传入应用市场上的APP名字-实际商品名称，
	// 例如：天天爱消除-游戏充值。
	// 示例值：腾讯充值中心-QQ会员充值
	Description string `json:"description"`

	// 结算信息 可选
	SettleInfo *CombineOrderSettleInfo `json:"settle_info,omitempty"`

	/////////// 查询结果字段 ///////////

	// 交易类型 查询必填
	// NATIVE：扫码支付
	// JSAPI：公众号支付
	// APP：APP支付
	// MWEB：H5支付
	TradeType string `json:"trade_type,omitempty"`

	// 交易状态 查询必填
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// USERPAYING：用户支付中
	// PAYERROR：支付失败(其他原因，如银行返回失败)
	TradeState string `json:"trade_state,omitempty"`

	// 付款银行 查询可选
	// 银行类型，采用字符串类型的银行标识。
	// 示例值：CMC
	BankType string `json:"bank_type,omitempty"`

	// 支付完成时间 查询必填
	SuccessTime time.Time `json:"success_time,omitempty"`

	// 微信订单号 查询必填
	TransactionID string `json:"transaction_id,omitempty"`
}

// CombineOrderAmount 订单金额
type CombineOrderAmount struct {
	// 标价金额 必填
	// 子单金额，单位为分。
	// 示例值：100
	TotalAmount int64 `json:"total_amount"`
	// 标价币种 必填
	// 符合ISO 4217标准的三位字母代码，人民币：CNY。
	// 示例值：CNY
	Currency string `json:"currency"`

	// 现金支付金额 查询必填
	PayerAmount int64 `json:"payer_amount,omitempty"`
	// 现金支付币种 查询可选
	PayerCurrency int64 `json:"payer_currency,omitempty"`
}

// CombineOrderSettleInfo 结算信息
type CombineOrderSettleInfo struct {
	// 是否指定分账 可选
	// 是否分账，与外层profit_sharing同时存在时，以本字段为准
	ProfitSharing *bool `json:"profit_sharing,omitempty"`

	// 补差金额 可选
	// SettleInfo.profit_sharing为true时，该金额才生效
	SubsidyAmount int64 `json:"subsidy_amount"`
}

// CombinePayerInfo 支付者信息
type CombinePayerInfo struct {
	// 用户标识 必填
	// 使用合单appid获取的对应用户openid。
	// 是用户在商户appid下的唯一标识。
	OpenID string `json:"openid"`
}

// APP 合单下单-APP支付API
//
// 使用合单支付接口，用户只输入一次密码，即可完成多个订单的支付。目前最多一次可支持50笔订单进行合单支付。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_1.shtml
func (s *PayCombineService) APP(ctx context.Context, order *CombineOrder) (string, error) {
	api := "combine-transactions/app"

	req, err := s.client.NewRequest(ctx, "POST", api, order)
	if err != nil {
		return "", err
	}

	var data struct {
		PrepayID string `json:"prepay_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.PrepayID, err
}

// JSAPI 合单下单-JS支付API
//
// 使用合单支付接口，用户只输入一次密码，即可完成多个订单的支付。目前最多一次可支持50笔订单进行合单支付。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_2.shtml
func (s *PayCombineService) JSAPI(ctx context.Context, order *CombineOrder) (string, error) {
	api := "combine-transactions/jsapi"

	req, err := s.client.NewRequest(ctx, "POST", api, order)
	if err != nil {
		return "", err
	}

	var data struct {
		PrepayID string `json:"prepay_id"`
	}
	_, err = s.client.Do(req, &data)
	return data.PrepayID, err
}

// H5 合单下单-H5支付API
//
// 使用合单支付接口，用户只输入一次密码，即可完成多个订单的支付。目前最多一次可支持50笔订单进行合单支付。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_9.shtml
func (s *PayCombineService) H5(ctx context.Context, order *CombineOrder) (string, error) {
	api := "combine-transactions/h5"

	req, err := s.client.NewRequest(ctx, "POST", api, order)
	if err != nil {
		return "", err
	}

	var data struct {
		H5URL string `json:"h5_url"`
	}
	_, err = s.client.Do(req, &data)
	return data.H5URL, err
}

// Native 合单下单-Native支付API
//
// 使用合单支付接口，用户只输入一次密码，即可完成多个订单的支付。目前最多一次可支持50笔订单进行合单支付。
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_9.shtml
func (s *PayCombineService) Native(ctx context.Context, order *CombineOrder) (string, error) {
	api := "combine-transactions/h5"

	req, err := s.client.NewRequest(ctx, "POST", api, order)
	if err != nil {
		return "", err
	}

	var data struct {
		CodeURL string `json:"code_url"`
	}
	_, err = s.client.Do(req, &data)
	return data.CodeURL, err
}

// OutTradeNo 合单查询订单API
//
// 电商平台通过合单查询订单API查询订单状态，完成下一步的业务逻辑
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_3.shtml
func (s *PayCombineService) OutTradeNo(ctx context.Context, combineOutTradeNo string) (*CombineOrder, error) {
	api := fmt.Sprintf("combine-transactions/out-trade-no/%s", combineOutTradeNo)
	req, err := s.client.NewRequest(ctx, "GET", api, nil)
	if err != nil {
		return nil, err
	}
	order := new(CombineOrder)
	_, err = s.client.Do(req, order)
	return order, err
}

// CombineSubOrderClose 子单关闭
type CombineSubOrderClose struct {
	// 子单商户号 必填
	// 子单发起方商户号，必须与发起方appid有绑定关系。
	// 示例值：1900000109
	MCHID string `json:"mchid"`
	// 子单商户订单号 必填
	// 商户系统内部订单号，要求32个字符内，
	// 只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	// 特殊规则：最小字符长度为6
	// 示例值：20150806125346
	OutTradeNo string `json:"out_trade_no"`
	// 二级商户号 必填(电商平台 服务商)
	// 二级商户商户号，由微信支付生成并下发。
	// 注意：仅适用于电商平台 服务商
	// 示例值：1900000109
	SubMCHID string `json:"sub_mchid,omitempty"`
}

// OutTradeNoClose 合单关闭订单API
//
// 合单支付订单只能使用此合单关单api完成关单
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/combine/chapter3_4.shtml
func (s *PayCombineService) OutTradeNoClose(ctx context.Context, appid, combineOutTradeNo string, subOrders []CombineSubOrderClose) error {
	api := fmt.Sprintf("combine-transactions/out-trade-no/%s/close", combineOutTradeNo)
	body := map[string]interface{}{
		"combine_appid": appid,
		"sub_orders":    subOrders,
	}
	req, err := s.client.NewRequest(ctx, "POST", api, body)
	if err != nil {
		return err
	}
	_, err = s.client.Do(req, nil)
	return err
}
