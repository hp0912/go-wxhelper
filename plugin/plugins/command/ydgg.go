package command

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// Ydgg
// @description: 御弟哥哥 随机帅哥图
// @param userId string 发信人
func Ydgg(userId string) {
	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
		return
	}

	res := resty.New()
	resp, err := res.R().
		Get("https://api.suyanw.cn/api/boy.php?type=text")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取随机帅哥图片失败: %v", err)
		return
	}
	imgurl := resp.String()
	if imgurl == "" {
		log.Printf("获取随机帅哥图片失败: 图片地址为空")
		return
	}

	filename := fmt.Sprintf("%d.jpg", time.Now().Nanosecond())
	response, err := http.Get(imgurl)
	if err != nil {
		log.Printf("下载帅哥图片(%s)失败: %v", imgurl, err)
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Println("下载帅哥图片失败，状态码不为 200")
		return
	}
	reader := response.Body
	defer reader.Close()
	filePath := fmt.Sprintf("/app/temp/%s", filename)
	file, err := os.Create(filePath)
	defer func() {
		file.Close()
	}()

	if err != nil {
		log.Printf("打开帅哥图片文件(%s)失败: %v", filePath, err)
		return
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Printf("保存帅哥图片文件(%s)失败: %v", filePath, err)
		return
	}

	// 发送图片
	wxPath := fmt.Sprintf(conf.Path, filename)
	log.Printf("发送对象ID: %s  图片路径: %s~", userId, wxPath)
	utils.SendImage(userId, wxPath, 0)
}
