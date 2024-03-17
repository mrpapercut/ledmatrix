package main

import (
	"time"
)

func convertColorToRGB(color int) (int, int, int) {
	r := (color >> 16) & 0xff
	g := (color >> 8) & 0xff
	b := color & 0xff

	return r, g, b
}

func reverseBinaryString(binaryString string) string {
	runes := []rune(binaryString)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func getSheetWidthHeight(sheet [][]int) (int, int) {
	maxSheetWidth := 0
	maxSheetHeight := len(sheet)

	for i := 0; i < len(sheet); i++ {
		if len(sheet[i]) > maxSheetWidth {
			maxSheetWidth = len(sheet[i])
		}
	}

	return maxSheetWidth, maxSheetHeight
}

func DuringWorkingHours() bool {
	config := getConfig()

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
