package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	SSL       bool    `json:"ssl"`
	Assets    string  `json:"assets"`
	TimeDelta float64 `json:"time_delta"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	json.Unmarshal(data, config)
	return config, nil
}

func DefaultConfig() *Config {
	return &Config{
		SSL:       false,
		Assets:    "./assets",
		TimeDelta: 0.1,
	}
}
