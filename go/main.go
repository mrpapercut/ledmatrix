package main

import (
	"os"
	"os/signal"

	"github.com/mrpapercut/ledmatrix/internals/canvas"
	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/scheduler"
)

func main() {
	canvas := canvas.GetCanvasInstance()
	config := config.GetConfig()

	schedulerInstance := scheduler.GetSchedulerInstance(canvas, config)
	defer schedulerInstance.Stop()

	// Prepare for cleanup
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for range signalChannel {
			schedulerInstance.Stop()
			os.Exit(0)
		}
	}()

	schedulerInstance.Start()

	select {}
}
