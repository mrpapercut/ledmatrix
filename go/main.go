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
		for range signalChannel {
			scheduler.Stop()
			os.Exit(0)
		}
	}()

	scheduler.Start()

	select {}
}
