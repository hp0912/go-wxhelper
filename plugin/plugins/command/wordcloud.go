package command

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/utils"
	"log"
	"time"
)

// WordCloud
// @description: 词云
// @param userId string 发信人
func WordCloud(userId string) {
	res, ok := config.Conf.Resource["wordcloud"]
	if !ok {
		log.Printf("获取词云路径失败~")
		return
	}
	// 获取昨日日期
	yd := time.Now().Local().AddDate(0, 0, -1).Format("20060102")
	fileName := fmt.Sprintf("%s_%s.png", yd, userId)

	// 发送词云
	log.Printf("群ID: %s  词云路径: %s~", userId, fmt.Sprintf(res.Path, fileName))
	utils.SendImage(userId, fmt.Sprintf(res.Path, fileName), 0)
}
