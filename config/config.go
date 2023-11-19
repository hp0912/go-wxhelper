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
