package emu

import (
	"sync"
	"time"
)

const (
	PpuLinesPerFrame = 154
	PpuTicksPerLine  = 456
	PpuYRes          = 144
	PpuXRes          = 160
	TargetFrameTime  = time.Second / 60
)

type Ppu struct {
	oam  *OamRam
	vram *RamBank

	prevFrameTime time.Time
	ticks         uint64
	nextFrame     []Pixel
	blankFrame    []Pixel
	currentFrame  []Pixel
	cfMux         sync.Mutex
	windowX       uint

	activeSprites []OamEntry

	ctx *Context
}

func NewPpu(ctx *Context) {
	ppu := &Ppu{
		oam: &OamRam{make([]byte, 160)},
		vram: &RamBank{
			offset: 0x8000,
			data:   make([]byte, 0x2000),
		},
		nextFrame:    make([]Pixel, PpuYRes*PpuXRes),
		currentFrame: make([]Pixel, PpuYRes*PpuXRes),
		blankFrame:   make([]Pixel, PpuYRes*PpuXRes),
		ctx:          ctx,
	}

	for i := range PpuXRes * PpuYRes {
		ppu.blankFrame[i] = Pixel{0xFF, 0x00, 0xFF}
	}

	ctx.ppu = ppu
}

func (p *Ppu) Read(address uint16) uint8 {
	if address < 0xA000 {
		return p.vram.Read(address)
	}

	if address < 0xFEA0 {
		return p.oam.Read(address)
	}

	return 0xFF
}

func (p *Ppu) Write(address uint16, value uint8) {
	if address < 0xA000 {
		p.vram.Write(address, value)
	} else if address < 0xFEA0 {
		p.oam.Write(address, value)
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
	case LcdModeDraw:
		p.doDraw()
	}
}

func (p *Ppu) doHblank() {
	if p.ticks < PpuTicksPerLine {
		return
	}

	p.ctx.lcd.IncrementLine()
	p.ticks = 0

	if p.ctx.lcd.Ly() < PpuYRes {
		p.ctx.lcd.SetMode(LcdModeOam)
		return
	}

	p.ctx.lcd.SetMode(LcdModeVblank)

	p.ctx.cpu.RequestInterrupt(InterruptVBlank)
	if p.ctx.lcd.GetStatus(LcdStatusVblank) {
		p.ctx.cpu.RequestInterrupt(InterruptLcdStat)
	}
}

func (p *Ppu) awaitNextFrame() {
	now := time.Now()
	if curFrameTime := now.Sub(p.prevFrameTime); curFrameTime < TargetFrameTime {
		time.Sleep(TargetFrameTime - curFrameTime)
	}

	p.prevFrameTime = time.Now()
}

func (p *Ppu) doVblank() {
	if p.ticks < PpuTicksPerLine {
		return
	}

	p.ctx.lcd.IncrementLine()
	if p.ctx.lcd.Ly() >= PpuLinesPerFrame {
		p.ctx.lcd.SetMode(LcdModeOam)
		p.ctx.lcd.ResetLy()
		p.ctx.pix.(*PixelFetcher).windowX = 0

		if !p.ctx.pix.(*PixelFetcher).done {
			p.cfMux.Lock()
			copy(p.currentFrame, p.nextFrame)

			// send to renderer
			frame := make([]Pixel, len(p.currentFrame))
			copy(frame, p.currentFrame)
			p.ctx.FrameCh <- frame
			p.cfMux.Unlock()

			copy(p.nextFrame, p.blankFrame)
			p.ctx.pix.(*PixelFetcher).done = true
		}

		p.awaitNextFrame()
	}

	p.ticks = 0
}

func (p *Ppu) doOam() {
	if p.ticks == 79 {
		p.activeSprites = p.oam.SelectObjects(
			p.ctx.lcd.Ly(),
			p.ctx.lcd.GetControl(LcdcObjecteDoubleHeight),
		)
	}
	if p.ticks < 80 {
		return
	}

	p.ctx.lcd.SetMode(LcdModeDraw)
	p.ctx.pix.(*PixelFetcher).Reset()
}

func (p *Ppu) doDraw() {
	p.ctx.pix.(*PixelFetcher).Process()

	if p.ctx.pix.(*PixelFetcher).pushed < PpuXRes {
		return
	}

	p.ctx.pix.(*PixelFetcher).pixFifo.Reset()

	p.ctx.lcd.SetMode(LcdModeHblank)

	if p.ctx.lcd.GetStatus(LcdStatusHblank) {
		p.ctx.cpu.RequestInterrupt(InterruptLcdStat)
	}
}
