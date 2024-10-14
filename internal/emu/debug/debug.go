package debug

import (
	"bytes"
	"log"

	"github.com/indeedhat/gb-emulator/internal/emu/context"
)

type Debug struct {
	enabled bool
	buf     *bytes.Buffer

	ctx *context.Context
}

func New(ctx *context.Context, enabled bool) {
	ctx.Debug = &Debug{
		buf:     new(bytes.Buffer),
		ctx:     ctx,
		enabled: enabled,
	}
}

func (d *Debug) Update() {
	if d.ctx.Bus.Read(0xFF02) != 0x81 || !d.enabled {
		return
	}

	d.buf.WriteByte(d.ctx.Bus.Read(0xFF01))
	d.ctx.Bus.Write(0xFF02, 0x00)
}

func (d *Debug) Print() {
	if d.buf.Len() == 0 || !d.enabled {
		return
	}

	log.Printf("[DEBUG]: %s", d.buf.String())
}

func (d *Debug) Enabled() bool {
	return d.enabled
}
