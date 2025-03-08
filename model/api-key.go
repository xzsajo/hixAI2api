package model

import (
	"gorm.io/gorm"
	"hixai2api/common"
	"time"
)

type ApiKey struct {
	Id         string    `json:"id" gorm:"type:varchar(64);not null;primaryKey"`
	Key        string    `json:"key" gorm:"type:varchar(255);not null;index"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;not null"`
}

func (c *ApiKey) Create(db *gorm.DB) error {
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
