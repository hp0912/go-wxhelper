package plugins

import (
	"fmt"
	"go-wechat/plugin"
	"go-wechat/plugin/plugins/command"
	"go-wechat/utils"
	"regexp"
	"strconv"
	"strings"
)

func GroupSummaryHandler(fromUser, groupUser string, args []string) {
	if len(args) == 2 {
		condition := args[1]
		switch condition {
		case "昨日":
			command.GroupSummary(fromUser, true, "")
		default:
			re := regexp.MustCompile(`\d+`)
			match := re.FindString(condition)
			if match != "" {
				num, err := strconv.Atoi(match)
				if err != nil {
					utils.SendMessage(fromUser, groupUser, fmt.Sprintf("指令错误: %s", err.Error()), 0)
					return
				}
				if strings.HasSuffix(condition, "时") {
					if num > 24 {
						utils.SendMessage(fromUser, groupUser, "指令错误，统计最近多少小时的消息时，小时不能大于 24，您可以使用尝试：/群聊总结 昨日", 0)
						return
					}
					command.GroupSummary(fromUser, false, fmt.Sprintf("%s HOUR", match))
				} else if strings.HasSuffix(condition, "分") {
					if num > 60 {
						utils.SendMessage(fromUser, groupUser, "指令错误，统计最近多少分钟的消息时，分钟不能大于 60，您可以使用尝试：/群聊总结 近1小时", 0)
						return
					}
					command.GroupSummary(fromUser, false, fmt.Sprintf("%s MINUTE", match))
				} else {
					utils.SendMessage(fromUser, groupUser, "指令错误，群聊总结只能统计最近多少分钟，最近多少小时，或者昨日的记录", 0)
				}
			} else {
				utils.SendMessage(fromUser, groupUser, "指令错误，下面是一个正确的例子：/群聊总结 近1小时，其中指令和统计时间之间有个空格", 0)
			}
		}
	} else {
		command.GroupSummary(fromUser, true, "")
	}
}

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
	case "/ai":
		command.AiCmd(m.FromUser, m.GroupUser, msgArray[1])
	case "/舔狗日记":
		command.DogLickingDiary(m.FromUser)
	case "/毒鸡汤":
		command.DogLickingDiary(m.FromUser)
	case "/暖男语录":
		command.ZhaNan(m.FromUser)
	case "/渣女语录":
		command.ZhaNv(m.FromUser)
	case "/昨日热词":
		command.WordCloud(m.FromUser)
	case "/神仙姐姐":
		command.Sxjj(m.FromUser)
	case "/御弟哥哥":
		command.Ydgg(m.FromUser)
	case "/绘图":
		command.Draw(m.FromUser, strings.Join(msgArray[1:], ""))
	case "/群聊总结":
		GroupSummaryHandler(m.FromUser, m.GroupUser, msgArray)
	default:
		utils.SendMessage(m.FromUser, m.GroupUser, "指令错误", 0)
	}

	// 中止后续消息处理
	m.Abort()
}
