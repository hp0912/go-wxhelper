package tasks

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"go-wechat/client"
	"go-wechat/constant"
	"go-wechat/entity"
	"go-wechat/model"
	"gorm.io/gorm"
	"log"
	"slices"
	"strings"
)

// 同步群成员

// syncFriends
// @description: 同步好友列表
func syncFriends() {
	var base model.Response[[]model.FriendItem]

	res := resty.New()
	resp, err := res.R().
		SetHeader("Content-Type", "application/json;chartset=utf-8").
		SetResult(&base).
		Post("http://10.0.0.73:19088/api/getContactList")
	if err != nil {
		log.Printf("获取好友列表失败: %s", err.Error())
		return
	}
	log.Printf("获取好友列表结果: %s", resp.String())

	tx := client.MySQL.Begin()
	defer tx.Commit()

	for _, friend := range base.Data {
		if strings.Contains(friend.Wxid, "gh_") || strings.Contains(friend.Wxid, "@openim") {
			continue
		}
		// 特殊Id跳过
		if slices.Contains(constant.SpecialId, friend.Wxid) {
			continue
		}
		log.Printf("昵称: %s -> 类型: %d -> 微信号: %s -> 微信原始Id: %s", friend.Nickname, friend.Type, friend.CustomAccount, friend.Wxid)

		// 判断是否存在，不存在的话就新增，存在就修改一下名字
		var count int64
		err = tx.Model(&entity.Friend{}).Where("wxid = ?", friend.Wxid).Count(&count).Error
		if err != nil {
			continue
		}
		if count == 0 {
			// 新增
			err = tx.Create(&entity.Friend{
				CustomAccount: friend.CustomAccount,
				Nickname:      friend.Nickname,
				Pinyin:        friend.Pinyin,
				PinyinAll:     friend.PinyinAll,
				Wxid:          friend.Wxid,
			}).Error
			if err != nil {
				log.Printf("新增好友失败: %s", err.Error())
				continue
			}
		} else {
			pm := map[string]any{
				"nickname":       friend.Nickname,
				"custom_account": friend.CustomAccount,
				"pinyin":         friend.Pinyin,
				"pinyin_all":     friend.PinyinAll,
			}
			err = tx.Model(&entity.Friend{}).Where("wxid = ?", friend.Wxid).Updates(pm).Error
			if err != nil {
				log.Printf("修改好友失败: %s", err.Error())
				continue
			}
		}

		// 群成员，同步一下成员信息
		if strings.Contains(friend.Wxid, "@chatroom") {
			syncGroupUsers(tx, friend.Wxid)
		}
	}

	log.Println("同步好友列表完成")
}

// syncGroupUsers
// @description: 同步群成员
// @param gid
func syncGroupUsers(tx *gorm.DB, gid string) {
	var baseResp model.Response[model.GroupUser]

	// 组装参数
	param := map[string]any{
		"chatRoomId": gid, // 群Id
	}
	pbs, _ := json.Marshal(param)

	res := resty.New()
	_, err := res.R().
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

	// 修改不在数组的群成员状态为不在
	err = tx.Model(&entity.GroupUser{}).Where("wxid NOT IN (?)", wxIds).Update("is_member", false).Error
	if err != nil {
		log.Printf("修改群成员状态失败: %s", err.Error())
		return
	}

	for _, wxid := range wxIds {
		// 获取成员信息
		cp, _ := getContactProfile(wxid)
		if cp.Wxid != "" {
			log.Printf("            微信Id: %s -> 昵称: %s -> 微信号: %s", wxid, cp.Nickname, cp.Account)
			// 查询成员是否存在，不在就新增，否则修改
			var count int64
			err = tx.Model(&entity.GroupUser{}).Where("group_id = ?", gid).Where("wxid = ?", wxid).Count(&count).Error
			if err != nil {
				log.Printf("查询群成员失败: %s", err.Error())
				continue
			}
			if count == 0 {
				// 新增
				err = tx.Create(&entity.GroupUser{
					GroupId:   gid,
					Account:   cp.Account,
					HeadImage: cp.HeadImage,
					Nickname:  cp.Nickname,
					Wxid:      cp.Wxid,
					IsMember:  true,
				}).Error
				if err != nil {
					log.Printf("新增群成员失败: %s", err.Error())
					continue
				}
			} else {
				// 修改
				pm := map[string]any{
					"account":    cp.Account,
					"head_image": cp.HeadImage,
					"nickname":   cp.Nickname,
				}
				err = tx.Model(&entity.GroupUser{}).Where("group_id = ?", gid).Where("wxid = ?", wxid).Updates(pm).Error
				if err != nil {
					log.Printf("修改群成员失败: %s", err.Error())
					continue
				}
			}
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
