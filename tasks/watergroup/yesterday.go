package watergroup

import (
	"fmt"
	"go-wechat/client"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/service"
	"go-wechat/utils"
	"log"
	"strings"
	"time"
)

// 水群排行榜

// Yesterday
// @description: 昨日排行榜
func Yesterday() {
	groups, err := service.GetAllEnableChatRank()
	if err != nil {
		log.Printf("获取启用了聊天排行榜的群组失败, 错误信息: %v", err)
		return
	}

	for _, group := range groups {
		// 消息统计
		dealYesterday(group.Wxid)

		res, ok := config.Conf.Resource["wordcloud"]
		if !ok {
			continue
		}

		// 获取昨日日期
		yd := time.Now().Local().AddDate(0, 0, -1).Format("20060102")
		// 发送词云
		fileName := fmt.Sprintf("%s_%s.png", yd, group.Wxid)
		utils.SendImage(group.Wxid, fmt.Sprintf(res.Path, fileName), 0)
	}
}

// dealYesterday
// @description: 处理请求
// @param gid
func dealYesterday(gid string) {
	notifyMsgs := []string{"#昨日水群排行榜"}

	// 获取昨日消息总数
	records, err := getRankData(gid, "yesterday")
	if err != nil {
		log.Printf("获取昨日消息排行失败, 错误信息: %v", err)
		return
	}
	log.Printf("昨日消息总数: %+v", records)
	// 莫得消息，直接返回
	if len(records) == 0 {
		log.Printf("昨日群[%s]无对话记录", gid)
		return
	}

	// 查询群成员总数
	var groupUsers int64
	err = client.MySQL.Model(&entity.GroupUser{}).
		Where("group_id = ?", gid).
		Where("is_member IS TRUE").
		Count(&groupUsers).Error
	if err != nil {
		log.Printf("查询群成员总数失败, 错误信息: %v", err)
	}
	// 计算活跃度
	showActivity := err == nil && groupUsers > 0
	activity := "0.00"
	if groupUsers > 0 {
		activity = fmt.Sprintf("%.2f", (float64(len(records))/float64(groupUsers))*100)
	}

	// 计算消息总数
	var msgCount int64
	for _, v := range records {
		msgCount += v.Count
	}
	// 组装消息总数推送信息
	notifyMsgs = append(notifyMsgs, " ")
	notifyMsgs = append(notifyMsgs, fmt.Sprintf("🗣️ 昨日本群 %d 位朋友共产生 %d 条发言", len(records), msgCount))
	if showActivity {
		notifyMsgs = append(notifyMsgs, fmt.Sprintf("🎭 活跃度: %s%%", activity))
	}
	notifyMsgs = append(notifyMsgs, "\n🏵 活跃用户排行榜 🏵")

	notifyMsgs = append(notifyMsgs, " ")
	for i, r := range records {
		// 只取前十条
		if i >= 10 {
			break
		}
		log.Printf("账号: %s[%s] -> %d", r.Nickname, r.GroupUser, r.Count)
		badge := "🏆"
		switch i {
		case 0:
			badge = "🥇"
		case 1:
			badge = "🥈"
		case 2:
			badge = "🥉"
		}
		notifyMsgs = append(notifyMsgs, fmt.Sprintf("%s %s -> %d条", badge, r.Nickname, r.Count))
	}

	notifyMsgs = append(notifyMsgs, " \n🎉感谢以上群友昨日对群活跃做出的卓越贡献，也请未上榜的群友多多反思。")

	log.Printf("排行榜: \n%s", strings.Join(notifyMsgs, "\n"))
	go utils.SendMessage(gid, "", strings.Join(notifyMsgs, "\n"), 0)
}
