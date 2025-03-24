package job

import (
	"fmt"
	"github.com/deanxv/CycleTLS/cycletls"
	"hixai2api/common/config"
	logger "hixai2api/common/loggger"
	"hixai2api/database"
	"hixai2api/hixapi"
	"hixai2api/model"
	"time"
)

func DelChatTask() {
	if config.ChatMaxDays < 0 {
		return
	}

	client := cycletls.Init()
	defer safeClose(client)
	for {
		logger.SysLog("hixai2api Scheduled DelChatTask Task Job Start!")

		chat := &model.Chat{}
		chats, err := chat.FindOlderThan(database.DB, config.ChatMaxDays)
		if err != nil {
			logger.SysError(fmt.Sprintf("FindOlderThan err: %v Id: %s", err, chat.Id))
		}

		for _, chat := range chats {
			err := chat.DeleteById(database.DB)
			if err != nil {
				logger.SysError(fmt.Sprintf("DeleteById err: %v Id: %s", err, chat.Id))
			}
			err = hixapi.MakeDelChatRequest(client, chat.Cookie, chat.HixChatId)
			if err != nil {
				logger.SysError(fmt.Sprintf("MakeDelChatRequest err: %v Id: %s", err, chat.Id))
			}
		}

		logger.SysLog("hixai2api Scheduled DelChatTask Task Job  End!")
		time.Sleep(60 * time.Minute)
	}
}
