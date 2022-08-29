package focus_go

import (
	"GoFocusMicroService/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func SetUp() {
	newLog := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		IgnoreRecordNotFoundError: true,
	})
	db, _ = gorm.Open(mysql.Open(conf.Conf.DataBase.FocusGoDB), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		DisableAutomaticPing: true,
		Logger:               newLog,
	})
	// 设置连接最大超时时间
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetConnMaxLifetime(800 * time.Second)
	}
	return
}
