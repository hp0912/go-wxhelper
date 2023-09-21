package service

import (
	"go-wechat/client"
	"go-wechat/entity"
	"log"
)

// SaveMessage
// @description: 消息入库
// @param msg
func SaveMessage(msg entity.Message) {
	// 检查消息是否存在，存在就跳过
	var count int64
	err := client.MySQL.Model(&entity.Message{}).Where("msg_id = ?", msg.MsgId).Count(&count).Error
	if err != nil {
		log.Printf("检查消息是否存在失败, 错误信息: %v", err)
		return
	}
	if count > 0 {
		return
	}
	err = client.MySQL.Create(&msg).Error
	if err != nil {
		log.Printf("消息入库失败, 错误信息: %v", err)
	}
	log.Printf("消息入库成功，消息Id: %d", msg.MsgId)
}
