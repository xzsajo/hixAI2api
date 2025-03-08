package model

import (
	"gorm.io/gorm"
	"hixai2api/common"
	"time"
)

type Cookie struct {
	Id         string    `json:"id" gorm:"type:varchar(64);not null;primaryKey"`
	Cookie     string    `json:"cookie" gorm:"type:text"`
	CookieHash string    `json:"cookie_hash" gorm:"type:varchar(255);not null;index"`
	Credit     int       `json:"credit" gorm:"type:bigint(20);not null"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;not null"`
}

func (c *Cookie) Create(db *gorm.DB) error {
	if c.Id == "" {
		id, err := common.NextID()
		if err != nil {
			return err
		}
		c.Id = id
		c.CreateTime = time.Now()
	}
	result := db.Create(c)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *Cookie) FindAllCookies(db *gorm.DB) ([]Cookie, error) {
	var cookies []Cookie
	result := db.Find(&cookies)
	if result.Error != nil {
		return nil, result.Error
	}
	return cookies, nil
}

func (t *Cookie) FindByMinimumCredit(db *gorm.DB, minCredit int) ([]Cookie, error) {
	var cookies []Cookie
	result := db.Where("credit >= ?", minCredit).Find(&cookies)
	if result.Error != nil {
		return nil, result.Error
	}
	return cookies, nil
}

// 定义用于接收查询结果的结构体
type ChatResult struct {
	Cookie    string
	HixChatId string
}

func QueryCookiesByChatHashAndModelAndCredit(db *gorm.DB, lastMessagesPairSha256Hash, modelName string, creditLimit int) (string, string, error) {
	var result ChatResult

	// 从Chat表出发进行查询，关联Cookie表
	err := db.Model(&Chat{}).
		Select("chats.cookie, chats.hix_chat_id").
		Joins("JOIN cookies ON chats.cookie_hash = cookies.cookie_hash").
		Where("chats.last_messages_pair_sha256_hash = ?  AND chats.model = ? AND cookies.credit >= ?",
					lastMessagesPairSha256Hash, modelName, creditLimit).
		First(&result).Error // 使用First获取第一条记录

	if err != nil {
		return "", "", err
	}

	return result.Cookie, result.HixChatId, nil
}

func (c *Cookie) UpdateCreditByCookieHash(db *gorm.DB, cookieHash string, newCredit int) error {
	// 使用 GORM 的 Model 方法指定模型，并使用 Where 方法指定条件
	result := db.Model(&Cookie{}).Where("cookie_hash = ?", cookieHash).Update("credit", newCredit)
	if result.Error != nil {
		return result.Error
	}

	// 检查是否有记录被更新
	//if result.RowsAffected == 0 {
	//	return fmt.Errorf("no record found with cookie_hash: %s", cookieHash)
	//}

	return nil
}
