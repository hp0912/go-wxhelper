package initialization

import (
	"go-wechat/common/current"
	"go-wechat/model"
	plugin "go-wechat/plugin"
	"go-wechat/plugin/plugins"
	"go-wechat/service"
)

// Plugin
// @description: 初始化插件
func Plugin() {
	// 定义一个处理器
	dispatcher := plugin.NewMessageMatchDispatcher()
	// 设置为异步处理
	dispatcher.SetAsync(true)

	// 注册插件

	// 保存消息进数据库
	dispatcher.RegisterHandler(func(*model.Message) bool {
		return true
	}, plugins.SaveToDb)

	// 私聊指令消息
	dispatcher.RegisterHandler(func(m *model.Message) bool {
		// 私聊消息 或 群聊艾特机器人并且以/开头的消息
		return (m.IsPrivateText() || (m.IsAt() && !m.IsAtAll() && m.CleanContentStartWith("/"))) && service.CheckIsEnableCommand(m.FromUser)
	}, plugins.Command)

	// AI消息插件
	dispatcher.RegisterHandler(func(m *model.Message) bool {
		// 群内@或者私聊文字消息
		return m.IsAt() && !m.IsAtAll() || m.IsPrivateText()
	}, plugins.AI)

	// 欢迎新成员
	dispatcher.RegisterHandler(func(m *model.Message) bool {
		return m.IsNewUserJoin()
	}, plugins.WelcomeNew)

	// 注册消息处理器
	current.SetRobotMessageHandler(plugin.DispatchMessage(dispatcher))
}
