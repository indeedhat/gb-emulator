package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// colors taken from https://github.com/mitxela/swotGB/blob/master/gbjs.htm
var ColorPallet = []Pixel{
	/* White */ Pixel{0xE0, 0xF8, 0xD0},
	/* Light */ Pixel{0x88, 0xC0, 0x70},
	/* Dark  */ Pixel{0x34, 0x68, 0x56},
	/* Black */ Pixel{0x08, 0x18, 0x20},
}

type Pixel struct {
	R byte
	G byte
	B byte
}

type LcdRenderer struct {
	displayTileData   bool
	buffer            []byte
	blankScreenBuffer []byte

	screenX int
	screenY int

	ctx *Context
}

func NewLcdRenderer(ctx *Context) *LcdRenderer {
	l := &LcdRenderer{
		ctx:             ctx,
		displayTileData: true,
	}
	l.init()

	return l
}

func (l *LcdRenderer) init() {
	l.screenX = 160
	l.screenY = 144

	if l.displayTileData {
		l.screenX = 386
	}

	l.blankScreenBuffer = make([]byte, l.screenX*l.screenY*4)
	l.buffer = make([]byte, l.screenX*l.screenY*4)
}

// Draw implements ebiten.Game.
func (l *LcdRenderer) Draw(screen *ebiten.Image) {
	copy(l.buffer, l.blankScreenBuffer)

	l.drawGame(screen)
	l.drawTileData(screen)

	screen.WritePixels(l.buffer)
}

func (l *LcdRenderer) drawGame(_ *ebiten.Image) {
	if !l.displayTileData {
		l.buffer = l.ctx.ppu.videoBuffer
		return
	}

	lineOffset := 386 * 4

	for i := 0; i < 144; i++ {
		copy(l.buffer[i*lineOffset:], l.ctx.ppu.videoBuffer[i*640:(i+1)*640])
	}
}

func (l *LcdRenderer) drawTileData(_ *ebiten.Image) {
	if !l.displayTileData {
		return
	}

	tileMap1 := l.ctx.ppu.vram.data

	for tileId := 0; tileId < 384; tileId++ {
		xOffset := 160 + (tileId%24)*9
		yOffset := (tileId / 24) * 9

		for y := 0; y < 8; y++ {
			hb := tileMap1[tileId*16+(y*2)]
			lb := tileMap1[tileId*16+(y*2)+1]

			for x := 7; x >= 0; x-- {
				i := (xOffset+x)*4 + (yOffset+y)*386*4
				color := getColor(hb, lb, uint8(x))

				l.buffer[i] = color.R
				l.buffer[i+1] = color.G
				l.buffer[i+2] = color.B
				l.buffer[i+3] = 0xFF
			}
		}
	}
}

func getColor(hb, lb, bit uint8) Pixel {
	i := 0
	if hb&(1<<bit) != 0 {
		i += 2
	}

	if lb&(1<<bit) != 0 {
		i++
	}

	return ColorPallet[i]
}

// Layout implements ebiten.Game.
func (l *LcdRenderer) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return l.screenX, l.screenY
}

// Update implements ebiten.Game.
func (l *LcdRenderer) Update() error {
	// TODO: implement io
	return nil
}

var _ ebiten.Game = (*LcdRenderer)(nil)
