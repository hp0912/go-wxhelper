package config

// Conf 配置
var Conf conf

// Config
// @description: 配置
type conf struct {
	Task     task                    `json:"task" yaml:"task"`         // 定时任务配置
	MySQL    mysql                   `json:"mysql" yaml:"mysql"`       // MySQL 配置
	Wechat   wechat                  `json:"wechat" yaml:"wechat"`     // 微信助手
	Ai       ai                      `json:"ai" yaml:"ai"`             // AI配置
	Resource map[string]resourceItem `json:"resource" yaml:"resource"` // 资源配置
}
