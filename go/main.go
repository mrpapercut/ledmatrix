package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/mrpapercut/ledmatrix/internals/canvas"
	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/scheduler"
)

func main() {
	// Setup logger
	logfile, err := os.OpenFile("ledmatrix.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatalf("Failed to open ledmatrix.log: %v", err)
	}
	defer logfile.Close()

	logger := slog.New(slog.NewJSONHandler(logfile, nil))

	slog.SetDefault(logger)

	config := config.GetConfig()

	// Specify canvas by initializing instance
	canvas.GetCanvasInstance()

	schedulerInstance := scheduler.GetSchedulerInstance(config)
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
