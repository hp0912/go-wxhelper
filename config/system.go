package config

// 系统配置
type system struct {
	NewFriendNotify newFriendNotify `json:"newFriendNotify" yaml:"newFriendNotify"` // 新好友通知
	DefaultRule     defaultRule     `json:"defaultRule" yaml:"defaultRule"`         // 默认规则
}

// 添加新好友或群之后通知给指定的人
type newFriendNotify struct {
	Enable bool     `json:"enable" yaml:"enable"` // 是否启用
	ToUser []string `json:"toUser" yaml:"toUser"` // 通知给谁
}

// 默认规则
type defaultRule struct {
	Ai       bool `json:"ai" yaml:"ai"`             // 是否启用AI
	ChatRank bool `json:"chatRank" yaml:"chatRank"` // 是否启用聊天排行榜
	Summary  bool `json:"summary" yaml:"summary"`   // 是否启用聊天总结
	Welcome  bool `json:"welcome" yaml:"welcome"`   // 是否启用欢迎新成员
}
