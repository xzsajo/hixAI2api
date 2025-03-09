package controller

import (
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"github.com/gin-gonic/gin"
	"hixai2api/common"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/hixapi"
	"hixai2api/model"
	"net/http"
	"net/url"
)

func SaveCookie(c *gin.Context) {

	var req model.CookieSaveReq
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
	var err error
	req.Cookie, err = url.QueryUnescape(req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), fmt.Sprintf("cookie QueryUnescape err:"), err)
		return
	}

	cookie := model.Cookie{
		Cookie:     req.Cookie,
		Remark:     req.Remark,
		CookieHash: common.StringToSHA256(req.Cookie),
	}

	exist, err := cookie.Exist(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 0, "COOKIE already exists", "")
		return
	}

	client := cycletls.Init()
	defer safeClose(client)

	// 校验cookie
	credit, err := hixapi.MakeSubUsageRequest(client, req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	cookie.Credit = credit

	err = cookie.Create(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func DeleteCookie(c *gin.Context) {
	id := c.Param("id") // 获取 URL 中的 id 参数
	if id == "" {
		common.SendResponse(c, http.StatusBadRequest, 0, "Invalid ID", "")
		return
	}

	cookie := model.Cookie{
		Id: id,
	}
	err := cookie.DeleteById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func UpdateCookie(c *gin.Context) {
	var req model.CookieUpdateReq
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusBadRequest, 0, err.Error(), "")
		return
	}

	if req.Id == "" {
		common.SendResponse(c, http.StatusBadRequest, 0, "Invalid ID", "")
		return
	}
	var err error
	req.Cookie, err = url.QueryUnescape(req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), fmt.Sprintf("cookie QueryUnescape err:"), err)
		return
	}

	cookie := model.Cookie{
		Id:         req.Id,
		Cookie:     req.Cookie,
		Remark:     req.Remark,
		CookieHash: common.StringToSHA256(req.Cookie),
	}

	exist, err := cookie.ExistsNotMe(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 0, "COOKIE already exists", "")
		return
	}

	client := cycletls.Init()
	defer safeClose(client)
	// 校验cookie
	credit, err := hixapi.MakeSubUsageRequest(client, req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	cookie.Credit = credit

	err = cookie.UpdateKeyById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

func GetAllCookie(c *gin.Context) {
	var resp []model.CookieResp

	cookie := model.Cookie{}
	cookies, err := cookie.GetAll(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}
	if len(cookies) > 0 {
		for _, k := range cookies {
			resp = append(resp, model.CookieResp{
				Id:         k.Id,
				Cookie:     k.Cookie,
				Credit:     k.Credit,
				Remark:     k.Remark,
				CreateTime: k.CreateTime.Format("2006-01-02 15:04:05"),
			})
		}
	}
	common.SendResponse(c, http.StatusOK, 0, "success", resp)
	return
}

func RefreshCookieCredit(c *gin.Context) {

	cookie := model.Cookie{}
	cookies, err := cookie.GetAll(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 0, err.Error(), "")
		return
	}

	if len(cookies) > 0 {
		go func() {
			client := cycletls.Init()
			defer safeClose(client)
			for _, k := range cookies {
				credit, err := hixapi.MakeSubUsageRequest(client, k.Cookie)
				if err != nil {
					logger.Errorf(c.Request.Context(), err.Error())
					continue
				}
				k.Credit = credit
				err = k.UpdateCreditByCookieHash(database.DB)
				if err != nil {
					logger.Errorf(c.Request.Context(), err.Error())
					continue
				}
			}
		}()
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}
