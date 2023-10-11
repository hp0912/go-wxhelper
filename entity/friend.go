package entity

import "time"

// Friend
// @description: 好友列表
type Friend struct {
	CustomAccount string `json:"customAccount"` // 微信号
	Nickname      string `json:"nickname"`      // 昵称
	Pinyin        string `json:"pinyin"`        // 昵称拼音大写首字母
	PinyinAll     string `json:"pinyinAll"`     // 昵称全拼
	Wxid          string `json:"wxid"`          // 微信原始Id
}

func (Friend) TableName() string {
	return "t_friend"
}

// GroupUser
// @description: 群成员
type GroupUser struct {
	GroupId   string    `json:"groupId"`                         // 群Id
	Account   string    `json:"account"`                         // 账号
	HeadImage string    `json:"headImage"`                       // 头像
	Nickname  string    `json:"nickname"`                        // 昵称
	Wxid      string    `json:"wxid"`                            // 微信Id
	IsMember  bool      `json:"isMember" gorm:"type:tinyint(1)"` // 是否群成员
	LeaveTime time.Time `json:"leaveTime"`                       // 离开时间
}

func (GroupUser) TableName() string {
	return "t_group_user"
}
