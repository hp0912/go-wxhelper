package main

import (
	"github.com/gin-gonic/gin"
	"go-wechat/config"
	"go-wechat/initialization"
	"go-wechat/router"
	"go-wechat/tasks"
	"go-wechat/tcpserver"
	"go-wechat/utils"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	initialization.InitConfig()          // 初始化配置
	initialization.InitWechatRobotInfo() // 初始化机器人信息
	initialization.Plugin()              // 注册插件
	tasks.InitTasks()                    // 初始化定时任务
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

	// 自定义模板引擎函数
	app.SetFuncMap(template.FuncMap{
		"checkSwap": func(flag bool) string {
			if flag {
				return "swap-active"
			}
			return ""
		},
	})

	app.LoadHTMLGlob("views/*.html")
	app.Static("/assets", "./views/static")
	app.StaticFile("/favicon.ico", "./views/wechat.ico")
	// 404返回数据
	app.NoRoute(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/api") {
			ctx.String(404, "接口不存在")
			return
		}
		// 404直接跳转到首页
		ctx.Redirect(302, "/index.html")
	})
	app.NoMethod(func(ctx *gin.Context) {
		ctx.String(http.StatusMethodNotAllowed, "不支持的请求方式")
	})
	// 初始化路由
	router.Init(app)
	if err := app.Run(":8080"); err != nil {
		log.Panicf("服务启动失败：%v", err)
	}
}
