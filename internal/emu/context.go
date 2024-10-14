package emu

type Context struct {
	ticks uint64

	cart interface {
		ReadWriter
		SaveLoader

		Mbc() MBC
	}

	cpu interface {
		RequestInterrupt(itype uint8)
		Step() error
		InterruptFlags() uint8
		SetInterruptFlags(value uint8)
		InterruptRegister() uint8
		SetInterruptRegister(value uint8)
	}
	debug interface {
		Update()
		Print()
		Enabled() bool
	}
	dma interface {
		Ticker

		Active() bool
		Start(value uint8)
	}
	jpad ReadWriter
	lcd  interface {
		ReadWriter

		GetMode() LcdMode
		SetMode(mode LcdMode)
		IncrementLine()
		GetStatus(code LcdStatus) bool
		GetControl(code Lcdc) bool
		String(pc uint16) string
		Ly() uint8
		ResetLy()
		ScrollX() uint8
		ScrollY() uint8
		WindowX() uint8
		WindowY() uint8

		BgWinTileAddress(address uint16) uint16
		WinTileAddress(address uint16) uint16
		BgTileAddress(address uint16) uint16
		BackgroundPallet() uint8
		ObjectPallet(i uint8) uint8
	}
	membus ReadWriter16
	pix    interface {
		WindowVisible() bool
		IncrementWindowX()
	}
	ppu interface {
		ReadWriter
		Ticker
	}
	timer interface {
		ReadWriter
		Ticker
	}
	io ReadWriter

	FrameCh  chan []Pixel
	JoypadCh chan KeyEvent
}

func NewContext() *Context {
	return &Context{
		FrameCh:  make(chan []Pixel, 2),
		JoypadCh: make(chan KeyEvent, 2),
	}
}

func (c *Context) EmuCycle(i uint8) {
	for range i {
		for range 4 {
			c.ticks++
			c.timer.Tick()
			c.ppu.Tick()
		}

		c.dma.Tick()
	}
}
