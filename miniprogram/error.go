package miniprogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrorResponse 响应错误结构
type ErrorResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`

	Response *http.Response `json:"-"`
}

// Error 实现 error 接口
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: <http %d> [%v] %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Code, r.Message)
}

// CheckResponse 检查响应是否有错误
func CheckResponse(r *http.Response) error {

	errResp := &ErrorResponse{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errResp)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if r.StatusCode == http.StatusOK && errResp.Code == 0 {
		return nil
	}

	return errResp
}
