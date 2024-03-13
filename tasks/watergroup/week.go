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

// Week
// @description: 周排行榜
func Week() {
	groups, err := service.GetAllEnableChatRank()
	if err != nil {
		log.Printf("获取启用了聊天排行榜的群组失败, 错误信息: %v", err)
		return
	}

	for _, group := range groups {
		// 消息统计
		dealWeek(group.Wxid)

		res, ok := config.Conf.Resource["wordcloud"]
		if !ok {
			continue
		}

		// 获取上周周数
		year, weekNo := time.Now().Local().AddDate(0, 0, -1).ISOWeek()
		// 发送词云
		fileName := fmt.Sprintf("%d%d_%s.png", year, weekNo, group.Wxid)
		utils.SendImage(group.Wxid, fmt.Sprintf(res.Path, fileName), 0)
	}
}

// dealWeek
// @description: 处理请求
// @param gid
func dealWeek(gid string) {
	notifyMsgs := []string{"#上周水群排行榜"}

	// 获取上周消息总数
	records, err := getRankData(gid, "week")
	if err != nil {
		log.Printf("获取上周消息排行失败, 错误信息: %v", err)
		return
	}
	log.Printf("上周消息总数: %+v", records)
	// 莫得消息，直接返回
	if len(records) == 0 {
		log.Printf("上周群[%s]无对话记录", gid)
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

	// 计算消息总数、中位数
	var msgCount int64
	var medianCount int64
	for idx, v := range records {
		msgCount += v.Count
		if idx == (len(records)/2)-1 {
			medianCount = v.Count
		}
	}
	// 计算活跃用户人均消息条数
	avgMsgCount := int(float64(msgCount) / float64(len(records)))

	// 组装消息总数推送信息
	notifyMsgs = append(notifyMsgs, " ")
	notifyMsgs = append(notifyMsgs, fmt.Sprintf("🗣️ 上周本群 %d 位朋友共产生 %d 条发言", len(records), msgCount))
	if showActivity {
		notifyMsgs = append(notifyMsgs, fmt.Sprintf("🎭 活跃度: %s%%，人均消息条数: %d，中位数: %d", activity, avgMsgCount, medianCount))
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

	notifyMsgs = append(notifyMsgs, " \n🎉感谢以上群友上周对群活跃做出的卓越贡献，也请未上榜的群友多多反思。")

	log.Printf("排行榜: \n%s", strings.Join(notifyMsgs, "\n"))
	go utils.SendMessage(gid, "", strings.Join(notifyMsgs, "\n"), 0)
}
