package main

/*
#cgo LDFLAGS: -L. -lrgbmatrix -lstdc++ -lm
#include "../include/led-matrix-c.h"
#include <stdlib.h>
*/
import "C"
import (
	"sync"
	"unsafe"
)

var canvasLock = &sync.Mutex{}

type Canvas struct {
	matrix  *C.struct_RGBLedMatrix
	canvas  *C.struct_LedCanvas
	options C.struct_RGBLedMatrixOptions
	config  *Config
}

var canvasInstance *Canvas

func getCanvasInstance() *Canvas {
	if canvasInstance == nil {
		canvasLock.Lock()
		defer canvasLock.Unlock()

		if canvasInstance == nil {
			canvasInstance = &Canvas{}
			canvasInstance.init()
		}
	}

	return canvasInstance
}

func (c *Canvas) SetDefaultOptions() {
	c.options.hardware_mapping = C.CString("regular")
	c.options.cols = C.int(c.config.Canvas.ScreenWidth)
	c.options.rows = C.int(c.config.Canvas.ScreenHeight)
	c.options.brightness = C.int(c.config.Canvas.Brightness)
	c.options.disable_hardware_pulsing = true

	// For whatever reason the lib swaps G and B; this corrects that
	c.options.led_rgb_sequence = C.CString("RBG")
}

func (c *Canvas) init() {
	c.config = getConfig()

	c.SetDefaultOptions()

	c.matrix = C.led_matrix_create_from_options(&c.options, nil, nil)
	c.canvas = C.led_matrix_get_canvas(c.matrix)
}

func (c *Canvas) SetPixel(_x int, _y int, _r int, _g int, _b int) {
	x := C.int(_x)
	y := C.int(_y)

	r := C.uchar(_r)
	g := C.uchar(_g)
	b := C.uchar(_b)

	C.led_canvas_set_pixel(c.canvas, x, y, r, g, b)
}

func (c *Canvas) Clear() {
	C.led_canvas_clear(c.canvas)
}

func (c *Canvas) Close() {
	C.free(unsafe.Pointer(c.options.hardware_mapping))
	C.led_matrix_delete(c.matrix)
}

func (c *Canvas) DrawScreen(pixeldata [][]int, colors []int, offsetX int, offsetY int) {
	for y := 0; y < len(pixeldata); y++ {
		for x := 0; x < len(pixeldata[y]); x++ {
			colorIndex := pixeldata[y][x]
			if colorIndex == -1 {
				continue
			}

			color := colors[colorIndex]
			if color == 0 {
				continue
			}

			r, g, b := convertColorToRGB(color)

			c.SetPixel(x+offsetX, y+offsetY, r, g, b)
		}
	}
}
