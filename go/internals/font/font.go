package font

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/spritesheet"
	"github.com/mrpapercut/ledmatrix/internals/types"
	"github.com/mrpapercut/ledmatrix/internals/utils"
)

type Font struct {
	Name       string           `json:"name"`
	Characters map[string][]int `json:"characters"`
}

func getFontFromJson(filename string) (*Font, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening config file:", err)
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading config file:", err)
		return nil, err
	}

	var font Font

	err = json.Unmarshal(jsonData, &font)
	if err != nil {
		log.Println("Error parsing config file:", err)
		return nil, err
	}

	return &font, nil
}

func GetDefaultFont() *Font {
	font, _ := getFontFromJson("./fonts/default.json")

	return font
}

func GetSMWFont() *Font {
	font, _ := getFontFromJson("./fonts/smw.json")

	return font
}

func GetSegmentedDisplayFont() *Font {
	font, _ := getFontFromJson("./fonts/segmented-display.json")

	return font
}

func GetMinimalNumbersFont() *Font {
	font, _ := getFontFromJson("./fonts/minimal-numbers.json")

	return font
}

func GetFontByName(name string) *Font {
	fontfns := map[string]func() *Font{
		"default":         GetDefaultFont,
		"smw":             GetSMWFont,
		"segmented":       GetSegmentedDisplayFont,
		"minimal-numbers": GetMinimalNumbersFont,
	}

	if fontfn, ok := fontfns[name]; ok {
		return fontfn()
	} else {
		log.Printf("Font not found: %s", name)
		return GetDefaultFont()
	}
}

func (f *Font) ConvertTextToSpritesheet(text string, convertOptions types.ConvertOptions) *spritesheet.Spritesheet {
	// Convert text to font-characters
	characters := make([][]int, 0, len(text))
	for _, char := range text {
		fontChar := f.Characters[string(char)]

		if fontChar != nil {
			characters = append(characters, fontChar)
		}
	}

	// Convert font-characters to sheets
	sheetsPixelData := types.PixelData{}
	for i := 0; i < len(characters); i++ {
		characterPixels := make([][]int, 0)

		for j := 0; j < len(characters[i]); j++ {
			binaryString := strconv.FormatInt(int64(characters[i][j]), 2)
			reversedString := utils.ReverseBinaryString(binaryString)
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

	return f.convertCharactersToSingleSpritesheets(convertOptions, sheetsPixelData)
}

func (f *Font) convertCharactersToSingleSpritesheets(convertOptions types.ConvertOptions, sheetsPixelData types.PixelData) *spritesheet.Spritesheet {
	config := config.GetConfig()

	sheet := &spritesheet.Spritesheet{
		Width:     0,
		Height:    0,
		FPS:       10,
		Animation: []int{0},
		Colors:    []int{0, config.Canvas.TextColor},
	}

	if convertOptions.CharacterSpacing == 0 {
		convertOptions.CharacterSpacing = 2
	}

	// Get max width & height of characters for alignment
	maxWidth := 0
	maxHeight := 0
	for _, sheet := range sheetsPixelData {
		width, height := utils.GetSheetWidthHeight(sheet)
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
			width, _ := utils.GetSheetWidthHeight(sheetsPixelData[sheet])

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

	sheet.PixelData = types.PixelData{singleSheet}

	sheet.Width = maxSheetWidth
	sheet.Height = maxHeight

	return sheet
}

func (f *Font) PrependLogoToTextSpritesheet(logo, textSheet *spritesheet.Spritesheet) *spritesheet.Spritesheet {
	config := config.GetConfig()

	sheet := &spritesheet.Spritesheet{
		Width:     0,
		Height:    0,
		FPS:       10,
		Animation: []int{0},
		Colors:    []int{0, config.Canvas.TextColor},
	}

	// Add logo colors
	sheet.Colors = append(sheet.Colors, logo.Colors...)
	// Update color indexes
	for i := 0; i < len(logo.PixelData[0]); i++ {
		for j := 0; j < len(logo.PixelData[0][i]); j++ {
			if logo.PixelData[0][i][j] == -1 {
				logo.PixelData[0][i][j] = 1
			}
			logo.PixelData[0][i][j] += 2
		}
	}

	// Create new sheet for pixeldata
	singleSheet := make([][]int, 0)

	logoIsLargest := logo.Height >= textSheet.Height
	if !logoIsLargest {
		// Increase the height of the logo to match the text
		newRows := make([][]int, 0)
		newRows = append(newRows, logo.PixelData[0]...)

		halvedDifference := int(math.Floor(float64(textSheet.Height-logo.Height) / 2))

		for i := 0; i < halvedDifference; i++ {
			newRows = append(newRows, make([]int, logo.Width+2))
		}

		logo.Height = len(newRows)
		logo.PixelData = types.PixelData{newRows}
	}

	maxSheetWidth := 0
	for i := 0; i < logo.Height; i++ {
		combinedRow := make([]int, 0)
		combinedRow = append(combinedRow, logo.PixelData[0][i]...)

		if len(logo.PixelData[0][i]) < logo.Width {
			// Make sure logo draws up to max width
			combinedRow = append(combinedRow, make([]int, logo.Width-len(logo.PixelData[0][i]))...)
		}

		if i < len(textSheet.PixelData[0]) {
			combinedRow = append(combinedRow, make([]int, 3)...) // Spacing between logo and text
			combinedRow = append(combinedRow, textSheet.PixelData[0][i]...)
		}

		singleSheet = append(singleSheet, combinedRow)

		if len(combinedRow) > maxSheetWidth {
			maxSheetWidth = len(combinedRow)
		}
	}

	sheet.Width = maxSheetWidth
	sheet.Height = logo.Height

	sheet.PixelData = types.PixelData{singleSheet}

	return sheet
}
