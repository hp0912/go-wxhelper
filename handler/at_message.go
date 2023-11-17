package handler

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/utils"
	"log"
	"regexp"
	"strings"
)

// handleAtMessage
// @description: 处理At机器人的消息
// @param m
func handleAtMessage(m entity.Message) {
	if !config.Conf.Ai.Enable {
		return
	}

	// 预处理一下发送的消息，用正则去掉@机器人的内容
	re := regexp.MustCompile(`@([^ ]+)`)
	matches := re.FindStringSubmatch(m.Content)
	if len(matches) > 0 {
		// 过滤掉第一个匹配到的
		m.Content = strings.Replace(m.Content, matches[0], "", 1)
	}

	// 组装消息体
	messages := make([]openai.ChatCompletionMessage, 0)
	if config.Conf.Ai.Personality != "" {
		// 填充人设
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.Conf.Ai.Personality,
		})
	}
	// 填充用户消息
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: m.Content,
	})

	// 配置模型
	model := openai.GPT3Dot5Turbo0613
	if config.Conf.Ai.Model != "" {
		model = config.Conf.Ai.Model
	}

	// 默认使用AI回复
	conf := openai.DefaultConfig(config.Conf.Ai.ApiKey)
	if config.Conf.Ai.BaseUrl != "" {
		conf.BaseURL = fmt.Sprintf("%s/v1", config.Conf.Ai.BaseUrl)
	}
	client := openai.NewClientWithConfig(conf)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("OpenAI聊天发起失败: %v", err.Error())
		utils.SendMessage(m.FromUser, m.GroupUser, "AI炸啦~", 0)
		return
	}

	// 发送消息
	utils.SendMessage(m.FromUser, m.GroupUser, "\n"+resp.Choices[0].Message.Content, 0)
}
