package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	return viper.ReadInConfig()
}
