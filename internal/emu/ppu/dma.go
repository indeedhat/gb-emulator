package ppu

import (
	"bytes"
	"encoding/binary"

	"github.com/indeedhat/gb-emulator/internal/emu/context"
)

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

func (d *Dma) LoadState(data []byte) {
	r := bytes.NewReader(data)

	binary.Read(r, binary.BigEndian, &d.active)
	binary.Read(r, binary.BigEndian, &d.startDelay)
	binary.Read(r, binary.BigEndian, &d.byteIdx)
	binary.Read(r, binary.BigEndian, &d.addr)
}

func (d *Dma) SaveState() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, d.active)
	binary.Write(&buf, binary.BigEndian, d.startDelay)
	binary.Write(&buf, binary.BigEndian, d.byteIdx)
	binary.Write(&buf, binary.BigEndian, d.addr)

	return buf.Bytes()
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
