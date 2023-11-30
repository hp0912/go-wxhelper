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
}
