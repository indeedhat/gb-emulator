package context

import (
	"bytes"
	"encoding/binary"

	. "github.com/indeedhat/gb-emulator/internal/emu/enum"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

type Context struct {
	ticks uint64

	Cart interface {
		ReadWriter
		SaveLoader

		Mbc() MBC
	}

	Cpu interface {
		RequestInterrupt(itype uint8)
		Step() error
		InterruptFlags() uint8
		SetInterruptFlags(value uint8)
		InterruptRegister() uint8
		SetInterruptRegister(value uint8)
	}
	Debug interface {
		Update()
		Print()
		Enabled() bool
	}
	Dma interface {
		Ticker

		Active() bool
		Start(value uint8)
	}
	Lcd interface {
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
	Bus ReadWriter16
	Pix interface {
		WindowVisible() bool
		IncrementWindowX()
	}
	Ppu interface {
		ReadWriter
		Ticker
	}
	Timer interface {
		ReadWriter
		Ticker
	}
	Io ReadWriter

	FrameCh  chan []Pixel
	JoypadCh chan KeyEvent
}

func NewContext() *Context {
	return &Context{
		FrameCh:  make(chan []Pixel, 2),
		JoypadCh: make(chan KeyEvent, 2),
	}
}

func (c *Context) LoadState(data []byte) {
	r := bytes.NewReader(data)

	var size int64
	binary.Read(r, binary.BigEndian, &size)
	tmp := make([]byte, size)
	r.Read(tmp)
	c.Cart.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Cpu.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Dma.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Lcd.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Bus.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Pix.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Ppu.(Stator).LoadState(tmp)

	binary.Read(r, binary.BigEndian, &size)
	tmp = make([]byte, size)
	r.Read(tmp)
	c.Timer.(Stator).LoadState(tmp)
}

func (c *Context) SaveState() []byte {
	var buf bytes.Buffer

	tmp := c.Cart.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Cpu.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Dma.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Lcd.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Bus.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Pix.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Ppu.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	tmp = c.Timer.(Stator).SaveState()
	binary.Write(&buf, binary.BigEndian, int64(len(tmp)))
	buf.Write(tmp)

	return buf.Bytes()
}

func (c *Context) Ticks() uint64 {
	return c.ticks
}

func (c *Context) EmuCycle(i uint8) {
	for range i {
		for range 4 {
			c.ticks++
			c.Timer.Tick()
			c.Ppu.Tick()
		}

		c.Dma.Tick()
	}
}
