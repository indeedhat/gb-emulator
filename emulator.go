package main

import (
	"time"
)

type Emulator struct {
	running bool
	paused  bool
	ticks   uint64
}

func (e *Emulator) Run(romPath string) error {
	e.running = true

	cartridge, err := LoadCartridge(romPath)
	if err != nil {
		return err
	}

	memory := &MemoryBus{Cart: cartridge}
	cpu := NewCpu(memory)

	for i := 0; i < 20; i++ {
		if e.paused {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if err := cpu.Step(); err != nil {
			return err
		}

		e.ticks++
	}

	return nil
}

func emu_cycle(i int) {

}
