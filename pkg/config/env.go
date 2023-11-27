package config

import (
	"errors"
	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(err.Error())
		}
	}

}
