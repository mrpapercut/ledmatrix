package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Font struct {
	Name       string         `json:"name"`
	Characters map[rune][]int `json:"characters"`
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
	config := getConfig()

	spritesheet := &Spritesheet{
		Width:     0,
		Height:    0,
		NumSheets: 1,
		FPS:       1,
		Animation: []int{1},
		Colors:    []int{0, config.Canvas.TextColor},
		PixelData: PixelData{},
	}

	// Convert text to font-characters
	characters := make([][]int, len(text))
	for _, char := range text {
		if f.Characters[char] != nil {
			characters = append(characters, f.Characters[char])
		}
	}

	fmt.Println(characters)

	// Convert font-characters to single sheet

	return spritesheet
}
