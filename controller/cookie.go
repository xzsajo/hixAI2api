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
	"strings"
)

// SaveCookie @Summary 保存COOKIE
// @Description 保存COOKIE
// @Tags COOKIE
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization BACKEND_SECRET"
// @Param req body model.CookieSaveReq true "COOKIE信息"
// @Router /api/cookie [put]
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

	if !strings.Contains(req.Cookie, "__Secure-next-auth.session-token=") {
		req.Cookie = "__Secure-next-auth.session-token=" + req.Cookie
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
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 1, "COOKIE already exists", "")
		return
	}

	client := cycletls.Init()
	defer safeClose(client)

	// 校验cookie
	isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	cookie.Credit = credit
	cookie.IsActiveSub = isActiveSub
	cookie.AdvancedCredit = advancedCredit

	err = cookie.Create(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

// DeleteCookie @Summary 删除COOKIE
// @Description 删除COOKIE
// @Tags COOKIE
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization BACKEND_SECRET"
// @Param id path string true "COOKIE ID"
// @Router /api/cookie/{id} [delete]
func DeleteCookie(c *gin.Context) {
	id := c.Param("id") // 获取 URL 中的 id 参数
	if id == "" {
		common.SendResponse(c, http.StatusBadRequest, 1, "Invalid ID", "")
		return
	}

	cookie := model.Cookie{
		Id: id,
	}
	err := cookie.DeleteById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

// UpdateCookie @Summary 更新COOKIE
// @Description 更新COOKIE
// @Tags COOKIE
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization BACKEND_SECRET"
// @Param req body model.CookieUpdateReq true "COOKIE信息"
// @Router /api/cookie/update [post]
func UpdateCookie(c *gin.Context) {
	var req model.CookieUpdateReq
	if err := c.BindJSON(&req); err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusBadRequest, 1, err.Error(), "")
		return
	}

	if req.Id == "" {
		common.SendResponse(c, http.StatusBadRequest, 1, "Invalid ID", "")
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
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	if exist {
		common.SendResponse(c, http.StatusBadRequest, 1, "COOKIE already exists", "")
		return
	}

	client := cycletls.Init()
	defer safeClose(client)
	// 校验cookie
	isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, req.Cookie)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	cookie.Credit = credit
	cookie.IsActiveSub = isActiveSub
	cookie.AdvancedCredit = advancedCredit

	err = cookie.UpdateKeyById(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}
	common.SendResponse(c, http.StatusOK, 0, "success", "")
	return
}

// GetAllCookie @Summary 查询全量COOKIE
// @Description 查询全量COOKIE
// @Tags COOKIE
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization BACKEND_SECRET"
// @Success 200 {object} common.ResponseResult{data=[]model.CookieResp} "成功"
// @Router /api/cookie/all [get]
func GetAllCookie(c *gin.Context) {
	var resp []model.CookieResp

	cookie := model.Cookie{}
	cookies, err := cookie.GetAll(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
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

// RefreshCookieCredit @Summary 同步更新全量COOKIE额度
// @Description 同步更新全量COOKIE额度
// @Tags COOKIE
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization BACKEND_SECRET"
// @Router /api/cookie/credit/refresh [post]
func RefreshCookieCredit(c *gin.Context) {

	cookie := model.Cookie{}
	cookies, err := cookie.GetAll(database.DB)
	if err != nil {
		logger.Errorf(c.Request.Context(), err.Error())
		common.SendResponse(c, http.StatusInternalServerError, 1, err.Error(), "")
		return
	}

	if len(cookies) > 0 {
		go func() {
			client := cycletls.Init()
			defer safeClose(client)
			for _, k := range cookies {
				isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, k.Cookie)
				if err != nil {
					logger.Errorf(c.Request.Context(), err.Error())
					continue
				}
				k.Credit = credit
				k.IsActiveSub = isActiveSub
				k.AdvancedCredit = advancedCredit
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
