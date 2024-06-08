package command

import (
	"go-wechat/utils"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ZhaNan
// @description: 渣男语录
// @param userId string 发信人
func ZhaNan(userId string) {
	type result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Text string `json:"text"`
	}

	var resData result

	msg := ""

	res := resty.New()
	resp, err := res.R().
		SetResult(&resData).
		Get("https://api.pearktrue.cn/api/random/zhanan?type=json")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取渣男语录失败: %v", err)
		msg = "获取渣男语录失败"
	} else if resData.Text != "" {
		msg = resData.Text
	} else {
		msg = resp.String()
	}

	// 发送消息
	utils.SendMessage(userId, "", msg, 0)
}
