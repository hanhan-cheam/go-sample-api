package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Config struct {
	DB *gorm.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "sample"
)

func Connect() *Config {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database")
	} else {
		fmt.Println("connect successful")
	}
	return &Config{
		DB: db,
	}
}
