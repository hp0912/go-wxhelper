package tasks

import (
	"github.com/go-co-op/gocron"
	"go-wechat/config"
	"log"
	"time"
)

// InitTasks
// @description: 初始化定时任务
func InitTasks() {
	// 没启用直接返回
	if !config.Conf.Task.Enable {
		log.Println("未启用定时任务")
		return
	}
	// 定时任务发送消息
	s := gocron.NewScheduler(time.Local)

	// 水群排行
	if config.Conf.Task.WaterGroup.Enable {
		log.Printf("水群排行任务已启用，执行表达式: %s", config.Conf.Task.WaterGroup.Cron)
		_, _ = s.Cron(config.Conf.Task.WaterGroup.Cron).Do(yesterday)
	}

	// 更新好友列表
	if config.Conf.Task.SyncFriends.Enable {
		log.Printf("更新好友列表任务已启用，执行表达式: %s", config.Conf.Task.SyncFriends.Cron)
		_, _ = s.Cron(config.Conf.Task.SyncFriends.Cron).Do(syncFriends)
	}

	// 开启定时任务
	s.StartAsync()
	log.Println("定时任务初始化成功")
}
