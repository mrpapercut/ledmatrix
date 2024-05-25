package clock

import (
	"fmt"
	"time"

	"github.com/mrpapercut/ledmatrix/internals/font"
	"github.com/mrpapercut/ledmatrix/internals/types"
)

type Clock struct {
}

func (c *Clock) formatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func (c *Clock) DrawDigitalClock() {
	currentTime := time.Now()

	startTime := c.formatTime(currentTime)

	f := font.GetFontByName("minimal-numbers")
	convertOptions := types.ConvertOptions{
		CharacterSpacing: 1,
	}
	timeSprite := f.ConvertTextToSpritesheet(startTime, convertOptions)
	timeSprite.Animation = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	timeSprite.FPS = 1

	for i := 0; i < 10; i++ {
		futureTime := currentTime.Add(time.Duration(i+1) * time.Second)

		sprite := f.ConvertTextToSpritesheet(c.formatTime(futureTime), convertOptions)

		timeSprite.PixelData = append(timeSprite.PixelData, sprite.PixelData[0])
	}

	drawOptions := types.DrawOptions{
		SpriteType: types.StaticSprite,
	}
	timeSprite.Draw(drawOptions)
}

func (c *Clock) DrawColoredClock() {

}
