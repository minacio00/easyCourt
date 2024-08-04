package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func Init() {
	var err error

	// Build the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
	)

	// Connect to the database
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	DB.AutoMigrate(&model.Tenant{})

}
