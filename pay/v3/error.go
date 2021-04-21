package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrorResponse 微信支付接口公共错误
//
// 文档: https://wechatpay-api.gitbook.io/wechatpay-api-v3/wei-xin-zhi-fu-api-v3-jie-kou-gui-fan#cuo-wu-xin-xi
type ErrorResponse struct {
	// 详细错误码
	Code string `json:"code"`

	// 错误描述，使用易理解的文字表示错误的原因。
	Message string `json:"message"`

	// 错误详情
	Detail struct {
		// 指示错误参数的位置。当错误参数位于请求body的JSON时，
		// 填写指向参数的JSON Pointer。
		// 当错误参数位于请求的url或者querystring时，填写参数的变量名。
		Field string `json:"field"`

		// 错误的值
		Value string `json:"value"`

		// 具体错误原因
		Issue string `json:"issue"`

		Location string `json:"location"`
	} `json:"detail"`

	Response *http.Response `json:"-"`
}

// Error 错误信息
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%s %v: <http %d> [%v] %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Code, r.Message)
}

// CheckResponse 检查接口调用是否成功
// 包含业务错误
//
// 状态码: https://wechatpay-api.gitbook.io/wechatpay-api-v3/wei-xin-zhi-fu-api-v3-jie-kou-gui-fan#http-zhuang-tai-ma
func CheckResponse(r *http.Response) error {

	switch r.StatusCode {
	case http.StatusOK, // 200
		//http.StatusAccepted,  // 202
		http.StatusNoContent: // 204
		return nil
	}

	errResp := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errResp)
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	return errResp
}

// IsErrorResponse 是否为 *ErrorResponse
func IsErrorResponse(err error) (*ErrorResponse, bool) {
	e, ok := err.(*ErrorResponse)
	return e, ok
}
