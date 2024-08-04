package main

import (
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
	oam *OamRam
	// oam  OamRam
	vram *RamBank

	prevFrameTime time.Time
	ticks         uint64
	nextFrame     []byte
	currentFrame  []byte
	blankFrame    []byte

	ctx *Context
}

func NewPpu(ctx *Context) {
	ctx.ppu = &Ppu{
		oam: &OamRam{make([]byte, 160)},
		vram: &RamBank{
			offset: 0x8000,
			data:   make([]byte, 0x2000),
		},
		nextFrame:    make([]byte, PpuYRes*PpuXRes*4),
		currentFrame: make([]byte, PpuYRes*PpuXRes*4),
		blankFrame:   make([]byte, PpuYRes*PpuXRes*4),
		ctx:          ctx,
	}

	for i := range PpuXRes * PpuYRes {
		ctx.ppu.blankFrame[i*4] = 0xFF
		ctx.ppu.blankFrame[i*4+1] = 0x00
		ctx.ppu.blankFrame[i*4+2] = 0xFF
		ctx.ppu.blankFrame[i*4+3] = 0xFF
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

	copy(p.currentFrame, p.nextFrame)
	copy(p.nextFrame, p.blankFrame)
	p.awaitNextFrame()
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
	}

	p.ticks = 0
}

func (p *Ppu) doOam() {
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

	p.ctx.pix.bgFifo.Reset()

	p.ctx.lcd.SetMode(LcdModeHblank)

	if p.ctx.lcd.GetStatus(LcdStatusHblank) {
		p.ctx.cpu.requestInterrupt(InterruptLcdStat)
	}
}

type OamRam struct {
	data []byte
}

func (o *OamRam) Read(address uint16) uint8 {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	return o.data[address]
}

func (o *OamRam) Write(address uint16, value uint8) {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	o.data[address] = value
}
