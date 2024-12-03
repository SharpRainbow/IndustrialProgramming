package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func GetMapper(cfg *Config) (*gorm.DB, error) {
	if db != nil {
		dbase, err := db.DB()
		if err == nil && dbase.Ping() == nil {
			return db, nil
		}
	}
	var err error
	db, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbName, cfg.DbPassword)), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func CloseMapper() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get SQL DB instance: %v", err)
			return
		}
		sqlDB.Close()
	}
}
