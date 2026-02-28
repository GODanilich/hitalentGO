package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func Open(databaseURL string, gormLogLevel string) (*gorm.DB, error) {
	level := parseGormLogLevel(gormLogLevel)

	baseLogger := log.New(os.Stdout, "gorm ", log.LstdFlags)

	lg := glogger.New(
		baseLogger,
		glogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: lg,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func parseGormLogLevel(s string) glogger.LogLevel {
	switch s {
	case "silent":
		return glogger.Silent
	case "error":
		return glogger.Error
	case "info":
		return glogger.Info
	default:
		return glogger.Warn
	}
}
