package message

// Message 消息
type Message interface {
	// Type 消息类型
	Type() string
}

// Text 文本消息
type Text struct {
	Content string `json:"content"` // 文本消息内容
}

// Type 消息类型
func (m Text) Type() string {
	return TypeText
}

// Link 图文链接
type Link struct {
	// 消息标题
	Title string `json:"title"`
	// 图文链接消息
	Description string `json:"description"`
	// 图文链接消息被点击后跳转的链接
	URL string `json:"url"`
	// 图文链接消息的图片链接，支持 JPG、PNG 格式，较好的效果为大图 640 X 320，小图 80 X 80
	ThumbURL string `json:"thumb_url"`
}

// Type 消息类型
func (m Link) Type() string {
	return TypeLink
}

// Image 图片消息
type Image struct {
	MediaID string `json:"media_id"` // 发送的图片的媒体ID，通过 新增素材接口 上传图片文件获得。
}

// Type 消息类型
func (m Image) Type() string {
	return TypeImage
}

// Voice 语音消息
type Voice struct {
	MediaID string `json:"media_id"`
}

// Type 消息类型
func (m Voice) Type() string {
	return TypeVoice
}

// Video 视频消息
type Video struct {
	MediaID      string `json:"media_id"`
	ThumbMediaID string `json:"thumb_media_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

// Type 消息类型
func (m Video) Type() string {
	return TypeVideo
}

// Music 音乐消息
type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicURL     string `json:"musicurl"`
	HQMusicURL   string `json:"hqmusicurl"`
	ThumbMediaID string `json:"thumb_media_id"`
}

// News 图文消息（点击跳转到外链）
type News struct {
	Articles []NewsArticle `json:"articles"`
}

// Type 消息类型
func (m News) Type() string {
	return TypeNews
}

// NewsArticle 图文消息内容
type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
}

// MPNews 图文消息（点击跳转到图文消息页面）
type MPNews struct {
	MediaID string `json:"media_id"`
}

// Type 消息类型
func (m MPNews) Type() string {
	return TypeMPNews
}

// MsgMenu 菜单消息
type MsgMenu struct {
	HeadContent string        `json:"head_content"`
	List        []MsgMenuItem `json:"list"`
	TailContent string        `json:"tail_content"`
}

// Type 消息类型
func (m MsgMenu) Type() string {
	return TypeMsgMenu
}

// MsgMenuItem 菜单项
type MsgMenuItem struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// WXCard 卡券消息
type WXCard struct {
	CardID string `json:"card_id"`
}

// Type 消息类型
func (m WXCard) Type() string {
	return TypeWXCard
}

// MiniProgramPage 小程序卡片
type MiniProgramPage struct {
	// 消息标题
	Title string `json:"title"`
	APPID string `json:"appid"`
	// 小程序的页面路径，跟app.json对齐，支持参数，比如pages/index/index?foo=bar
	PagePath string `json:"page_path"`
	// 小程序消息卡片的封面， image 类型的 media_id，通过 新增素材接口 上传图片文件获得，建议大小为 520*416
	ThumbMediaID string `json:"thumb_media_id"`
}

// Type 消息类型
func (m MiniProgramPage) Type() string {
	return TypeMiniProgramPage
}
