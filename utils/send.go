package utils

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

// SendMessage
// @description: 发送消息
// @param toId
// @param atId
// @param msg
func SendMessage(toId, atId, msg string) {
	// 组装参数
	param := map[string]any{
		"wxid": toId, // 群或好友Id
		"msg":  msg,  // 消息
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post("http://10.0.0.73:19088/api/sendTextMsg")
	if err != nil {
		log.Printf("发送文本消息失败: %s", err.Error())
		return
	}
	log.Printf("发送文本消息结果: %s", resp.String())
}
