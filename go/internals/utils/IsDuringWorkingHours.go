package utils

import (
	"time"

	"github.com/mrpapercut/ledmatrix/internals/config"
)

func DuringWorkingHours() bool {
	config := config.GetConfig()

	if config.Debug {
		return true
	}

	currentTime := time.Now()

	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.WorkingHours.Start, 0, 0, 0, currentTime.Location())
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.WorkingHours.End, 0, 0, 0, currentTime.Location())

	if config.WorkingHours.End < config.WorkingHours.Start {
		endTime = endTime.AddDate(0, 0, 1)
	}

	return currentTime.After(startTime) && currentTime.Before(endTime)
}
