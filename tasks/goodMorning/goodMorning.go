package goodMorning

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/service"
	"go-wechat/utils"
	"image"
	"unicode/utf8"

	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

// GoodMorning
// @description: 早安书
func GoodMorning() {
	// 获取当前文件的绝对路径
	absPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Printf("无法获取当前文件路径: %v", err)
		return
	}

	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
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

	groups, err := service.GetAllEnableGoodMorning()
	if err != nil {
		log.Printf("获取启用了早安书的群组失败, 错误信息: %v", err)
		return
	}

	// 每日一言
	dailyWords := "早上好，今天接口挂了，没有早安语。"
	res := resty.New()
	resp, err := res.R().
		Post("https://api.pearktrue.cn/api/hitokoto/")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取随机一言失败: %v", err)
	}
	respText := resp.String()
	if respText != "" {
		dailyWords = respText
	}

	// 循环发送图片
	for _, group := range groups {
		summary, err := getGroupSummary(group.Wxid)
		if err != nil {
			log.Printf("统计群[%s]信息失败: %v", group.Nickname, err)
			continue
		}

		summary.Year = year
		summary.Month = int(month)
		summary.Date = day
		summary.Week = weekdays[weekday]

		err = Draw(bgFilePath, fontFilePath, dailyWords, summary)
		if err != nil {
			log.Printf("生成群[%s]图片失败: %v", group.Nickname, err)
			continue
		}
		// 发送图片
		wxPath := fmt.Sprintf(conf.Path, fmt.Sprintf("morning_%s.png", group.Wxid))
		log.Printf("早安书 -> 发送对象ID: %s  图片路径: %s~", group.Wxid, wxPath)
		utils.SendImage(group.Wxid, wxPath, 0)
		// 休眠20秒，防止频繁发送
		time.Sleep(20 * time.Second)
	}
}

func Draw(bgPath, fontPath, dailyWords string, summary groupSummary) error {
	// 打开背景图片
	bgFile, err := os.Open(bgPath)
	if err != nil {
		return err
	}
	defer bgFile.Close()

	// 解码
	bgFileImage, err := png.Decode(bgFile)
	if err != nil {
		return err
	}
	// 新建一张和模板文件一样大小的画布
	newBgImage := image.NewRGBA(bgFileImage.Bounds())
	// 将模板图片画到新建的画布上
	draw.Draw(newBgImage, bgFileImage.Bounds(), bgFileImage, bgFileImage.Bounds().Min, draw.Over)

	// 加载字体文件  这里我们加载两种字体文件
	fontKai, err := loadFont(fontPath)
	if err != nil {
		return err
	}

	// 向图片中写入文字
	// 在写入之前有一些准备工作
	content := freetype.NewContext()
	content.SetClip(newBgImage.Bounds())
	content.SetDst(newBgImage)
	content.SetSrc(image.Black) // 设置字体颜色
	content.SetDPI(72)          // 设置字体分辨率

	content.SetFontSize(24)  // 设置字体大小
	content.SetFont(fontKai) // 设置字体样式，就是我们上面加载的字体

	drawLeft := 16
	drawTop := 30

	// 参数1：要写入的文字
	// 参数2：文字坐标
	txt := fmt.Sprintf("%d年%d月%d日  %s", summary.Year, summary.Month, summary.Date, summary.Week)
	content.DrawString(txt, freetype.Pt(drawLeft, drawTop))

	content.SetSrc(image.Opaque)
	content.SetFontSize(18) // 设置字体大小
	drawTop += 44
	content.DrawString("早上好，", freetype.Pt(drawLeft, drawTop))

	drawTop += 30
	drawContent := "群内成员一共"
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
	drawLeft += (utf8.RuneCountInString(drawContent) * 18)
	drawContent = fmt.Sprintf("%d", summary.UserTotalCount)
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	content.SetSrc(image.Opaque)
	drawLeft += (len(drawContent) * 10)
	drawContent = "名，"
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft = 16
	drawTop += 30
	if summary.UserJoinCount == 0 && summary.UserLeaveCount == 0 {
		drawContent = "没有人加入，也没有人离开，"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))
	} else if summary.UserJoinCount == 0 {
		drawContent = "没有人加入，有"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (utf8.RuneCountInString(drawContent) * 18)
		content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
		drawContent = fmt.Sprintf("%d", summary.UserLeaveCount)
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (len(drawContent) * 10)
		content.SetSrc(image.Opaque)
		drawContent = "人离开了我们，"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))
	} else if summary.UserLeaveCount == 0 {
		drawContent = "有"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (utf8.RuneCountInString(drawContent) * 18)
		content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
		drawContent = fmt.Sprintf("%d", summary.UserJoinCount)
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (len(drawContent) * 10)
		content.SetSrc(image.Opaque)
		drawContent = "人加入，没有人离开，"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))
	} else {
		drawContent = "有"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (utf8.RuneCountInString(drawContent) * 18)
		content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
		drawContent = fmt.Sprintf("%d", summary.UserJoinCount)
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (len(drawContent) * 10)
		content.SetSrc(image.Opaque)
		drawContent = "人加入，但也有"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (utf8.RuneCountInString(drawContent) * 18)
		content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
		drawContent = fmt.Sprintf("%d", summary.UserLeaveCount)
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

		drawLeft += (len(drawContent) * 10)
		content.SetSrc(image.Opaque)
		drawContent = "人离开了我们，"
		content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))
	}

	drawLeft = 16
	drawTop += 30
	content.SetSrc(image.Opaque)
	drawContent = "共有"
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft += (utf8.RuneCountInString(drawContent) * 18)
	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
	drawContent = fmt.Sprintf("%d", summary.UserChatCount)
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft += (len(drawContent) * 10)
	content.SetSrc(image.Opaque)
	drawContent = "人侃侃而谈"
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft += (utf8.RuneCountInString(drawContent) * 18)
	content.SetSrc(image.NewUniform(color.RGBA{R: 237, G: 39, B: 90, A: 255})) // 设置字体颜色
	drawContent = fmt.Sprintf("%d", summary.MessageCount)
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft += (len(drawContent) * 10)
	content.SetSrc(image.Opaque)
	drawContent = "句。"
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	drawLeft = 16
	drawTop = 220
	content.SetFontSize(13)
	content.SetSrc(image.Opaque)
	drawContent = dailyWords
	content.DrawString(drawContent, freetype.Pt(drawLeft, drawTop))

	// 保存图片
	err = saveFile(newBgImage, summary.GroupID)
	if err != nil {
		return err
	}

	return nil
}

// 根据路径加载字体文件
// path 字体的路径
func loadFont(path string) (font *truetype.Font, err error) {
	var fontBytes []byte
	fontBytes, err = os.ReadFile(path) // 读取字体文件
	if err != nil {
		err = fmt.Errorf("加载字体文件出错:%s", err.Error())
		return
	}
	font, err = freetype.ParseFont(fontBytes) // 解析字体文件
	if err != nil {
		err = fmt.Errorf("解析字体文件出错,%s", err.Error())
		return
	}
	return
}

func saveFile(pic *image.RGBA, groupID string) error {
	dstFile, err := os.Create(fmt.Sprintf("/app/temp/morning_%s.png", groupID))
	if err != nil {
		return err
	}
	defer dstFile.Close()
	png.Encode(dstFile, pic)
	return nil
}
