package config

import "strings"

// wxHelper
// @description: 微信助手
type wechat struct {
	Host            string   `json:"host" yaml:"host"`                       // 接口地址
	VncUrl          string   `json:"vncUrl" yaml:"vncUrl"`                   // vnc页面地址
	AutoSetCallback bool     `json:"autoSetCallback" yaml:"autoSetCallback"` // 是否自动设置回调地址
	Callback        string   `json:"callback" yaml:"callback"`               // 回调地址
	Forward         []string `json:"forward" yaml:"forward"`                 // 转发地址
}

// Check
// @description: 检查配置是否可用
// @receiver w
// @return bool
func (w wechat) Check() bool {
	if w.Host == "" {
		return false
	}
	if w.AutoSetCallback && w.Callback == "" {
		return false
	}
	return true
}

func (w wechat) GetURL(uri string) string {
	host := w.Host
	if !strings.HasPrefix(w.Host, "http://") {
		host = "http://" + w.Host
	}
	return host + uri
}
