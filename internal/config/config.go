package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    EVMR_PORT   string `mapstructure:"EVMR_PORT"`
    EVMR_SERVER string `mapstructure:"EVMR_SERVER"`
    EVMR_AUTH	string `mapstructure:"EVMR_AUTH"`
}


func LoadConfig() (Config, error) {
    config := Config{EVMR_PORT: "3000", EVMR_SERVER: "http://localhost", EVMR_AUTH: ""}
	
	viper.SetConfigFile(".env")

    err := viper.ReadInConfig()
    if err != nil {
		return config, err
	}

	config.EVMR_AUTH = viper.GetString("EVMR_AUTH")

    return config, err
}