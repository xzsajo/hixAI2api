package database

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"hixai2api/common"
	"hixai2api/common/config"
	"hixai2api/common/env"

	logger "hixai2api/common/loggger"
	"hixai2api/model"
	"log"
	"os"
	"time"
)

var (
	DB     *gorm.DB
	LOG_DB *gorm.DB
)

func InitDB() {
	var err error
	//DB, err = ConfigDB(config.MysqlDsn)

	DB, err = chooseDB("MYSQL_DSN")
	if err != nil {
		logger.FatalLog("failed to initialize database: " + err.Error())
		return
	}

	setDBConns(DB)

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

func chooseDB(envName string) (*gorm.DB, error) {
	dsn := os.Getenv(envName)

	switch {
	//case strings.HasPrefix(dsn, "postgres://"):
	//	// Use PostgreSQL
	//	return openPostgreSQL(dsn)
	case dsn != "":
		// Use MySQL
		return openMySQL(dsn)
	default:
		// Use SQLite
		return openSQLite()
	}
}

func openSQLite() (*gorm.DB, error) {
	logger.SysLog("SQL_DSN not set, using SQLite as database")
	common.UsingSQLite = true
	dsn := fmt.Sprintf("%s?_busy_timeout=%d", common.SQLitePath, common.SQLiteBusyTimeout)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
}

func openMySQL(dsn string) (*gorm.DB, error) {
	newLogger := gormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlog.Config{
			SlowThreshold:             time.Second,    // Slow SQL threshold
			LogLevel:                  gormlog.Silent, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,           // Don't include params in the SQL log
			Colorful:                  false,          // Disable color
		},
	)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logrus.WithField("err", err).Fatal("连接数据库失败")
		return nil, err
	}

	return db, nil
}

func setDBConns(db *gorm.DB) *sql.DB {
	if config.DebugSQLEnabled {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.FatalLog("failed to connect database: " + err.Error())
		return nil
	}

	sqlDB.SetMaxIdleConns(env.Int("SQL_MAX_IDLE_CONNS", 100))
	sqlDB.SetMaxOpenConns(env.Int("SQL_MAX_OPEN_CONNS", 1000))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(env.Int("SQL_MAX_LIFETIME", 60)))
	return sqlDB
}

func CloseDB() error {
	if LOG_DB != DB {
		err := closeDB(LOG_DB)
		if err != nil {
			return err
		}
	}
	return closeDB(DB)
}

func closeDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
