package plugins

import (
	"go-wechat/plugin"
	"go-wechat/plugin/plugins/command"
	"go-wechat/utils"
	"regexp"
	"strings"
)

// Command
// @description: 自定义指令
// @param m
func Command(m *plugin.MessageContext) {
	// 如果是群聊，提取出消息
	content := m.Content

	if m.IsGroup() {
		re := regexp.MustCompile(`@([^ | ]+)`)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 0 {
			// 过滤掉第一个匹配到的
			content = strings.Replace(content, matches[0], "", 1)
		}
		// 去掉最前面的空格
		content = strings.TrimLeft(content, " ")
		content = strings.TrimLeft(content, " ")
	}
	// 判断是不是指令
	if !strings.HasPrefix(content, "/") {
		return
	}

	// 用空格分割消息，下标0表示指令
	msgArray := strings.Split(content, " ")
	cmd := msgArray[0]

	switch cmd {
	case "/帮助", "/h", "/help", "/?", "/？":
		command.HelpCmd(m)
	case "/雷神", "/ls":
		command.LeiGodCmd(m.FromUser, msgArray[1], msgArray[2:]...)
	case "/肯德基", "/kfc":
		command.KfcCrazyThursdayCmd(m.FromUser)
	default:
		utils.SendMessage(m.FromUser, m.GroupUser, "指令错误", 0)
	}

	// 中止后续消息处理
	m.Abort()
}
