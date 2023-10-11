package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// 配置管理工具
var vp *viper.Viper

// InitConfig
// @description: 初始化配置
func InitConfig() {
	vp = viper.New()
	vp.AddConfigPath(".")      // 设置配置文件路径
	vp.SetConfigName("config") // 设置配置文件名
	vp.SetConfigType("yaml")   // 设置配置文件类型
	// 读取配置文件
	if err := vp.ReadInConfig(); err != nil {
		log.Panicf("读取配置文件失败: %v", err)
	}
	// 绑定配置文件
	if err := vp.Unmarshal(&Conf); err != nil {
		log.Panicf("配置文件解析失败: %v", err)
	}
	log.Printf("配置文件解析完成: %+v", Conf)
	// 初始化数据库连接
	//db.Init()
	//redis.Init()

	// 下面的代码是配置变动之后自动刷新的
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		// 绑定配置文件
		if err := vp.Unmarshal(&Conf); err != nil {
			log.Printf("配置文件更新失败: %v", err)
		} else {
			// 初始化数据库连接
			//db.Init()
			//redis.Init()
		}
	})
}
