package handler

import (
	"go-wechat/model"
	"go-wechat/utils"
	"strings"
)

// handleSysMessage
// @description: 系统消息处理
// @param m
func handleSysMessage(m model.Message) {
	// 有人进群
	if strings.Contains(m.Content, "\"邀请\"") && strings.Contains(m.Content, "\"加入了群聊") {
		// 发一张图乐呵乐呵
		// 自己欢迎自己图片地址  D:\Share\emoticon\welcome-yourself.gif
		utils.SendImage(m.FromUser, "D:\\Share\\emoticon\\welcome-yourself.gif", 0)
	}
}
