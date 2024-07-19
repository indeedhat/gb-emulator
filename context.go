package main

type Context struct {
	ticks uint64

	cart   *Cartridge
	cpu    *Cpu
	debug  *Debug
	membus *MemoryBus
	timer  *Timer
	io     *IO
}

func (c *Context) EmuCycle(i int) {
	for range i * 4 {
		c.ticks++
		c.timer.Tick()
	}
}
