package utils

import (
	"encoding/json"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/go-resty/resty/v2"
	"go-wechat/config"
	"log"
	"net"
	"strconv"
	"strings"
)

// ClearCallback
// @description: 清理微信HOOK回调
func ClearCallback() {
	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		Post(config.Conf.Wechat.GetURL("/api/unhookSyncMsg"))
	if err != nil {
		log.Panicf("清理微信HOOK回调失败: %s", err.Error())
	}
	log.Printf("清理微信HOOK回调结果: %s", resp.String())
}

// SetCallback
// @description: 设置微信HOOK回调
// @param host
func SetCallback(userHost string) {
	// 获取本机IP地址
	host := userHost
	if userHost == "auto" {
		host = net.ParseIP(netutil.GetInternalIp()).String()
	}

	port := 19099
	if userHost != "" {
		uh := strings.Split(strings.TrimSpace(userHost), ":")
		host = uh[0]
		if len(uh) == 2 {
			port, _ = strconv.Atoi(uh[1])
		}
	}

	// 组装参数
	param := map[string]any{
		"port":       port, // socket端口
		"ip":         host, // socketIP
		"url":        "",   // http接口地址
		"timeout":    3000, // 超时毫秒数
		"enableHttp": 0,    // 是否使用http接口
	}
	pbs, _ := json.Marshal(param)
	log.Printf("设置微信HOOK回调参数: %s", string(pbs))

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		Post(config.Conf.Wechat.GetURL("/api/hookSyncMsg"))
	if err != nil {
		log.Panicf("设置微信HOOK回调失败: %s", err.Error())
	}
	log.Printf("设置微信HOOK回调结果: %s", resp.String())
}
