package miniprogram

import (
	"context"
)

// AnalysisService 数据分析服务
type AnalysisService service

// VisitRetain 访问留存信息
type VisitRetain struct {
	RefDate    string    `json:"ref_date"`     // 日期。格式为 yyyymmdd
	VisitUVNew []VisitUV `json:"visit_uv_new"` // 新增用户留存
	VisitUV    []VisitUV `json:"visit_uv"`     // 活跃用户留存
}

// VisitUV 访问用户数据
type VisitUV struct {
	// 标识
	// 0开始，表示当天，1表示1天后。依此类推，key取值分别是：0,1,2,3,4,5,6,7,14,30
	// 0开始，表示当月，1表示1月后。key取值分别是：0,1
	// 0开始，表示当周，1表示1周后。依此类推，取值分别是：0,1,2,3,4
	Key int64 `json:"key"`
	// key对应日期的新增用户数/活跃用户数（key=0时）或留存用户数（k>0时）
	Value int64 `json:"value"`
}

// GetDailyRetain 获取用户访问小程序日留存
//
// beginDate 开始日期。endDate 结束日期，限定查询1天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-retain/analysis.getDailyRetain.html
func (s *AnalysisService) GetDailyRetain(ctx context.Context, beginDate, endDate string) (*VisitRetain, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappiddailyretaininfo", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	visitRetain := new(VisitRetain)
	_, err = s.client.Do(req, visitRetain)
	return visitRetain, err
}

// GetMonthlyRetain 获取用户访问小程序月留存
//
// beginDate 开始日期，为自然月第一天。endDate 结束日期，为自然月最后一天，限定查询一个月数据。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-retain/analysis.getMonthlyRetain.html
func (s *AnalysisService) GetMonthlyRetain(ctx context.Context, beginDate, endDate string) (*VisitRetain, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidmonthlyretaininfo", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	visitRetain := new(VisitRetain)
	_, err = s.client.Do(req, visitRetain)
	return visitRetain, err
}

// GetWeeklyRetain 获取用户访问小程序周留存
//
// beginDate 开始日期，为周一日期。endDate 结束日期，为周日日期，限定查询一周数据。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-retain/analysis.getWeeklyRetain.html
func (s *AnalysisService) GetWeeklyRetain(ctx context.Context, beginDate, endDate string) (*VisitRetain, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidweeklyretaininfo", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	visitRetain := new(VisitRetain)
	_, err = s.client.Do(req, visitRetain)
	return visitRetain, err
}

// DailySummary 用户访问小程序数据概况
type DailySummary struct {
	RefDate    string `json:"ref_date"`    // 日期。格式为 yyyymmdd
	VisitTotal int64  `json:"visit_total"` // 累计用户数
	SharePV    int64  `json:"share_pv"`    // 转发次数
	ShareUV    int64  `json:"share_uv"`    // 转发人数
}

// GetDailySummary 获取用户访问小程序数据概况
//
// beginDate 开始日期。endDate 结束日期，限定查询1天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/analysis.getDailySummary.html
func (s *AnalysisService) GetDailySummary(ctx context.Context, beginDate, endDate string) ([]*DailySummary, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappiddailysummarytrend", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	var data struct {
		List []*DailySummary `json:"list"`
	}

	_, err = s.client.Do(req, &data)
	return data.List, err
}

// VisitTrend 访问趋势数据
type VisitTrend struct {
	RefDate         string `json:"ref_date"`          // 日期。格式为 yyyymmdd
	SessionCnt      int64  `json:"session_cnt"`       // 打开次数
	VisitPV         int64  `json:"visit_pv"`          // 访问次数
	VisitUV         int64  `json:"visit_uv"`          // 访问人数
	VisitUVNew      int64  `json:"visit_uv_new"`      // 新用户数
	StayTimeUV      int64  `json:"stay_time_uv"`      // 人均停留时长 (浮点型，单位：秒)
	StayTimeSession int64  `json:"stay_time_session"` // 次均停留时长 (浮点型，单位：秒)
	VisitDepth      int64  `json:"visit_depth"`       // 平均访问深度 (浮点型)
}

// GetDailyVisitTrend 获取用户访问小程序数据日趋势
//
// beginDate 开始日期。endDate 结束日期，限定查询1天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-trend/analysis.getDailyVisitTrend.html
func (s *AnalysisService) GetDailyVisitTrend(ctx context.Context, beginDate, endDate string) ([]*VisitTrend, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappiddailyvisittrend", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	var data struct {
		List []*VisitTrend `json:"list"`
	}

	_, err = s.client.Do(req, &data)
	return data.List, err
}

