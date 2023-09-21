package entity

// Friend
// @description: 好友列表
type Friend struct {
	CustomAccount string `json:"customAccount"` // 微信号
	EncryptName   string `json:"encryptName"`   // 不知道
	Nickname      string `json:"nickname"`      // 昵称
	Pinyin        string `json:"pinyin"`        // 昵称拼音大写首字母
	PinyinAll     string `json:"pinyinAll"`     // 昵称全拼
	Reserved1     int    `json:"reserved1"`     // 未知
	Reserved2     int    `json:"reserved2"`     // 未知
	Type          int    `json:"type"`          // 类型
	VerifyFlag    int    `json:"verifyFlag"`    // 未知
	Wxid          string `json:"wxid"`          // 微信原始Id
}

func (Friend) TableName() string {
	return "t_friend"
}
