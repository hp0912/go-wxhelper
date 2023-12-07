package initialization

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go-wechat/client"
	"go-wechat/config"
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
	if err := vp.Unmarshal(&config.Conf); err != nil {
		log.Panicf("配置文件解析失败: %v", err)
	}
	log.Printf("配置文件解析完成: %+v", config.Conf)
	if !config.Conf.Wechat.Check() {
		log.Panicf("微信HOOK配置缺失")
	}
	// 初始化数据库连接
	client.InitMySQLClient()
	//redis.Init()

	// 下面的代码是配置变动之后自动刷新的
	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		// 绑定配置文件
		if err := vp.Unmarshal(&config.Conf); err != nil {
			log.Printf("配置文件更新失败: %v", err)
		} else {
			if !config.Conf.Wechat.Check() {
				log.Panicf("微信HOOK配置缺失")
			}
			// 初始化数据库连接
			client.InitMySQLClient()
			//redis.Init()
		}
	})
}
