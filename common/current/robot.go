package current

import "go-wechat/model"

var robotInfo model.RobotUserInfo

// SetRobotInfo
// @description: 设置机器人信息
// @param info
func SetRobotInfo(info model.RobotUserInfo) {
	robotInfo = info
}

// GetRobotInfo
// @description: 获取机器人信息
// @return model.RobotUserInfo
func GetRobotInfo() model.RobotUserInfo {
	return robotInfo
}
