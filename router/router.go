package router

import (
	"github.com/gin-gonic/gin"
	"go-wechat/app"
)

// Init
// @description: 初始化路由
// @param g
func Init(g *gin.Engine) {
	g.GET("/", func(ctx *gin.Context) {
		// 重定向到index.html
		ctx.Redirect(302, "/index.html")
	})

	g.GET("/index.html", app.Index) // 首页
	g.GET("/test.html", func(ctx *gin.Context) {
		ctx.HTML(200, "test.html", nil)
	})

	// 接口
	api := g.Group("/api")
	api.PUT("/ai/status", app.ChangeEnableAiStatus)               // 修改是否开启AI状态
	api.PUT("/welcome/status", app.ChangeEnableWelcomeStatus)     // 修改是否开启迎新状态
	api.PUT("/grouprank/status", app.ChangeEnableGroupRankStatus) // 修改是否开启水群排行榜状态
	api.PUT("/grouprank/skip", app.ChangeSkipGroupRankStatus)     // 修改是否跳过水群排行榜状态
	api.GET("/group/users", app.GetGroupUsers)                    // 获取群成员列表
}
