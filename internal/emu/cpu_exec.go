package emu

func (c *Cpu) execJP(instruction CpuInstriction, data uint16) {
	if c.registers.CheckFlag(instruction.Condition) {
		c.registers.PC = data
		c.ctx.EmuCycle(1)
	}
}

func (c *Cpu) execJR(instruction CpuInstriction, data uint16) {
	offset := int8(uint8(data & 0xFF))
	// NB: this mess handles the int/uint conversions properly without having to resort to unsafe pointers
	c.execJP(instruction, uint16(int16((int32(c.registers.PC)+int32(offset))&0xFFFF)))
}

func (c *Cpu) execXOR(data uint16) {
	var isZero uint8
	c.registers.A ^= uint8(data)
	if c.registers.A == 0 {
		isZero = 1
	}
	c.registers.SetFlags(isZero, 0, 0, 0)
}

func (c *Cpu) execDI() {
	c.ime = false
}

func (c *Cpu) execEI() {
	c.enablingIME = true
}

func (c *Cpu) execLD(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	if nil != destAddress {
		if instruction.Register2.Is16bit() {
			c.ctx.EmuCycle(1)
			c.ctx.membus.Write16(destAddress.Address, data)
		} else {
			c.ctx.membus.Write(destAddress.Address, uint8(data&0xFF))
		}

		c.ctx.EmuCycle(1)
		return
	}

	if instruction.AddressMode == AddressModeHL_SPR {
		sp := c.readFromRegister(RegisterTypeSP)
		offset := int32(int8(uint8(data & 0xFF)))
		final := uint16(int16((int32(sp) + offset) & 0xFFFF))

		var hflag, cflag uint8
		if sp&0xF+data&0xF >= 0x10 {
			hflag = 1
		}
		if sp&0xFF+data&0xFF >= 0x100 {
			cflag = 1
		}
		c.registers.SetFlags(0, 0, hflag, cflag)
		c.writeToRegister(RegisterTypeHL, final)
		return
	}

	c.writeToRegister(instruction.Register1, data)
}

func (c *Cpu) execLDH(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	if instruction.Register1 == RegisterTypeA {
		c.writeToRegister(RegisterTypeA, uint16(c.ctx.membus.Read(0xFF00|data)))
	} else {
		c.ctx.membus.Write(destAddress.Address, c.registers.A)
	}

	c.ctx.EmuCycle(1)
}

func (c *Cpu) execINC(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	data++
	if instruction.Register1.Is16bit() {
		c.ctx.EmuCycle(1)
	}

	if destAddress != nil || !instruction.Register1.Is16bit() {
		data &= 0x00FF
	}

	var zflag uint8
	if data == 0 {
		zflag = 1
	}
	hflag := halfCarry(data-1, 1, data)

	if nil != destAddress {
		c.ctx.membus.Write(destAddress.Address, uint8(data))
	} else {
		c.writeToRegister(instruction.Register1, data)
	}

	if !instruction.Register1.Is16bit() || destAddress != nil {
		c.registers.SetFlags(zflag, 0, hflag, 0xFF)
	}
}

func (c *Cpu) execDEC(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	data--
	if instruction.Register1.Is16bit() {
		c.ctx.EmuCycle(1)
	}

	if destAddress != nil || !instruction.Register1.Is16bit() {
		data &= 0x00FF
	}

	var zflag uint8
	if data == 0 {
		zflag = 1
	}
	hflag := halfCarry(data+1, 1, data)

	if nil != destAddress {
		c.ctx.membus.Write(destAddress.Address, uint8(data))
	} else {
		c.writeToRegister(instruction.Register1, data)
	}

	if !instruction.Register1.Is16bit() || destAddress != nil {
		c.registers.SetFlags(zflag, 1, hflag, 0xFF)
	}
}

