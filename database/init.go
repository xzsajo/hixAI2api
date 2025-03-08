package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hixai2api/common/config"
	"hixai2api/model"
)

var (
	DB *gorm.DB
)

func InitDB() {
	var err error
	DB, err = ConfigDB(config.MysqlDsn)
	if err != nil {
		logrus.Fatal("连接数据库失败:", err)
	}

	err = DB.AutoMigrate(&model.Chat{})
	if err != nil {
		logrus.Fatal("自动迁移表结构失败:", err)
	}
	err = DB.AutoMigrate(&model.ApiKey{})
	if err != nil {
		logrus.Fatal("自动迁移表结构失败:", err)
	}
	err = DB.AutoMigrate(&model.Cookie{})
	if err != nil {
		logrus.Fatal("自动迁移表结构失败:", err)
	}

}
