package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Font struct {
	name string	`json:"name"`
	characters map[string][]int `json:"characters"`
}

func getFontFromJson(filename string) (*Font, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	var font Font

	err = json.Unmarshal(jsonData, &font)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return nil, err
	}

	return &font, nil
}

func (f *Font) ConvertTextToSpritesheet(text string) *Spritesheet {
	spritesheet := &Spritesheet{
		Width: 0,
		Height: 0,
	}

	return spritesheet
}
