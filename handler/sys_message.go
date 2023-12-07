package handler

import (
	"go-wechat/client"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/model"
	"go-wechat/utils"
	"log"
)

// handleNewUserJoin
// @description: 欢迎新成员
// @param m
func handleNewUserJoin(m model.Message) {
	// 判断是否开启迎新
	var count int64
	err := client.MySQL.Model(&entity.Friend{}).
		Where("enable_welcome IS TRUE").
		Where("wxid = ?", m.FromUser).
		Count(&count).Error
	if err != nil {
		log.Printf("查询是否开启迎新失败: %s", err.Error())
		return
	}
	if count < 1 {
		return
	}

	// 读取欢迎新成员配置
	conf, ok := config.Conf.Resource["welcomeNew"]
	if !ok {
		// 未配置，跳过
		return
	}
	switch conf.Type {
	case "text":
		// 文字类型
		utils.SendMessage(m.FromUser, "", conf.Path, 0)
	case "image":
		// 图片类型
		utils.SendImage(m.FromUser, conf.Path, 0)
	case "emotion":
		// 表情类型
		utils.SendEmotion(m.FromUser, conf.Path, 0)
	}
}
