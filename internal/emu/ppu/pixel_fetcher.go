package ppu

import (
	"bytes"
	"encoding/binary"

	"github.com/indeedhat/gb-emulator/internal/emu/config"
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	. "github.com/indeedhat/gb-emulator/internal/emu/enum"
	"github.com/indeedhat/gb-emulator/internal/emu/palette"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

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

	windowX uint8

	frame int
	done  bool

	bgTileId uint8
	bgLoBit  uint8
	bgHiBit  uint8

	spriteLoBit []uint8
	spriteHiBit []uint8
	fetchedOam  []OamEntry

	pixFifo *PixelFifo

	ctx *context.Context
}

func NewPixelFetcher(ctx *context.Context) {
	ctx.Pix = &PixelFetcher{
		ctx: ctx,

		pixFifo: &PixelFifo{pixels: make([]Pixel, 16)},
	}
}

func (p *PixelFetcher) LoadState(data []byte) {
	r := bytes.NewReader(data)

	binary.Read(r, binary.BigEndian, &p.pushed)
	binary.Read(r, binary.BigEndian, &p.fetched)
	binary.Read(r, binary.BigEndian, &p.fifoX)
	binary.Read(r, binary.BigEndian, &p.lineX)
	binary.Read(r, binary.BigEndian, &p.tileX)
	binary.Read(r, binary.BigEndian, &p.tileY)
	binary.Read(r, binary.BigEndian, &p.mapX)
	binary.Read(r, binary.BigEndian, &p.mapY)
	binary.Read(r, binary.BigEndian, &p.windowX)
	binary.Read(r, binary.BigEndian, &p.frame)
	binary.Read(r, binary.BigEndian, &p.done)
	binary.Read(r, binary.BigEndian, &p.bgTileId)
	binary.Read(r, binary.BigEndian, &p.bgLoBit)
	binary.Read(r, binary.BigEndian, &p.bgHiBit)

	var l int64
	binary.Read(r, binary.BigEndian, &l)
	p.spriteLoBit = make([]uint8, l)
	r.Read(p.spriteLoBit)
	p.spriteHiBit = make([]uint8, l)
	r.Read(p.spriteHiBit)

	binary.Read(r, binary.BigEndian, &l)
	p.fetchedOam = make([]OamEntry, l)
	for i := range p.fetchedOam {
		binary.Read(r, binary.BigEndian, &p.fetchedOam[i].x)
		binary.Read(r, binary.BigEndian, &p.fetchedOam[i].y)
		binary.Read(r, binary.BigEndian, &p.fetchedOam[i].tileIdx)
		binary.Read(r, binary.BigEndian, &p.fetchedOam[i].flags)
	}

	binary.Read(r, binary.BigEndian, &p.pixFifo.head)
	binary.Read(r, binary.BigEndian, &p.pixFifo.tail)
	binary.Read(r, binary.BigEndian, &p.pixFifo.fill)
	for i := range p.pixFifo.pixels {
		binary.Read(r, binary.BigEndian, &p.pixFifo.pixels[i].R)
		binary.Read(r, binary.BigEndian, &p.pixFifo.pixels[i].G)
		binary.Read(r, binary.BigEndian, &p.pixFifo.pixels[i].B)
	}
}

