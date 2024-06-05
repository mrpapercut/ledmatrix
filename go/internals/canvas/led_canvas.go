package canvas

/*
#cgo LDFLAGS: -L../../../lib -lrgbmatrix -lstdc++ -lm
#include "../../../include/led-matrix-c.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/mrpapercut/ledmatrix/internals/config"
	"github.com/mrpapercut/ledmatrix/internals/utils"
)

type LEDCanvas struct {
	matrix  *C.struct_RGBLedMatrix
	canvas  *C.struct_LedCanvas
	options C.struct_RGBLedMatrixOptions
	config  *config.Config
}

func (c *LEDCanvas) SetDefaultOptions() {
	c.options.hardware_mapping = C.CString("regular")
	c.options.cols = C.int(c.config.Canvas.ScreenWidth)
	c.options.rows = C.int(c.config.Canvas.ScreenHeight)
	c.options.brightness = C.int(c.config.Canvas.Brightness)
	c.options.disable_hardware_pulsing = true

	// For whatever reason the lib swaps G and B; this corrects that
	c.options.led_rgb_sequence = C.CString("RBG")
}

func (c *LEDCanvas) init() {
	c.config = config.GetConfig()

	c.SetDefaultOptions()

	c.matrix = C.led_matrix_create_from_options(&c.options, nil, nil)
	c.canvas = C.led_matrix_get_canvas(c.matrix)
}

func (c *LEDCanvas) setPixel(_x int, _y int, _r int, _g int, _b int) error {
	x := C.int(_x)
	y := C.int(_y)

	r := C.uchar(_r)
	g := C.uchar(_g)
	b := C.uchar(_b)

	C.led_canvas_set_pixel(c.canvas, x, y, r, g, b)

	return nil
}

func (c *LEDCanvas) Clear() error {
	C.led_canvas_clear(c.canvas)

	return nil
}

func (c *LEDCanvas) Close() error {
	C.free(unsafe.Pointer(c.options.hardware_mapping))
	C.led_matrix_delete(c.matrix)

	return nil
}

func (c *LEDCanvas) DrawScreen(pixeldata [][]int, colors []int, offsetX int, offsetY int) error {
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

			r, g, b := utils.ConvertColorToRGB(color)

			c.setPixel(x+offsetX, y+offsetY, r, g, b)
		}
	}

	return nil
}
