package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo/mutable"
	"gorm.io/gorm"
	"hixai2api/common"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"

	"hixai2api/database"
	"hixai2api/hixapi"
	"hixai2api/model"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	errNoValidCookies = "No valid cookies available"
	responseIDFormat  = "chatcmpl-%s"
)

// ChatForOpenAI @Summary OpenAI对话接口
// @Description OpenAI对话接口
// @Tags OpenAI
// @Accept json
// @Produce json
// @Param req body model.OpenAIChatCompletionRequest true "OpenAI对话请求"
// @Param Authorization header string true "Authorization API-KEY"
// @Router /v1/chat/completions [post]
func ChatForOpenAI(c *gin.Context) {
	client := cycletls.Init()
	defer safeClose(client)

	var openAIReq model.OpenAIChatCompletionRequest
	if err := c.BindJSON(&openAIReq); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		c.JSON(http.StatusInternalServerError, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Invalid request parameters",
				Type:    "request_error",
				Code:    "500",
			},
		})
		return
	}

	if openAIReq.Stream {
		handleStreamRequest(c, client, openAIReq)
	} else {
		handleNonStreamRequest(c, client, openAIReq)
	}
}

func handleNonStreamRequest(c *gin.Context, client cycletls.CycleTLS, openAIReq model.OpenAIChatCompletionRequest) {
	var err error
	var cookies []model.Cookie
	searchType := ""
	if strings.HasSuffix(openAIReq.Model, "-search") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-search", "", 1)
		searchType = "internet"
	} else if strings.HasSuffix(openAIReq.Model, "-news") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-news", "", 1)
		searchType = "news"
	} else if strings.HasSuffix(openAIReq.Model, "-academic") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-academic", "", 1)
		searchType = "academic"
	}
	modelInfo, ok := common.GetHixModelInfo(openAIReq.Model)
	// 1. 先获取该模型的Credit
	if ok {
		// 2. 从token表中获取该credit大于该模型的token值
		cookieRecord := &model.Cookie{
			Credit: modelInfo.Credit,
		}
		if modelInfo.Type == "STANDARD" {
			cookies, err = cookieRecord.FindByMinimumCredit(database.DB)
		} else {
			cookies, err = cookieRecord.FindByMinimumCreditAdvanced(database.DB)
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "no token"})
			return
		}
		if len(cookies) == 0 {
			c.JSON(500, gin.H{"error": errNoValidCookies})
			return
		}
	} else {
		c.JSON(500, gin.H{"error": "no model"})
		return
	}

	newChatFlag := false
	var hixChatId string
	// 1. 获取符合messagehash的tokens
	pair, b, err := openAIReq.GetPreviousMessagePair()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !b {
		// 没有last问答对 视为新会话
		newChatFlag = true
	} else {
		msgPairSha256 := common.StringToSHA256(strings.TrimSpace(pair))
		cookie, chatId, err := model.QueryCookiesByChatHashAndModelAndCredit(database.DB, msgPairSha256, openAIReq.Model, modelInfo.Credit)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newChatFlag = true
		} else if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		} else {
			cookies = []model.Cookie{
				{
					Cookie: cookie,
				},
			}
			hixChatId = chatId
		}
	}

	responseId := fmt.Sprintf(responseIDFormat, time.Now().Format("20060102150405"))
	ctx := c.Request.Context()

	mutable.Shuffle(cookies)

	maxRetries := len(cookies)

	var messagesPair []model.OpenAIChatMessage
	messagesPair = append(messagesPair, openAIReq.Messages[len(openAIReq.Messages)-1])
	for attempt := 0; attempt < maxRetries; attempt++ {
		cookie := cookies[attempt]
		if newChatFlag {
			chatId, err := hixapi.MakeCreateChatRequest(client, cookie.Cookie, modelInfo.ModelID)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			hixChatId = chatId
		}

		requestBody, err := createRequestBody(c, hixChatId, &openAIReq, searchType, modelInfo)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to marshal request body"})
			return
		}
		sseChan, err := hixapi.MakeStreamChatRequest(c, client, openAIReq.Model, hixChatId, jsonData, cookie.Cookie)
		if err != nil {
			logger.Errorf(ctx, "MakeStreamChatRequest err on attempt %d: %v", attempt+1, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		isRateLimit := false
		var assistantMsgContent string
		var delta string
		var shouldContinue bool
		thinkStartType := new(bool) // 初始值为false
		for response := range sseChan {
			if response.Status == 403 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Forbidden"})
				return
			}
			if response.Done {
				logger.Warnf(ctx, response.Data)
				return
			}

			data := response.Data
			if data == "" {
				continue
			}

			logger.Debug(ctx, strings.TrimSpace(data))

			switch {
			case common.IsCloudflareChallenge(data):
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cf challenge"})
				return
			case common.IsNotLogin(data):
				isRateLimit = true
				logger.Warnf(ctx, "Cookie Not Login, switching to next cookie, attempt %d/%d, COOKIE:%s", attempt+1, maxRetries, cookie)
				// 删除cookie
				//config.RemoveCookie(cookie)
				break
			}

			streamDelta, streamShouldContinue := processNoStreamData(c, data, responseId, openAIReq.Model, jsonData, thinkStartType)
			delta = streamDelta
			shouldContinue = streamShouldContinue
			// 处理事件流数据
			if !shouldContinue {
				// 保存chat记录
				messagesPair = append(messagesPair, model.OpenAIChatMessage{
					Role:    "assistant",
					Content: strings.TrimSpace(assistantMsgContent),
				})
				bytes, err := json.Marshal(messagesPair)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				messagesPairStr := strings.NewReplacer(
					`\n`, "",
					`\t`, "",
					`\r`, "",
				).Replace(string(bytes))

				chat := model.Chat{
					Cookie:                     cookie.Cookie,
					Model:                      openAIReq.Model,
					CookieHash:                 cookie.CookieHash,
					HixChatId:                  hixChatId,
					LastMessagesPair:           string(bytes),
					LastMessagesPairSha256Hash: common.StringToSHA256(messagesPairStr),
				}
				if newChatFlag {
					if err := chat.Create(database.DB); err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
					}
				} else {
					// 更新chat记录
					if err := chat.UpdateLastMessages(database.DB); err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
					}
				}
				promptTokens := model.CountTokenText(string(jsonData), openAIReq.Model)
				completionTokens := model.CountTokenText(assistantMsgContent, openAIReq.Model)
				finishReason := "stop"

				c.JSON(http.StatusOK, model.OpenAIChatCompletionResponse{
					ID:      fmt.Sprintf(responseIDFormat, time.Now().Format("20060102150405")),
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   openAIReq.Model,
					Choices: []model.OpenAIChoice{{
						Message: model.OpenAIMessage{
							Role:    "assistant",
							Content: assistantMsgContent,
						},
						FinishReason: &finishReason,
					}},
					Usage: model.OpenAIUsage{
						PromptTokens:     promptTokens,
						CompletionTokens: completionTokens,
						TotalTokens:      promptTokens + completionTokens,
					},
				})

				go func() {
					isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, cookie.Cookie)
					if err != nil {
						logger.Errorf(ctx, "MakeSubUsageRequest err: %v", err)
					}
					cookieRecord := &model.Cookie{
						CookieHash:     cookie.CookieHash,
						Credit:         credit,
						AdvancedCredit: advancedCredit,
						IsActiveSub:    isActiveSub,
					}
					err = cookieRecord.UpdateCreditByCookieHash(database.DB)
					if err != nil {
						logger.Errorf(ctx, "UpdateCreditByCookieHash err: %v", err)
					}
				}()
				return
			} else {
				//if strings.TrimSpace(delta) != "" {
				assistantMsgContent = assistantMsgContent + delta

				//}
			}

		}
		if !isRateLimit {
			return
		}

	}
	logger.Errorf(ctx, "All cookies exhausted after %d attempts", maxRetries)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "All cookies are temporarily unavailable."})
	return
}

