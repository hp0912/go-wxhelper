package service

import (
	"go-wechat/client"
	"go-wechat/entity"
)

// GetAllEnableAI
// @description: 取出所有启用了AI的好友或群组
// @return []entity.Friend
func GetAllEnableAI() (records []entity.Friend, err error) {
	err = client.MySQL.Where("enable_ai = ?", 1).Find(&records).Error
	return
}

// GetAllEnableChatRank
// @description: 取出所有启用了聊天排行榜的群组
// @return records
// @return err
func GetAllEnableChatRank() (records []entity.Friend, err error) {
	err = client.MySQL.Where("enable_chat_rank = ?", 1).Where("wxid LIKE '%@chatroom'").Find(&records).Error
	return
}
