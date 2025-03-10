package hixapi

import (
	"encoding/json"
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"github.com/gin-gonic/gin"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"
)

const (
	baseURL            = "https://hix.ai"
	chatEndpoint       = baseURL + "/api/hix/chat"
	createChatEndpoint = baseURL + "/api/trpc/hixChat.createChat?batch=1"
	subUsageEndpoint   = baseURL + "/api/trpc/subUsage.getSubUsage?batch=1"
	deleteEndpoint     = baseURL + "/api/project/delete?project_id=%s"
	uploadEndpoint     = baseURL + "/api/get_upload_personal_image_url"
	chatType           = "COPILOT_MOA_CHAT"
	imageType          = "COPILOT_MOA_IMAGE"
)

type CreateChatResponse struct {
	Result struct {
		Data struct {
			JSON struct {
				ID string `json:"id"`
			} `json:"json"`
		} `json:"data"`
	} `json:"result"`
}

func MakeCreateChatRequest(client cycletls.CycleTLS, cookie string, modelId int) (string, error) {
	createChatBody := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"title": "Untitled",
				"botId": modelId,
			},
		},
	}
	bytes, err := json.Marshal(createChatBody)
	if err != nil {
		return "", err
	}
	accept := "application/json"

	response, err := client.Do(fmt.Sprintf(createChatEndpoint), cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Method:  "POST",
		Body:    string(bytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "POST")

	if err != nil {
		return "", err
	}

	var responses []CreateChatResponse
	err = json.Unmarshal([]byte(response.Body), &responses)
	if err != nil {
		return "", err
	}

	// 检查数组是否非空并提取ID
	if len(responses) > 0 {
		id := responses[0].Result.Data.JSON.ID
		return id, nil
	} else {
		return "", fmt.Errorf("MakeCreateChatRequest err")
	}
}

type SubUsageResponse struct {
	Result struct {
		Data struct {
			JSON struct {
				UsageList []struct {
					ID             int    `json:"id"`
					SubscriptionID *int   `json:"subscription_id"`
					TotalCount     int    `json:"total_count"`
					UseCount       int    `json:"use_count"`
					Status         string `json:"status"`
					UsageType      string `json:"usage_type"`
					AppName        string `json:"app_name"`
					DateStart      string `json:"date_start"`
					DateEnd        string `json:"date_end"`
					PriceID        *int   `json:"price_id"`
				} `json:"usageList"`
				IsActiveSub bool `json:"isActiveSub"`
			} `json:"json"`
			Meta struct {
				Values struct {
					DateStart []string `json:"usageList.0.date_start"`
					DateEnd   []string `json:"usageList.0.date_end"`
				} `json:"values"`
			} `json:"meta"`
		} `json:"data"`
	} `json:"result"`
}

func MakeSubUsageRequest(client cycletls.CycleTLS, cookie string) (int, error) {
	subUsageReqParam := map[string]interface{}{
		"0": map[string]interface{}{
			"json": map[string]interface{}{
				"appName": "HIXChat",
			},
		},
		"1": map[string]interface{}{
			"json": map[string]interface{}{
				"appName": "HIXChat",
			},
		},
	}
	bytes, err := json.Marshal(subUsageReqParam)
	if err != nil {
		return 0, err
	}
	accept := "application/json"

	response, err := client.Do(fmt.Sprintf(subUsageEndpoint+"&input=%s", string(bytes)), cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Method:  "GET",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "GET")

	if err != nil {
		return 0, err
	}
	bodyBytes := []byte(response.Body)
	var responses []SubUsageResponse
	err = json.Unmarshal(bodyBytes, &responses)
	if err != nil {
		return 0, err
	}

	// 检查数组是否非空并提取ID
	if len(responses) > 0 {
		if len(responses[0].Result.Data.JSON.UsageList) > 0 {
			totalCount := responses[0].Result.Data.JSON.UsageList[0].TotalCount
			useCount := responses[0].Result.Data.JSON.UsageList[0].UseCount
			return totalCount - useCount, nil
		} else {
			return 0, fmt.Errorf("MakeSubUsageRequest err ResqBody: %s Cookie: %s", string(bodyBytes), cookie)

		}

	} else {
		return 0, fmt.Errorf("MakeSubUsageRequest err ResqBody: %s Cookie: %s", string(bodyBytes), cookie)
	}
}

// makeRequest 发送HTTP请求
func makeChatRequest(client cycletls.CycleTLS, jsonData []byte, cookie string, isStream bool) (cycletls.Response, error) {
	accept := "application/json"
	if isStream {
		accept = "text/event-stream"
	}

	return client.Do(chatEndpoint, cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Body:    string(jsonData),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       accept,
			"Origin":       baseURL,
			"Referer":      baseURL + "/",
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}, "POST")
}

func MakeStreamChatRequest(c *gin.Context, client cycletls.CycleTLS, modelName, hixChatId string, jsonData []byte, cookie string) (<-chan cycletls.SSEResponse, error) {

	options := cycletls.Options{
		Timeout: 10 * 60 * 60,
		Proxy:   config.ProxyUrl, // 在每个请求中设置代理
		Body:    string(jsonData),
		Method:  "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "text/event-stream",
			"Origin":       baseURL,
			"Referer":      baseURL + "/" + modelName + "?/id=" + hixChatId,
			"Cookie":       cookie,
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome",
		},
	}

	logger.Debug(c.Request.Context(), fmt.Sprintf("cookie: %v", cookie))

	sseChan, err := client.DoSSE(chatEndpoint, options, "POST")
	if err != nil {
		logger.Errorf(c, "Failed to make stream request: %v", err)
		return nil, fmt.Errorf("Failed to make stream request: %v", err)
	}
	return sseChan, nil
}
