package main

import (
	"time"
)

type Emulator struct {
	running bool
	paused  bool

	ctx *Context
}

func (e *Emulator) Run(romPath string) error {
	e.running = true

	cartridge, err := LoadCartridge(romPath)
	if err != nil {
		return err
	}

	e.ctx = &Context{cart: cartridge}

	NewMemoryBus(e.ctx)
	NewCpu(e.ctx)
	NewDebug(e.ctx)
	NewTimer(e.ctx)
	NewIO(e.ctx)
	NewPpu(e.ctx)
	NewLcd(e.ctx)

	e.ctx.debug.enbled = true

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
