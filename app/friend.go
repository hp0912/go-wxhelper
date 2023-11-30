package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// changeStatusParam
// @description: 修改状态用的参数集
type changeStatusParam struct {
	WxId   string `json:"wxId" binding:"required"`
	UserId string `json:"userId"`
}

// ChangeEnableAiStatus
// @description: 修改是否开启AI
// @param ctx
func ChangeEnableAiStatus(ctx *gin.Context) {
	var p changeStatusParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("待修改的微信Id：%s", p.WxId)

	ctx.String(http.StatusOK, "操作成功")
}

// ChangeEnableGroupRankStatus
// @description: 修改是否开启水群排行榜
// @param ctx
func ChangeEnableGroupRankStatus(ctx *gin.Context) {
	var p changeStatusParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("待修改的群Id：%s", p.WxId)

	ctx.String(http.StatusOK, "操作成功")
}

// ChangeSkipGroupRankStatus
// @description: 修改是否跳过水群排行榜
// @param ctx
func ChangeSkipGroupRankStatus(ctx *gin.Context) {
	var p changeStatusParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("待修改的群Id：%s -> %s", p.WxId, p.UserId)

	ctx.String(http.StatusOK, "操作成功")
}
