package config

var Conf Config

// Config
// @description: 配置
type Config struct {
	Task   task   `json:"task" yaml:"task"`     // 定时任务配置
	MySQL  mysql  `json:"mysql" yaml:"mysql"`   // MySQL 配置
	Wechat wechat `json:"wechat" yaml:"wechat"` // 微信助手
	Ai     ai     `json:"ai" yaml:"ai"`         // AI配置
}

// task
// @description: 定时任务
type task struct {
	Enable      bool `json:"enable" yaml:"enable"` // 是否启用
	SyncFriends struct {
		Enable bool   `json:"enable" yaml:"enable"` // 是否启用
		Cron   string `json:"cron" yaml:"cron"`     // 定时任务表达式
	} `json:"syncFriends" yaml:"syncFriends"` // 同步好友
	WaterGroup struct {
		Enable    bool     `json:"enable" yaml:"enable"`       // 是否启用
		Cron      string   `json:"cron" yaml:"cron"`           // 定时任务表达式
		Groups    []string `json:"groups" yaml:"groups"`       // 启用的群Id
		Blacklist []string `json:"blacklist" yaml:"blacklist"` // 黑名单
	} `json:"waterGroup" yaml:"waterGroup"` // 水群排行榜
}
