package cpu

func (c *Cpu) execCB(_ CpuInstriction, cbyte uint16) bool {
	c.ctx.EmuCycle(1)

	if bitOp := uint8(cbyte >> 6 & 0x03); bitOp != 0 {
		c.execCB_BitOp(bitOp, cbyte)
		return true
	}

	opCode := uint8(cbyte >> 3 & 0x07)
	reg := c.cbRegLookup(cbyte & 0x07)

	if reg == RegisterTypeHL {
		c.ctx.EmuCycle(2)
	}

	var zflag, cflag, final uint8
	rbyte := c.cbReadRegister(reg)

	switch opCode {
	case 0: // RLC
		cflag = (rbyte >> 7) & 0x1
		final = (rbyte << 1) | cflag
	case 1: // RRC
		cflag = rbyte & 0x01
		final = (rbyte >> 1) | (cflag << 7)
	case 2: // RL
		cflag = uint8(rbyte >> 7)
		final = (rbyte << 1) | c.registers.GetFlag(CpuFlagC)
	case 3: // RR
		cflag = uint8(rbyte & 0x01)
		final = (rbyte >> 1) | c.registers.GetFlag(CpuFlagC)<<7&0xFF
	case 4: // SLA
		cflag = uint8(rbyte >> 7)
		final = (rbyte << 1) & 0xFF
	case 5: // SRA
		msb := rbyte & 0x80
		cflag = uint8(rbyte & 0x01)
		final = (rbyte >> 1) | msb
	case 6: // SWAP
		final = (rbyte << 4 & 0xF0) | (rbyte >> 4 & 0x0F)
	case 7: // SRL
		cflag = uint8(rbyte & 0x01)
		final = (rbyte >> 1) & 0x7F
	}

	if reg == RegisterTypeHL {
		c.ctx.Bus.Write(c.readFromRegister(RegisterTypeHL), final)
	} else {
		c.writeToRegister(reg, uint16(final))
	}

	if final == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0, 0, cflag)

	return true
}

func (c *Cpu) execCB_BitOp(bitOp uint8, cbyte uint16) {
	idx := uint8(cbyte >> 3 & 0x07)
	reg := c.cbRegLookup(cbyte & 0x07)

	if reg == RegisterTypeHL {
		c.ctx.EmuCycle(2)
	}

	switch bitOp {
	case 0x1: // BIT
		var zflag uint8
		if c.cbReadRegister(reg)>>idx&0x01 == 0 {
			zflag = 1
		}
		c.registers.SetFlags(zflag, 0, 1, 0xFF)
	case 0x2: // RES
		mask := ^uint8(1 << idx)
		c.cbWriteRegister(reg, c.cbReadRegister(reg)&mask)
	case 0x3: // SET
		c.cbWriteRegister(reg, c.cbReadRegister(reg)|1<<idx)
	}
}

func (c *Cpu) cbRegLookup(i uint16) RegisterType {
	return []RegisterType{
		RegisterTypeB,
		RegisterTypeC,
		RegisterTypeD,
		RegisterTypeE,
		RegisterTypeH,
		RegisterTypeL,
		RegisterTypeHL,
		RegisterTypeA,
	}[i]
}

func (c *Cpu) cbReadRegister(reg RegisterType) uint8 {
	switch reg {
	case RegisterTypeA:
		return c.registers.A
	case RegisterTypeB:
		return c.registers.B
	case RegisterTypeC:
		return c.registers.C
	case RegisterTypeD:
		return c.registers.D
	case RegisterTypeE:
		return c.registers.E
	case RegisterTypeH:
		return c.registers.H
	case RegisterTypeL:
		return c.registers.L
	case RegisterTypeHL:
		return c.ctx.Bus.Read(c.readFromRegister(RegisterTypeHL))
	default:
		return 0
	}
}

func (c *Cpu) cbWriteRegister(reg RegisterType, val uint8) {
	switch reg {
	case RegisterTypeA:
		c.registers.A = val
	case RegisterTypeB:
		c.registers.B = val
	case RegisterTypeC:
		c.registers.C = val
	case RegisterTypeD:
		c.registers.D = val
	case RegisterTypeE:
		c.registers.E = val
	case RegisterTypeH:
		c.registers.H = val
	case RegisterTypeL:
		c.registers.L = val
	case RegisterTypeHL:
		c.ctx.Bus.Write(c.readFromRegister(RegisterTypeHL), val)
	}
}
