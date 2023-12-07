package initialization

import (
	"github.com/go-resty/resty/v2"
	"go-wechat/common/current"
	"go-wechat/config"
	"go-wechat/model"
	"log"
)

// InitWechatRobotInfo
// @description: 初始化微信机器人信息
func InitWechatRobotInfo() {
	// 获取数据
	var base model.Response[model.RobotUserInfo]
	_, err := resty.New().R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetResult(&base).
		Post(config.Conf.Wechat.GetURL("/api/userInfo"))
	if err != nil {
		log.Printf("获取机器人信息失败: %s", err.Error())
		return
	}

	log.Printf("机器人Id: %s", base.Data.WxId)
	log.Printf("机器人微信号: %s", base.Data.Account)
	log.Printf("机器人名称: %s", base.Data.Name)

	// 设置为单例
	current.SetRobotInfo(base.Data)
}
