package emu

// colors taken from https://github.com/mitxela/swotGB/blob/master/gbjs.htm
var ColorPallet = []Pixel{
	/* White */ {0xE0, 0xF8, 0xD0},
	/* Light */ {0x88, 0xC0, 0x70},
	/* Dark  */ {0x34, 0x68, 0x56},
	/* Black */ {0x08, 0x18, 0x20},
}

type Pixel struct {
	R byte
	G byte
	B byte
}

func getColor(palette, hb, lb, bit uint8) Pixel {
	i := getColorIdx(hb, lb, bit)
	return ColorPallet[(palette>>(i*2))&0b11]
}

func getColorIdx(hb, lb, bit uint8) int {
	i := 0
	if hb&(1<<bit) != 0 {
		i += 2
	}

	if lb&(1<<bit) != 0 {
		i++
	}

	return i
}
