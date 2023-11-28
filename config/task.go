package config

// task
// @description: 定时任务
type task struct {
	Enable      bool        `json:"enable" yaml:"enable"`           // 是否启用
	SyncFriends syncFriends `json:"syncFriends" yaml:"syncFriends"` // 同步好友
	WaterGroup  waterGroup  `json:"waterGroup" yaml:"waterGroup"`   // 水群排行榜
}

// syncFriends
// @description: 同步好友
type syncFriends struct {
	Enable bool   `json:"enable" yaml:"enable"` // 是否启用
	Cron   string `json:"cron" yaml:"cron"`     // 定时任务表达式
}

// waterGroup
// @description: 水群排行榜
type waterGroup struct {
	Enable bool           `json:"enable" yaml:"enable"` // 是否启用
	Cron   waterGroupCron `json:"cron" yaml:"cron"`     // 定时任务表达式
}

// waterGroupCron
// @description: 水群排行榜定时任务
type waterGroupCron struct {
	Yesterday string `json:"yesterday" yaml:"yesterday"` // 昨日排行榜
	Week      string `json:"week" yaml:"week"`           // 周排行榜
	Month     string `json:"month" yaml:"month"`         // 月排行榜
}
