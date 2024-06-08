package command

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/utils"
	"time"
)

// WordCloud
// @description: 词云
// @param userId string 发信人
func WordCloud(userId string) {
	res, ok := config.Conf.Resource["wordcloud"]
	if !ok {
		return
	}
	// 获取昨日日期
	yd := time.Now().Local().AddDate(0, 0, -1).Format("20060102")
	// 发送词云
	fileName := fmt.Sprintf("%s_%s.png", yd, userId)
	utils.SendImage(userId, fmt.Sprintf(res.Path, fileName), 0)
}
