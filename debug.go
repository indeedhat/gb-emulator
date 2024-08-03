package main

import (
	"bytes"
	"log"
)

type Debug struct {
	enbled bool
	buf    *bytes.Buffer

	ctx *Context
}

func NewDebug(ctx *Context) {
	ctx.debug = &Debug{
		buf: new(bytes.Buffer),
		ctx: ctx,
	}
}

func (d *Debug) Update() {
	if d.ctx.membus.Read(0xFF02) != 0x81 || !d.enbled {
		return
	}

	d.buf.WriteByte(d.ctx.membus.Read(0xFF01))
	d.ctx.membus.Write(0xFF02, 0x00)
}

func (d *Debug) Print() {
	if d.buf.Len() == 0 || !d.enbled {
		return
	}

	log.Printf("[DEBUG]: %s", d.buf.String())
}
