package config

import "github.com/spf13/viper"

func LoadConfig() {
	// viper.Debug()
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read the config file")
	}
}
