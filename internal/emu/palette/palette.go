package palette

import . "github.com/indeedhat/gb-emulator/internal/emu/types"

// colors taken from https://github.com/mitxela/swotGB/blob/master/gbjs.htm
var ColorPallet = []Pixel{
	/* White */ {R: 0xE0, G: 0xF8, B: 0xD0},
	/* Light */ {R: 0x88, G: 0xC0, B: 0x70},
	/* Dark  */ {R: 0x34, G: 0x68, B: 0x56},
	/* Black */ {R: 0x08, G: 0x18, B: 0x20},
}

func GetColor(palette, hb, lb, bit uint8) Pixel {
	i := GetColorIdx(hb, lb, bit)
	return ColorPallet[(palette>>(i*2))&0b11]
}

func GetColorIdx(hb, lb, bit uint8) int {
	i := 0
	if hb&(1<<bit) != 0 {
		i += 2
	}

	if lb&(1<<bit) != 0 {
		i++
	}

	return i
}
