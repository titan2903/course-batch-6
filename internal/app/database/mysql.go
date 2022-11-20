package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewConnDatabase() *gorm.DB {
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWD")
	dbUser := os.Getenv("DB_USER")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