// GetMonthlyVisitTrend 获取用户访问小程序数据月趋势
//
// beginDate 开始日期，为自然月第一天。endDate 结束日期，为自然月最后一天，限定查询一个月的数据。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-trend/analysis.getMonthlyVisitTrend.html
func (s *AnalysisService) GetMonthlyVisitTrend(ctx context.Context, beginDate, endDate string) ([]*VisitTrend, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidmonthlyvisittrend", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	var data struct {
		List []*VisitTrend `json:"list"`
	}

	_, err = s.client.Do(req, &data)
	return data.List, err
}

// GetWeeklyVisitTrend 获取用户访问小程序数据周趋势
//
// beginDate 开始日期，为周一日期。endDate 结束日期，为周日日期，限定查询一周数据。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/visit-trend/analysis.getWeeklyVisitTrend.html
func (s *AnalysisService) GetWeeklyVisitTrend(ctx context.Context, beginDate, endDate string) ([]*VisitTrend, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidweeklyvisittrend", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	var data struct {
		List []*VisitTrend `json:"list"`
	}

	_, err = s.client.Do(req, &data)
	return data.List, err
}

// UserPortrait 用户画像分布信息
type UserPortrait struct {
	RefDate    string   `json:"ref_date"`     // 时间范围，如："20170611-20170617"
	VisitUVNew Portrait `json:"visit_uv_new"` // 新用户画像
	VisitUV    Portrait `json:"visit_uv"`     // 活跃用户画像
}

// Portrait 画像数据
type Portrait struct {
	Index     int            `json:"index"`     // 分布类型
	Province  []PortraitAttr `json:"province"`  // 省份，如北京、广东等
	City      []PortraitAttr `json:"city"`      // 城市，如北京、广州等
	Genders   []PortraitAttr `json:"genders"`   // 性别，包括男、女、未知
	Platforms []PortraitAttr `json:"platforms"` // 终端类型，包括 iPhone，android，其他
	Devices   []PortraitAttr `json:"devices"`   // 机型，如苹果 iPhone 6，OPPO R9 等
	Ages      []PortraitAttr `json:"ages"`      // 年龄，包括17岁以下、18-24岁等区间
}

// PortraitAttr 画像属性
type PortraitAttr struct {
	ID                  int   `json:"id"`                     // 属性值id
	Name                int   `json:"name"`                   // 属性值名称，与id对应。如属性为 province 时，返回的属性值名称包括「广东」等。
	AccessSourceVisitUV int64 `json:"access_source_visit_uv"` // 该场景访问uv
}

// GetUserPortrait 获取小程序新增或活跃用户的画像分布数据。
// 时间范围支持昨天、最近7天、最近30天。其中，新增用户数为时间范围内首次访问小程序的去重用户数，
// 活跃用户数为时间范围内访问过小程序的去重用户数。
//
// beginDate 开始日期。endDate 结束日期，开始日期与结束日期相差的天数限定为0/6/29，
// 分别表示查询最近1/7/30天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/analysis.getUserPortrait.html
func (s *AnalysisService) GetUserPortrait(ctx context.Context, beginDate, endDate string) (*UserPortrait, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappiduserportrait", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	userPortrait := new(UserPortrait)
	_, err = s.client.Do(req, userPortrait)
	return userPortrait, err
}

// VisitDistribution 访问分布
type VisitDistribution struct {
	RefDate string                  `json:"ref_date"` // 日期。格式为 yyyymmdd
	List    []VisitDistributionData `json:"list"`     // 数据列表
}

// VisitDistributionData 分布数据
type VisitDistributionData struct {
	// 分布类型
	// access_source_session_cnt 访问来源分布
	// access_staytime_info	访问时长分布
	// access_depth_info 访问深度的分布
	Index string `json:"index"`

	// 分布数据列表
	ItemList []VisitDistributionItem `json:"item_list"`
}

