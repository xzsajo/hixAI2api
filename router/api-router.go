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
	router.Use(middleware.IPBlacklistMiddleware())
	router.Use(middleware.RequestRateLimit())

	v1Router := router.Group(fmt.Sprintf("%s/v1", ProcessPath(config.RoutePrefix)))
	v1Router.Use(middleware.OpenAIAuth())
	v1Router.POST("/chat/completions", controller.ChatForOpenAI)
	//v1Router.POST("/images/generations", controller.ImagesForOpenAI)
	v1Router.GET("/models", controller.OpenaiModels)

	if config.BackendApiEnable == 1 {
		apiRouter := router.Group(fmt.Sprintf("/api"))
		apiRouter.Use(middleware.BackendAuth())
		//apiRouter.POST("/auth/verify", controller.AuthVerify)
		apiRouter.PUT("/key", controller.SaveApiKey)
		apiRouter.DELETE("/key/:id", controller.DeleteApiKey)
		apiRouter.POST("/key/update", controller.UpdateApiKey)
		apiRouter.GET("/key/all", controller.GetAllApiKey)

		apiRouter.PUT("/cookie", controller.SaveCookie)
		apiRouter.DELETE("/cookie/:id", controller.DeleteCookie)
		apiRouter.POST("/cookie/update", controller.UpdateCookie)
		apiRouter.GET("/cookie/all", controller.GetAllCookie)
		apiRouter.POST("/cookie/credit/refresh", controller.RefreshCookieCredit)
	}

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
