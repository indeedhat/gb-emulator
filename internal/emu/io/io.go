package io

import "github.com/indeedhat/gb-emulator/internal/emu/context"

type IO struct {
	serial []uint8

	ctx *context.Context
}

func New(ctx *context.Context) {
	ctx.Io = &IO{
		serial: make([]uint8, 2),
		ctx:    ctx,
	}
}

func (i *IO) Read(addr uint16) uint8 {
	switch true {
	case addr == 0xFF00:
		return i.ctx.Jpad.Read(0)
	case addr == 0xFF01:
		return i.serial[0]
	case addr == 0xFF02:
		return i.serial[1]
	case addr >= 0xFF04 && addr <= 0xFF07:
		return i.ctx.Timer.Read(addr)
	case addr == 0xFF0F:
		return i.ctx.Cpu.InterruptFlags()
	case addr >= 0xFF40 && addr <= 0xFF4B:
		return i.ctx.Lcd.Read(addr)
	default:
		// log.Printf("unsupported mem.read (IO) 0x%X", addr)
		return 0
	}
}

func (i *IO) Write(addr uint16, value uint8) {
	switch true {
	case addr == 0xFF00:
		i.ctx.Jpad.Write(0, value)
	case addr == 0xFF01:
		i.serial[0] = value
	case addr == 0xFF02:
		i.serial[1] = value
	case addr >= 0xFF04 && addr <= 0xFF07:
		i.ctx.Timer.Write(addr, value)
	case addr == 0xFF0F:
		i.ctx.Cpu.SetInterruptFlags(value)
	case addr >= 0xFF40 && addr <= 0xFF4B:
		i.ctx.Lcd.Write(addr, value)
	default:
		// log.Printf("unsupported mem.write (IO) 0x%X", addr)
	}
}