// VisitDistributionItem 分布数据项
type VisitDistributionItem struct {
	// 场景 id，定义在各个 index 下不同，具体参见下方表格
	//
	// 访问来源 key 对应关系（index="access_source_session_cnt")
	// key	访问来源				对应场景值
	// 1	小程序历史列表			1001
	// 2	搜索					1005 1006 1027 1042 1053
	// 3	会话					1007 1008 1044 1096
	// 4	扫一扫二维码			1011 1047
	// 5	公众号主页				1020
	// 6	聊天顶部				1022
	// 7	系统桌面				1023
	// 8	小程序主页				1024
	// 9	附近的小程序			1026 1068
	// 11	模板消息				1014 1043
	// 13	公众号菜单				1035
	// 14	APP分享					1036
	// 15	支付完成页				1034
	// 16	长按识别二维码			1012 1048
	// 17	相册选取二维码			1013 1049
	// 18	公众号文章				1058
	// 19	钱包					1019
	// 20	卡包					1028
	// 21	小程序内卡券			1029
	// 22	其他小程序				1037
	// 23	其他小程序返回			1038
	// 24	卡券适用门店列表		1052
	// 25	搜索框快捷入口			1054
	// 26	小程序客服消息			1073 1081
	// 27	公众号下发				1074 1082
	// 29	任务栏-最近使用			1089
	// 30	长按小程序菜单圆点		1090
	// 31	连wifi成功页			1078
	// 32	城市服务				1092
	// 33	微信广告				1045 1046 1067 1084
	// 34	其他移动应用			1069
	// 35	发现入口-我的小程序（基础库2.2.4版本起1103场景值废弃，不影响此处统计结果）	1103
	// 36	任务栏-我的小程序（基础库2.2.4版本起1104场景值废弃，不影响此处统计结果）	1104
	// 10	其他	除上述外其余场景值
	//
	// 访问时长 key 对应关系（index="access_staytime_info"）
	// key	访问时长
	// 1	0-2s
	// 2	3-5s
	// 3	6-10s
	// 4	11-20s
	// 5	20-30s
	// 6	30-50s
	// 7	50-100s
	// 8	>100s
	//
	// 平均访问深度 key 对应关系（index="access_depth_info"）
	// key	访问时长
	// 1	1 页
	// 2	2 页
	// 3	3 页
	// 4	4 页
	// 5	5 页
	// 6	6-10 页
	// 7	>10 页
	//
	Key int `json:"key"`

	Value               int64 `json:"value"`                  // 该场景 id 访问 pv
	AccessSourceVisitUV int64 `json:"access_source_visit_uv"` // 该场景访问uv

}

// GetVisitDistribution 获取用户小程序访问分布数据
//
// beginDate 开始日期。endDate 结束日期，限定查询 1 天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/analysis.getVisitDistribution.html
func (s *AnalysisService) GetVisitDistribution(ctx context.Context, beginDate, endDate string) (*VisitDistribution, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidvisitdistribution", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	vd := new(VisitDistribution)
	_, err = s.client.Do(req, vd)
	return vd, err
}

// VisitPage 访问页面数据
type VisitPage struct {
	PagePath       string `json:"page_path"`        // 页面路径
	PageVisitPV    string `json:"page_visit_pv"`    // 访问次数
	PageVisitUV    string `json:"page_visit_uv"`    // 访问人数
	PageStaytimePV string `json:"page_staytime_pv"` // 次均停留时长
	EntryPagePV    string `json:"entry_page_pv"`    // 进入页次数
	ExitPagePV     string `json:"exit_page_pv"`     // 退出页次数
	PageSharePV    string `json:"page_share_pv"`    // 转发次数
	PageShareUV    string `json:"page_share_uv"`    // 转发人数
}

// GetVisitPage 访问页面。
// 目前只提供按 page_visit_pv 排序的 top200。
//
// beginDate 开始日期。endDate 结束日期，限定查询 1 天数据，允许设置的最大值为昨日。格式为 yyyymmdd
//
// 微信文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/data-analysis/analysis.getVisitPage.html
func (s *AnalysisService) GetVisitPage(ctx context.Context, beginDate, endDate string) ([]*VisitPage, error) {

	u, err := s.client.apiURL(ctx, "datacube/getweanalysisappidvisitpage", nil)
	if err != nil {
		return nil, err
	}

	rdata := map[string]string{
		"begin_date": beginDate,
		"end_date":   endDate,
	}
	req, err := s.client.NewRequest(ctx, "POST", u.String(), rdata)
	if err != nil {
		return nil, err
	}

	var data struct {
		List []*VisitPage `json:"list"` // 数据列表
	}

	_, err = s.client.Do(req, &data)
	return data.List, err
}
