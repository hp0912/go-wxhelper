package goodMorning

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDraw(t *testing.T) {
	// 获取当前文件的绝对路径
	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Errorf("无法获取当前文件路径: %v", err)
		return
	}

	// 获取项目根目录路径
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(absPath)))
	// 拼接背景图片的路径
	bgFilePath := filepath.Join(projectRoot, "assets", "background2.png")
	// 拼接字体路径
	fontFilePath := filepath.Join(projectRoot, "assets", "simkai.ttf")

	// 获取当前时间
	now := time.Now()
	// 获取年、月、日
	year := now.Year()
	month := now.Month()
	day := now.Day()
	// 获取星期
	weekday := now.Weekday()
	// 定义中文星期数组
	weekdays := [...]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}

	summary := groupSummary{}
	summary.GroupID = "xxxxx"
	summary.Year = year
	summary.Month = int(month)
	summary.Date = day
	summary.Week = weekdays[weekday]
	summary.UserTotalCount = 490
	summary.UserJoinCount = 0
	summary.UserLeaveCount = 1
	summary.UserChatCount = 60
	summary.MessageCount = 1490

	err = Draw(bgFilePath, fontFilePath, "知识是很美的，它们可以让你不出家门就了解这世上的许多事。", summary)
	if err != nil {
		t.Errorf("绘图失败: %v", err)
		return
	}
}