func (c *Cpu) execADD(instruction CpuInstriction, data uint16) {
	var zflag, hflag, cflag uint8
	rval := c.readFromRegister(instruction.Register1)
	final := rval + data

	if instruction.Register1 == RegisterTypeSP {
		c.ctx.EmuCycle(1)
		final = uint16(int16(rval) + int16(int8(data&0xFF)))
		hflag = halfCarry(data, rval, rval+data)
		if rval&0xF+data&0xF >= 0x10 {
			hflag = 1
		}
		if int16(rval&0xFF)+int16(data&0xFF) >= 0x100 {
			cflag = 1
		}
	} else if instruction.Register1.Is16bit() {
		c.ctx.EmuCycle(1)
		zflag = 0xFF
		hflag = halfCarry16(rval, data)
		cflag = carry16(rval, data)
	} else {
		final &= 0xFF
		if final == 0 {
			zflag = 1
		}

		hflag = halfCarry(data, rval, rval+data)
		cflag = carry(data, rval, rval+data)
	}

	c.writeToRegister(instruction.Register1, final)
	c.registers.SetFlags(zflag, 0, hflag, cflag)
}

func (c *Cpu) execADC(instruction CpuInstriction, data uint16) {
	rval := c.readFromRegister(instruction.Register1)
	cval := uint16(c.registers.GetFlag(CpuFlagC))
	final := uint8((rval + data + cval) & 0x00FF)

	var zflag, hflag, cflag uint8
	if rval&0xF+data&0xF+cval >= 0x10 {
		hflag = 1
	}
	if rval&0xFF+data&0xFF+cval >= 0x100 {
		cflag = 1
	}
	if final == 0 {
		zflag = 1
	}

	c.writeToRegister(instruction.Register1, uint16(final))

	c.registers.SetFlags(zflag, 0, hflag, cflag)
}

func (c *Cpu) execSUB(instruction CpuInstriction, data uint16) {
	rval := c.readFromRegister(instruction.Register1)
	final := rval - data

	var zflag uint8
	if final == 0 {
		zflag = 1
	}
	hflag := halfCarry(rval, data, final)
	cflag := carry(rval, data, final)

	c.writeToRegister(instruction.Register1, final)
	c.registers.SetFlags(zflag, 1, hflag, cflag)
}

func (c *Cpu) execSBC(instruction CpuInstriction, data uint16) {
	var zflag, cflag, hflag uint8
	rval := c.readFromRegister(instruction.Register1)
	cval := uint16(c.registers.GetFlag(CpuFlagC))
	final := (rval - data - cval) & 0xFF
	if data+cval > rval {
		cflag = 1
	}
	if int8(rval&0xF)-int8(data&0xF)-int8(cval) < 0 {
		hflag = 1
	}

	c.writeToRegister(instruction.Register1, final)

	if final == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 1, hflag, cflag)
}

func (c *Cpu) execAND(_ CpuInstriction, data uint16) {
	c.registers.A &= uint8(data)

	var zflag uint8
	if c.registers.A == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0, 1, 0)
}

func (c *Cpu) execOR(_ CpuInstriction, data uint16) {
	c.registers.A |= uint8(data)
	var zflag uint8
	if c.registers.A == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0, 0, 0)
}

func (c *Cpu) execCALL(instruction CpuInstriction, data uint16) {
	if instruction.Condition != ConditionTypeNone && !c.registers.CheckFlag(instruction.Condition) {
		return
	}

	c.ctx.EmuCycle(2)
	c.stackPush(c.registers.PC)
	c.registers.PC = data
	c.ctx.EmuCycle(1)
}

func (c *Cpu) execCCF() {
	var cflag uint8
	if c.registers.GetFlag(CpuFlagC) == 0 {
		cflag = 1
	}
	c.registers.SetFlags(0xFF, 0, 0, cflag)
}

func (c *Cpu) execCP(instruction CpuInstriction, data uint16) {
	if instruction.Register2 == RegisterTypeA {
		c.registers.SetFlags(1, 1, 0, 0)
		return
	}

	rval := c.readFromRegister(instruction.Register1)
	final := rval - data

	var zflag uint8
	if final == 0 {
		zflag = 1
	}

	hflag := halfCarry(data, rval, final)
	cflag := carry(data, rval, final)
	c.registers.SetFlags(zflag, 1, hflag, cflag)
}

