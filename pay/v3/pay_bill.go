package pay

import "context"

// BillService 基础支付-账单
type BillService service

// TradeBill 申请交易账单
//
// 微信支付按天提供交易账单文件，商户可以通过该接口获取账单文件的下载地址。
// 文件内包含交易相关的金额、时间、营销等信息，供商户核对订单、退款、银行到账等情况
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/bill/chapter3_1.shtml
func (s *BillService) TradeBill(ctx context.Context) (string, string, string, error) {
	return "", "", "", nil
}

// FundFlowBill 申请资金账单
//
// 微信支付按天提供微信支付账户的资金流水账单文件，商户可以通过该接口获取账单文件的下载地址。
// 文件内包含该账户资金操作相关的业务单号、收支金额、记账时间等信息，供商户进行核对
//
// 文档: https://pay.weixin.qq.com/wiki/doc/apiv3/wxpay/pay/bill/chapter3_2.shtml
func (s *BillService) FundFlowBill(ctx context.Context) (string, string, string, error) {
	return "", "", "", nil
}

// DownloadBill 下载账单
func (s *BillService) DownloadBill(hashType, hashValue, downloadURL string) ([]byte, error) {
	return nil, nil
}
