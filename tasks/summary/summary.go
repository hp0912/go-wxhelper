package summary

import (
	"context"
	"fmt"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/service"
	"go-wechat/utils"
	"go-wechat/vo"
	"log"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// GroupSummary
// @param group entity.Friend 群信息
// @param isCron bool 是否由定时任务触发
// @param condition string 额外查询条件
// @description: AI总结，单个群聊记录
func GroupSummary(group entity.Friend, isCron bool, condition string) {
	var err error
	// 获取对话记录
	var records []vo.TextMessageItem
	if records, err = service.GetTextMessagesById(group.Wxid, service.GetMessageOption{IsCron: isCron, Condition: condition}); err != nil {
		log.Printf("获取群[%s]对话记录失败, 错误信息: %v", group.Wxid, err)
		return
	}
	if isCron {
		if len(records) < 100 {
			log.Printf("群[%s]对话记录不足100条，跳过总结", group.Wxid)
			return
		}
	} else {
		if len(records) < 10 {
			log.Printf("群[%s]对话记录不足10条，跳过总结", group.Wxid)
			return
		}
	}

	// 组装对话记录为字符串
	var content []string
	for _, record := range records {
		content = append(content, fmt.Sprintf(`{"%s": "%s"}--end--`, record.Nickname, strings.ReplaceAll(record.Message, "\n", "。。")))
	}

	msgTmp := `请帮我总结一下一下的群聊内容的梗概，生成的梗概需要尽可能详细，需要带上一些聊天关键信息，并且带上群友名字。
注意，他们可能是多个话题，请仔细甄别。
每一行代表一个人的发言，每一行的的格式为： {"{nickname}": "{content}"}--end--

群名称: %s
聊天记录如下:
%s
`

	msg := fmt.Sprintf(msgTmp, group.Nickname, strings.Join(content, "\n"))

	// AI总结
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}

	// 默认使用AI回复
	conf := openai.DefaultConfig(config.Conf.Ai.ApiKey)
	if config.Conf.Ai.BaseUrl != "" {
		conf.BaseURL = fmt.Sprintf("%s/v1", config.Conf.Ai.BaseUrl)
	}
	ai := openai.NewClientWithConfig(conf)
	var resp openai.ChatCompletionResponse
	resp, err = ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    config.Conf.Ai.SummaryModel,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("群聊记录总结失败: %v", err.Error())
		return
	}

	// 返回消息为空
	if resp.Choices[0].Message.Content == "" {
		return
	}
	tips := "#昨日消息总结\n又是一天过去了，让我们一起来看看昨儿群友们都聊了什么有趣的话题吧~"
	if !isCron {
		tips = "#消息总结\n一会儿不看群，你们就群聊消息999+，让我们一起来看看群友们都聊了什么有趣的话题吧~"
	}
	replyMsg := fmt.Sprintf("%s\n\n%s", tips, resp.Choices[0].Message.Content)
	// log.Printf("群[%s]对话记录总结成功，总结内容: %s", group.Wxid, replyMsg)
	utils.SendMessage(group.Wxid, "", replyMsg, 0)
}

// AiSummary
// @description: AI总结群聊记录
func AiSummary() {
	groups, err := service.GetAllEnableSummary()
	if err != nil {
		log.Printf("获取启用了聊天排行榜的群组失败, 错误信息: %v", err)
		return
	}

	for _, group := range groups {
		GroupSummary(group, true, "")
		// 休眠一秒，防止频繁发送
		time.Sleep(time.Second)
	}
}
