package main

import (
	"time"
)

type Emulator struct {
	running bool
	paused  bool

	ctx *Context
}

func NewEmulator(romPath string, debug bool) (*Emulator, error) {
	e := &Emulator{}

	cartridge, err := LoadCartridge(romPath)
	if err != nil {
		return nil, err
	}

	e.ctx = &Context{cart: cartridge}

	NewMemoryBus(e.ctx)
	NewCpu(e.ctx)
	NewDebug(e.ctx)
	NewTimer(e.ctx)
	NewIO(e.ctx)
	NewPixelFetcher(e.ctx)
	NewPpu(e.ctx)
	NewJoypad(e.ctx)
	NewLcd(e.ctx)
	NewDma(e.ctx)

	e.ctx.debug.enbled = debug

	return e, nil
}

func (e *Emulator) Run() error {
	e.running = true

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
