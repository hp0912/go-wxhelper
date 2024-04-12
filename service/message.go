package service

import (
	"go-wechat/client"
	"go-wechat/entity"
	"go-wechat/vo"
	"log"
	"os"
	"strconv"
)

// SaveMessage
// @description: 消息入库
// @param msg
func SaveMessage(msg entity.Message) {
	if flag, _ := strconv.ParseBool(os.Getenv("DONT_SAVE")); flag {
		return
	}

	// 检查消息是否存在，存在就跳过
	var count int64
	err := client.MySQL.Model(&entity.Message{}).Where("msg_id = ?", msg.MsgId).Count(&count).Error
	if err != nil {
		log.Printf("检查消息是否存在失败, 错误信息: %v", err)
		return
	}
	if count > 0 {
		//log.Printf("消息已存在，消息Id: %d", msg.MsgId)
		return
	}
	err = client.MySQL.Create(&msg).Error
	if err != nil {
		log.Printf("消息入库失败, 错误信息: %v", err)
	}
	log.Printf("消息入库成功，消息Id: %d", msg.MsgId)

	// 更新最后活跃时间
	// 只处理收到的消息
	if msg.MsgId > 1 {
		go updateLastActive(msg)
	}
}

// GetTextMessagesById
// @description: 根据群id或者用户Id获取消息
// @param id
// @return records
// @return err
func GetTextMessagesById(id string) (records []vo.TextMessageItem, err error) {
	tx := client.MySQL.
		Table("`t_message` AS tm").
		Joins("LEFT JOIN t_group_user AS tgu ON tm.group_user = tgu.wxid AND tgu.group_id = tm.from_user").
		Select("tgu.nickname", "IF( tm.type = 49, EXTRACTVALUE ( tm.content, \"/msg/appmsg/title\" ), tm.content ) AS message").
		Where("tm.`from_user` = ?", id).
		Where(`(tm.type = 1 OR ( tm.type = 49 AND EXTRACTVALUE ( tm.content, "/msg/appmsg/type" ) = '57' ))`).
		Where("DATE ( tm.create_at ) = DATE ( CURDATE() - INTERVAL 1 DAY )").
		Order("tm.create_at ASC")

	err = tx.Find(&records).Error
	return
}
