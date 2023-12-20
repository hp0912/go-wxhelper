package plugins

import (
	"go-wechat/plugin"
	"go-wechat/utils"
	"strings"
)

// Command
// @description: 自定义指令
// @param m
func Command(m *plugin.MessageContext) {
	// 判断是不是指令
	if !strings.HasPrefix(m.Content, "/") {
		return
	}

	// 用空格分割消息，下标0表示指令
	msgArray := strings.Split(m.Content, " ")
	cmd := msgArray[0]

	switch cmd {
	case "/帮助", "/h", "/help", "/?", "/？":
		helpCmd(m)
	case "/ls", "/雷神":
		leiGodCmd(m.FromUser, msgArray[1], msgArray[2:]...)
	default:
		utils.SendMessage(m.FromUser, m.GroupUser, "指令错误", 0)
	}

	// 中止后续消息处理
	m.Abort()
}

// helpCmd
// @description: 帮助指令
func helpCmd(m *plugin.MessageContext) {
	str := `帮助菜单:
指令消息必须以'/'开头，比如: '/帮助'。
支持的指令:

#1. 雷神加速器
/ls option args
option: 指令选项，可选值: 
  绑定账户：'绑定'、'b'，参数: 账户名 密码 [-f]，-f表示强制绑定，非必传项
  详情: '详情'、'i'
  暂停: '暂停'、'p'
示例: 绑定:
/ls 绑定 123456 123456 或者 /ls b 123456 123456
`
	utils.SendMessage(m.FromUser, m.GroupUser, str, 0)

}

// leiGodCmd
// @description: 雷神加速器指令
// @param userId
// @param cmd
// @param args
// @return string
func leiGodCmd(userId, cmd string, args ...string) {
	lg := newLeiGod(userId)

	var replyMsg string
	switch cmd {
	case "绑定", "b":
		var force bool
		if len(args) == 3 && args[2] == "-f" {
			force = true
		}
		replyMsg = lg.binding(args[0], args[1], force)
	case "详情", "i":
		replyMsg = lg.info()
	case "暂停", "p":
		replyMsg = lg.pause()
	default:
		replyMsg = "指令错误"
	}

	// 返回消息
	if strings.TrimSpace(replyMsg) != "" {
		utils.SendMessage(userId, "", replyMsg, 0)
	}
}
