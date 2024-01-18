package util

import "image/color"

func Dim(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r - 50), uint8(g - 50), uint8(b - 50), uint8(a)}
}

type ballColor color.Color

var (
	_colorYellow ballColor = color.RGBA{250, 216, 74, 255}
	_colorBlue   ballColor = color.RGBA{0x27, 0x5A, 0xFF, 0xFF}
	_colorRed    ballColor = color.RGBA{0xa0, 0x0, 0x0, 0xFF}
	_colorPurple ballColor = color.RGBA{0x6E, 0x38, 0x80, 0xFF}
	_colorOrange ballColor = color.RGBA{242, 169, 60, 255}
	_colorGreen  ballColor = color.RGBA{0x9, 0x8B, 0x0, 0xFF}
	_colorMaroon ballColor = color.RGBA{0x7E, 0x3C, 0x0, 0xFF}
	_colorBlack  ballColor = color.RGBA{0x0, 0x0, 0x0, 0xFF}
	_colorWhite  ballColor = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	_colorGrey   ballColor = color.RGBA{0x0, 0x4E, 0x41, 0xFF}
)

var (
	_boardColor = color.RGBA{0x44, 0x96, 0x80, 0xFF}
	_edgeColor  = color.RGBA{0x20, 0x70, 0x60, 0xFF}
	_arrowColor = color.RGBA{0x0, 0x0, 0x0, 0xFF}
)
