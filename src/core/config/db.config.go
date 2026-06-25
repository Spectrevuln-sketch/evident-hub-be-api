package config

import (
	"evidence-hub-be/src/core/config/envs"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(cfg *envs.Config) {
	var err error
	log.Printf("ENV DB %s", cfg.DbConfig.DbHost)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DbConfig.DbHost,
		cfg.DbConfig.DbUser,
		cfg.DbConfig.DbPassword,
		cfg.DbConfig.DbName,
		cfg.DbConfig.DbPort,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Cannot Connect to database !")
		os.Exit(2)
	}

	db, _ := DB.DB()
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	err = db.Ping()
	if err != nil {
		log.Fatal("Unable to ping database")
		os.Exit(1)
	}
	log.Println("Database Connected 1")
}
