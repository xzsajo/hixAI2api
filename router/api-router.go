package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hixai2api/common/config"
	"hixai2api/controller"
	"hixai2api/middleware"
	"strings"
)

func SetApiRouter(router *gin.Engine) {
	router.Use(middleware.CORS())
	//router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.IPBlacklistMiddleware())
	router.Use(middleware.RequestRateLimit())

	router.GET("/")

	//https://api.openai.com/v1/images/generations
	v1Router := router.Group(fmt.Sprintf("%s/v1", ProcessPath(config.RoutePrefix)))
	v1Router.Use(middleware.OpenAIAuth())
	v1Router.POST("/chat/completions", controller.ChatForOpenAI)
	//v1Router.POST("/images/generations", controller.ImagesForOpenAI)
	v1Router.GET("/models", controller.OpenaiModels)
}

func ProcessPath(path string) string {
	// 判断字符串是否为空
	if path == "" {
		return ""
	}

	// 判断开头是否为/，不是则添加
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// 判断结尾是否为/，是则去掉
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	return path
}
