package miniprogram

import (
	"context"

	"github.com/google/go-querystring/query"
)

// OperationService 运维中心
type OperationService service

// OperationRealtimeLogSearchRequest 实时日志查询请求
type OperationRealtimeLogSearchRequest struct {
	// YYYYMMDD格式的日期，仅支持最近7天
	Date string `url:"date"`
	// 开始时间，必须是date指定日期的时间
	Begintime int64 `url:"begintime"`
	// 结束时间，必须是date指定日期的时间
	Endtime int64 `url:"endtime"`
	// 开始返回的数据下标，用作分页，默认为0
	Start int64 `url:"start,omitempty"`
	// 返回的数据条数，用作分页，默认为20
	Limit int64 `url:"limit,omitempty"`
	// 小程序启动的唯一ID，按TraceId查询会展示该次小程序启动过程的所有页面的日志。
	TraceID string `url:"traceId,omitempty"`
	// 小程序页面路径，例如pages/index/index
	URL string `url:"url,omitempty"`
	// 用户微信号或者OpenId
	ID string `url:"id,omitempty"`
	// 开发者通过setFileterMsg/addFilterMsg指定的filterMsg字段
	FilterMsg string `url:"filterMsg,omitempty"`
	// 日志等级，返回大于等于level等级的日志，
	// level的定义为2（Info）、4（Warn）、8（Error），
	// 如果指定为4，则返回大于等于4的日志，即返回Warn和Error日志。
	Level int `url:"level,omitempty"`
}

// OperationRealtimeLog 实时日志
type OperationRealtimeLog struct {
	// 日志等级，是msg数组里面的所有level字段的或操作得到的结果。例如msg数组里有两条日志，Info（值为2）和Warn（值为4），则level值为6
	Level int `json:"level,omitempty"`
	// 基础库版本
	LibraryVersion string `json:"libraryVersion"`
	// 客户端版本
	ClientVersion string `json:"clientVersion"`
	// 微信用户OpenID
	ID string `json:"id"`
	// 打日志的Unix时间戳
	Timestamp int64 `json:"timestamp"`
	// 1 安卓 2 IOS
	Platform int `json:"platform"`
	// 小程序页面链接
	URL string `json:"url"`
	// 小程序启动的唯一ID，按TraceId查询会展示该次小程序启动过程的所有页面的日志。
	TraceID string `json:"traceid"`
	// 开发者通过setFileterMsg/addFilterMsg指定的filterMsg字段
	FilterMsg string `json:"filterMsg"`
	// 日志内容数组，log.info等的内容存在这里
	Msg     []OperationRealtimeLogMsg `json:"msg"`
	Device  string                    `json:"device"`
	MediaID uint64                    `json:"media_id"`
}

// OperationRealtimeLogMsg 实时日志信息
type OperationRealtimeLogMsg struct {
	Time  int64    `json:"time"`
	Level int      `json:"level"`
	Msg   []string `json:"msg"`
}

// RealtimeLogSearch 实时日志查询
//
// 文档: https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/operation/operation.realtimelogSearch.html
func (s *OperationService) RealtimeLogSearch(ctx context.Context, r *OperationRealtimeLogSearchRequest) (int64, []OperationRealtimeLog, error) {

	v, err := query.Values(r)
	if err != nil {
		return 0, nil, err
	}
	u, err := s.client.apiURL(ctx, "wxaapi/userlog/userlog_search", v)
	if err != nil {
		return 0, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, nil, err
	}
	var data struct {
		Data struct {
			Total int64                  `json:"total"`
			List  []OperationRealtimeLog `json:"list"`
		} `json:"data"`
	}
	_, err = s.client.Do(req, &data)
	return data.Data.Total, data.Data.List, err
}
