package message

// 发送消息类型
const (
	TypeText            = "text"            // 文本消息
	TypeImage           = "image"           // 图片消息
	TypeLink            = "link"            // 图文链接
	TypeVoice           = "voice"           // 语音消息
	TypeVideo           = "video"           // 视频消息
	TypeMusic           = "music"           // 音乐消息
	TypeNews            = "news"            // 图文消息（点击跳转到外链）
	TypeMPNews          = "mpnews"          // 图文消息（点击跳转到图文消息页面）
	TypeMsgMenu         = "msgmenu"         // 菜单消息
	TypeWXCard          = "wxcard"          // 卡券
	TypeMiniProgramPage = "miniprogrampage" // 小程序卡片
)

var typeList = []string{
	TypeText,
	TypeImage,
	TypeLink,
	TypeVoice,
	TypeVideo,
	TypeMusic,
	TypeNews,
	TypeMPNews,
	TypeMsgMenu,
	TypeWXCard,
	TypeMiniProgramPage,
}

// IsMessageType 是否是消息类型
func IsMessageType(t string) bool {
	for _, v := range typeList {
		if t == v {
			return true
		}
	}
	return false
}

// 客服输入状态
type StatusCommand string

// 客服输入状态
const (
	StatusCommandTyping       StatusCommand = "Typing"       // 正在输入
	StatusCommandCancelTyping StatusCommand = "CancelTyping" // 取消对用户的”正在输入"状态
)
