package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	SSL       bool
	Assets    string
	TimeDelta float64
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config *Config
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
