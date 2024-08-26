package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
		displayTileData: false,
	}
	l.init()

	return l
}

func (l *LcdRenderer) init() {
	l.screenX = PpuXRes
	l.screenY = PpuYRes

	if l.displayTileData {
		l.screenX += 226
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
	l.ctx.ppu.cfMux.Lock()
	defer l.ctx.ppu.cfMux.Unlock()

	if l.displayTileData {
		lineOffset := 386 * 4

		for i, p := range l.ctx.ppu.currentFrame {
			y := i / PpuXRes
			x := i % PpuXRes
			window := y*lineOffset + x*4

			l.buffer[window] = p.R
			l.buffer[window+1] = p.G
			l.buffer[window+2] = p.B
			l.buffer[window+3] = 0xFF
		}
		return
	}

	for i, p := range l.ctx.ppu.currentFrame {
		l.buffer[i*4] = p.R
		l.buffer[i*4+1] = p.G
		l.buffer[i*4+2] = p.B
		l.buffer[i*4+3] = 0xFF
	}
	return
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
				color := getColor(l.ctx.lcd.ctx.lcd.objectPallet0, hb, lb, uint8(x))

				l.buffer[i] = color.R
				l.buffer[i+1] = color.G
				l.buffer[i+2] = color.B
				l.buffer[i+3] = 0xFF
			}
		}
	}
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

// Layout implements ebiten.Game.
func (l *LcdRenderer) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return l.screenX, l.screenY
}

// Update implements ebiten.Game.
func (l *LcdRenderer) Update() error {
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		switch k {
		case ebiten.KeyComma:
			l.ctx.jpad.Up = true
		case ebiten.KeyO:
			l.ctx.jpad.Down = true
		case ebiten.KeyE:
			l.ctx.jpad.Right = true
		case ebiten.KeyA:
			l.ctx.jpad.Left = true
		case ebiten.KeyEnter:
			l.ctx.jpad.A = true
		case ebiten.KeyJ:
			l.ctx.jpad.B = true
		case ebiten.KeyQuote:
			l.ctx.jpad.Select = true
		case ebiten.KeyPeriod:
			l.ctx.jpad.Start = true
		}
	}
	for _, k := range inpututil.AppendJustReleasedKeys(nil) {
		switch k {
		case ebiten.KeyComma:
			l.ctx.jpad.Up = false
		case ebiten.KeyO:
			l.ctx.jpad.Down = false
		case ebiten.KeyE:
			l.ctx.jpad.Right = false
		case ebiten.KeyA:
			l.ctx.jpad.Left = false
		case ebiten.KeyEnter:
			l.ctx.jpad.A = false
		case ebiten.KeyJ:
			l.ctx.jpad.B = false
		case ebiten.KeyQuote:
			l.ctx.jpad.Select = false
		case ebiten.KeyPeriod:
			l.ctx.jpad.Start = false
		}
	}
	return nil
}

var _ ebiten.Game = (*LcdRenderer)(nil)
