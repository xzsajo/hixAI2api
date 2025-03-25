package job

import (
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/hixapi"
	"hixai2api/model"
	"time"
)

func UpdateCookieCreditTask() {
	client := cycletls.Init()
	defer safeClose(client)
	for {
		logger.SysLog("hixai2api Scheduled UpdateCookieCreditTask Task Job Start!")
		cookieRecord := model.Cookie{}
		cookies, err := cookieRecord.FindAllCookies(database.DB)
		if err != nil {
			logger.SysError(fmt.Sprintf("FindAllCookies err: %v", err))
		}
		if len(cookies) != 0 {
			for _, cookie := range cookies {
				isActiveSub, credit, advancedCredit, err := hixapi.MakeSubUsageRequest(client, cookie.Cookie)
				if err != nil {
					logger.SysError(fmt.Sprintf("UpdateCookieCreditTask err: %v", err))
				}
				cookieRecord := &model.Cookie{
					CookieHash:     cookie.CookieHash,
					Credit:         credit,
					AdvancedCredit: advancedCredit,
					IsActiveSub:    isActiveSub,
				}
				err = cookieRecord.UpdateCreditByCookieHash(database.DB)
				if err != nil {
					logger.SysError(fmt.Sprintf("UpdateCreditByCookieHash err: %v cookie: %s", err, cookie.Cookie))
				}
			}
		}
		logger.SysLog("hixai2api Scheduled UpdateCookieCreditTask Task Job End!")

		// 计算到下一个整点的时间
		now := time.Now()
		next := now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		time.Sleep(next.Sub(now))
	}
}
func safeClose(client cycletls.CycleTLS) {
	if client.ReqChan != nil {
		close(client.ReqChan)
	}
	if client.RespChan != nil {
		close(client.RespChan)
	}
}
