package model

import "go-wechat/types"

// Message
// @description: 消息
type Message struct {
	MsgId              int64             `json:"msgId" gorm:"primarykey"`
	CreateTime         int               `json:"createTime"`
	Content            string            `json:"content"`
	DisplayFullContent string            `json:"displayFullContent" gorm:"-"`
	FromUser           string            `json:"fromUser"`
	MsgSequence        int               `json:"msgSequence"`
	Pid                int               `json:"pid"`
	Signature          string            `json:"signature"`
	ToUser             string            `json:"toUser"`
	Type               types.MessageType `json:"type"`
}
