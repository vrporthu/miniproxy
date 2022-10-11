package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerUrl    string
	Replacements map[string]string
}

func loadConfig() (Config, error) {
	b, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
