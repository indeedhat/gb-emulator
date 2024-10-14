package ppu

import "github.com/indeedhat/gb-emulator/internal/emu/context"

type Dma struct {
	active bool

	startDelay uint8
	byteIdx    uint8
	addr       uint16

	ctx *context.Context
}

func NewDma(ctx *context.Context) {
	ctx.Dma = &Dma{
		ctx: ctx,
	}
}

func (d *Dma) Active() bool {
	return d.active
}

func (d *Dma) Start(value uint8) {
	d.active = true
	d.byteIdx = 0
	d.startDelay = 2
	d.addr = uint16(value)
}

func (d *Dma) Tick() {
	if !d.active {
		return
	}

	if d.startDelay > 0 {
		d.startDelay--
		return
	}

	d.ctx.Ppu.Write(
		uint16(d.byteIdx)+0xFE00,
		d.ctx.Bus.Read((d.addr*0x100)+uint16(d.byteIdx)),
	)

	d.byteIdx++
	d.active = d.byteIdx < 0xA0
}
