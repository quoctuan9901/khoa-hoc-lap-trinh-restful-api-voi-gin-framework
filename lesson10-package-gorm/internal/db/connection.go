package db

import (
	"context"
	"fmt"
	"hoc-gin/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	connStr := config.NewConfig().DNS()

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: connStr,
	}), config)

	if err != nil {
		return fmt.Errorf("error opening DB connection: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return fmt.Errorf("DB ping error: %w", err)
	}

	log.Println("Connected")

	return nil
}
