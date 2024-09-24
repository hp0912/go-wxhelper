package command

import (
	"go-wechat/utils"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// PoisonousChickenSoup
// @description: 毒鸡汤
// @param userId string 发信人
func PoisonousChickenSoup(userId string) {
	msg := ""

	res := resty.New()
	resp, err := res.R().
		Get("https://api.pearktrue.cn/api/dujitang")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取毒鸡汤失败: %v", err)
		msg = "获取毒鸡汤失败"
	}
	msg = resp.String()

	// 发送消息
	utils.SendMessage(userId, "", msg, 0)
}
