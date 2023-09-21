package tasks

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

// InitTasks
// @description: 初始化定时任务
func InitTasks() {
	// 定时任务发送消息
	s := gocron.NewScheduler(time.Local)

	// 每天早上九点半发送前一天的水群排行
	_, _ = s.Every(1).Day().At("09:30").Do(yesterday)

	// 每小时更新一次好友列表
	_, _ = s.Every(1).Hour().Do(syncFriends)

	// 开启定时任务
	s.StartAsync()
	log.Println("定时任务初始化成功")
}
