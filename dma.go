package main

type Dma struct {
	Active bool

	startDelay uint8
	byte       uint8
	addr       uint16

	ctx *Context
}

func NewDma(ctx *Context) {
	ctx.dma = &Dma{
		ctx: ctx,
	}
}

func (d *Dma) Start(value uint8) {
	d.Active = true
	d.byte = 0
	d.startDelay = 2
	d.addr = uint16(value)
}

func (d *Dma) Tick() {
	if !d.Active {
		return
	}

	if d.startDelay > 0 {
		d.startDelay--
		return
	}

	d.ctx.ppu.oam.Write(
		uint16(d.byte),
		d.ctx.membus.Read(d.addr*0x100+uint16(d.byte)),
	)

	d.byte++
	d.Active = d.byte < 0xA0
}
