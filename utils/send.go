package utils

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"go-wechat/config"
	"log"
	"time"
)

// SendMessage
// @description: 发送消息
// @param toId
// @param atId
// @param msg
func SendMessage(toId, atId, msg string, retryCount int) {
	if retryCount > 5 {
		log.Printf("重试五次失败，停止发送")
		return
	}
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
		Post(config.Conf.Wechat.GetURL("/api/sendTextMsg"))
	if err != nil {
		log.Printf("发送文本消息失败: %s", err.Error())
		// 休眠五秒后重新发送
		time.Sleep(5 * time.Second)
		SendMessage(toId, atId, msg, retryCount+1)
	}
	log.Printf("发送文本消息结果: %s", resp.String())
}
