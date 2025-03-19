// @title HIX-AI-2API
// @version 1.0.0
// @description HIX-AI-2API
// @BasePath
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hixai2api/check"
	"hixai2api/common"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/job"
	"hixai2api/middleware"
	"hixai2api/model"
	"hixai2api/router"
	"os"
	"strconv"
)

//var buildFS embed.FS

func main() {
	logger.SetupLogger()
	logger.SysLog(fmt.Sprintf("hixai2api %s starting...", common.Version))
	database.InitDB()
	defer func() {
		err := database.CloseDB()
		if err != nil {
			logger.FatalLog("failed to close database: " + err.Error())
		}
	}()
	check.CheckEnvVariable()

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	var err error

	model.InitTokenEncoders()
	go job.UpdateCookieCreditTask()
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)

	router.SetRouter(server)
	//router.SetRouter(server, buildFS)

	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}

	if config.DebugEnabled {
		logger.SysLog("running in DEBUG mode.")
	}

	logger.SysLog("hixai2api start success. enjoy it! ^_^\n")

	err = server.Run(":" + port)

	if err != nil {
		logger.FatalLog("failed to start HTTP server: " + err.Error())
	}

}
