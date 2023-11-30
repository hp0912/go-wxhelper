package app

import (
	"github.com/gin-gonic/gin"
	"go-wechat/service"
	"net/http"
)

type getGroupUser struct {
	GroupId string `json:"groupId" form:"groupId" binding:"required"` // 群Id
}

// GetGroupUsers
// @description: 获取群成员列表
// @param ctx
func GetGroupUsers(ctx *gin.Context) {
	var p getGroupUser
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	// 查询数据
	records, err := service.GetGroupUsersByGroupId(p.GroupId)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "查询失败: %s", err.Error())
		return
	}
	// 暂时先就这样写着，跑通了再改
	ctx.JSON(http.StatusOK, records)
}
