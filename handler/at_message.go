package handler

import (
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/sashabaranov/go-openai"
	"go-wechat/client"
	"go-wechat/common/current"
	"go-wechat/config"
	"go-wechat/entity"
	"go-wechat/model"
	"go-wechat/service"
	"go-wechat/types"
	"go-wechat/utils"
	"log"
	"regexp"
	"strings"
	"time"
)

// handleAtMessage
// @description: 处理At机器人的消息
// @param m
func handleAtMessage(m model.Message) {
	if !config.Conf.Ai.Enable {
		return
	}

	// 取出所有启用了AI的好友或群组
	var count int64
	client.MySQL.Model(&entity.Friend{}).Where("enable_ai IS TRUE").Where("wxid = ?", m.FromUser).Count(&count)
	if count < 1 {
		return
	}

	// 预处理一下发送的消息，用正则去掉@机器人的内容
	re := regexp.MustCompile(`@([^ | ]+)`)
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

	// 查询发信人前面几条文字信息，组装进来
	var oldMessages []entity.Message
	client.MySQL.Model(&entity.Message{}).
		Where("msg_id != ?", m.MsgId).
		Where("create_at >= DATE_SUB(NOW(),INTERVAL 30 MINUTE)").
		Where("from_user = ? AND group_user = ? AND display_full_content LIKE ?", m.FromUser, m.GroupUser, "%在群聊中@了你").
		Or("to_user = ? AND group_user = ?", m.FromUser, m.GroupUser).
		Order("create_at desc").
		Limit(4).Find(&oldMessages)
	// 翻转数组
	slice.Reverse(oldMessages)
	// 循环填充消息
	for _, message := range oldMessages {
		// 剔除@机器人的内容
		msgStr := message.Content
		matches = re.FindStringSubmatch(msgStr)
		if len(matches) > 0 {
			// 过滤掉第一个匹配到的
			msgStr = strings.Replace(msgStr, matches[0], "", 1)
		}
		// 填充消息
		role := openai.ChatMessageRoleUser
		if message.ToUser != current.GetRobotInfo().WxId {
			// 如果收信人不是机器人，表示这条消息是 AI 发的
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: msgStr,
		})
	}

	// 填充用户消息
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: m.Content,
	})

	// 配置模型
	chatModel := openai.GPT3Dot5Turbo0613
	if config.Conf.Ai.Model != "" {
		chatModel = config.Conf.Ai.Model
	}

	// 默认使用AI回复
	conf := openai.DefaultConfig(config.Conf.Ai.ApiKey)
	if config.Conf.Ai.BaseUrl != "" {
		conf.BaseURL = fmt.Sprintf("%s/v1", config.Conf.Ai.BaseUrl)
	}
	ai := openai.NewClientWithConfig(conf)
	resp, err := ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    chatModel,
			Messages: messages,
		},
	)

	if err != nil {
		log.Printf("OpenAI聊天发起失败: %v", err.Error())
		utils.SendMessage(m.FromUser, m.GroupUser, "AI炸啦~", 0)
		return
	}

	// 保存一下AI 返回的消息，消息 Id 使用传入 Id 的负数
	var replyMessage entity.Message
	replyMessage.MsgId = -m.MsgId
	replyMessage.CreateTime = int(time.Now().Local().Unix())
	replyMessage.CreateAt = time.Now().Local()
	replyMessage.Content = resp.Choices[0].Message.Content
	replyMessage.FromUser = current.GetRobotInfo().WxId // 发信人是机器人
	replyMessage.GroupUser = m.GroupUser                // 群成员
	replyMessage.ToUser = m.FromUser                    // 收信人是发信人
	replyMessage.Type = types.MsgTypeText
	service.SaveMessage(replyMessage) // 保存消息

	// 发送消息
	utils.SendMessage(m.FromUser, m.GroupUser, "\n"+resp.Choices[0].Message.Content, 0)
}
