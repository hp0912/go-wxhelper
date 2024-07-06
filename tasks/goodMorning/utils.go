package goodMorning

import "go-wechat/client"

// groupSummary
// @description: 群动态
type groupSummary struct {
	GroupID        string
	Year           int
	Month          int
	Date           int
	Week           string
	UserTotalCount int // 当前群成员总数
	UserJoinCount  int // 昨天入群数
	UserLeaveCount int // 昨天离群数
	UserChatCount  int // 昨天聊天人数
	MessageCount   int // 昨天消息数
}

// getGroupSummary
// @description: 获取群动态
// @param groupId string 群Id
// @return summary
// @return err
func getGroupSummary(groupId string) (groupSummary, error) {
	summary := groupSummary{}

	// 当前群总人数
	userTotalCount := []groupSummary{}
	tx := client.MySQL.Table("t_group_user").Select("count( 1 ) AS `user_total_count`").Where("group_id = ?", groupId).Where("is_member = 1")
	err := tx.Find(&userTotalCount).Error
	if err != nil {
		return summary, err
	}

	// 昨天入群人数
	userJoinCount := []groupSummary{}
	tx = client.MySQL.Table("t_group_user").Select("count( 1 ) AS `user_join_count`").Where("group_id = ?", groupId).Where("DATEDIFF(join_time, NOW()) = -1")
	err = tx.Find(&userJoinCount).Error
	if err != nil {
		return summary, err
	}

	// 昨天离群人数
	userLeaveCount := []groupSummary{}
	tx = client.MySQL.Table("t_group_user").Select("count( 1 ) AS `user_leave_count`").Where("group_id = ?", groupId).Where("DATEDIFF(leave_time, NOW()) = -1")
	err = tx.Find(&userLeaveCount).Error
	if err != nil {
		return summary, err
	}

	// 昨天聊天人数和消息数
	userChatCount := []groupSummary{}
	tx = client.MySQL.Table("t_message AS tm").
		Joins("LEFT JOIN t_group_user AS tgu ON tgu.wxid = tm.group_user AND tm.from_user = tgu.group_id").
		Select("count( 1 ) AS `user_chat_count`").
		Where("tm.from_user = ?", groupId).
		Where("tm.type < 10000").
		Where("DATEDIFF(tm.create_at, NOW()) = -1").
		Group("tm.group_user, tgu.nickname")
	err = tx.Find(&userChatCount).Error
	if err != nil {
		return summary, err
	}
	messageCount := 0
	for _, item := range userChatCount {
		messageCount += item.UserChatCount
	}

	summary.GroupID = groupId
	if len(userTotalCount) > 0 {
		summary.UserTotalCount = userTotalCount[0].UserTotalCount
	}
	if len(userJoinCount) > 0 {
		summary.UserJoinCount = userJoinCount[0].UserJoinCount
	}
	if len(userLeaveCount) > 0 {
		summary.UserLeaveCount = userLeaveCount[0].UserLeaveCount
	}
	summary.UserChatCount = len(userChatCount)
	summary.MessageCount = messageCount

	return summary, nil
}
