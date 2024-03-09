package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var configLock = &sync.Mutex{}

type Config struct {
	Log struct {
		Level string `json:"level"`
		LogToFile bool `json:"log_to_file"`
		Filename string `json:"filename"`
	} `json:"log"`
	Canvas struct {
		ScreenWidth int `json:"screen_width"`
		ScreenHeight int `json:"screen_height"`
		Brightness int `json:"brightness"`
		TextColor int `json:"text_color"`
	} `json:"canvas"`
	Database struct {
		Filename string `json:"filename"`
	} `json:"db"`
	Youtube struct {
		ApiKey string `json:"apikey"`
		Channels map[string]string `json:"channels"`
	} `json:"youtube"`
	Weather struct {
		ApiKey string `json:"apikey"`
		Location int `json:"location"`
	} `json:"weather"`
}

var configInstance *Config

func getConfig() *Config {
	if configInstance == nil {
		configLock.Lock()
		defer configLock.Unlock()

		if configInstance == nil {
			configInstance = &Config{}
			configInstance.init()
		}
	}

	return configInstance
}

func (c *Config) init() {
	file, err := os.Open("../config.json")
	defer file.Close()

	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	err = json.Unmarshal(jsonData, &c)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}
}