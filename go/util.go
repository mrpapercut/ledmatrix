package main

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
