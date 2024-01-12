package app

import (
	"github.com/gin-gonic/gin"
	"go-wechat/client"
	"go-wechat/entity"
	"gorm.io/gorm"
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

	err := client.MySQL.Model(&entity.Friend{}).
		Where("wxid = ?", p.WxId).
		Update("`enable_ai`", gorm.Expr(" !`enable_ai`")).Error
	if err != nil {
		log.Printf("修改是否开启AI失败：%s", err)
		ctx.String(http.StatusInternalServerError, "操作失败: %s", err)
		return
	}

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

	err := client.MySQL.Model(&entity.Friend{}).
		Where("wxid = ?", p.WxId).
		Update("`enable_chat_rank`", gorm.Expr(" !`enable_chat_rank`")).Error
	if err != nil {
		log.Printf("修改开启水群排行榜失败：%s", err)
		ctx.String(http.StatusInternalServerError, "操作失败: %s", err)
		return
	}

	ctx.String(http.StatusOK, "操作成功")
}

// ChangeEnableWelcomeStatus
// @description: 修改是否开启迎新
// @param ctx
func ChangeEnableWelcomeStatus(ctx *gin.Context) {
	var p changeStatusParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("待修改的群Id：%s", p.WxId)

	err := client.MySQL.Model(&entity.Friend{}).
		Where("wxid = ?", p.WxId).
		Update("`enable_welcome`", gorm.Expr(" !`enable_welcome`")).Error
	if err != nil {
		log.Printf("修改开启迎新失败：%s", err)
		ctx.String(http.StatusInternalServerError, "操作失败: %s", err)
		return
	}

	ctx.String(http.StatusOK, "操作成功")
}

// ChangeEnableCommandStatus
// @description: 修改是否开启指令
// @param ctx
func ChangeEnableCommandStatus(ctx *gin.Context) {
	var p changeStatusParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}
	log.Printf("待修改的群Id：%s", p.WxId)

	err := client.MySQL.Model(&entity.Friend{}).
		Where("wxid = ?", p.WxId).
		Update("`enable_command`", gorm.Expr(" !`enable_command`")).Error
	if err != nil {
		log.Printf("修改指令启用状态失败：%s", err)
		ctx.String(http.StatusInternalServerError, "操作失败: %s", err)
		return
	}

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

	err := client.MySQL.Model(&entity.GroupUser{}).
		Where("group_id = ?", p.WxId).
		Where("wxid = ?", p.UserId).
		Update("`skip_chat_rank`", gorm.Expr(" !`skip_chat_rank`")).Error
	if err != nil {
		log.Printf("修改跳过水群排行榜失败：%s", err)
		ctx.String(http.StatusInternalServerError, "操作失败: %s", err)
		return
	}

	ctx.String(http.StatusOK, "操作成功")
}
