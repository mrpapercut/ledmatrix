package main

import (
	"os"
	"os/signal"
)

func main() {
	canvas := getCanvasInstance()
	config := getConfig()

	scheduler := getSchedulerInstance(canvas, config)
	defer scheduler.Stop()

	// Prepare for cleanup
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for {
			select {
			case <- signalChannel:
				scheduler.Stop()
				os.Exit(0)
			}
		}
	}()

	scheduler.Start()

	select {}
}

/*
func runJob() {
	fmt.Println("Job executed at", time.Now())

	rand.Seed(time.Now().UnixNano())

	callableFunctions := []func(){
		DrawKirbyAnimation,
		DrawClock,
	}

	randomIndex := rand.Intn(len(callableFunctions))
	callableFunctions[randomIndex]()
}

func DrawText(textMessage string) {
	font := getSMWFont()
	convertOptions := ConvertOptions{
		CharacterSpacing: 2,
	}
	textsprite := font.ConvertTextToSpritesheet(textMessage, convertOptions)

	textDrawOptions := DrawOptions{
		ScrollSpeed: 3,
		SpriteType: TextSprite,
	}
	textsprite.Draw(textDrawOptions)
}

func DrawLogo() {
	spritesheet, _ := getSpritesheetFromJson("./sprites/youtubeLogo.json")

	drawOptions := DrawOptions{
		SpriteType: StaticSprite,
		Loop: true,
		Duration: 10,
	}

	spritesheet.Draw(drawOptions)
}
*/
