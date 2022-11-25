package database

import (
	"exercise/internal/app/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnDatabasePostgres() *gorm.DB {

	// https://github.com/go-gorm/postgres
	dsnMaster := os.Getenv("DATABASE_URL")
	dbMaster, errMaster := gorm.Open(postgres.Open(dsnMaster), &gorm.Config{})
	if errMaster != nil {
		log.Panic(errMaster)
	}

	if errMaster = dbMaster.AutoMigrate(
		&domain.User{},
		&domain.Exercise{},
		&domain.Answer{},
		&domain.Question{},
	); errMaster != nil {
		log.Panic(errMaster)
	}

	return dbMaster
}
