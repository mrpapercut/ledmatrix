package main

import (
	"time"
	"fmt"
)

type Clock struct {

}

func (c *Clock) formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func (c *Clock) formatDate(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func (c *Clock) formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func (c *Clock) DrawDigitalClock() {
	currentTime := time.Now()

	startTime := c.formatTime(currentTime)

	font := getMinimalNumbersFont()
	convertOptions := ConvertOptions{
		CharacterSpacing: 1,
	}
	timeSprite := font.ConvertTextToSpritesheet(startTime, convertOptions)
	timeSprite.NumSheets = 10
	timeSprite.Animation = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	timeSprite.FPS = 1

	for i := 0; i < 9; i++ {
		futureTime := currentTime.Add(time.Duration(i + 1) * time.Second)

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
