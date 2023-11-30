package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-wechat/service"
	"net/http"
)

// Index
// @description: 首页
// @param ctx
func Index(ctx *gin.Context) {
	var result = gin.H{
		"msg": "success",
	}
	// 取出所有好友列表
	friends, groups, err := service.GetAllFriend()
	if err != nil {
		result["msg"] = fmt.Sprintf("数据获取失败: %s", err.Error())
	}
	result["friends"] = friends
	result["groups"] = groups
	// 渲染页面
	ctx.HTML(http.StatusOK, "index.html", result)
}
