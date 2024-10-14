package emu

import (
	"log"
	"time"
)

type Emulator struct {
	running bool
	paused  bool

	ctx *Context
}

func NewEmulator(romPath string, debug bool) (*Emulator, *Context, error) {
	e := &Emulator{}

	cartridge, err := LoadCartridge(romPath)
	if err != nil {
		return nil, nil, err
	}

	e.ctx = NewContext()
	e.ctx.cart = cartridge

	NewMemoryBus(e.ctx)
	NewCpu(e.ctx)
	NewDebug(e.ctx, debug)
	NewTimer(e.ctx)
	NewIO(e.ctx)
	NewPixelFetcher(e.ctx)
	NewPpu(e.ctx)
	NewJoypad(e.ctx)
	NewLcd(e.ctx)
	NewDma(e.ctx)

	return e, e.ctx, nil
}

func (e *Emulator) Run() error {
	e.running = true

	go e.saveBatteryRam()

	for {
		if e.paused {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if err := e.ctx.cpu.Step(); err != nil {
			return err
		}
	}
}

func (e *Emulator) saveBatteryRam() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if mbc3, ok := e.ctx.cart.Mbc().(Ticker); ok {
			mbc3.Tick()
		}

		if err := e.ctx.cart.Save(); err != nil {
			log.Printf("failed to save battery backed ram: %s", err)
		}
	}
}
