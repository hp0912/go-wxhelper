package handler

import (
	"go-wechat/client"
	"go-wechat/entity"
	"go-wechat/model"
	"go-wechat/utils"
)

// handleNewUserJoin
// @description: 欢迎新成员
// @param m
func handleNewUserJoin(m model.Message) {
	// 判断是否开启迎新
	var count int64
	_ = client.MySQL.Model(&entity.Friend{}).Where("enable_welcome IS TRUE").Where("wxid = ?", m.FromUser).Count(&count).Error
	if count < 1 {
		return
	}

	// 发一张图乐呵乐呵

	// 自己欢迎自己图片地址  D:\Share\emoticon\welcome-yourself.gif
	utils.SendImage(m.FromUser, "D:\\Share\\emoticon\\welcome-yourself.gif", 0)
}
