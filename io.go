package main

import (
	"log"
)

type IO struct {
	serial []uint8

	ctx *Context
}

func NewIO(ctx *Context) {
	ctx.io = &IO{
		serial: make([]uint8, 2),
		ctx:    ctx,
	}
}

func (i *IO) Read(addr uint16) uint8 {
	switch true {
	case addr == 0xFF01:
		return i.serial[0]
	case addr == 0xFF02:
		return i.serial[1]
	case addr >= 0xFF04 && addr <= 0xFF07:
		return i.ctx.timer.Read(addr)
	case addr == 0xFF0F:
		return i.ctx.cpu.interuptFlags
	default:
		log.Printf("unsupported mem.read (IO) 0x%X", addr)
		return 0
	}
}

func (i *IO) Write(addr uint16, value uint8) {
	switch true {
	case addr == 0xFF01:
		i.serial[0] = value
	case addr == 0xFF02:
		i.serial[1] = value
	case addr >= 0xFF04 && addr <= 0xFF07:
		i.ctx.timer.Write(addr, value)
	case addr == 0xFF0F:
		i.ctx.cpu.interuptFlags = value
	default:
		log.Printf("unsupported mem.write (IO) 0x%X", addr)
	}
}
