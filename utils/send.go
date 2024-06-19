package utils

import (
	"encoding/json"
	"fmt"
	"go-wechat/common/current"
	"go-wechat/config"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
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

	// 接口地址
	apiUrl := config.Conf.Wechat.GetURL("/api/sendTextMsg")
	if atId != "" {
		apiUrl = config.Conf.Wechat.GetURL("/api/sendAtText")
		param = map[string]any{
			"chatRoomId": toId,
			"wxids":      atId,
			"msg":        msg, // 消息
		}
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post(apiUrl)
	if err != nil {
		log.Printf("发送文本消息失败: %s", err.Error())
		// 休眠五秒后重新发送
		time.Sleep(5 * time.Second)
		SendMessage(toId, atId, msg, retryCount+1)
	}
	log.Printf("发送文本消息结果: %s", resp.String())
}

// SendImage
// @description: 发送图片
// @param toId string 群或者好友Id
// @param imgPath string 图片路径
// @param retryCount int 重试次数
func SendImage(toId, imgPath string, retryCount int) {
	if retryCount > 5 {
		log.Printf("重试五次失败，停止发送")
		return
	}

	// 组装参数
	param := map[string]any{
		"wxid":      toId,    // 群或好友Id
		"imagePath": imgPath, // 图片地址
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post(config.Conf.Wechat.GetURL("/api/sendImagesMsg"))
	if err != nil {
		log.Printf("发送图片消息失败: %s", err.Error())
		// 休眠五秒后重新发送
		time.Sleep(5 * time.Second)
		SendImage(toId, imgPath, retryCount+1)
	}
	log.Printf("发送图片消息结果: %s", resp.String())
}

// SendEmotion
// @description: 发送自定义表情包
// @param toId string 群或者好友Id
// @param emotionHash string 表情包hash(md5值)
// @param retryCount int 重试次数
func SendEmotion(toId, emotionHash string, retryCount int) {
	if retryCount > 5 {
		log.Printf("重试五次失败，停止发送")
		return
	}

	// 组装表情包本地地址
	// 规则：机器人数据目录\FileStorage\CustomEmotion\表情包hash前两位\表情包hash
	emotionPath := fmt.Sprintf("%sFileStorage\\CustomEmotion\\%s\\%s",
		current.GetRobotInfo().CurrentDataPath, emotionHash[:2], emotionHash)

	// 组装参数
	param := map[string]any{
		"wxid":     toId,        // 群或好友Id
		"filePath": emotionPath, // 图片地址
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post(config.Conf.Wechat.GetURL("/api/sendCustomEmotion"))
	if err != nil {
		log.Printf("发送表情包消息失败: %s", err.Error())
		// 休眠五秒后重新发送
		time.Sleep(5 * time.Second)
		SendImage(toId, emotionHash, retryCount+1)
	}
	log.Printf("发送表情包消息结果: %s", resp.String())
}

// DeleteGroupMember
// @description: 删除群成员
// @param chatRoomId 群Id
// @param memberIds 成员id,用','分隔
func DeleteGroupMember(chatRoomId, memberIds string, retryCount int) {
	if retryCount > 5 {
		log.Printf("重试五次失败，停止发送")
		return
	}

	// 组装参数
	param := map[string]any{
		"chatRoomId": chatRoomId, // 群Id
		"memberIds":  memberIds,  // 成员id
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post(config.Conf.Wechat.GetURL("/api/delMemberFromChatRoom"))
	if err != nil {
		log.Printf("删除群成员失败: %s", err.Error())
		// 休眠五秒后重新发送
		time.Sleep(5 * time.Second)
		DeleteGroupMember(chatRoomId, memberIds, retryCount+1)
	}
	log.Printf("删除群成员结果: %s", resp.String())
	// 这个逼接口要调用两次，第一次调用成功，第二次调用才会真正删除
	DeleteGroupMember(chatRoomId, memberIds, 5)
}
