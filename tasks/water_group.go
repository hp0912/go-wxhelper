package tasks

import (
	"fmt"
	"go-wechat/client"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/utils"
	"log"
	"strings"
	"time"
)

// 水群排行榜

// yesterday
// @description: 昨日排行榜
func yesterday() {
	for _, id := range config.Conf.Task.WaterGroup.Groups {
		// 消息统计
		dealYesterday(id)
		// 获取昨日日期
		yd := time.Now().Local().AddDate(0, 0, -1).Format("20060102")
		// 发送词云
		fileName := fmt.Sprintf("%s_%s.png", yd, id)
		utils.SendImage(id, "D:\\Share\\wordcloud\\"+fileName, 0)
	}
}

// dealYesterday
// @description: 处理请求
// @param gid
func dealYesterday(gid string) {
	notifyMsgs := []string{"#昨日水群排行榜"}

	// 获取昨日消息总数
	var yesterdayMsgCount int64
	err := client.MySQL.Model(&entity.Message{}).
		Where("from_user = ?", gid).
		Where("`type` < 10000").
		Where("DATEDIFF(create_at,NOW()) = -1").
		Count(&yesterdayMsgCount).Error
	if err != nil {
		log.Printf("获取昨日消息总数失败, 错误信息: %v", err)
		return
	}
	log.Printf("昨日消息总数: %d", yesterdayMsgCount)
	if yesterdayMsgCount == 0 {
		return
	}

	notifyMsgs = append(notifyMsgs, " ")
	notifyMsgs = append(notifyMsgs, fmt.Sprintf("昨日消息总数: %d", yesterdayMsgCount))

	// 返回数据
	type record struct {
		GroupUser string
		Nickname  string
		Count     int64
	}

	var records []record
	tx := client.MySQL.Table("t_message AS tm").
		Joins("LEFT JOIN t_group_user AS tgu ON tgu.wxid = tm.group_user AND tm.from_user = tgu.group_id").
		Select("tm.group_user", "tgu.nickname", "count( 1 ) AS `count`").
		Where("tm.from_user = ?", gid).
		Where("tm.type < 10000").
		Where("DATEDIFF(tm.create_at,NOW()) = -1").
		Group("tm.group_user, tgu.nickname").Order("`count` DESC").
		Limit(10)

	// 黑名单
	blacklist := config.Conf.Task.WaterGroup.Blacklist
	// 如果有黑名单，过滤掉
	if len(blacklist) > 0 {
		tx.Where("tm.group_user NOT IN (?)", blacklist)
	}

	err = tx.Find(&records).Error

	if err != nil {
		log.Printf("获取昨日消息失败, 错误信息: %v", err)
		return
	}
	notifyMsgs = append(notifyMsgs, " ")
	for i, r := range records {
		log.Printf("账号: %s[%s] -> %d", r.Nickname, r.GroupUser, r.Count)
		notifyMsgs = append(notifyMsgs, fmt.Sprintf("#%d: %s -> %d条", i+1, r.Nickname, r.Count))
	}

	notifyMsgs = append(notifyMsgs, " \n请未上榜的群友多多反思。")

	log.Printf("排行榜: \n%s", strings.Join(notifyMsgs, "\n"))
	go utils.SendMessage(gid, "", strings.Join(notifyMsgs, "\n"), 0)
}