func (p *PixelFetcher) SaveState() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, p.pushed)
	binary.Write(&buf, binary.BigEndian, p.fetched)
	binary.Write(&buf, binary.BigEndian, p.fifoX)
	binary.Write(&buf, binary.BigEndian, p.lineX)
	binary.Write(&buf, binary.BigEndian, p.tileX)
	binary.Write(&buf, binary.BigEndian, p.tileY)
	binary.Write(&buf, binary.BigEndian, p.mapX)
	binary.Write(&buf, binary.BigEndian, p.mapY)
	binary.Write(&buf, binary.BigEndian, p.windowX)
	binary.Write(&buf, binary.BigEndian, p.frame)
	binary.Write(&buf, binary.BigEndian, p.done)
	binary.Write(&buf, binary.BigEndian, p.bgTileId)
	binary.Write(&buf, binary.BigEndian, p.bgLoBit)
	binary.Write(&buf, binary.BigEndian, p.bgHiBit)

	binary.Write(&buf, binary.BigEndian, int64(len(p.spriteLoBit)))
	buf.Write(p.spriteLoBit)
	buf.Write(p.spriteHiBit)

	binary.Write(&buf, binary.BigEndian, int64(len(p.fetchedOam)))
	for i := range p.fetchedOam {
		binary.Write(&buf, binary.BigEndian, p.fetchedOam[i].x)
		binary.Write(&buf, binary.BigEndian, p.fetchedOam[i].y)
		binary.Write(&buf, binary.BigEndian, p.fetchedOam[i].tileIdx)
		binary.Write(&buf, binary.BigEndian, p.fetchedOam[i].flags)
	}

	binary.Write(&buf, binary.BigEndian, p.pixFifo.head)
	binary.Write(&buf, binary.BigEndian, p.pixFifo.tail)
	binary.Write(&buf, binary.BigEndian, p.pixFifo.fill)
	for i := range p.pixFifo.pixels {
		binary.Write(&buf, binary.BigEndian, p.pixFifo.pixels[i].R)
		binary.Write(&buf, binary.BigEndian, p.pixFifo.pixels[i].G)
		binary.Write(&buf, binary.BigEndian, p.pixFifo.pixels[i].B)
	}

	return buf.Bytes()
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
	p.mapX = p.fetched + p.ctx.Lcd.ScrollX()
	p.mapY = p.ctx.Lcd.Ly() + p.ctx.Lcd.ScrollY()
	p.tileY = (p.mapY % 8) * 2

	if p.ctx.Ppu.(*Ppu).ticks%2 == 0 {
		p.fetch()
	}

	p.pushPixel()
}

func (p *PixelFetcher) fetch() {
	switch p.mode {
	case PixFetchModeTile:
		p.doFetchModeTile()

	case PixFetchModeDataHigh:
		p.mode = PixFetchModeDataLow
		p.bgHiBit = p.ctx.Bus.Read(
			p.ctx.Lcd.BgWinTileAddress(uint16(p.bgTileId)*16 + uint16(p.tileY) + 1),
		)
		p.loadSpriteTileData(true)

	case PixFetchModeDataLow:
		p.mode = PixFetchModeSleep
		p.bgLoBit = p.ctx.Bus.Read(
			p.ctx.Lcd.BgWinTileAddress(uint16(p.bgTileId)*16 + uint16(p.tileY)),
		)
		p.loadSpriteTileData(false)

	case PixFetchModeSleep:
		p.mode = PixFetchModePush

	case PixFetchModePush:
		p.fetchPixels()
	}

}

func (p *PixelFetcher) doFetchModeTile() {
	p.fetchedOam = nil

	if p.ctx.Lcd.GetControl(LcdcBgwEnable) {
		if p.WindowVisible() &&
			(p.fetched+7 >= p.ctx.Lcd.WindowX() && p.fetched+7 < p.ctx.Lcd.WindowX()+config.PpuYRes+14) &&
			(p.ctx.Lcd.Ly() >= p.ctx.Lcd.WindowY() && p.ctx.Lcd.Ly() < p.ctx.Lcd.WindowY()+config.PpuXRes) {

			p.bgTileId = p.ctx.Bus.Read(p.ctx.Lcd.WinTileAddress(
				(uint16(p.fetched+7-p.ctx.Lcd.WindowX()) / 8) +
					(uint16(p.windowX)/8)*32,
			))
		} else {
			p.bgTileId = p.ctx.Bus.Read(
				p.ctx.Lcd.BgTileAddress(uint16(p.mapX/8) + uint16(p.mapY/8)*32),
			)
		}

		if !p.ctx.Lcd.GetControl(LcdcBgwTileArea) {
			p.bgTileId += 128
		}
	}

	if p.ctx.Lcd.GetControl(LcdcObjecteEnable) && len(p.ctx.Ppu.(*Ppu).activeSprites) > 0 {
		for _, entry := range p.ctx.Ppu.(*Ppu).activeSprites {
			x := (entry.x - 8) + p.ctx.Lcd.ScrollX()%8

			if (x >= p.fetched && x < p.fetched+8) ||
				(x+8 >= p.fetched && x+8 < p.fetched+8) {

				p.fetchedOam = append(p.fetchedOam, entry)
			}

			if len(p.fetchedOam) == 3 {
				break
			}
		}
	}

	p.fetched += 8
	p.mode = PixFetchModeDataHigh

	p.spriteHiBit = make([]uint8, len(p.fetchedOam))
	p.spriteLoBit = make([]uint8, len(p.fetchedOam))
}

