package router

import (
	"go-wechat/app"

	"github.com/gin-gonic/gin"
)

// Init
// @description: 初始化路由
// @param g
func Init(g *gin.Engine) {
	g.GET("/", func(ctx *gin.Context) {
		// 重定向到index.html
		ctx.Redirect(302, "/index.html")
	})

	g.GET("/index.html", app.Index)         // 首页
	g.GET("/friend.html", app.Friend)       // 好友列表
	g.GET("/group.html", app.Group)         // 群组列表
	g.GET("/assistant.html", app.Assistant) // AI角色

	g.GET("/404.html", app.PageNotFound) // 群组列表

	// 接口
	api := g.Group("/api")
	api.PUT("/ai/status", app.ChangeEnableAiStatus)               // 修改是否开启AI状态
	api.POST("/ai/model", app.ChangeUseAiModel)                   // 修改使用的AI模型
	api.POST("/ai/assistant", app.ChangeUseAiAssistant)           // 修改使用的AI助手
	api.PUT("/welcome/status", app.ChangeEnableWelcomeStatus)     // 修改是否开启迎新状态
	api.PUT("/command/status", app.ChangeEnableCommandStatus)     // 修改是否开启指令状态
	api.PUT("/news/status", app.ChangeEnableNewsStatus)           // 修改是否开启早报状态
	api.PUT("/grouprank/status", app.ChangeEnableGroupRankStatus) // 修改是否开启水群排行榜状态
	api.PUT("/grouprank/skip", app.ChangeSkipGroupRankStatus)     // 修改是否跳过水群排行榜状态
	api.GET("/group/users", app.GetGroupUsers)                    // 获取群成员列表
	api.PUT("/summary/status", app.ChangeEnableSummaryStatus)     // 修改是否开启群聊总结状态
	api.PUT("/clearmembers", app.AutoClearMembers)                // 自动清理群成员

	api.POST("/assistant", app.SaveAssistant) // 保存AI助手
}
