package watergroup

import "go-wechat/client"

// rankUser
// @description: 排行榜用户
type rankUser struct {
	GroupUser string // 微信Id
	Nickname  string // 昵称
	Count     int64  // 消息数
}

// getRankData
// @description: 获取消息排行榜
// @param groupId string 群Id
// @param d string 模式(yesterday | week | month)
// @return rank
// @return err
func getRankData(groupId, date string) (rank []rankUser, err error) {
	tx := client.MySQL.Table("t_message AS tm").
		Joins("LEFT JOIN t_group_user AS tgu ON tgu.wxid = tm.group_user AND tm.from_user = tgu.group_id AND tgu.skip_chat_rank = 0 AND is_member = 1").
		Select("tm.group_user", "tgu.nickname", "count( 1 ) AS `count`").
		Where("tm.from_user = ?", groupId).
		Where("tm.type < 10000").
		Group("tm.group_user, tgu.nickname").
		Order("`count` DESC")

	// 根据参数获取不同日期的数据
	switch date {
	case "yesterday":
		tx.Where("DATEDIFF(tm.create_at,NOW()) = -1")
	case "week":
		tx.Where("YEARWEEK(date_format(tm.create_at, '%Y-%m-%d')) = YEARWEEK(now()) - 1")
	case "month":
		tx.Where("PERIOD_DIFF(date_format(now(), '%Y%m'), date_format(create_at, '%Y%m')) = 1")
	case "year":
		tx.Where("YEAR(tm.create_at) = YEAR(NOW()) - 1")
	}

	// 查询指定时间段全部数据
	err = tx.Find(&rank).Error
	return
}
