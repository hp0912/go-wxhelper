package plugins

// Message
// @description: 插件消息
type Message struct {
	GroupId string // 消息来源群Id
	UserId  string // 消息来源用户Id
	Message string // 消息内容
	IsBreak bool   // 是否中断消息传递
}
