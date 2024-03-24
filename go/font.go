package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ConvertOptions struct {
	CharacterSpacing int
}

type Font struct {
	Name       string           `json:"name"`
	Characters map[string][]int `json:"characters"`
}

func getFontFromJson(filename string) (*Font, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
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

func getDefaultFont() *Font {
	font, _ := getFontFromJson("./fonts/default.json")

	return font
}

func getSMWFont() *Font {
	font, _ := getFontFromJson("./fonts/smw.json")

	return font
}

func getSegmentedDisplayFont() *Font {
	font, _ := getFontFromJson("./fonts/segmented-display.json")

	return font
}

func getMinimalNumbersFont() *Font {
	font, _ := getFontFromJson("./fonts/minimal-numbers.json")

	return font
}

func (f *Font) ConvertTextToSpritesheet(text string, convertOptions ConvertOptions) *Spritesheet {
	config := getConfig()

	spritesheet := &Spritesheet{
		Width:     0,
		Height:    0,
		NumSheets: 1,
		FPS:       10,
		Animation: []int{0},
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
			binaryString := strconv.FormatInt(int64(characters[i][j]), 2)
			reversedString := reverseBinaryString(binaryString)
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
	var singleSheetPixelData PixelData

	if convertOptions.CharacterSpacing == 0 {
		convertOptions.CharacterSpacing = 2
	}
	maxWidth := 0
	maxHeight := 0

	for _, sheet := range sheetsPixelData {
		width, height := getSheetWidthHeight(sheet)
		if width > maxWidth {
			maxWidth = width
		}
		if height > maxHeight {
			maxHeight = height
		}
	}

	singleSheet := make([][]int, 0)
	maxSheetWidth := 0

	for y := 0; y < maxHeight; y++ {
		combinedRow := make([]int, 0)

		for sheet := 0; sheet < len(sheetsPixelData); sheet++ {
			width, _ := getSheetWidthHeight(sheetsPixelData[sheet])

			if y < len(sheetsPixelData[sheet]) {
				combinedRow = append(combinedRow, sheetsPixelData[sheet][y]...)

				padding := width - len(sheetsPixelData[sheet][y])
				if padding > 0 {
					combinedRow = append(combinedRow, make([]int, padding)...)
				}
			} else {
				combinedRow = append(combinedRow, make([]int, width)...)
			}

			if sheet < len(sheetsPixelData)-1 {
				combinedRow = append(combinedRow, make([]int, convertOptions.CharacterSpacing)...)
			}
		}

		if len(combinedRow) > maxSheetWidth {
			maxSheetWidth = len(combinedRow)
		}

		singleSheet = append(singleSheet, combinedRow)
	}

	singleSheetPixelData = append(singleSheetPixelData, singleSheet)
	spritesheet.PixelData = singleSheetPixelData

	spritesheet.Width = maxSheetWidth
	spritesheet.Height = maxHeight

	return spritesheet
}
