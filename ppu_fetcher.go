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

	frame int
	done  bool

	bgTileId uint8
	bgLoBit  uint8
	bgHiBit  uint8

	spriteTileId []uint8
	spriteLoBit  []uint8
	spriteHiBit  []uint8
	fetchedOam   []OamEntry

	pixFifo *PixelFifo

	ctx *Context
}

func NewPixelFetcher(ctx *Context) {
	ctx.pix = &PixelFetcher{
		ctx: ctx,

		pixFifo: &PixelFifo{pixels: make([]Pixel, 16)},
	}
}

func (p *PixelFetcher) Reset() {
	p.mode = PixFetchModeTile
	p.pushed = 0
	p.fetched = 0
	p.lineX = 0
	p.fifoX = 0
	p.done = false
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
		p.fetchedOam = nil
		p.fetched += 8
		p.mode = PixFetchModeDataHigh

		if p.ctx.lcd.GetControl(LcdcBgwEnable) {
			p.bgTileId = p.ctx.membus.Read(
				p.ctx.lcd.BgTileAddress(uint16(p.mapX/8) + uint16(p.mapY/8)*32),
			)
			if !p.ctx.lcd.GetControl(LcdcBgwTileArea) {
				p.bgTileId += 128
			}
		}

		if p.ctx.lcd.GetControl(LcdcObjecteEnable) && len(p.ctx.ppu.activeSprites) > 0 {
			for _, entry := range p.ctx.ppu.activeSprites {
				x := (entry.x - 8) + p.ctx.lcd.scrollX%8

				if (x >= p.fetched && x < p.fetched+8) ||
					(x+8 >= p.fetched && x+8 < p.fetched+8) {

					p.fetchedOam = append(p.fetchedOam, entry)
				}

				if len(p.fetchedOam) == 3 {
					break
				}
			}
		}

		p.spriteHiBit = make([]uint8, len(p.fetchedOam))
		p.spriteLoBit = make([]uint8, len(p.fetchedOam))

	case PixFetchModeDataHigh:
		p.mode = PixFetchModeDataLow
		p.bgHiBit = p.ctx.membus.Read(
			p.ctx.lcd.BgWinTileAddress(uint16(p.bgTileId)*16 + uint16(p.tileY)),
		)
		p.loadSpriteTileData(true)

	case PixFetchModeDataLow:
		p.mode = PixFetchModeSleep
		p.bgLoBit = p.ctx.membus.Read(
			p.ctx.lcd.BgWinTileAddress(uint16(p.bgTileId)*16 + uint16(p.tileY) + 1),
		)
		p.loadSpriteTileData(false)

	case PixFetchModeSleep:
		p.mode = PixFetchModePush

	case PixFetchModePush:
		p.fetchPixels()
	}
}

func (p *PixelFetcher) loadSpriteTileData(hi bool) {
	var height uint8 = 8
	if p.ctx.lcd.GetControl(LcdcObjecteDoubleHeight) {
		height = 16
	}

	for i, entry := range p.fetchedOam {
		y := (p.ctx.lcd.ly + 16 - entry.y) * 2
		if entry.Check(OamFlagYflip) {
			y = height*2 - 2 - y
		}

		tileId := entry.tileIdx
		if p.ctx.lcd.GetControl(LcdcObjecteDoubleHeight) {
			tileId &= 0xFE
		}

		if hi {
			p.spriteHiBit[i] = p.ctx.membus.Read(0x8000 + uint16(tileId)*16 + uint16(y))
		} else {
			p.spriteLoBit[i] = p.ctx.membus.Read(0x8000 + uint16(tileId)*16 + uint16(y) + 1)
		}
	}
}

func (p *PixelFetcher) fetchPixels() {
	if p.pixFifo.fill > 8 {
		return
	}

	p.mode = PixFetchModeTile
	xPos := p.fetched - (8 - (p.ctx.lcd.scrollX % 8))

	for i := 7; i >= 0; i-- {
		cid := getColorIdx(p.bgLoBit, p.bgHiBit, uint8(i))
		c := getColor(p.ctx.lcd.backgroundPallet, p.bgLoBit, p.bgHiBit, uint8(i))

		if !p.ctx.lcd.GetControl(LcdcBgwEnable) {
			c = ColorPallet[p.ctx.lcd.backgroundPallet&0b11]
		}

		if p.ctx.lcd.GetControl(LcdcObjecteEnable) {
			if sp := p.fetchSpritePixel(cid); sp != nil {
				c = *sp
			}
		}

		if xPos >= 0 {
			p.pixFifo.Enqueue(c)
			p.fifoX++
		}
	}
}

func (p *PixelFetcher) fetchSpritePixel(bgColorId int) *Pixel {
	for i, entry := range p.fetchedOam {
		x := (entry.x - 8) + p.ctx.lcd.scrollX%8
		if x+8 < p.ctx.pix.fifoX {
			continue
		}

		offset := p.fifoX - x
		if offset < 0 || offset > 7 {
			continue
		}

		bit := 7 - offset
		if entry.Check(OamFlagXFlip) {
			bit = offset
		}

		hiBit := p.spriteHiBit[i]
		loBit := p.spriteLoBit[i]

		if getColorIdx(hiBit, loBit, bit) == 0 {
			continue
		}

		if entry.Check(OamFlagPriority) && bgColorId != 0 {
			continue
		}

		palette := p.ctx.lcd.objectPallet0
		if entry.Check(OamFlagDmgPalette) {
			palette = p.ctx.lcd.objectPallet1
		}

		pix := getColor(palette, hiBit, loBit, bit)
		return &pix
	}

	return nil
}

func (p *PixelFetcher) pushPixel() {
	if p.pixFifo.fill <= 8 {
		return
	}

	pix, ok := p.pixFifo.Dequeue()
	if !ok {
		panic("failed to get pixel data")
	}

	if p.lineX >= (p.ctx.lcd.scrollX % 8) {
		i := uint16(p.pushed) + (uint16(p.ctx.lcd.ly) * PpuXRes)

		p.ctx.ppu.nextFrame[i] = pix
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
