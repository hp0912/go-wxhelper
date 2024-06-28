package command

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Ydgg
// @description: 御弟哥哥 随机帅哥图
// @param userId string 发信人
func Ydgg(userId string) {
	type result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Url  string `json:"url"`
	}

	var resData result

	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
		return
	}

	res := resty.New()
	resp, err := res.R().
		SetResult(&resData).
		Get("https://api.52vmy.cn/api/img/tu/boy")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取随机帅哥图片失败: %v", err)
		return
	}
	if resData.Url == "" {
		log.Printf("获取随机帅哥图片失败: 图片地址为空")
		return
	}

	urlPathArr := strings.Split(resData.Url, "/")
	filename := fmt.Sprintf("%d_%s", time.Now().Nanosecond(), urlPathArr[len(urlPathArr)-1])

	response, err := http.Get(resData.Url)
	if err != nil || response.StatusCode != http.StatusOK {
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
