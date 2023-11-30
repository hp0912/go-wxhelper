package vo

import (
	"go-wechat/common/types"
)

// FriendItem
// @description: 好友列表数据
type FriendItem struct {
	CustomAccount  string         // 微信号
	Nickname       string         // 昵称
	Pinyin         string         // 昵称拼音大写首字母
	PinyinAll      string         // 昵称全拼
	Wxid           string         // 微信原始Id
	EnableAi       bool           // 是否使用AI
	EnableChatRank bool           // 是否使用聊天排行
	IsOk           bool           // 是否还在通讯库(群聊是要还在群里也算)
	LastActiveTime types.DateTime // 最后活跃时间
}
