package utils

func GetSheetWidthHeight(sheet [][]int) (int, int) {
	maxSheetWidth := 0
	maxSheetHeight := len(sheet)

	for i := 0; i < len(sheet); i++ {
		if len(sheet[i]) > maxSheetWidth {
			maxSheetWidth = len(sheet[i])
		}
	}

	return maxSheetWidth, maxSheetHeight
}
