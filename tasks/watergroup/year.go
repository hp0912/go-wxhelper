package watergroup

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/service"
	"go-wechat/utils"
	"log"
	"strings"
	"time"
)

// Year
// @description: 年排行榜
func Year() {
	groups, err := service.GetAllEnableChatRank()
	if err != nil {
		log.Printf("获取启用了聊天排行榜的群组失败, 错误信息: %v", err)
		return
	}

	for _, group := range groups {
		// 消息统计
		dealYear(group.Wxid)

		res, ok := config.Conf.Resource["wordcloud"]
		if !ok {
			continue
		}

		// 获取上周周数
		year := time.Now().Local().AddDate(0, 0, -1).Year()
		// 发送词云
		fileName := fmt.Sprintf("%d_%s.png", year, group.Wxid)
		utils.SendImage(group.Wxid, fmt.Sprintf(res.Path, fileName), 0)
	}
}

// dealYear
// @description: 处理年度排行榜
// @param gid
func dealYear(gid string) {
	notifyMsgs := []string{"#年度水群排行榜"}

	// 获取上周消息总数
	records, err := getRankData(gid, "year")
	if err != nil {
		log.Printf("获取去年消息排行失败, 错误信息: %v", err)
		return
	}
	log.Printf("去年消息总数: %+v", records)
	// 莫得消息，直接返回
	if len(records) == 0 {
		log.Printf("去年本群[%s]无对话记录", gid)
		return
	}
	// 计算消息总数
	var msgCount int64
	for _, v := range records {
		msgCount += v.Count
	}
	// 组装消息总数推送信息
	notifyMsgs = append(notifyMsgs, " ")
	notifyMsgs = append(notifyMsgs, "亲爱的群友们，新年已经悄悄来临，让我们一起迎接这充满希望和美好的时刻。在这个特殊的日子里，我要向你们致以最真挚的祝福。")
	notifyMsgs = append(notifyMsgs, "首先，我想对去年在群中表现出色、积极参与的成员们表示衷心的祝贺和感谢！你们的活跃与奉献让群聊更加充满了生机和活力。你们的贡献不仅仅是为了自己，更是为了我们整个群体的进步与成长。")
	notifyMsgs = append(notifyMsgs, "特此给去年年度活跃成员排行榜上的朋友们送上真诚的祝福。你们的热情、智慧和参与度，令我们很是钦佩。愿新的一年中，你们继续保持着你们的活力和激情，为群中带来更多的惊喜和启迪。")
	notifyMsgs = append(notifyMsgs, "对于那些未上榜的朋友们，我要说，你们也是我们群聊中非常重要的一部分。你们或许没有在排行榜上留下痕迹，但你们的存在和参与同样不可或缺。你们为群聊注入了新的思维和观点，为我们提供了不同的视角和见解。")
	notifyMsgs = append(notifyMsgs, "因此，我想特别鼓励未上榜的朋友们，继续发扬你们的热情和积极性。无论是在分享知识、讨论问题、还是互相支持鼓励，你们的贡献都是宝贵的。让我们共同创造一个更加活跃和有意义的群聊环境。")
	notifyMsgs = append(notifyMsgs, "最后，让我们一起迈向新的一年，相信自己的潜力和可能性，用我们的友谊和互助支持彼此。愿新的一年给我们带来更多的快乐、成功和成长。")
	notifyMsgs = append(notifyMsgs, fmt.Sprintf("祝福你们新年快乐！让我们一起迎接%d年的到来！", time.Now().Local().Year()))
	notifyMsgs = append(notifyMsgs, " ")
	notifyMsgs = append(notifyMsgs, fmt.Sprintf("🗣️ 去年本群 %d 位朋友共产生 %d 条发言", len(records), msgCount))
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

	log.Printf("排行榜: \n%s", strings.Join(notifyMsgs, "\n"))
	go utils.SendMessage(gid, "", strings.Join(notifyMsgs, "\n"), 0)
}
