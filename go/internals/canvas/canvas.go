package canvas

import (
	"sync"
)

type CanvasInterface interface {
	Clear() error
	Close() error
	DrawScreen(pixeldata [][]int, colors []int, offsetX int, offsetY int) error
	init()
}

var canvasLock = &sync.Mutex{}
var canvasInstance CanvasInterface

func GetCanvasInstance(canvasInterface ...CanvasInterface) CanvasInterface {
	if canvasInstance == nil {
		canvasLock.Lock()
		defer canvasLock.Unlock()

		if canvasInstance == nil {
			if len(canvasInterface) > 0 {
				canvasInstance = canvasInterface[0]
			} else {
				canvasInstance = &LEDCanvas{}
			}

			canvasInstance.init()
		}
	}

	return canvasInstance
}
