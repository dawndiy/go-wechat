package token

// AccessToken 接口调用凭据
type AccessToken struct {
	// 获取到的凭证
	Value string `json:"access_token"`
	// 凭证有效时间，单位：秒。目前是7200秒之内的值。
	ExpiresIn int64 `json:"expires_in"`
}
