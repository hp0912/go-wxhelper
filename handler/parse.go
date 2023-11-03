package handler

import (
	"encoding/json"
	"go-wechat/entity"
	"go-wechat/model"
	"go-wechat/service"
	"go-wechat/types"
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
	groupUser := ""
	msgStr := m.Content
	if strings.Contains(m.FromUser, "@") {
		switch m.Type {
		case types.MsgTypeRecalled:
			// 消息撤回
		case types.MsgTypeSys:
			// 系统消息
			go handleSysMessage(m)
		default:
			// 默认消息处理
			groupUser = strings.Split(m.Content, "\n")[0]
			groupUser = strings.ReplaceAll(groupUser, ":", "")

			// 文字消息单独提出来处理一下
			msgStr = strings.Join(strings.Split(m.Content, "\n")[1:], "\n")
		}
	}
	log.Printf("%s\n消息来源: %s\n群成员: %s\n消息类型: %v\n消息内容: %s", remoteAddr, m.FromUser, groupUser, m.Type, msgStr)

	// 转换为结构体之后入库
	var ent entity.Message
	ent.MsgId = m.MsgId
	ent.CreateTime = m.CreateTime
	ent.CreateAt = time.Unix(int64(m.CreateTime), 0)
	ent.Content = msgStr
	ent.FromUser = m.FromUser
	ent.GroupUser = groupUser
	ent.ToUser = m.ToUser
	ent.Type = m.Type
	ent.DisplayFullContent = m.DisplayFullContent
	ent.Raw = string(msg)

	go service.SaveMessage(ent)
}