func createRequestBody(c *gin.Context, chatId string, openAIReq *model.OpenAIChatCompletionRequest, searchType string, modelInfo common.HixModelInfo) (map[string]interface{}, error) {
	if config.PRE_MESSAGES_JSON != "" {
		err := openAIReq.PrependMessagesFromJSON(config.PRE_MESSAGES_JSON)
		if err != nil {
			return nil, fmt.Errorf("PrependMessagesFromJSON err: %v JSON:%s", err, config.PRE_MESSAGES_JSON)
		}
	}
	openAIReq.FilterUserMessage()
	var question string
	switch content := openAIReq.Messages[0].Content.(type) {
	case string:
		runeCountInString := utf8.RuneCountInString(content)
		if runeCountInString > modelInfo.MaxTokens {
			return nil, fmt.Errorf("input text too long: %d", runeCountInString)
		}
		question = content
	default:
		return nil, fmt.Errorf("Invalid message content type: %T", content)
	}
	requestBody := map[string]interface{}{
		"chatId":   chatId,
		"fileUrl":  "",
		"question": question,
	}
	if searchType != "" &&
		(searchType == "internet" || searchType == "news" || searchType == "academic") {
		requestBody["search"] = true
		requestBody["searchType"] = searchType
	}
	// 创建请求体
	logger.Debug(c.Request.Context(), fmt.Sprintf("RequestBody: %v", requestBody))

	return requestBody, nil
}

