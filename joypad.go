package main

import "log"

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

	ctx *Context
}

func NewJoypad(ctx *Context) *Joypad {
	ctx.jpad = &Joypad{ctx: ctx}
	return ctx.jpad
}

func (j *Joypad) Read() uint8 {
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

	mode := ""
	if !j.ModeDpad {
		mode += "D"
	}
	if !j.ModeActions {
		mode += "A"
	}
	log.Printf("Read %2s: %02x", mode, value)
	return value
}

func (j *Joypad) Write(value uint8) {
	j.ModeDpad = JpadModeDpad&value == 0
	j.ModeActions = JpadModeActions&value == 0
}
