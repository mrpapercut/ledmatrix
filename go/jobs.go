package main

import (
	"fmt"
	"log"
	"math/rand"
)

type Jobs struct{}

func (j *Jobs) DrawClock() {
	clock := Clock{}
	clock.DrawDigitalClock()
}

func (j *Jobs) DrawKirbyAnimation() {
	spritesheets := []string{
		"./sprites/kirbyWalking.json",
		"./sprites/kirbyRunning.json",
		"./sprites/kirbyTumbling.json",
	}
	randomIndex := rand.Intn(len(spritesheets))
	jsonFile := spritesheets[randomIndex]

	spritesheet, _ := getSpritesheetFromJson(jsonFile)

	animationDrawOptions := DrawOptions{
		ScrollSpeed: 3,
		SpriteType:  AnimationSprite,
		Reverse:     rand.Intn(2) == 0,
	}

	spritesheet.Draw(animationDrawOptions)
}

func (j *Jobs) DrawLogo(logo string) {
	spritesheet, err := getSpritesheetFromJson(fmt.Sprintf("./sprites/%sLogo.json", logo))
	if err != nil {
		log.Fatalf("Error creating spritesheet: %v", err)
	}

	drawOptions := DrawOptions{
		SpriteType: StaticSprite,
		Duration:   5,
		Loop:       true,
	}

	spritesheet.Draw(drawOptions)
}

// Screen in case there's nothing more important to show
func (j *Jobs) DrawIdleScreen() {
	callableFunctions := []func(){
		j.DrawKirbyAnimation,
		j.DrawKirbyAnimation,
		j.DrawKirbyAnimation,
		j.DrawClock,
	}

	randomIndex := rand.Intn(len(callableFunctions))
	callableFunctions[randomIndex]()
}

func (j *Jobs) GetHighPriorityMessage() (bool, FeedMessage) {
	sqlite := getSQLiteInstance()

	priority := 1
	message, err := sqlite.GetLatestFeedMessage(priority)
	if err != nil {
		log.Println("Error retrieving latest feed message:", err)
		return false, FeedMessage{}
	}

	return true, message
}

func (j *Jobs) GetLowPriorityMessage() (bool, FeedMessage) {
	sqlite := getSQLiteInstance()

	priority := 2
	message, err := sqlite.GetLatestFeedMessage(priority)
	if err != nil {
		log.Println("Error retrieving latest feed message:", err)
		return false, FeedMessage{}
	}

	return true, message
}

func (j *Jobs) DrawFeedMessage(message FeedMessage) {
	font := getDefaultFont()
	convertOptions := ConvertOptions{}
	messageSprite := font.ConvertTextToSpritesheet(message.Message, convertOptions)

	drawOptions := DrawOptions{
		SpriteType:  TextSprite,
		Loop:        false,
		Scroll:      true,
		ScrollSpeed: 3,
		Direction:   Left,
	}
	messageSprite.Draw(drawOptions)
}
