package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

var (
	cleanupOnce sync.Once
)

func createCleanupFunc(canvas *Canvas) func() {
	return func() {
		fmt.Println("\nCTRL+C signal received, cleaning up...")
		canvas.Close()
	}
}

func main() {
	canvas := getCanvasInstance()

	// Prepare for cleanup
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for {
			select {
			case <- signalChannel:
				cleanupOnce.Do(createCleanupFunc(canvas))
				os.Exit(0)
			}
		}
	}()

	font, err := getFontFromJson("./fonts/default.json")
	if err != nil {
		fmt.Println("Error reading font: ", err)
	}

	textsprite := font.ConvertTextToSpritesheet("Hello, world!")

	fmt.Println("\n\n", textsprite)

	spritesheet, err := getSpritesheetFromJson("./sprites/kirbyWalking.json")
	if err != nil {
		fmt.Println(err)
	}
	drawOptions := DrawOptions{Reverse: true}
	spritesheet.Draw(drawOptions)
}
