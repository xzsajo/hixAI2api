package controller

import (
	"github.com/gin-gonic/gin"
	"hixai2api/common"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/model"
	"net/http"
)

func AuthVerify(c *gin.Context) {

	var req model.AuthVerifyReq
	if err := c.BindJSON(&req); err != nil {
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
	common.SendResponse(c, http.StatusOK, 0, "success", "")
}

func SaveApiKey(c *gin.Context) {

	var req model.ApiKeySaveReq
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}

	apiKey := model.ApiKey{
		ApiKey: req.ApiKey,
		Remark: req.Remark,
	}

	exist, err := apiKey.Exist(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 0, "API-KEY already exists", "")
		return
	}

	err = apiKey.Create(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func DeleteApiKey(c *gin.Context) {
	id := c.Param("id") // 获取 URL 中的 id 参数
	if id == "" {
		common.SendResponse(c, http.StatusBadRequest, 0, "Invalid ID", "")
		return
	}

	apiKey := model.ApiKey{
		Id: id,
	}
	err := apiKey.DeleteById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func UpdateApiKey(c *gin.Context) {
	var req model.ApiKeyUpdateReq
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusBadRequest, 0, err.Error(), "")
		return
	}

	if req.Id == "" {
		common.SendResponse(c, http.StatusBadRequest, 0, "Invalid ID", "")
		return
	}

	apiKey := model.ApiKey{
		Id:     req.Id,
		ApiKey: req.ApiKey,
		Remark: req.Remark,
	}

	exist, err := apiKey.ExistsNotMe(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 0, "API-KEY already exists", "")
		return
	}

	err = apiKey.UpdateKeyById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func GetAllApiKey(c *gin.Context) {
	var resp []model.ApiKeyResp

	apiKey := model.ApiKey{}
	apiKeys, err := apiKey.GetAll(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if len(apiKeys) > 0 {
		for _, k := range apiKeys {
			resp = append(resp, model.ApiKeyResp{
				Id:         k.Id,
				ApiKey:     k.ApiKey,
				CreateTime: k.CreateTime.Format("2006-01-02 15:04:05"),
			})
		}
	}
	common.SendResponse(c, http.StatusOK, 0, "success", resp)
	return
}
