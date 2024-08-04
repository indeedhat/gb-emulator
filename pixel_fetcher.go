package main

type PixFetchMode uint8

const (
	PixFetchModeTile PixFetchMode = iota
	PixFetchModeDataHigh
	PixFetchModeDataLow
	PixFetchModeSleep
	PixFetchModePush
)

type PixelFetcher struct {
	mode PixFetchMode

	pushed  uint8
	fetched uint8
	fifoX   uint8
	lineX   uint8
	tileX   uint8
	tileY   uint8
	mapX    uint8
	mapY    uint8

	dataTileId  uint8
	dataLowBit  uint8
	dataHighBit uint8

	bgFifo  *PixelFifo
	objFifo *PixelFifo

	ctx *Context
}

func NewPixelFetcher(ctx *Context) {
	ctx.pix = &PixelFetcher{
		ctx: ctx,

		bgFifo:  &PixelFifo{pixels: make([]Pixel, 16)},
		objFifo: &PixelFifo{pixels: make([]Pixel, 16)},
	}
}

func (p *PixelFetcher) Reset() {
	p.mode = PixFetchModeTile
	p.pushed = 0
	p.fetched = 0
	p.lineX = 0
}

func (p *PixelFetcher) Process() {
	p.mapX = p.fetched + p.ctx.lcd.scrollX
	p.mapY = p.ctx.lcd.ly + p.ctx.lcd.scrollY
	p.tileY = (p.mapY % 8) * 2

	if p.ctx.ppu.ticks%2 == 0 {
		p.fetch()
	}

	p.pushPixel()
}

func (p *PixelFetcher) fetch() {
	switch p.mode {
	case PixFetchModeTile:
		p.fetched += 8
		p.mode = PixFetchModeDataHigh

		if !p.ctx.lcd.GetControl(LcdcBgwEnable) {
			return
		}

		p.dataTileId = p.ctx.membus.Read(
			p.ctx.lcd.BgTileAddress(uint16(p.mapX/8) + uint16(p.mapY/8)*32),
		)
		if !p.ctx.lcd.GetControl(LcdcBgwTileArea) {
			p.dataTileId += 128
		}
	case PixFetchModeDataHigh:
		p.mode = PixFetchModeDataLow
		p.dataHighBit = p.ctx.membus.Read(
			p.ctx.lcd.BgWinTileAddress(uint16(p.dataTileId)*16 + uint16(p.tileY)),
		)
	case PixFetchModeDataLow:
		p.mode = PixFetchModeSleep
		p.dataLowBit = p.ctx.membus.Read(
			p.ctx.lcd.BgWinTileAddress(uint16(p.dataTileId)*16 + uint16(p.tileY) + 1),
		)
	case PixFetchModeSleep:
		p.mode = PixFetchModePush
	case PixFetchModePush:
		if p.bgFifo.fill > 8 {
			return
		}

		p.mode = PixFetchModeTile
		xPos := p.fetched - (8 - (p.ctx.lcd.scrollX % 8))

		for i := 7; i >= 0; i-- {
			c := getColor(p.dataLowBit, p.dataHighBit, uint8(i))

			if xPos >= 0 {
				p.bgFifo.Enqueue(c)
			}
		}
	}
}

func (p *PixelFetcher) pushPixel() {
	if p.bgFifo.fill <= 8 {
		return
	}

	pix, ok := p.bgFifo.Dequeue()
	if !ok {
		panic("failed to get pixel")
	}

	if p.lineX >= (p.ctx.lcd.scrollX % 8) {
		i := uint16(p.pushed) + (uint16(p.ctx.lcd.ly) * PpuXRes)

		p.ctx.ppu.nextFrame[i*4] = pix.R
		p.ctx.ppu.nextFrame[i*4+1] = pix.G
		p.ctx.ppu.nextFrame[i*4+2] = pix.B
		p.ctx.ppu.nextFrame[i*4+3] = 0xFF

		p.pushed++
	}

	p.lineX++
}

type PixelFifo struct {
	pixels []Pixel
	head   int
	tail   int
	fill   int
}

func (p *PixelFifo) Enqueue(val Pixel) bool {
	if p.fill >= len(p.pixels) {
		return false
	}

	p.pixels[p.tail] = val
	p.tail = (p.tail + 1) % len(p.pixels)
	p.fill++

	return true
}

func (p *PixelFifo) Dequeue() (Pixel, bool) {
	if p.fill == 0 {
		return Pixel{}, false
	}

	val := p.pixels[p.head]
	p.head = (p.head + 1) % len(p.pixels)
	p.fill--

	return val, true
}

func (p *PixelFifo) Reset() {
	p.head = 0
	p.tail = 0
	p.fill = 0
}
