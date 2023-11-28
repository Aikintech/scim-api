package config

import (
	"errors"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.SetConfigName(".env")
	viper.AddConfigPath("/app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(err.Error())
		}
	}

}
