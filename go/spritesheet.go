package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type PixelData [][][]int
type FontData [][]int

type DrawOptions struct {
	Reverse bool
	Loop    bool
	Scroll  bool
}

type Spritesheet struct {
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	NumSheets int       `json:"num_sheets"`
	FPS       int       `json:"fps"`
	Animation []int     `json:"animation"`
	Colors    []int     `json:"colors"`
	PixelData PixelData `json:"pixeldata"`
}

func getSpritesheetFromJson(filename string) (*Spritesheet, error) {
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

	var spritesheet Spritesheet

	err = json.Unmarshal(jsonData, &spritesheet)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return nil, err
	}

	return &spritesheet, nil
}

func (s *Spritesheet) reverseSheet(sheet [][]int) [][]int {
	for _, row := range sheet {
		length := len(row)
		for i := 0; i < length/2; i++ {
			row[i], row[length-i-1] = row[length-i-1], row[i]
		}
	}

	return sheet
}

func (s *Spritesheet) Draw(drawOptions DrawOptions) {
	config := getConfig()
	canvas := getCanvasInstance()

	fps := s.FPS
	frameDuration := time.Duration(int(1000/fps)) * time.Millisecond

	animationFrames := s.Animation
	animationIndex := 0

	maxSpriteWidth := s.Width
	maxSpriteHeight := s.Height

	colors := s.Colors

	offsetX := 0 - maxSpriteWidth
	offsetY := (config.Canvas.ScreenHeight - maxSpriteHeight) / 2

	if drawOptions.Reverse {
		offsetX = config.Canvas.ScreenWidth + maxSpriteWidth

		for i, sheet := range s.PixelData {
			s.PixelData[i] = s.reverseSheet(sheet)
		}
	}

	for {
		canvas.Clear()

		currentSprite := s.PixelData[animationFrames[animationIndex]]
		animationIndex = (animationIndex + 1) % len(s.Animation)

		canvas.DrawScreen(currentSprite, colors, offsetX, offsetY)

		if drawOptions.Reverse {
			offsetX = offsetX - 3
			if offsetX < (0 - maxSpriteWidth) {
				if drawOptions.Loop {
					offsetX = config.Canvas.ScreenWidth + maxSpriteWidth
				} else {
					return
				}
			}
		} else {
			offsetX = offsetX + 3
			if offsetX > (config.Canvas.ScreenWidth + maxSpriteWidth) {
				if drawOptions.Loop {
					offsetX = 0 - maxSpriteWidth
				} else {
					return
				}
			}
		}

		time.Sleep(frameDuration)
	}
}
