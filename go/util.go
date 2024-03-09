package main

func convertColorToRGB(color int) (int, int, int) {
	r := (color >> 16) & 0xff
	g := (color >> 8) & 0xff
	b := color & 0xff

	return r, g, b
}
