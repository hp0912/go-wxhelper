package main

import (
	"github.com/gin-gonic/gin"
	"go-wechat/config"
	"go-wechat/initialization"
	"go-wechat/router"
	"go-wechat/tasks"
	"go-wechat/tcpserver"
	"go-wechat/utils"
	"log"
	"time"
)

func init() {
	initialization.InitConfig()
	tasks.InitTasks()
}

func main() {
	// 如果启用了自动配置回调，就设置一下
	if config.Conf.Wechat.AutoSetCallback {
		utils.ClearCallback()
		time.Sleep(500 * time.Millisecond) // 休眠五百毫秒再设置
		utils.SetCallback(config.Conf.Wechat.Callback)
	}

	// 启动TCP服务
	go tcpserver.Start()

	// 启动HTTP服务
	app := gin.Default()
	app.LoadHTMLGlob("views/*.html")
	app.Static("/assets", "./views/static")
	app.StaticFile("/favicon.ico", "./views/wechat.ico")
	// 404返回数据
	app.NoRoute(func(ctx *gin.Context) {
		// 404直接跳转到首页
		ctx.Redirect(302, "/index.html")
	})
	// 初始化路由
	router.Init(app)
	if err := app.Run(":8080"); err != nil {
		log.Panicf("服务启动失败：%v", err)
	}
}
