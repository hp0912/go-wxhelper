package tasks

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"go-wechat/constant"
	"go-wechat/entity"
	"go-wechat/model"
	"log"
	"slices"
	"strings"
)

// 同步群成员

// syncFriends
// @description: 同步好友列表
func syncFriends() {
	var base model.Response[[]entity.Friend]

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetResult(&base).
		Post("http://10.0.0.73:19088/api/getContactList")
	if err != nil {
		log.Printf("获取好友列表失败: %s", err.Error())
		return
	}
	log.Printf("获取好友列表结果: %s", resp.String())
	for _, friend := range base.Data {
		if strings.Contains(friend.Wxid, "gh_") || strings.Contains(friend.Wxid, "@openim") {
			continue
		}
		// 特殊Id跳过
		if slices.Contains(constant.SpecialId, friend.Wxid) {
			continue
		}
		log.Printf("昵称: %s -> 类型: %d -> 微信号: %s -> 微信原始Id: %s", friend.Nickname, friend.Type, friend.CustomAccount, friend.Wxid)

		// 群成员，同步一下成员信息
		if strings.Contains(friend.Wxid, "@chatroom") {
			syncGroupUsers(friend.Wxid)
		}

	}
}

// syncGroupUsers
// @description: 同步群成员
// @param gid
func syncGroupUsers(gid string) {
	var baseResp model.Response[model.GroupUser]

	// 组装参数
	param := map[string]any{
		"chatRoomId": gid, // 群Id
	}
	pbs, _ := json.Marshal(param)

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		SetResult(&baseResp).
		Post("http://10.0.0.73:19088/api/getMemberFromChatRoom")
	if err != nil {
		log.Printf("获取群成员信息失败: %s", err.Error())
		return
	}

	// 昵称Id
	wxIds := strings.Split(baseResp.Data.Members, "^G")

	log.Printf("      群成员数: %d", len(wxIds))
	for _, wxid := range wxIds {
		// 获取成员信息
		cp, _ := getContactProfile(wxid)
		if cp.Wxid != "" {
			log.Printf("            微信Id: %s -> 昵称: %s -> 微信号: %s", wxid, cp.Nickname, cp.Account)
		}
	}
}

// getContactProfile
// @description: 获取成员详情
// @param wxid
// @return ent
// @return err
func getContactProfile(wxid string) (ent model.ContactProfile, err error) {
	var baseResp model.Response[model.ContactProfile]

	// 组装参数
	param := map[string]any{
		"wxid": wxid, // 群Id
	}
	pbs, _ := json.Marshal(param)

	client := resty.New()
	_, err = client.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetBody(string(pbs)).
		SetResult(&baseResp).
		Post("http://10.0.0.73:19088/api/getContactProfile")
	if err != nil {
		log.Printf("获取成员详情失败: %s", err.Error())
		return
	}
	ent = baseResp.Data
	return
}
