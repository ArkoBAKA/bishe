package storage

import (
	"fmt"
	"time"

	"server/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg config.Config) (*gorm.DB, error) {
	dsn := cfg.MySQL.DSN()
	if dsn == "" {
		return nil, fmt.Errorf("mysql dsn is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
