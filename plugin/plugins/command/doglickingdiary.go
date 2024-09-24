package command

import (
	"go-wechat/utils"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// DogLickingDiary
// @description: 舔狗日记
// @param userId string 发信人
func DogLickingDiary(userId string) {
	msg := ""

	res := resty.New()
	resp, err := res.R().
		Get("https://api.pearktrue.cn/api/jdyl/tiangou.php")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取舔狗日记失败: %v", err)
		msg = "获取舔狗日记失败"
	}
	msg = resp.String()

	// 发送消息
	utils.SendMessage(userId, "", msg, 0)
}
