package command

import (
	"fmt"
	"go-wechat/config"
	"go-wechat/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
)

type Cogview3PlusRequest struct {
	Model  string  `json:"model"`
	Prompt string  `json:"prompt"`
	Size   *string `json:"size"`
	UserID *string `json:"user_id"`
}

type Cogview3PlusDataItem struct {
	Url string `json:"url"`
}

type Cogview3PlusResponse struct {
	Created       int64                  `json:"created"`        // 请求创建时间，是以秒为单位的Unix时间戳
	Data          []Cogview3PlusDataItem `json:"data"`           // 数组，包含生成的图片 URL。目前数组中只包含一张图片
	Url           *string                `json:"url"`            // 图片链接。图片的临时链接有效期为 30天，请及时转存图片
	ContentFilter interface{}            `json:"content_filter"` // 返回内容安全的相关信息
	Role          *string                `json:"role"`           // 安全生效环节，包括 role = assistant 模型推理，role = user 用户输入，role = history 历史上下文
	Level         *int                   `json:"level"`          // 严重程度 level 0-3，level 0表示最严重，3表示轻微
}

// Draw
// @description: 智谱绘图接口 cogview-3-plus
// @param userId string 发信人
func Draw(userId, prompt string) {
	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
		return
	}

	var respData Cogview3PlusResponse

	res := resty.New()
	resp, err := res.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.Conf.Ai.DrawApiKey)).
		SetHeader("Content-Type", "application/json").
		SetBody(Cogview3PlusRequest{
			Model:  config.Conf.Ai.DrawModel,
			Prompt: prompt,
			UserID: &userId,
		}).
		SetResult(&respData).
		Post("https://open.bigmodel.cn/api/paas/v4/images/generations")
	if err != nil {
		log.Printf("调用绘图接口失败: %v", err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("调用绘图接口失败，返回状态码不是 200: %d", resp.StatusCode())
		return
	}

	if len(respData.Data) == 0 || respData.Data[0].Url == "" {
		log.Printf("调用绘图接口失败: 图片地址为空")
		return
	}

	imgurl := respData.Data[0].Url
	filename := fmt.Sprintf("%s-%d%s", userId, time.Now().Nanosecond(), filepath.Ext(imgurl))
	response, err := http.Get(imgurl)
	if err != nil {
		log.Printf("下载图片(%s)失败: %v", imgurl, err)
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Println("下载图片失败，状态码不为 200")
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
		log.Printf("打开图片文件(%s)失败: %v", filePath, err)
		return
	}

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Printf("保存图片文件(%s)失败: %v", filePath, err)
		return
	}

	// 发送图片
	wxPath := fmt.Sprintf(conf.Path, filename)
	log.Printf("发送对象ID: %s  图片路径: %s~", userId, wxPath)
	utils.SendImage(userId, wxPath, 1)
}
