package main

type Timer struct {
	div  uint16
	tima uint8
	tma  uint8
	tac  uint8

	ctx *Context
}

func NewTimer(ctx *Context) {
	ctx.timer = &Timer{
		div: 0xABCC,
		ctx: ctx,
	}
}

func (t *Timer) Tick() {
	pdiv := t.div
	t.div++

	var (
		clockDivider uint16
		tacDivider   uint8 = 0x4
		clockSelect        = t.tac & 0x3
	)

	switch clockSelect {
	case 0x0:
		clockDivider = 1 << 9
	case 0x1:
		clockDivider = 1 << 3
	case 0x2:
		clockDivider = 1 << 5
	case 0x3:
		clockDivider = 1 << 7
	}

	if (pdiv&clockDivider) == clockDivider &&
		(t.div&clockDivider) != clockDivider &&
		(t.tac&tacDivider) == tacDivider {

		t.tima++
		if t.tima == 0xFF {
			t.tima = t.tma

			t.ctx.cpu.requestInterrupt(InterruptTimer)
		}
	}
}

func (t *Timer) Write(addr uint16, value uint8) {
	switch addr {
	case 0xFF04:
		t.div = 0
	case 0xFF05:
		t.tima = value
	case 0xFF06:
		t.tma = value
	case 0xFF07:
		t.tac = value
	}
}

func (t *Timer) Read(addr uint16) uint8 {
	switch addr {
	case 0xFF04:
		return uint8(t.div >> 8)
	case 0xFF05:
		return t.tima
	case 0xFF06:
		return t.tma
	case 0xFF07:
		return t.tac
	default:
		panic("bad timer address")
	}
}
