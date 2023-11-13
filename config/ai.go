package config

// ai
// @description: AI配置
type ai struct {
	Enable  bool   `json:"enable" yaml:"enable"`   // 是否启用AI
	ApiKey  string `json:"apiKey" yaml:"apiKey"`   // API Key
	BaseUrl string `json:"baseUrl" yaml:"baseUrl"` // API地址
}