func (p *PixelFetcher) loadSpriteTileData(hi bool) {
	var height uint8 = 8
	if p.ctx.Lcd.GetControl(LcdcObjecteDoubleHeight) {
		height = 16
	}

	for i, entry := range p.fetchedOam {
		y := (p.ctx.Lcd.Ly() + 16 - entry.y) * 2
		if entry.Check(OamFlagYflip) {
			y = height*2 - 2 - y
		}

		tileId := entry.tileIdx
		if p.ctx.Lcd.GetControl(LcdcObjecteDoubleHeight) {
			tileId &= 0xFE
		}

		if hi {
			p.spriteHiBit[i] = p.ctx.Bus.Read(0x8000 + uint16(tileId)*16 + uint16(y) + 1)
		} else {
			p.spriteLoBit[i] = p.ctx.Bus.Read(0x8000 + uint16(tileId)*16 + uint16(y))
		}
	}
}

func (p *PixelFetcher) fetchPixels() {
	if p.pixFifo.fill > 8 {
		return
	}

	p.mode = PixFetchModeTile
	xPos := p.fetched - (8 - (p.ctx.Lcd.ScrollX() % 8))

	for i := 7; i >= 0; i-- {
		cid := palette.GetColorIdx(p.bgHiBit, p.bgLoBit, uint8(i))
		c := palette.GetColor(p.ctx.Lcd.BackgroundPallet(), p.bgHiBit, p.bgLoBit, uint8(i))

		if !p.ctx.Lcd.GetControl(LcdcBgwEnable) {
			c = palette.ColorPallet[p.ctx.Lcd.BackgroundPallet()&0b11]
		}

		if p.ctx.Lcd.GetControl(LcdcObjecteEnable) {
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
		x := (entry.x - 8) + p.ctx.Lcd.ScrollX()%8
		if x+8 < p.fifoX {
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

		if palette.GetColorIdx(hiBit, loBit, bit) == 0 {
			continue
		}

		if entry.Check(OamFlagPriority) && bgColorId != 0 {
			continue
		}

		activePalette := p.ctx.Lcd.ObjectPallet(0)
		if entry.Check(OamFlagDmgPalette) {
			activePalette = p.ctx.Lcd.ObjectPallet(1)
		}

		pix := palette.GetColor(activePalette, hiBit, loBit, bit)
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

	if p.lineX >= (p.ctx.Lcd.ScrollX() % 8) {
		i := uint16(p.pushed) + (uint16(p.ctx.Lcd.Ly()) * config.PpuXRes)

		p.ctx.Ppu.(*Ppu).nextFrame[i] = pix
		p.pushed++
	}

	p.lineX++
}

func (p *PixelFetcher) WindowVisible() bool {
	if !p.ctx.Lcd.GetControl(LcdcWindowEnable) {
		return false
	}

	return p.ctx.Lcd.WindowX() >= 0 && p.ctx.Lcd.WindowX() <= 166 &&
		p.ctx.Lcd.WindowY() >= 0 && p.ctx.Lcd.WindowY() < config.PpuYRes
}

func (p *PixelFetcher) IncrementWindowX() {
	p.windowX++
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
