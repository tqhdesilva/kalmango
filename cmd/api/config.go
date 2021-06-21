package main

import (
	"github.com/spf13/viper"
)

// API_SSL, API_ASSETS, API_TIME_DELTA
type Config struct {
	SSL       bool    `json:"ssl"`
	Assets    string  `json:"assets"`
	TimeDelta float64 `json:"time_delta"`
	Port      int     `json:"port"`
}

func LoadConfig() Config {
	viper.SetEnvPrefix("API")
	viper.SetDefault("SSL", false)
	viper.SetDefault("ASSETS", "./assets")
	viper.SetDefault("TIME_DELTA", 0.1)
	viper.SetDefault("PORT", 8080)
	viper.AutomaticEnv()
	ssl := viper.GetBool("SSL")
	assets := viper.GetString("ASSETS")
	timeDelta := viper.GetFloat64("TIME_DELTA")
	port := viper.GetInt("PORT")
	return Config{ssl, assets, timeDelta, port}
}
