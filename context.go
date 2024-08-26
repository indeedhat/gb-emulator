package main

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
}

func (c *Context) EmuCycle(i int) {
	for range i {
		for range 4 {
			c.ticks++
			c.timer.Tick()
			c.ppu.Tick()
		}

		c.dma.Tick()
	}
}
