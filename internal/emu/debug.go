package emu

import (
	"bytes"
	"log"
)

type Debug struct {
	enabled bool
	buf     *bytes.Buffer

	ctx *Context
}

func NewDebug(ctx *Context, enabled bool) {
	ctx.debug = &Debug{
		buf:     new(bytes.Buffer),
		ctx:     ctx,
		enabled: enabled,
	}
}

func (d *Debug) Update() {
	if d.ctx.membus.Read(0xFF02) != 0x81 || !d.enabled {
		return
	}

	d.buf.WriteByte(d.ctx.membus.Read(0xFF01))
	d.ctx.membus.Write(0xFF02, 0x00)
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
