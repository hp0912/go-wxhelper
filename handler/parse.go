package handler

import (
	"encoding/json"
	"go-wechat/entity"
	"go-wechat/model"
	"go-wechat/service"
	"go-wechat/types"
	"go-wechat/utils"
	"log"
	"net"
	"strings"
	"time"
)

// Parse
// @description: 解析消息
// @param msg
func Parse(remoteAddr net.Addr, msg []byte) {
	var m model.Message
	if err := json.Unmarshal(msg, &m); err != nil {
		log.Printf("[%s]消息解析失败： %v", remoteAddr, err)
		log.Printf("[%s]消息内容： %d -> %v", remoteAddr, len(msg), string(msg))
		return
	}
	// 提取出群成员信息
	//groupUser := ""
	//msgStr := m.Content
	if strings.Contains(m.FromUser, "@") {
		// 群消息，处理一下消息和发信人
		groupUser := strings.Split(m.Content, "\n")[0]
		groupUser = strings.ReplaceAll(groupUser, ":", "")
		// 如果两个id一致，说明是系统发的
		if m.FromUser != groupUser {
			m.GroupUser = groupUser
		}
		// 用户的操作单独提出来处理一下
		m.Content = strings.Join(strings.Split(m.Content, "\n")[1:], "\n")
	}
	log.Printf("%s\n消息来源: %s\n群成员: %s\n消息类型: %v\n消息内容: %s", remoteAddr, m.FromUser, m.GroupUser, m.Type, m.Content)

	// 异步处理消息
	go func() {
		if m.IsNewUserJoin() {
			log.Printf("%s -> 开始迎新 -> %s", m.FromUser, m.Content)
			// 欢迎新成员
			go handleNewUserJoin(m)
		} else if m.IsAt() {
			// @机器人的消息
			go handleAtMessage(m)
		} else if !strings.Contains(m.FromUser, "@") && m.Type == types.MsgTypeText {
			// 私聊消息处理
			utils.SendMessage(m.FromUser, "", "暂未开启私聊AI", 0)
		}
	}()

	// 转换为结构体之后入库
	var ent entity.Message
	ent.MsgId = m.MsgId
	ent.CreateTime = m.CreateTime
	ent.CreateAt = time.Unix(int64(m.CreateTime), 0)
	ent.Content = m.Content
	ent.FromUser = m.FromUser
	ent.GroupUser = m.GroupUser
	ent.ToUser = m.ToUser
	ent.Type = m.Type
	ent.DisplayFullContent = m.DisplayFullContent
	ent.Raw = string(msg)

	go service.SaveMessage(ent)
}