// createStreamResponse 创建流式响应
func createStreamResponse(responseId, modelName string, jsonData []byte, delta model.OpenAIDelta, finishReason *string) model.OpenAIChatCompletionResponse {
	promptTokens := model.CountTokenText(string(jsonData), modelName)
	completionTokens := model.CountTokenText(delta.Content, modelName)
	return model.OpenAIChatCompletionResponse{
		ID:      responseId,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   modelName,
		Choices: []model.OpenAIChoice{
			{
				Index:        0,
				Delta:        delta,
				FinishReason: finishReason,
			},
		},
		Usage: model.OpenAIUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}
}

// handleDelta 处理消息字段增量
func handleDelta(c *gin.Context, delta string, responseId, modelName string, jsonData []byte) error {
	// 创建基础响应
	createResponse := func(content string) model.OpenAIChatCompletionResponse {
		return createStreamResponse(
			responseId,
			modelName,
			jsonData,
			model.OpenAIDelta{Content: content, Role: "assistant"},
			nil,
		)
	}

	// 发送基础事件
	var err error
	if err = sendSSEvent(c, createResponse(delta)); err != nil {
		return err
	}

	return err
}

// handleMessageResult 处理消息结果
func handleMessageResult(c *gin.Context, responseId, modelName string, jsonData []byte) bool {
	finishReason := "stop"
	var delta string

	streamResp := createStreamResponse(responseId, modelName, jsonData, model.OpenAIDelta{Content: delta, Role: "assistant"}, &finishReason)
	if err := sendSSEvent(c, streamResp); err != nil {
		logger.Warnf(c.Request.Context(), "sendSSEvent err: %v", err)
		return false
	}
	c.SSEvent("", " [DONE]")
	return false
}

// sendSSEvent 发送SSE事件
func sendSSEvent(c *gin.Context, response model.OpenAIChatCompletionResponse) error {
	jsonResp, err := json.Marshal(response)
	if err != nil {
		logger.Errorf(c.Request.Context(), "Failed to marshal response: %v", err)
		return err
	}
	c.SSEvent("", " "+string(jsonResp))
	c.Writer.Flush()
	return nil
}

