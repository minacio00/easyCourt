package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/minacio00/easyCourt/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Init() {
	once.Do(func() {
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
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to the database: %v", err)
		}
		if err := DB.AutoMigrate(
			&model.User{},
			&model.Location{},
			&model.Court{}, &model.Timeslot{}, &model.Booking{},
		); err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}

	})
}

func GetDB() *gorm.DB {
	return DB
}
