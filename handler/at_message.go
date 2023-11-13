package handler

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/utils"
)

// handleAtMessage
// @description: 处理At机器人的消息
// @param m
func handleAtMessage(m entity.Message) {
	if !config.Conf.Ai.Enable {
		return
	}
	// 默认使用AI回复
	client := openai.NewClient(config.Conf.Ai.ApiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0613,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: m.Content,
				},
			},
		},
	)

	if err != nil {
		utils.SendMessage(m.FromUser, m.GroupUser, "AI炸啦~", 0)
		return
	}

	// 发送消息
	utils.SendMessage(m.FromUser, m.GroupUser, resp.Choices[0].Message.Content, 0)
}