func (c *Cpu) execCPL() {
	c.registers.A = ^c.registers.A
	c.registers.SetFlags(0xFF, 1, 1, 0xFF)
}

func (c *Cpu) execDAA() {
	var cflag, zflag uint8
	var adjustment uint8

	if c.registers.GetFlag(CpuFlagH) == 1 ||
		(c.registers.GetFlag(CpuFlagN) == 0 && c.registers.A&0xF > 9) {

		adjustment = 6
	}

	if c.registers.GetFlag(CpuFlagC) == 1 ||
		(c.registers.GetFlag(CpuFlagN) == 0 && c.registers.A > 0x99) {

		adjustment |= 0x60
		cflag = 1
	}

	if c.registers.GetFlag(CpuFlagN) == 1 {
		adjustment = -adjustment
	}

	c.registers.A += adjustment
	if c.registers.A == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0xFF, 0, cflag)
}

func (c *Cpu) execPOP(Instruction CpuInstriction) {
	data := c.stackPop()
	c.ctx.EmuCycle(2)

	if Instruction.Register1 == RegisterTypeAF {
		c.writeToRegister(Instruction.Register1, data&0xFFF0)
	} else {
		c.writeToRegister(Instruction.Register1, data)
	}
}

func (c *Cpu) execPUSH(Instruction CpuInstriction) {
	data := c.readFromRegister(Instruction.Register1)
	c.ctx.EmuCycle(2)

	c.stackPush(data)
	c.ctx.EmuCycle(1)
}

func (c *Cpu) execRET(instruction CpuInstriction) {
	if instruction.Condition != ConditionTypeNone {
		c.ctx.EmuCycle(1)
	}

	if !c.registers.CheckFlag(instruction.Condition) {
		return
	}

	data := c.stackPop()
	c.ctx.EmuCycle(2)

	c.registers.PC = data
	c.ctx.EmuCycle(1)
}

func (c *Cpu) execRETI(instruction CpuInstriction) {
	c.ime = true
	c.execRET(instruction)
}

func (c *Cpu) execRLA() {
	cflag := c.registers.GetFlag(CpuFlagC)
	c.registers.SetFlags(0, 0, 0, c.registers.A&0x80>>7)
	c.registers.A = (c.registers.A << 1) | cflag
}

func (c *Cpu) execRLCA() {
	msb := c.registers.A & 0x80 >> 7
	c.registers.SetFlags(0, 0, 0, msb)
	c.registers.A = (c.registers.A << 1) | msb
}

func (c *Cpu) execRRA() {
	cflag := c.registers.GetFlag(CpuFlagC)
	c.registers.SetFlags(0, 0, 0, c.registers.A&0x01)
	c.registers.A = (c.registers.A >> 1) | cflag<<7
}

func (c *Cpu) execRRCA() {
	lsb := c.registers.A & 0x01
	c.registers.SetFlags(0, 0, 0, lsb)
	c.registers.A = (c.registers.A >> 1) | lsb<<7
}

func (c *Cpu) execRST(instruction CpuInstriction) {
	if instruction.Condition != ConditionTypeNone && !c.registers.CheckFlag(instruction.Condition) {
		return
	}

	c.ctx.EmuCycle(2)
	c.stackPush(c.registers.PC)
	c.registers.PC = uint16(instruction.Param) & 0xFF
	c.ctx.EmuCycle(1)
}

func (c *Cpu) execSCF() {
	c.registers.SetFlags(0xFF, 0, 0, 1)
}

func (c *Cpu) execSTOP(_ uint16) {
	// writing to the timers div register resets it
	// c.ctx.membus.Write(0xFF04, 0x01)
	// panic("stop not implemented")
}

func (c *Cpu) execHALT() {
	c.halted = true
}
