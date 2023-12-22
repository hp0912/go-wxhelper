package model

import (
	"encoding/xml"
	"go-wechat/types"
	"regexp"
	"strings"
)

// Message
// @description: 消息
type Message struct {
	MsgId              int64             `json:"msgId"`
	CreateTime         int               `json:"createTime"`
	Content            string            `json:"content"`
	DisplayFullContent string            `json:"displayFullContent"`
	FromUser           string            `json:"fromUser"`
	GroupUser          string            `json:"-"`
	MsgSequence        int               `json:"msgSequence"`
	Pid                int               `json:"pid"`
	Signature          string            `json:"signature"`
	ToUser             string            `json:"toUser"`
	Type               types.MessageType `json:"type"`
	Raw                string            `json:"raw"`
}

// systemMsgDataXml
// @description: 微信系统消息的xml结构
type systemMsgDataXml struct {
	SysMsg sysMsg `xml:"sysmsg"`
	Type   string `xml:"type,attr"`
}

// sysMsg
// @description: 消息主体
type sysMsg struct{}

func (m Message) IsGroup() bool {
	return strings.HasSuffix(m.FromUser, "@chatroom")
}

// IsPat
// @description: 是否是拍一拍消息
// @receiver m
// @return bool
func (m Message) IsPat() bool {
	// 解析xml
	var d systemMsgDataXml
	if err := xml.Unmarshal([]byte(m.Content), &d); err != nil {
		return false
	}

	return m.Type == types.MsgTypeRecalled && d.Type == "pat"
}

// IsRevokeMsg
// @description: 是否是撤回消息
// @receiver m
// @return bool
func (m Message) IsRevokeMsg() bool {
	// 解析xml
	var d systemMsgDataXml
	if err := xml.Unmarshal([]byte(m.Content), &d); err != nil {
		return false
	}

	return m.Type == types.MsgTypeRecalled && d.Type == "revokemsg"
}

// IsNewUserJoin
// @description: 是否是新人入群
// @receiver m
// @return bool
func (m Message) IsNewUserJoin() bool {
	sysFlag := m.Type == types.MsgTypeSys && strings.Contains(m.Content, "\"邀请\"") && strings.Contains(m.Content, "\"加入了群聊")
	if sysFlag {
		return true
	}
	// 解析另一种情况
	var d systemMsgDataXml
	if err := xml.Unmarshal([]byte(m.Content), &d); err != nil {
		return false
	}
	return m.Type == types.MsgTypeSys && d.Type == "delchatroommember"
}

// IsAt
// @description: 是否是At机器人的消息
// @receiver m
// @return bool
func (m Message) IsAt() bool {
	return strings.HasSuffix(m.DisplayFullContent, "在群聊中@了你")
}

// IsPrivateText
// @description: 是否是私聊消息
// @receiver m
// @return bool
func (m Message) IsPrivateText() bool {
	// 发信人不以@chatroom结尾且消息类型为文本
	return !strings.HasSuffix(m.FromUser, "chatroom") && m.Type == types.MsgTypeText
}

// CleanContentStartWith
// @description: 判断是否包含指定消息前缀
// @receiver m
// @param prefix
// @return bool
func (m Message) CleanContentStartWith(prefix string) bool {
	content := m.Content

	// 如果是@消息，过滤掉@的内容
	if m.IsAt() {
		re := regexp.MustCompile(`@([^ | ]+)`)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 0 {
			// 过滤掉第一个匹配到的
			content = strings.Replace(content, matches[0], "", 1)
		}
	}

	// 去掉最前面的空格
	content = strings.TrimLeft(content, " ")
	content = strings.TrimLeft(content, " ")

	return strings.HasPrefix(content, prefix)
}
