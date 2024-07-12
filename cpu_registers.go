package main

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

func (r cpuRegisters) CheckFlag(cond ConditionType) bool {
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

func (r cpuRegisters) GetFlag(flag uint8) uint8 {
	if r.F&flag == flag {
		return 1
	}
	return 0
}

func (r cpuRegisters) SetFlags(z, n, h, c uint8) {
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
