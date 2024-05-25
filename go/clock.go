package main

import (
	"fmt"
	"time"
)

type Clock struct {
}

func (c *Clock) formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func (c *Clock) DrawDigitalClock() {
	currentTime := time.Now()

	startTime := c.formatTime(currentTime)

	font := getFontByName("minimal-numbers")
	convertOptions := ConvertOptions{
		CharacterSpacing: 1,
	}
	timeSprite := font.ConvertTextToSpritesheet(startTime, convertOptions)
	timeSprite.Animation = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	timeSprite.FPS = 1

	for i := 0; i < 10; i++ {
		futureTime := currentTime.Add(time.Duration(i+1) * time.Second)

		sprite := font.ConvertTextToSpritesheet(c.formatTime(futureTime), convertOptions)

		timeSprite.PixelData = append(timeSprite.PixelData, sprite.PixelData[0])
	}

	drawOptions := DrawOptions{
		SpriteType: StaticSprite,
	}
	timeSprite.Draw(drawOptions)
}

func (c *Clock) DrawColoredClock() {

}
