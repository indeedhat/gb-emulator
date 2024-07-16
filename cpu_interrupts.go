package main

const (
	InterruptVBlank uint8 = iota
	InterruptLcdStat
	InterruptTimer
	InterruptSerial
	InterruptJoyPad
)

func (c *Cpu) requestInterrupts() {

}

func (c *Cpu) handleInterrupts() {
	_ = c.attempInterrupt(InterruptVBlank, 0x40) ||
		c.attempInterrupt(InterruptLcdStat, 0x48) ||
		c.attempInterrupt(InterruptTimer, 0x50) ||
		c.attempInterrupt(InterruptSerial, 0x58) ||
		c.attempInterrupt(InterruptJoyPad, 0x60)
}

func (c *Cpu) attempInterrupt(interrupt uint8, addr uint16) bool {
	if c.interuptFlags&c.interruptRegister&interrupt != interrupt {
		return false
	}

	c.stackPush(c.registers.PC)
	c.registers.PC = addr
	c.interuptFlags &= ^interrupt
	c.halted = false
	c.ime = false

	return true
}
