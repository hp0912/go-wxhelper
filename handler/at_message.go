package handler

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/utils"
	"log"
)

// handleAtMessage
// @description: 处理At机器人的消息
// @param m
func handleAtMessage(m entity.Message) {
	if !config.Conf.Ai.Enable {
		return
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
		log.Printf("OpenAI聊天发起失败: %v", err.Error())
		utils.SendMessage(m.FromUser, m.GroupUser, "AI炸啦~", 0)
		return
	}

	// 发送消息
	utils.SendMessage(m.FromUser, m.GroupUser, resp.Choices[0].Message.Content, 0)
}
