package main

import "time"

const (
	PpuLinesPerFrame = 154
	PpuTicksPerLine  = 456
	PpuYRes          = 144
	PpuXRes          = 160
	TargetFrameTime  = time.Second / 60
)

type Ppu struct {
	oam  OamRam
	vram *RamBank

	prevFrameTime time.Time
	frameIdx      uint64
	ticks         uint64
	videoBuffre   []byte

	ctx *Context
}

func NewPpu(ctx *Context) {
	ctx.ppu = &Ppu{
		oam: make(OamRam, 40),
		vram: &RamBank{
			offset: 0x8000,
			data:   make([]byte, 0x2000),
		},
		videoBuffre: make([]byte, PpuYRes*PpuXRes*4),
		ctx:         ctx,
	}
}

func (p *Ppu) Tick() {
	p.ticks++

	switch p.ctx.lcd.GetMode() {
	case LcdModeHblank:
		p.doHblank()
	case LcdModeVblank:
		p.doVblank()
	case LcdModeOam:
		p.doOam()
	case LcdModeDrawLine:
		p.doDrawLine()
	}
}

func (p *Ppu) doHblank() {
	if p.ticks < PpuTicksPerLine {
		return
	}

	p.ctx.lcd.IncrementLine()
	p.ticks = 0

	if p.ctx.lcd.ly < PpuYRes {
		p.ctx.lcd.SetMode(LcdModeOam)
	}

	p.ctx.lcd.SetMode(LcdModeVblank)

	p.ctx.cpu.requestInterrupt(InterruptVBlank)
	p.frameIdx++

	p.awaitNextFrame()

}

func (p *Ppu) awaitNextFrame() {
	now := time.Now()
	if curFrameTime := now.Sub(p.prevFrameTime); curFrameTime < TargetFrameTime {
		time.Sleep(TargetFrameTime - curFrameTime)
	}

	p.prevFrameTime = now
}

func (p *Ppu) doVblank() {
	if p.ticks < PpuTicksPerLine {
		return
	}

	p.ctx.lcd.IncrementLine()
	if p.ctx.lcd.ly >= PpuLinesPerFrame {
		p.ctx.lcd.ly = 0
		p.ctx.lcd.SetMode(LcdModeOam)
	}

	p.ticks = 0
}

func (p *Ppu) doOam() {
	if p.ticks < 80 {
		return
	}

	p.ctx.lcd.SetMode(LcdModeDrawLine)
}

func (p *Ppu) doDrawLine() {
	if p.ticks >= 252 {
		return
	}

	p.ctx.lcd.SetMode(LcdModeHblank)
}

type OamEntry struct {
	x       uint8
	y       uint8
	tileIdx uint8
	// 7   Priority:    0 = No, 1 = BG and Window colors 1–3 are drawn over this OBJ
	// 6   Y flip:      0 = Normal, 1 = Entire OBJ is vertically mirrored
	// 5   X flip:      0 = Normal, 1 = Entire OBJ is horizontally mirrored
	// 4   DMG palette: 0 = OBP0, 1 = OBP1
	// 3   Bank:        0 = Fetch tile from VRAM bank 0, 1 = Fetch tile from VRAM bank 1
	// 0-2 CGB palette: Which of OBP0–7 to use
	flags uint8
}

func (e *OamEntry) Read(address uint8) uint8 {
	switch address {
	case 0:
		return e.x
	case 1:
		return e.y
	case 2:
		return e.tileIdx
	case 3:
		return e.flags
	default:
		panic("invalid address for oam entry")
	}
}

func (e *OamEntry) Write(address uint8, value uint8) {
	switch address {
	case 0:
		e.x = value
	case 1:
		e.y = value
	case 2:
		e.tileIdx = value
	case 3:
		e.flags = value
	default:
		panic("invalid address for oam entry")
	}
}

type OamRam []*OamEntry

func (o *OamRam) Read(address uint16) uint8 {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	i := address % 4
	return (*o)[(address-i)/4].Read(uint8(i))
}

func (o *OamRam) Write(address uint16, value uint8) {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	i := address % 4
	(*o)[(address-i)/4].Write(uint8(i), value)
}
