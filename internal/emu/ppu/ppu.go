package ppu

import (
	"sync"
	"time"

	"github.com/indeedhat/gb-emulator/internal/emu/config"
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	. "github.com/indeedhat/gb-emulator/internal/emu/enum"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
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

	ctx *context.Context
}

func New(ctx *context.Context) {
	ppu := &Ppu{
		oam:          &OamRam{make([]byte, 160)},
		vram:         NewRamBank(0x8000, 0x2000),
		nextFrame:    make([]Pixel, config.PpuYRes*config.PpuXRes),
		currentFrame: make([]Pixel, config.PpuYRes*config.PpuXRes),
		blankFrame:   make([]Pixel, config.PpuYRes*config.PpuXRes),
		ctx:          ctx,
	}

	for i := range config.PpuXRes * config.PpuYRes {
		ppu.blankFrame[i] = Pixel{R: 0xFF, G: 0x00, B: 0xFF}
	}

	ctx.Ppu = ppu
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

	switch p.ctx.Lcd.GetMode() {
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
	if p.ticks < config.PpuTicksPerLine {
		return
	}

	p.ctx.Lcd.IncrementLine()
	p.ticks = 0

	if p.ctx.Lcd.Ly() < config.PpuYRes {
		p.ctx.Lcd.SetMode(LcdModeOam)
		return
	}

	p.ctx.Lcd.SetMode(LcdModeVblank)

	p.ctx.Cpu.RequestInterrupt(InterruptVBlank)
	if p.ctx.Lcd.GetStatus(LcdStatusVblank) {
		p.ctx.Cpu.RequestInterrupt(InterruptLcdStat)
	}
}

func (p *Ppu) awaitNextFrame() {
	now := time.Now()
	if curFrameTime := now.Sub(p.prevFrameTime); curFrameTime < config.TargetFrameTime {
		time.Sleep(config.TargetFrameTime - curFrameTime)
	}

	p.prevFrameTime = time.Now()
}

func (p *Ppu) doVblank() {
	if p.ticks < config.PpuTicksPerLine {
		return
	}

	p.ctx.Lcd.IncrementLine()
	if p.ctx.Lcd.Ly() >= config.PpuLinesPerFrame {
		p.ctx.Lcd.SetMode(LcdModeOam)
		p.ctx.Lcd.ResetLy()
		p.ctx.Pix.(*PixelFetcher).windowX = 0

		if !p.ctx.Pix.(*PixelFetcher).done {
			p.cfMux.Lock()
			copy(p.currentFrame, p.nextFrame)

			// send to renderer
			frame := make([]Pixel, len(p.currentFrame))
			copy(frame, p.currentFrame)
			p.ctx.FrameCh <- frame
			p.cfMux.Unlock()

			copy(p.nextFrame, p.blankFrame)
			p.ctx.Pix.(*PixelFetcher).done = true
		}

		p.awaitNextFrame()
	}

	p.ticks = 0
}

func (p *Ppu) doOam() {
	if p.ticks == 79 {
		p.activeSprites = p.oam.SelectObjects(
			p.ctx.Lcd.Ly(),
			p.ctx.Lcd.GetControl(LcdcObjecteDoubleHeight),
		)
	}
	if p.ticks < 80 {
		return
	}

	p.ctx.Lcd.SetMode(LcdModeDraw)
	p.ctx.Pix.(*PixelFetcher).Reset()
}

func (p *Ppu) doDraw() {
	p.ctx.Pix.(*PixelFetcher).Process()

	if p.ctx.Pix.(*PixelFetcher).pushed < config.PpuXRes {
		return
	}

	p.ctx.Pix.(*PixelFetcher).pixFifo.Reset()

	p.ctx.Lcd.SetMode(LcdModeHblank)

	if p.ctx.Lcd.GetStatus(LcdStatusHblank) {
		p.ctx.Cpu.RequestInterrupt(InterruptLcdStat)
	}
}
