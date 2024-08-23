package main

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

	activeSprites []OamEntry

	ctx *Context
}

func NewPpu(ctx *Context) {
	ctx.ppu = &Ppu{
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
		ctx.ppu.blankFrame[i] = Pixel{0xFF, 0x00, 0xFF}
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

	if p.ctx.lcd.ly < PpuYRes {
		p.ctx.lcd.SetMode(LcdModeOam)
		return
	}

	p.ctx.lcd.SetMode(LcdModeVblank)

	p.ctx.cpu.requestInterrupt(InterruptVBlank)
	if p.ctx.lcd.GetStatus(LcdStatusVblank) {
		p.ctx.cpu.requestInterrupt(InterruptLcdStat)
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
	if p.ctx.lcd.ly >= PpuLinesPerFrame {
		p.ctx.lcd.SetMode(LcdModeOam)
		p.ctx.lcd.ly = 0

		if !p.ctx.pix.done {
			p.cfMux.Lock()
			copy(p.currentFrame, p.nextFrame)
			p.cfMux.Unlock()

			copy(p.nextFrame, p.blankFrame)
			p.ctx.pix.done = true
		}

		p.awaitNextFrame()
	}

	p.ticks = 0
}

func (p *Ppu) doOam() {
	if p.ticks == 79 {
		p.activeSprites = p.oam.SelectObjects(
			p.ctx.lcd.ly,
			p.ctx.lcd.GetControl(LcdcObjecteDoubleHeight),
		)
	}
	if p.ticks < 80 {
		return
	}

	p.ctx.lcd.SetMode(LcdModeDraw)
	p.ctx.pix.Reset()
}

func (p *Ppu) doDraw() {
	p.ctx.pix.Process()

	if p.ctx.pix.pushed < PpuXRes {
		return
	}

	p.ctx.pix.pixFifo.Reset()

	p.ctx.lcd.SetMode(LcdModeHblank)

	if p.ctx.lcd.GetStatus(LcdStatusHblank) {
		p.ctx.cpu.requestInterrupt(InterruptLcdStat)
	}
}
