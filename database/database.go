package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	fullDsn := fmt.Sprintf(
		"%s?charset=utf8mb4&parseTime=True&loc=Local",
		dsn,
	)

	db, err := gorm.Open(mysql.Open(fullDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// sqlDB.SetMaxIdleConns(100)
	// sqlDB.SetMaxOpenConns(100)

	return db, nil
}
