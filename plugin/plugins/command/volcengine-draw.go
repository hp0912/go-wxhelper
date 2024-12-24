package command

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"go-wechat/config"
	"go-wechat/utils"
	"go-wechat/utils/volcengine/drawing"
)

type VolDrawRequest struct {
	ReqKey          string `json:"req_key"`
	Prompt          string `json:"prompt"`
	ModelVersion    string `json:"model_version"`
	ReqScheduleConf string `json:"req_schedule_conf"`
	ReturnUrl       bool   `json:"return_url"`
}

type VolDrawResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AlgorithmBaseResp struct {
			StatusCode    int    `json:"status_code"`
			StatusMessage string `json:"status_message"`
		} `json:"algorithm_base_resp"`
		BinaryDataBase64  []string `json:"binary_data_base64"`
		ImageUrls         []string `json:"image_urls"`
		PeResult          string   `json:"pe_result"`
		PredictTagsResult string   `json:"predict_tags_result"`
		RephraserResult   string   `json:"rephraser_result"`
		RequestId         string   `json:"request_id"`
	} `json:"data"`
	RequestId   string `json:"request_id"`
	TimeElapsed string `json:"time_elapsed"`
	Status      int    `json:"status"`
}

// VolDraw
// @description: 豆包绘图接口
// @param toUserId string 发给谁（群）
// @param groupUserId string 艾特谁
func VolDraw(toUserId, groupUserId, prompt string) {
	conf, ok := config.Conf.Resource["temp"]
	if !ok {
		log.Printf("获取临时目录失败~")
		return
	}

	queries := make(url.Values)

	data := VolDrawRequest{
		ReqKey:          "high_aes_general_v21_L",
		Prompt:          prompt,
		ModelVersion:    "general_v2.1_L",
		ReqScheduleConf: config.Conf.Ai.DrawReqScheduleConf, // general_v20_9B_pe general_v20_9B_rephraser
		ReturnUrl:       true,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("调用绘图接口失败: %v", err)
		return
	}

	resp, err := drawing.DoRequest("POST", queries, body)
	if err != nil {
		utils.SendMessage(toUserId, groupUserId, err.Error(), 0)
		return
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendMessage(toUserId, groupUserId, err.Error(), 0)
		return
	}

	// 将响应体解析为结构体
	var respData VolDrawResponse
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		utils.SendMessage(toUserId, groupUserId, err.Error(), 0)
		return
	}

	if respData.Code != 10000 {
		msg := "触发敏感词，绘图失败"
		switch respData.Code {
		case 150411:
			msg = "输入图片前审核未通过"
		case 50511:
			msg = "输出图片后审核未通过"
		case 50412:
			msg = "输入文本前审核未通过"
		case 50512:
			msg = "输出文本后审核未通过"
		case 50413:
			msg = "输入文本NER、IP、Blocklist等拦截"
		}
		utils.SendMessage(toUserId, groupUserId, msg, 0)
		return
	}

	if len(respData.Data.ImageUrls) == 0 || respData.Data.ImageUrls[0] == "" {
		utils.SendMessage(toUserId, groupUserId, "调用绘图接口失败: 图片地址为空", 0)
		return
	}

	imgurl := respData.Data.ImageUrls[0]
	filename := fmt.Sprintf("%s-%d%s", toUserId, time.Now().Nanosecond(), ".jpeg")
	response, err := http.Get(imgurl)
	if err != nil {
		utils.SendMessage(toUserId, groupUserId, fmt.Sprintf("下载图片(%s)失败: %v", imgurl, err), 0)
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
	log.Printf("发送对象ID: %s  图片路径: %s~", toUserId, wxPath)
	utils.SendImage(toUserId, wxPath, 1)
}
