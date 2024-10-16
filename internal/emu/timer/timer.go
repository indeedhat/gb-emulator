package timer

import (
	"bytes"
	"encoding/binary"

	"github.com/indeedhat/gb-emulator/internal/emu/context"
	. "github.com/indeedhat/gb-emulator/internal/emu/enum"
)

type Timer struct {
	div  uint16
	tima uint8
	tma  uint8
	tac  uint8

	ctx *context.Context
}

func New(ctx *context.Context) {
	ctx.Timer = &Timer{
		div: 0xABCC,
		ctx: ctx,
	}
}

func (t *Timer) LoadState(data []byte) {
	r := bytes.NewReader(data)

	binary.Read(r, binary.BigEndian, &t.div)
	binary.Read(r, binary.BigEndian, &t.tima)
	binary.Read(r, binary.BigEndian, &t.tma)
	binary.Read(r, binary.BigEndian, &t.tac)
}

func (t *Timer) SaveState() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, t.div)
	binary.Write(&buf, binary.BigEndian, t.tima)
	binary.Write(&buf, binary.BigEndian, t.tma)
	binary.Write(&buf, binary.BigEndian, t.tac)

	return buf.Bytes()
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

			t.ctx.Cpu.RequestInterrupt(InterruptTimer)
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
