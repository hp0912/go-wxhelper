package types

import "fmt"

type MessageType int

const (
	MsgTypeText           MessageType = 1     // 文本消息
	MsgTypeImage          MessageType = 3     // 图片消息
	MsgTypeVoice          MessageType = 34    // 语音消息
	MsgTypeVerify         MessageType = 37    // 认证消息
	MsgTypePossibleFriend MessageType = 40    // 好友推荐消息
	MsgTypeShareCard      MessageType = 42    // 名片消息
	MsgTypeVideo          MessageType = 43    // 视频消息
	MsgTypeEmoticon       MessageType = 47    // 表情消息
	MsgTypeLocation       MessageType = 48    // 地理位置消息
	MsgTypeApp            MessageType = 49    // APP消息
	MsgTypeVoip           MessageType = 50    // VOIP消息
	MsgTypeVoipNotify     MessageType = 52    // VOIP结束消息
	MsgTypeVoipInvite     MessageType = 53    // VOIP邀请
	MsgTypeMicroVideo     MessageType = 62    // 小视频消息
	MsgTypeSys            MessageType = 10000 // 系统消息
	MsgTypeRecalled       MessageType = 10002 // 消息撤回
)

var MessageTypeMap = map[MessageType]string{
	MsgTypeText:           "文本消息",
	MsgTypeImage:          "图片消息",
	MsgTypeVoice:          "语音消息",
	MsgTypeVerify:         "认证消息",
	MsgTypePossibleFriend: "好友推荐消息",
	MsgTypeShareCard:      "名片消息",
	MsgTypeVideo:          "视频消息",
	MsgTypeEmoticon:       "表情消息",
	MsgTypeLocation:       "地理位置消息",
	MsgTypeApp:            "APP消息",
	MsgTypeVoip:           "VOIP消息",
	MsgTypeVoipNotify:     "VOIP结束消息",
	MsgTypeVoipInvite:     "VOIP邀请",
	MsgTypeMicroVideo:     "小视频消息",
	MsgTypeSys:            "系统消息",
	MsgTypeRecalled:       "消息撤回",
}

func (mt MessageType) String() string {
	if msg, ok := MessageTypeMap[mt]; ok {
		return msg
	}
	return fmt.Sprintf("未知消息类型(%d)", mt)
}
