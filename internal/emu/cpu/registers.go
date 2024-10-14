package cpu

const (
	CpuFlagZ uint8 = 1 << 7
	CpuFlagN uint8 = 1 << 6
	CpuFlagH uint8 = 1 << 5
	CpuFlagC uint8 = 1 << 4
)

type cpuRegisters struct {
	A  uint8
	F  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	H  uint8
	L  uint8
	SP uint16
	PC uint16
}

func (c *Cpu) readFromRegister(r RegisterType) uint16 {
	switch r {
	case RegisterTypeA:
		return uint16(c.registers.A)
	case RegisterTypeB:
		return uint16(c.registers.B)
	case RegisterTypeC:
		return uint16(c.registers.C)
	case RegisterTypeD:
		return uint16(c.registers.D)
	case RegisterTypeE:
		return uint16(c.registers.E)
	case RegisterTypeH:
		return uint16(c.registers.H)
	case RegisterTypeL:
		return uint16(c.registers.L)
	case RegisterTypeAF:
		return uint16(c.registers.A)<<8 | uint16(c.registers.F)
	case RegisterTypeBC:
		return uint16(c.registers.B)<<8 | uint16(c.registers.C)
	case RegisterTypeDE:
		return uint16(c.registers.D)<<8 | uint16(c.registers.E)
	case RegisterTypeHL:
		return uint16(c.registers.H)<<8 | uint16(c.registers.L)
	case RegisterTypeSP:
		return c.registers.SP
	case RegisterTypePC:
		return c.registers.PC
	}

	// NB: not really possible but keeps the compiler happy
	return 0
}

func (c *Cpu) writeToRegister(r RegisterType, val uint16) {
	switch r {
	case RegisterTypeA:
		c.registers.A = uint8(val)
	case RegisterTypeB:
		c.registers.B = uint8(val)
	case RegisterTypeC:
		c.registers.C = uint8(val)
	case RegisterTypeD:
		c.registers.D = uint8(val)
	case RegisterTypeE:
		c.registers.E = uint8(val)
	case RegisterTypeH:
		c.registers.H = uint8(val)
	case RegisterTypeL:
		c.registers.L = uint8(val)
	case RegisterTypeAF:
		c.registers.F = uint8(val) & 0xF0
		c.registers.A = uint8(val >> 8)
	case RegisterTypeBC:
		c.registers.C = uint8(val)
		c.registers.B = uint8(val >> 8)
	case RegisterTypeDE:
		c.registers.E = uint8(val)
		c.registers.D = uint8(val >> 8)
	case RegisterTypeHL:
		c.registers.L = uint8(val)
		c.registers.H = uint8(val >> 8)
	case RegisterTypeSP:
		c.registers.SP = val
	case RegisterTypePC:
		c.registers.PC = val
	}
}

func (r *cpuRegisters) CheckFlag(cond ConditionType) bool {
	c := r.F&CpuFlagC == CpuFlagC
	z := r.F&CpuFlagZ == CpuFlagZ

	switch cond {
	case ConditionTypeNone:
		return true
	case ConditionTypeC:
		return c
	case ConditionTypeNC:
		return !c
	case ConditionTypeZ:
		return z
	case ConditionTypeNZ:
		return !z
	default:
		return false
	}
}

func (r *cpuRegisters) GetFlag(flag uint8) uint8 {
	if r.F&flag == flag {
		return 1
	}
	return 0
}

func (r *cpuRegisters) SetFlags(z, n, h, c uint8) {
	if z == 1 {
		r.F |= CpuFlagZ
	} else if z == 0 {
		r.F &= ^CpuFlagZ
	}
	if n == 1 {
		r.F |= CpuFlagN
	} else if n == 0 {
		r.F &= ^CpuFlagN
	}
	if h == 1 {
		r.F |= CpuFlagH
	} else if h == 0 {
		r.F &= ^CpuFlagH
	}
	if c == 1 {
		r.F |= CpuFlagC
	} else if c == 0 {
		r.F &= ^CpuFlagC
	}
}

func halfCarry(a, b, result uint16) uint8 {
	if (a^b^result)&0x10 == 0x10 {
		return 1
	}
	return 0
}

func carry(a, b, result uint16) uint8 {
	if (a^b^result)&0x100 == 0x100 {
		return 1
	}
	return 0
}

func halfCarry16(a, b uint16) uint8 {
	if a&0xFFF+b&0xFFF >= 0x1000 {
		return 1
	}
	return 0
}

func carry16(a, b uint16) uint8 {
	if uint32(a&0xFFFF)+uint32(b&0xFFFF) >= 0x10000 {
		return 1
	}
	return 0
}
