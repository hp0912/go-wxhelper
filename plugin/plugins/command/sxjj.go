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

// Sxjj
// @description: 神仙姐姐
// @param userId string 发信人
func Sxjj(userId string) {
	type result struct {
		Code   string `json:"code"`
		Imgurl string `json:"imgurl"`
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
		Get("https://api.suyanw.cn/api/sjmv.php?return=json")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取随机美女图片失败: %v", err)
		return
	}
	if resData.Imgurl == "" {
		log.Printf("获取随机美女图片失败: 图片地址为空")
		return
	}

	urlPathArr := strings.Split(resData.Imgurl, "/")
	filename := fmt.Sprintf("%d_%s", time.Now().Nanosecond(), urlPathArr[len(urlPathArr)-1])

	response, err := http.Get(strings.Replace(resData.Imgurl, "/large/", "/middle/", 1))
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println("下载美女图片失败，状态码不为 200")
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
		log.Printf("打开美女图片文件(%s)失败: %v", filePath, err)
		return
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Printf("保存美女图片文件(%s)失败: %v", filePath, err)
		return
	}

	// 发送图片
	wxPath := fmt.Sprintf(conf.Path, filename)
	log.Printf("发送对象 ID: %s  图片路径: %s~", userId, wxPath)
	utils.SendImage(userId, wxPath, 0)
}
