package command

import (
	"go-wechat/plugin"
	"go-wechat/utils"
)

// HelpCmd
// @description: 帮助指令
func HelpCmd(m *plugin.MessageContext) {
	str := `帮助菜单:
指令消息必须以'/'开头，比如: '/帮助'。
支持的指令:

#1. 雷神加速器
/ls option args
option: 指令选项，可选值: 
  绑定账户：'绑定'、'b'，参数: 账户名 密码 [-f]，-f表示强制绑定，非必传项
  详情: '详情'、'i'
  暂停: '暂停'、'p'
示例:
示例1: /ls 绑定 123456 123456
示例2: /ls b 123456 123456

#2. 肯德基疯狂星期四文案
示例1: /kfc
示例2: /肯德基

#3. AI助手
/ai option
option: 指令选项，可选值:
  启用: '启用'、'打开'、'enable'
  停用: '停用'、'禁用'、'关闭'、'disable'
示例1: /ai 启用
示例2: /ai 禁用

#4. 舔狗日记
示例: /舔狗日记

#5. 毒鸡汤
示例: /毒鸡汤

#6. 渣男/女语录
示例1: /渣男语录
示例2: /渣女语录

#7. 昨日热词
示例: /昨日热词
`
	utils.SendMessage(m.FromUser, m.GroupUser, str, 0)

}