func handleStreamRequest(c *gin.Context, client cycletls.CycleTLS, openAIReq model.OpenAIChatCompletionRequest) {

	var err error

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	var cookies []model.Cookie
	searchType := ""
	if strings.HasSuffix(openAIReq.Model, "-search") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-search", "", 1)
		searchType = "internet"
	} else if strings.HasSuffix(openAIReq.Model, "-news") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-news", "", 1)
		searchType = "news"
	} else if strings.HasSuffix(openAIReq.Model, "-academic") {
		openAIReq.Model = strings.Replace(openAIReq.Model, "-academic", "", 1)
		searchType = "academic"
	}
	modelInfo, ok := common.GetHixModelInfo(openAIReq.Model)
	// 1. 先获取该模型的Credit
	if ok {
		// 2. 从token表中获取该credit大于该模型的token值
		cookieRecord := &model.Cookie{
			Credit: modelInfo.Credit,
		}
		if modelInfo.Type == "STANDARD" {
			cookies, err = cookieRecord.FindByMinimumCredit(database.DB)
		} else {
			cookies, err = cookieRecord.FindByMinimumCreditAdvanced(database.DB)
		}
		if err != nil {
			c.JSON(500, gin.H{"error": "no token"})
			return
		}
		if len(cookies) == 0 {
			c.JSON(500, gin.H{"error": errNoValidCookies})
			return
		}
	} else {
		c.JSON(500, gin.H{"error": "no model"})
		return
	}

	newChatFlag := false
	var hixChatId string
	// 1. 获取符合messagehash的tokens
	pair, b, err := openAIReq.GetPreviousMessagePair()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !b {
		// 没有last问答对 视为新会话
		newChatFlag = true
	} else {
		msgPairSha256 := common.StringToSHA256(strings.TrimSpace(pair))
		cookie, chatId, err := model.QueryCookiesByChatHashAndModelAndCredit(database.DB, msgPairSha256, openAIReq.Model, modelInfo.Credit)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newChatFlag = true
		} else if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		} else {
			cookies = []model.Cookie{
				{
					Cookie: cookie,
				},
			}
			hixChatId = chatId
		}
	}

	responseId := fmt.Sprintf(responseIDFormat, time.Now().Format("20060102150405"))
	ctx := c.Request.Context()

	mutable.Shuffle(cookies)

	maxRetries := len(cookies)

	var messagesPair []model.OpenAIChatMessage
	messagesPair = append(messagesPair, openAIReq.Messages[len(openAIReq.Messages)-1])
	c.Stream(func(w io.Writer) bool {
		for attempt := 0; attempt < maxRetries; attempt++ {
			cookie := cookies[attempt]
			if newChatFlag {
				chatId, err := hixapi.MakeCreateChatRequest(client, cookie.Cookie, modelInfo.ModelID)
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return false
				}
				hixChatId = chatId
			}

			requestBody, err := createRequestBody(c, hixChatId, &openAIReq, searchType, modelInfo)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return false
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to marshal request body"})
				return false
			}
			sseChan, err := hixapi.MakeStreamChatRequest(c, client, openAIReq.Model, hixChatId, jsonData, cookie.Cookie)
			if err != nil {
				logger.Errorf(ctx, "MakeStreamChatRequest err on attempt %d: %v", attempt+1, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return false
			}

			isRateLimit := false
			var assistantMsgContent string
			thinkStartType := new(bool) // 初始值为false
		SSELoop:
			for response := range sseChan {

				if response.Status == 403 {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Forbidden"})
					return false
				}

				if response.Done {
					logger.Warnf(ctx, response.Data)
					return false
				}

				data := response.Data
				if data == "" {
					continue
				}

				logger.Debug(ctx, strings.TrimSpace(data))

				switch {
				case common.IsCloudflareChallenge(data):
					c.JSON(http.StatusInternalServerError, gin.H{"error": "cf challenge"})
					return false
				case common.IsNotLogin(data):
					isRateLimit = true
					logger.Warnf(ctx, "Cookie Not Login, switching to next cookie, attempt %d/%d, COOKIE:%s", attempt+1, maxRetries, cookie)
					// 删除cookie
					//config.RemoveCookie(cookie)
					break SSELoop // 使用 label 跳出 SSE 循环
				}

				delta, shouldContinue := processStreamData(c, data, responseId, openAIReq.Model, jsonData, thinkStartType)
				// 处理事件流数据

				if !shouldContinue {
					// 保存chat记录
					messagesPair = append(messagesPair, model.OpenAIChatMessage{
						Role:    "assistant",
						Content: strings.TrimSpace(assistantMsgContent),
					})
					bytes, err := json.Marshal(messagesPair)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return false
					}
					messagesPairStr := strings.NewReplacer(
						`\n`, "",
						`\t`, "",
						`\r`, "",
					).Replace(string(bytes))

					chat := model.Chat{
						Cookie:                     cookie.Cookie,
						Model:                      openAIReq.Model,
						CookieHash:                 cookie.CookieHash,
						HixChatId:                  hixChatId,
						LastMessagesPair:           string(bytes),
						LastMessagesPairSha256Hash: common.StringToSHA256(messagesPairStr),
					}
					if newChatFlag {
						if err := chat.Create(database.DB); err != nil {
							c.JSON(500, gin.H{"error": err.Error()})
						}
					} else {
						// 更新chat记录
						if err := chat.UpdateLastMessages(database.DB); err != nil {
							c.JSON(500, gin.H{"error": err.Error()})
						}
					}

					go func() {
						isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, cookie.Cookie)
						if err != nil {
							logger.Errorf(ctx, "MakeSubUsageRequest err: %v", err)
						}
						cookieRecord := &model.Cookie{
							CookieHash:     cookie.CookieHash,
							Credit:         credit,
							AdvancedCredit: advancedCredit,
							IsActiveSub:    isActiveSub,
						}
						err = cookieRecord.UpdateCreditByCookieHash(database.DB)
						if err != nil {
							logger.Errorf(ctx, "UpdateCreditByCookieHash err: %v", err)
						}
					}()

					return false
				} else {
					//if strings.TrimSpace(delta) != "" {
					assistantMsgContent = assistantMsgContent + delta
					//}
				}
			}

			if !isRateLimit {
				return true
			}

		}

		logger.Errorf(ctx, "All cookies exhausted after %d attempts", maxRetries)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "All cookies are temporarily unavailable."})
		return false
	})
}

