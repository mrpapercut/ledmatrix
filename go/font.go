package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Font struct {
	Name       string           `json:"name"`
	Characters map[string][]int `json:"characters"`
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
	}

	// Convert text to font-characters
	characters := make([][]int, 0, len(text))
	for _, char := range text {
		fontChar := f.Characters[string(char)]

		if fontChar != nil {
			characters = append(characters, fontChar)
		}
	}

	// Convert font-characters to sheets
	sheetsPixelData := PixelData{}
	for i := 0; i < len(characters); i++ {
		characterPixels := make([][]int, 0)

		for j := 0; j < len(characters[i]); j++ {
			// Convert to binary string
			binaryString := strconv.FormatInt(int64(characters[i][j]), 2)

			// Reverse the binary string
			reversedString := reverseBinaryString(binaryString)

			// Remove trailing zeros
			reversedString = strings.TrimRight(reversedString, "0")

			boolList := make([]int, 0)
			for _, c := range reversedString {
				if c == '1' {
					boolList = append(boolList, 1)
				} else {
					boolList = append(boolList, 0)
				}
			}

			characterPixels = append(characterPixels, boolList)
		}

		sheetsPixelData = append(sheetsPixelData, characterPixels)
	}

	// Convert character-sheets to single sheet
	singleSheetPixelData := PixelData{}
	offsetX := 0
	// y := 0
	for y := 0; y < len(sheetsPixelData); y++ {
		if singleSheetPixelData[y] == nil {
			singleSheetPixelData[y] = make([][]int, 0)
		}

		for x := 0; x < len(sheetsPixelData[y]); x++ {
			singleSheetPixelData[y][x + offsetX] = sheetsPixelData[y][x]
		}

		width, _ := getSheetWidthHeight(sheetsPixelData[y])

		offsetX = offsetX + width + 2
	}

	spritesheet.PixelData = singleSheetPixelData

	return spritesheet
}
