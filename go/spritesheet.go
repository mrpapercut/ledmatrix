package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type PixelData [][][]int
type FontData [][]int

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type SpriteType int

const (
	AnimationSprite SpriteType = iota
	TextSprite
	StaticSprite
)

type DrawOptions struct {
	SpriteType  SpriteType
	Reverse     bool
	Loop        bool
	Scroll      bool
	ScrollSpeed int
	Direction   Direction
	Duration    int
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
	switch drawOptions.SpriteType {
	case AnimationSprite:
		s.drawAnimation(drawOptions)
	case TextSprite:
		s.drawText(drawOptions)
	case StaticSprite:
		s.drawStaticImage(drawOptions)
	}

	canvas := getCanvasInstance()
	canvas.Clear()
}

func (s *Spritesheet) drawAnimation(drawOptions DrawOptions) {
	config := getConfig()
	canvas := getCanvasInstance()

	frameDuration := time.Second / time.Duration(s.FPS)

	animationFrames := s.Animation
	animationIndex := 0

	maxSpriteWidth := s.Width
	maxSpriteHeight := s.Height

	colors := s.Colors

	offsetX := 0 - maxSpriteWidth
	offsetY := (config.Canvas.ScreenHeight - maxSpriteHeight) / 2

	if drawOptions.Reverse {
		offsetX = config.Canvas.ScreenWidth

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
			offsetX = offsetX - drawOptions.ScrollSpeed
			if offsetX < (0 - maxSpriteWidth) {
				if drawOptions.Loop {
					offsetX = config.Canvas.ScreenWidth + maxSpriteWidth
				} else {
					return
				}
			}
		} else {
			offsetX = offsetX + drawOptions.ScrollSpeed

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

func (s *Spritesheet) drawText(drawOptions DrawOptions) {
	config := getConfig()
	canvas := getCanvasInstance()

	frameDuration := time.Second / time.Duration(s.FPS)

	animationFrames := s.Animation
	animationIndex := 0

	maxSpriteWidth := s.Width
	maxSpriteHeight := s.Height

	colors := s.Colors

	offsetX := 0 - maxSpriteWidth
	offsetY := (config.Canvas.ScreenHeight - maxSpriteHeight) / 2

	if drawOptions.Direction == Left {
		offsetX = config.Canvas.ScreenWidth
	}

	if drawOptions.Reverse {
		offsetX = config.Canvas.ScreenWidth

		for i, sheet := range s.PixelData {
			s.PixelData[i] = s.reverseSheet(sheet)
		}
	}

	for {
		canvas.Clear()

		currentSprite := s.PixelData[animationFrames[animationIndex]]
		animationIndex = (animationIndex + 1) % len(s.Animation)

		canvas.DrawScreen(currentSprite, colors, offsetX, offsetY)

		if drawOptions.Reverse || drawOptions.Direction == Left {
			offsetX = offsetX - drawOptions.ScrollSpeed
			if offsetX < (0 - maxSpriteWidth) {
				if drawOptions.Loop {
					offsetX = config.Canvas.ScreenWidth + maxSpriteWidth
				} else {
					return
				}
			}
		} else {
			offsetX = offsetX + drawOptions.ScrollSpeed

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

func (s *Spritesheet) drawStaticImage(drawOptions DrawOptions) {
	config := getConfig()
	canvas := getCanvasInstance()

	if s.FPS < 1 {
		s.FPS = 1
	}

	if drawOptions.Duration == 0 {
		drawOptions.Duration = 1
	}

	frameDuration := time.Second / time.Duration(s.FPS)

	animationFrames := s.Animation
	animationIndex := 0

	maxSpriteWidth := s.Width
	maxSpriteHeight := s.Height

	colors := s.Colors

	offsetX := (config.Canvas.ScreenWidth - maxSpriteWidth) / 2
	offsetY := (config.Canvas.ScreenHeight - maxSpriteHeight) / 2

	timesRun := 0

	for {
		canvas.Clear()

		currentSprite := s.PixelData[animationFrames[animationIndex]]

		if animationIndex+1 >= len(s.Animation) && !drawOptions.Loop {
			return
		} else {
			animationIndex = (animationIndex + 1) % len(s.Animation)
		}

		canvas.DrawScreen(currentSprite, colors, offsetX, offsetY)

		time.Sleep(frameDuration)

		if drawOptions.Duration > 0 {
			timesRun += 1
			if timesRun >= drawOptions.Duration {
				return
			}
		}
	}
}
