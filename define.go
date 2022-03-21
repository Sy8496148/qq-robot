package main

import (
	"encoding/json"
	"log"
	"os"
)

const (
	defaultConfigFile = "config.json"
)

type RobotConfig struct {
	AppID uint64 `json:"app_id"`
	Token string `json:"token"`
}

func loadJONConfig(cfgFile string) (*RobotConfig, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("close file error ", cfgFile, err)
		}
	}()
	var config RobotConfig
	err = json.NewDecoder(f).Decode(&config)
	return &config, err
}
