package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/congziqi77/task-scheduling/global"
	"github.com/congziqi77/task-scheduling/internal/modules/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DB_PING_INTERVAL = 90 * time.Second
	DB_Max_LIFT_TIME = 2 * time.Hour
)

var DB *gorm.DB

//创建数据库实例
func NewDBEngine() (*gorm.DB, error) {
	dns := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		global.DbSetting.User,
		global.DbSetting.Password,
		global.DbSetting.Host,
		global.DbSetting.Port,
		global.DbSetting.Database)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(global.DbSetting.MaxOpenConns)
	sqlDB.SetMaxIdleConns(global.DbSetting.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(DB_Max_LIFT_TIME)

	//ping实例
	go keepPing(sqlDB)
	return db, nil
}

func keepPing(sqlDB *sql.DB) {
	t := time.Tick(DB_PING_INTERVAL)
	for {
		select {
		case <-t:
			err := sqlDB.Ping()
			if err != nil {
				logger.Printf("database ping: %s", err)
			}
		}
	}
}
