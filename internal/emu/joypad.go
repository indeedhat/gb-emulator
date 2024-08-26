package emu

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

type KeyCode uint8

const (
	KeyA KeyCode = iota
	KeyB
	KeySelect
	KeyStart
	KeyUp
	KeyRight
	KeyDown
	KeyLeft
	KeyUnknown
)

type KeyEvent struct {
	Key  KeyCode
	Down bool
}

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

	go func() {
		for press := range ctx.JoypadCh {
			switch press.Key {
			case KeyUp:
				ctx.jpad.Up = press.Down
			case KeyDown:
				ctx.jpad.Down = press.Down
			case KeyRight:
				ctx.jpad.Right = press.Down
			case KeyLeft:
				ctx.jpad.Left = press.Down
			case KeyA:
				ctx.jpad.A = press.Down
			case KeyB:
				ctx.jpad.B = press.Down
			case KeySelect:
				ctx.jpad.Select = press.Down
			case KeyStart:
				ctx.jpad.Start = press.Down
			}
		}
	}()

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

	return value
}

func (j *Joypad) Write(value uint8) {
	j.ModeDpad = JpadModeDpad&value == 0
	j.ModeActions = JpadModeActions&value == 0
}
