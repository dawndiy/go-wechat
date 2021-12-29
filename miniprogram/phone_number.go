package miniprogram

import "context"

// PhoneNumberService 手机号服务
type PhoneNumberService service

// PhoneNumberResult 手机号获取结果
type PhoneNumberResult struct {
	// 用户绑定的手机号（国外手机号会有区号）
	PhoneNumber string `json:"phoneNumber"`
	// 没有区号的手机号
	PurePhoneNumber string `json:"purePhoneNumber"`
	// 区号
	CountryCode string `json:"countryCode"`
	// 数据水印
	WaterMark struct {
		// 小程序appid
		APPID string `json:"appid"`
		// 用户获取手机号操作的时间戳
		Timestamp int64 `json:"timestamp"`
	} `json:"watermark"`
}

// GetPhoneNumber code换取用户手机号
//
// 每个code只能使用一次，code的有效期为5min
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/phonenumber/phonenumber.getPhoneNumber.html
func (s *PhoneNumberService) GetPhoneNumber(ctx context.Context, code string) (*PhoneNumberResult, error) {
	u, err := s.client.apiURL(ctx, "wxa/business/getuserphonenumber", nil)
	if err != nil {
		return nil, err
	}
	body := map[string]string{
		"code": code,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	var data struct {
		PhoneInfo PhoneNumberResult `json:"phone_info"`
	}
	_, err = s.client.Do(req, &data)
	return &data.PhoneInfo, err
}
