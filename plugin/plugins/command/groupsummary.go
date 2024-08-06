package command

import (
	"go-wechat/service"
	"go-wechat/tasks/summary"
	"go-wechat/utils"
	"log"
)

// GroupSummary
// @description: 群聊总结
// @param userId string 发信人
func GroupSummary(groupId string) {
	var msg string

	group, err := service.GetFriendInfoById(groupId)
	if err != nil {
		log.Printf("获取群[%s]信息失败, 错误信息: %v", groupId, err)
		utils.SendMessage(groupId, "获取群聊信息失败", msg, 0)
		return
	}

	summary.GroupSummary(group)
}
