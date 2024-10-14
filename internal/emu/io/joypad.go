package io

import (
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	. "github.com/indeedhat/gb-emulator/internal/emu/enum"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

const (
	JpadModeActions = uint8(1) << 5
	JpadModeDpad    = uint8(1) << 4
)

const (
	JpadDpadDown  = uint8(1) << 3
	JpadDpadUp    = uint8(1) << 2
	JpadDpadLeft  = uint8(1) << 1
	JpadDpadRight = uint8(1)
)

const (
	JpadActionStart  = uint8(1) << 3
	JpadActionSelect = uint8(1) << 2
	JpadActionB      = uint8(1) << 1
	JpadActionA      = uint8(1)
)

type Joypad struct {
	// data mode registers
	ModeDpad    bool
	ModeActions bool

	// d pad
	Up    bool
	Right bool
	Down  bool
	Left  bool

	// action buttons
	A      bool
	B      bool
	Start  bool
	Select bool

	ctx *context.Context
}

func NewJoypad(ctx *context.Context) ReadWriter {
	jpad := &Joypad{ctx: ctx}

	go func() {
		for press := range ctx.JoypadCh {
			switch press.Key {
			case KeyUp:
				jpad.Up = press.Down
			case KeyDown:
				jpad.Down = press.Down
			case KeyRight:
				jpad.Right = press.Down
			case KeyLeft:
				jpad.Left = press.Down
			case KeyA:
				jpad.A = press.Down
			case KeyB:
				jpad.B = press.Down
			case KeySelect:
				jpad.Select = press.Down
			case KeyStart:
				jpad.Start = press.Down
			}
		}
	}()

	ctx.Jpad = jpad
	return ctx.Jpad
}

func (j *Joypad) Read(_ uint16) uint8 {
	var value uint8 = 0xFF

	switch true {
	case j.ModeDpad:
		value &= ^JpadModeDpad
		if j.Down {
			value &= ^JpadDpadDown
		}
		if j.Up {
			value &= ^JpadDpadUp
		}
		if j.Left {
			value &= ^JpadDpadLeft
		}
		if j.Right {
			value &= ^JpadDpadRight
		}
	case j.ModeActions:
		value &= ^JpadModeActions
		if j.Start {
			value &= ^JpadActionStart
		}
		if j.Select {
			value &= ^JpadActionSelect
		}
		if j.B {
			value &= ^JpadActionB
		}
		if j.A {
			value &= ^JpadActionA
		}
	}

	return value
}

func (j *Joypad) Write(_ uint16, value uint8) {
	j.ModeDpad = JpadModeDpad&value == 0
	j.ModeActions = JpadModeActions&value == 0
}
