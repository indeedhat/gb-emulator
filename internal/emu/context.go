package emu

type Context struct {
	ticks uint64

	cart   *Cartridge
	cpu    *Cpu
	debug  *Debug
	dma    *Dma
	jpad   *Joypad
	lcd    *Lcd
	membus *MemoryBus
	pix    *PixelFetcher
	ppu    *Ppu
	timer  *Timer
	io     *IO

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