// 处理流式数据的辅助函数，返回bool表示是否继续处理
func processStreamData(c *gin.Context, data string, responseId, model string, jsonData []byte, thinkStartType *bool) (string, bool) {
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "data: ")

	if data == "[DONE]" {
		handleMessageResult(c, responseId, model, jsonData)
		return "", false
	}

	if !strings.HasPrefix(data, "{\"content\":") &&
		!strings.HasPrefix(data, "{\"reasoning_content\":") &&
		!strings.HasPrefix(data, "{\"thinking_time\":") {
		return "", true
	}

	var event map[string]interface{}
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		logger.Errorf(c.Request.Context(), "Failed to unmarshal event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", false
	}
	delta, ok := event["content"].(string)
	if ok {
		if err := handleDelta(c, delta, responseId, model, jsonData); err != nil {
			logger.Errorf(c.Request.Context(), "handleDelta err: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", false
		}
		return delta, true
	}
	delta, ok = event["reasoning_content"].(string)
	if ok {
		if !*thinkStartType {
			delta = "<think>\n" + delta
			*thinkStartType = true
		}
		if err := handleDelta(c, delta, responseId, model, jsonData); err != nil {
			logger.Errorf(c.Request.Context(), "handleDelta err: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", false
		}
		return delta, true
	}
	_, ok = event["thinking_time"].(float64)
	if ok {
		delta = "\n</think>"
		if err := handleDelta(c, delta, responseId, model, jsonData); err != nil {
			logger.Errorf(c.Request.Context(), "handleDelta err: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return "", false
		}
		return delta, true
	}

	return "", true

}

func processNoStreamData(c *gin.Context, data string, responseId, model string, jsonData []byte, thinkStartType *bool) (string, bool) {
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "data: ")

	if data == "[DONE]" {
		return "", false
	}

	if !strings.HasPrefix(data, "{\"content\":") &&
		!strings.HasPrefix(data, "{\"reasoning_content\":") &&
		!strings.HasPrefix(data, "{\"thinking_time\":") {
		return "", true
	}

	var event map[string]interface{}
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		logger.Errorf(c.Request.Context(), "Failed to unmarshal event: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", false
	}
	delta, ok := event["content"].(string)
	if ok {
		return delta, true
	}
	delta, ok = event["reasoning_content"].(string)
	if ok {
		if !*thinkStartType {
			delta = "<think>\n" + delta
			*thinkStartType = true
		}
		return delta, true
	}
	_, ok = event["thinking_time"].(float64)
	if ok {
		delta = "\n</think>"
		return delta, true
	}

	return "", true

}

func OpenaiModels(c *gin.Context) {
	var modelsResp []string

	maxCookies, err := (&model.Cookie{}).FindMaxCreditByActiveSub(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 提取两种状态的credit
	var standardCredit, advancedCredit int
	hasStandard := false
	hasAdvanced := false
	for _, cookie := range maxCookies {
		if !cookie.IsActiveSub {
			standardCredit = cookie.Credit
			hasStandard = true
		} else {
			advancedCredit = cookie.AdvancedCredit
			hasAdvanced = true
		}
	}

	// 遍历modelRegistry，收集符合条件的模型
	modelsResp = make([]string, 0)
	for modelName, info := range common.ModelRegistry {
		credit := info.Credit
		modelType := info.Type

		if modelType == "STANDARD" && hasStandard && standardCredit >= credit {
			modelsResp = append(modelsResp, modelName)
		}
		if modelType == "ADVANCED" && hasAdvanced && advancedCredit >= credit {
			modelsResp = append(modelsResp, modelName)
		}
	}

	var openaiModelListResponse model.OpenaiModelListResponse
	var openaiModelResponse []model.OpenaiModelResponse
	openaiModelListResponse.Object = "list"

	for _, modelResp := range modelsResp {
		openaiModelResponse = append(openaiModelResponse, model.OpenaiModelResponse{
			ID:     modelResp,
			Object: "model",
		})
	}
	openaiModelListResponse.Data = openaiModelResponse
	c.JSON(http.StatusOK, openaiModelListResponse)
	return
}

func safeClose(client cycletls.CycleTLS) {
	if client.ReqChan != nil {
		close(client.ReqChan)
	}
	if client.RespChan != nil {
		close(client.RespChan)
	}
}
