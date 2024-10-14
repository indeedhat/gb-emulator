package emu

import (
	"log"
	"time"

	"github.com/indeedhat/gb-emulator/internal/emu/cart"
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	"github.com/indeedhat/gb-emulator/internal/emu/cpu"
	"github.com/indeedhat/gb-emulator/internal/emu/debug"
	"github.com/indeedhat/gb-emulator/internal/emu/io"
	"github.com/indeedhat/gb-emulator/internal/emu/lcd"
	"github.com/indeedhat/gb-emulator/internal/emu/memory"
	"github.com/indeedhat/gb-emulator/internal/emu/ppu"
	"github.com/indeedhat/gb-emulator/internal/emu/timer"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

type Emulator struct {
	running bool
	paused  bool

	ctx *context.Context
}

func NewEmulator(romPath string, debugEnabled bool) (*Emulator, *context.Context, error) {
	e := &Emulator{}

	cartridge, err := cart.Load(romPath)
	if err != nil {
		return nil, nil, err
	}

	e.ctx = context.NewContext()
	e.ctx.Cart = cartridge

	memory.NewBus(e.ctx)
	cpu.New(e.ctx)
	debug.New(e.ctx, debugEnabled)
	timer.New(e.ctx)
	io.New(e.ctx)
	io.NewJoypad(e.ctx)
	ppu.New(e.ctx)
	ppu.NewPixelFetcher(e.ctx)
	ppu.NewDma(e.ctx)
	lcd.New(e.ctx)

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

		if err := e.ctx.Cpu.Step(); err != nil {
			return err
		}
	}
}

func (e *Emulator) saveBatteryRam() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if mbc3, ok := e.ctx.Cart.Mbc().(Ticker); ok {
			mbc3.Tick()
		}

		if err := e.ctx.Cart.Save(); err != nil {
			log.Printf("failed to save battery backed ram: %s", err)
		}
	}
}
