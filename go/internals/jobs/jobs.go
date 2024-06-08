package jobs

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/mrpapercut/ledmatrix/internals/clock"
	"github.com/mrpapercut/ledmatrix/internals/font"
	"github.com/mrpapercut/ledmatrix/internals/spritesheet"
	"github.com/mrpapercut/ledmatrix/internals/sqlite"
	"github.com/mrpapercut/ledmatrix/internals/types"
)

type Jobs struct{}

func (j *Jobs) DrawClock() {
	clock := clock.Clock{}
	clock.DrawDigitalClock()
}

func (j *Jobs) DrawKirbyDance() {
	j.DrawAnimationByFilename("kirbyDance01")
}

func (j *Jobs) DrawKirbyAnimation() {
	spritesheets := []string{
		"./sprites/kirbyWalking.json",
		"./sprites/kirbyRunning.json",
		"./sprites/kirbyTumbling.json",
	}
	randomIndex := rand.Intn(len(spritesheets))
	jsonFile := spritesheets[randomIndex]

	sheet, _ := spritesheet.GetSpritesheetFromJson(jsonFile)

	animationDrawOptions := types.DrawOptions{
		ScrollSpeed: 3,
		SpriteType:  types.AnimationSprite,
		Reverse:     rand.Intn(2) == 0,
	}

	sheet.Draw(animationDrawOptions)
}

func (j *Jobs) DrawAnimationByFilename(filename string) {
	sheet, err := spritesheet.GetSpritesheetFromJson("./sprites/" + filename + ".json")
	if err != nil {
		log.Fatal("Could not load spritesheet:", err)
	}

	animationDrawOptions := types.DrawOptions{
		SpriteType: types.StaticSprite,
	}

	sheet.Draw(animationDrawOptions)
}

func (j *Jobs) DrawLogo(logo string) {
	sheet, err := spritesheet.GetSpritesheetFromJson(fmt.Sprintf("./sprites/%sLogo.json", logo))
	if err != nil {
		log.Fatalf("Error creating spritesheet: %v", err)
	}

	drawOptions := types.DrawOptions{
		SpriteType: types.StaticSprite,
		Loop:       true,
	}

	sheet.Draw(drawOptions)
}

// Screen in case there's nothing more important to show
func (j *Jobs) DrawIdleScreen() {
	callableFunctions := []func(){
		j.DrawKirbyAnimation,
		j.DrawKirbyAnimation,
		j.DrawKirbyDance,
		j.DrawClock,
	}

	randomIndex := rand.Intn(len(callableFunctions))
	callableFunctions[randomIndex]()
}

func (j *Jobs) GetHighPriorityMessage() (bool, types.FeedMessage) {
	sql := sqlite.GetSQLiteInstance()

	priority := 1
	message, err := sql.GetLatestFeedMessage(priority)
	if err != nil {
		// log.Println("Error retrieving latest feed message:", err)
		return false, types.FeedMessage{}
	}

	return true, message
}

func (j *Jobs) GetLowPriorityMessage() (bool, types.FeedMessage) {
	sql := sqlite.GetSQLiteInstance()

	priority := 2
	message, err := sql.GetLatestFeedMessage(priority)
	if err != nil {
		log.Println("Error retrieving latest feed message:", err)
		return false, types.FeedMessage{}
	}

	return true, message
}

func (j *Jobs) DrawFeedMessage(message types.FeedMessage, logo ...*spritesheet.Spritesheet) {
	f := font.GetFontByName("default")
	convertOptions := types.ConvertOptions{}
	messageSprite := f.ConvertTextToSpritesheet(message.Message, convertOptions)

	if len(logo) > 0 {
		messageSprite = f.PrependLogoToTextSpritesheet(logo[0], messageSprite)
	}

	drawOptions := types.DrawOptions{
		SpriteType:  types.TextSprite,
		Loop:        false,
		Scroll:      true,
		ScrollSpeed: 3,
		Direction:   types.Left,
	}

	messageSprite.Draw(drawOptions)
}
