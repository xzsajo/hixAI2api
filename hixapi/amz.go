package hixapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"strings"
	"time"
)

func GetSignURL(client cycletls.CycleTLS, cookie, chatId, fileExtension string) (string, error) {
	// 配置第一个请求参数
	headers := map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"content-type":       "application/json",
		"cookie":             cookie,
		"origin":             "https://hix.ai",
		"priority":           "u=1, i",
		"referer":            "https://hix.ai/home?id=" + chatId,
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
	}

	timeStr := time.Now().Format("20060102150405")
	options := cycletls.Options{
		Body:    `{"filename":"` + timeStr + fileExtension + `","type":"file"}`,
		Headers: headers,
	}

	// 发送 POST 请求
	response, err := client.Do("https://hix.ai/api/upload/sign", options, "POST")
	if err != nil {
		return "", err
	}

	// 解析响应数据
	var result struct {
		Sign string `json:"sign"`
	}
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	return result.Sign, nil
}

func UploadToS3(client cycletls.CycleTLS, signURL, base64Data, mimeType string) error {

	// 去除可能的 MIME 类型前缀
	base64Str := strings.SplitN(base64Data, ",", 2)[1]

	// 解码 Base64 字符串
	fileContent, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return fmt.Errorf("base64 解码失败: %v", err)
	}

	// 配置第二个请求参数
	headers := map[string]string{
		"Accept":             "*/*",
		"Accept-Language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"Connection":         "keep-alive",
		"Content-Type":       mimeType,
		"Origin":             "https://hix.ai",
		"Referer":            "https://hix.ai/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "cross-site",
		"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"macOS"`,
	}
	options := cycletls.Options{
		Body:    string(fileContent),
		Headers: headers,
	}

	// 发送 PUT 请求
	response, err := client.Do(signURL, options, "PUT")
	if err != nil {
		return fmt.Errorf("上传请求失败: %v", err)
	}

	if response.Status != 200 {
		return fmt.Errorf("上传失败，状态码: %d", response.Status)
	}

	return nil
}
