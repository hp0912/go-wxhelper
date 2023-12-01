package service

import (
	"go-wechat/client"
	"go-wechat/vo"
)

// GetGroupUsersByGroupId
// @description: 根据群Id取出群成员列表
// @param groupId
// @return records
// @return err
func GetGroupUsersByGroupId(groupId string) (records []vo.GroupUserItem, err error) {
	err = client.MySQL.
		Table("t_group_user AS tgu").
		Joins("LEFT JOIN t_message AS tm ON tm.from_user = tgu.group_id AND tm.group_user = tgu.wxid").
		//Select("tgu.wxid", "tgu.nickname", "tgu.head_image", "tgu.is_member", "tgu.leave_time",
		//	"tgu.skip_chat_rank", "MAX(tm.create_at) AS last_active_time").
		Select("tgu.*", "MAX(tm.create_at) AS last_active_time").
		Where("tgu.group_id = ?", groupId).
		Group("tgu.group_id, tgu.wxid").
		Order("tgu.join_time DESC").
		Find(&records).Error
	return
}
