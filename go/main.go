package main

/*
#cgo LDFLAGS: -L. -lrgbmatrix -lstdc++ -lm
#include "../include/led-matrix-c.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func main() {
	var options C.struct_RGBLedMatrixOptions
	options.hardware_mapping = C.CString("regular")
	options.cols = 64
	options.rows = 32
	options.brightness = 50
	options.disable_hardware_pulsing = true

	matrix := C.led_matrix_create_from_options(&options, nil, nil)
	canvas := C.led_matrix_get_canvas(matrix)

	C.led_canvas_set_pixel(canvas, 0, 0, 0xff, 0xff, 0xff)

	// Don't forget to free the C string
	C.free(unsafe.Pointer(options.hardware_mapping))

	C.led_matrix_delete(matrix)
}
