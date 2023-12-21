package command

import (
	"github.com/go-resty/resty/v2"
	"go-wechat/utils"
	"log"
)

// KfcCrazyThursdayCmd
// @description: 肯德基疯狂星期四文案
// @param userId string 发信人
func KfcCrazyThursdayCmd(userId string) {
	// 接口调用
	str := kfcApi1()
	if str == "" {
		str = kfcApi2()
	}
	if str == "" {
		str = "文案获取失败"
	}

	// 发送消息
	utils.SendMessage(userId, "", str, 0)
}

// kfcApi1
// @description: 肯德基疯狂星期四文案接口1
// @return string
func kfcApi1() string {
	res := resty.New()
	resp, err := res.R().
		Post("https://api.jixs.cc/api/wenan-fkxqs/index.php")
	if err != nil {
		log.Panicf("KFC接口1文案获取失败: %s", err.Error())
	}
	log.Printf("KFC接口1文案获取结果: %s", resp.String())
	return resp.String()
}

// kfcApi2
// @description: 肯德基疯狂星期四文案接口2
// @return string
func kfcApi2() string {
	type result struct {
		Code int    `json:"code"`
		Text string `json:"text"`
		Data struct {
			Msg string `json:"msg"`
		} `json:"data"`
	}

	var resData result

	res := resty.New()
	resp, err := res.R().
		SetResult(&resData).
		Post("https://api.jixs.cc/api/wenan-fkxqs/index.php")
	if err != nil {
		log.Panicf("KFC接口2文案获取失败: %s", err.Error())
	}
	log.Printf("KFC接口2文案获取结果: %s", resp.String())
	if resData.Data.Msg != "" {
		return resData.Data.Msg
	}
	return resp.String()
}
