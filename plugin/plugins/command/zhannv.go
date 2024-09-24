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

// ZhaNan
// @description: 渣女语录
// @param userId string 发信人
func ZhaNv(userId string) {
	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
		return
	}

	type result struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		Audiopath string `json:"audiopath"`
	}

	var resData result

	Audiopath := ""

	res := resty.New()
	resp, err := res.R().
		SetResult(&resData).
		Get("https://api.pearktrue.cn/api/greentea/")
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Printf("获取渣女语录失败: %v", err)
		return
	} else if resData.Audiopath != "" {
		Audiopath = resData.Audiopath
	} else {
		Audiopath = resp.String()
	}

	if Audiopath == "" {
		log.Printf("获取渣女语录失败: 地址为空")
		return
	}

	filename := fmt.Sprintf("%d.mp3", time.Now().Nanosecond())
	response, err := http.Get(Audiopath)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println("下载渣女语录失败，状态码不为 200")
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
		log.Printf("打开渣女语录文件(%s)失败: %v", filePath, err)
		return
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Printf("保存渣女语录文件(%s)失败: %v", filePath, err)
		return
	}

	// 发送文件
	wxPath := fmt.Sprintf(conf.Path, filename)
	log.Printf("发送对象ID: %s  图片路径: %s~", userId, wxPath)
	utils.SendFile(userId, wxPath, 0)
}
