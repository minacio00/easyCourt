package config

import (
	"log"

	"github.com/spf13/viper"
)

var HOST = ""
var USER = ""
var PASSWORD = ""
var DBNAME = ""
var PORT = ""

func LdVars() {
	viper.SetConfigFile(".env")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error reading .env file, %s", err)
	}

	HOST = viper.GetString("HOST")
	USER = viper.GetString("USER")
	PASSWORD = viper.GetString("PASSWORD")
	DBNAME = viper.GetString("DBNAME")
	PORT = viper.GetString("PORT")
}
