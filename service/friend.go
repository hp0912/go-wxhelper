package service

import (
	"go-wechat/client"
	"go-wechat/entity"
	"go-wechat/vo"
	"strings"
)

// GetAllFriend
// @description: 取出所有好友
// @return friends
// @return groups
// @return err
func GetAllFriend() (friends, groups []vo.FriendItem, err error) {
	var records []vo.FriendItem
	err = client.MySQL.
		Table("t_friend AS tf").
		Joins("LEFT JOIN t_message AS tm ON tf.wxid = tm.from_user").
		Select("tf.*", "MAX(tm.create_at) AS last_active_time").
		Group("tf.wxid").
		Order("last_active_time DESC").
		Find(&records).Error
	if err != nil {
		return
	}
	for _, record := range records {
		if strings.HasSuffix(record.Wxid, "@chatroom") {
			groups = append(groups, record)
		} else {
			friends = append(friends, record)
		}
	}
	return
}

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
