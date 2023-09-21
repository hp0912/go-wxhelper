package tasks

import (
	"go-wechat/client"
	"go-wechat/entity"
	"log"
)

// 水群排行榜

// yesterday
// @description: 昨日排行榜
func yesterday() {
	// 获取昨日消息总数
	var yesterdayMsgCount int64
	err := client.MySQL.Model(&entity.Message{}).
		Where("from_user = ?", "18958257758@chatroom").
		Where("DATEDIFF(create_at,NOW()) = -1").
		Count(&yesterdayMsgCount).Error
	if err != nil {
		log.Printf("获取昨日消息总数失败, 错误信息: %v", err)
		return
	}
	log.Printf("昨日消息总数: %d", yesterdayMsgCount)

	// 返回数据
	type record struct {
		GroupUser string
		Count     int64
	}

	var records []record
	err = client.MySQL.Model(&entity.Message{}).
		Select("group_user", "count( 1 ) AS `count`").
		Where("from_user = ?", "18958257758@chatroom").
		Where("DATEDIFF(create_at,NOW()) = -1").
		Group("group_user").Order("`count` DESC").
		Limit(5).Find(&records).Error
	if err != nil {
		log.Printf("获取昨日消息失败, 错误信息: %v", err)
		return
	}
	for _, r := range records {
		log.Printf("账号: %s -> %d", r.GroupUser, r.Count)
	}
}
