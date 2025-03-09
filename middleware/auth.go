package middleware

import (
	"github.com/gin-gonic/gin"
	"hixai2api/common"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/model"
	"net/http"
	"strings"
)

func isValidSecret(secret string) (error, bool) {
	apiKey := model.ApiKey{
		ApiKey: secret,
	}

	apiKeys, err := apiKey.GetAll(database.DB)
	if err != nil {
		return err, false
	}

	if len(apiKeys) == 0 {
		return nil, true
	}

	count, err := apiKey.CountByKey(database.DB)
	if err != nil {
		return err, false
	}
	if count > 0 {
		return nil, true
	} else {
		return nil, false
	}
}

func isValidBackendSecret(secret string) bool {
	return config.BackendSecret != "" && !(config.BackendSecret == secret)
}

func authHelperForOpenai(c *gin.Context) {
	secret := c.Request.Header.Get("Authorization")
	secret = strings.Replace(secret, "Bearer ", "", 1)

	err, b := isValidSecret(secret)
	if err != nil {
		common.SendResponse(c, http.StatusUnauthorized, 1, "unauthorized", "")
		c.Abort()
		return
	}
	if !b {
		c.JSON(http.StatusUnauthorized, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "API-KEY校验失败",
				Type:    "invalid_request_error",
				Code:    "invalid_authorization",
			},
		})
		c.Abort()
		return
	}

	//if config.ApiSecret == "" {
	//	c.Request.Header.Set("Authorization", "")
	//}

	c.Next()
	return
}

func authHelperForBackend(c *gin.Context) {
	secret := c.Request.Header.Get("Authorization")
	secret = strings.Replace(secret, "Bearer ", "", 1)
	if isValidBackendSecret(secret) {
		logger.Debugf(c.Request.Context(), "BackendSecret is not empty, but not equal to %s", secret)
		common.SendResponse(c, http.StatusUnauthorized, 1, "unauthorized", "")
		c.Abort()
		return
	}

	if config.BackendSecret == "" {
		c.Request.Header.Set("Authorization", "")
	}

	c.Next()
	return
}

func OpenAIAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelperForOpenai(c)
	}
}

func BackendAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelperForBackend(c)
	}
}
