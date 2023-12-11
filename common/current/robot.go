package current

import (
	"go-wechat/model"
	plugin "go-wechat/plugin"
)

// robotInfo
// @description: 机器人信息
type robotInfo struct {
	info           model.RobotUserInfo
	MessageHandler plugin.MessageHandler // 启用的插件
}

// 当前接入的机器人信息
var ri robotInfo

// SetRobotInfo
// @description: 设置机器人信息
// @param info
func SetRobotInfo(info model.RobotUserInfo) {
	ri.info = info
}

// GetRobotInfo
// @description: 获取机器人信息
// @return model.RobotUserInfo
func GetRobotInfo() model.RobotUserInfo {
	return ri.info
}

// GetRobotMessageHandler
// @description: 获取机器人插件信息
// @return robotInfo
func GetRobotMessageHandler() plugin.MessageHandler {
	return ri.MessageHandler
}

// SetRobotMessageHandler
// @description: 设置机器人插件信息
// @param handler
func SetRobotMessageHandler(handler plugin.MessageHandler) {
	ri.MessageHandler = handler
}
