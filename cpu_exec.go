package main

import (
	"fmt"
	"regexp"
)

func (c *Cpu) execJP(instruction CpuInstriction, data uint16) {
	if c.registers.CheckFlag(instruction.Condition) {
		c.registers.PC = data
		emu_cycle(1)
	}
}

func (c *Cpu) execJR(instruction CpuInstriction, data uint16) {
	offset := int8(data & 0xFF00)
	// NB: this mess handles the int/uint conversions properly without having to resort to unsafe pointers
	c.execJP(instruction, uint16(int16(c.registers.PC)+int16(offset)))
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
	c.masterInterupt = false
}

func (c *Cpu) execEI() {
	c.masterInterupt = true
}

func (c *Cpu) execLD(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	if nil != destAddress {
		if instruction.Register2.Is16bit() {
			c.membus.Write16(destAddress.Address, data)
			emu_cycle(1)
		} else {
			c.membus.Write(destAddress.Address, uint8(data&0xFF))
		}

		return
	}

	if instruction.AddressMode == AddressModeHL_SPR {
		final := c.readFromRegister(instruction.Register2) + data
		hflag := halfCarry(c.readFromRegister(instruction.Register2), data, final)
		cflag := carry(c.readFromRegister(instruction.Register2), data, final)
		c.registers.SetFlags(0, 0, hflag, cflag)
		c.writeToRegister(instruction.Register1, final)
		return
	}

	c.writeToRegister(instruction.Register1, data)
}

func (c *Cpu) execLDH(instruction CpuInstriction, data uint16) {
	if instruction.Register1 == RegisterTypeA {
		c.writeToRegister(RegisterTypeA, uint16(c.membus.Read(0xFF00|data)&0xFF))
	} else {
		c.membus.Write(0xFF00|c.readFromRegister(RegisterTypeA), uint8(data))
	}
}

func (c *Cpu) execINC(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	if nil != destAddress {
		c.membus.Write(destAddress.Address, uint8(data&0xFF)+1)
		emu_cycle(1)
		return
	}

	if instruction.Register1 < RegisterTypeAF {
		var zflag uint8
		if 0 == uint8(data) {
			zflag = 1
		}

		hflag := halfCarry(data, 1, data+1)
		c.registers.SetFlags(zflag, 0, hflag, 0xFF)
	}

	c.writeToRegister(instruction.Register1, data+1)
}

func (c *Cpu) execDEC(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) {
	if nil != destAddress {
		c.membus.Write(destAddress.Address, uint8(data&0xFF)-1)
		emu_cycle(1)
		return
	}

	if !instruction.Register1.Is16bit() {
		var zflag uint8
		if 0 == uint8(data) {
			zflag = 1
		}

		hflag := halfCarry(data, 1, data-1)
		c.registers.SetFlags(zflag, 0, hflag, 0xFF)
	}

	c.writeToRegister(instruction.Register1, data-1)
}

func (c *Cpu) execADD(instruction CpuInstriction, data uint16) {
	var zflag, hflag, cflag uint8
	rval := c.readFromRegister(instruction.Register1)

	if instruction.Register1.Is16bit() {
		c.writeToRegister(instruction.Register1, rval+data)

		if instruction.Register1 != RegisterTypeSP {
			zflag = 0xFF
		}
	} else {
		final := uint8(rval + data)
		c.writeToRegister(instruction.Register1, uint16(final))

		if final == 0 {
			zflag = 1
		}
	}

	hflag = halfCarry(data, rval, rval+data)
	cflag = halfCarry(data, rval, rval+data)
	c.registers.SetFlags(zflag, 0, hflag, cflag)
}

func (c *Cpu) execSUB(instruction CpuInstriction, data uint16) {
	if instruction.Register2 == RegisterTypeA {
		c.writeToRegister(instruction.Register1, 0)
		c.registers.SetFlags(1, 1, 0, 0)
		return
	}

	rval := c.readFromRegister(instruction.Register1)
	final := uint8(rval - data)

	c.writeToRegister(instruction.Register1, uint16(data))

	var zflag uint8
	if final == 0 {
		zflag = 1
	}

	hflag := halfCarry(data, rval, rval-data)
	cflag := halfCarry(data, rval, rval-data)
	c.registers.SetFlags(zflag, 1, hflag, cflag)
}

func (c *Cpu) execADC(instruction CpuInstriction, data uint16) {
	var zflag uint8
	rval := c.readFromRegister(instruction.Register1)
	final := rval + data + uint16(c.registers.GetFlag(CpuFlagC))
	hflag := halfCarry(data, rval, final)
	cflag := halfCarry(data, rval, final)

	c.writeToRegister(instruction.Register1, final)

	if final == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0, hflag, cflag)
}

func (c *Cpu) execSBC(instruction CpuInstriction, data uint16) {
	var zflag uint8
	rval := c.readFromRegister(instruction.Register1) + uint16(c.registers.GetFlag(CpuFlagC))
	final := (rval - data) & 0xFF
	hflag := halfCarry(data, rval, final)
	cflag := halfCarry(data, rval, final)

	c.writeToRegister(instruction.Register1, final)

	if final == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 1, hflag, cflag)
}

func (c *Cpu) execAND(Instruction CpuInstriction, data uint16) {
	c.registers.A &= uint8(data)
	if Instruction.AddressMode != AddressModeR_R {
		emu_cycle(1)
	}

	var zflag uint8
	if c.registers.A == 0 {
		zflag = 1
	}

	c.registers.SetFlags(zflag, 0, 1, 0)
}

func (c *Cpu) execOR(Instruction CpuInstriction, data uint16) {
	c.registers.A |= uint8(data)
	if Instruction.AddressMode != AddressModeR_R {
		emu_cycle(1)
	}

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

	c.stackPush(c.registers.PC)
	c.registers.PC = data
	emu_cycle(3)
}

func (c *Cpu) execCCF() {
	var cflag uint8
	if c.registers.GetFlag(CpuFlagC) == 0 {
		cflag = 1
	}
	c.registers.SetFlags(0xFF, 0, 0, cflag)
}

func (c *Cpu) execCP(instruction CpuInstriction, data uint16) {
	if instruction.AddressMode != AddressModeR_R {
		emu_cycle(1)
	}
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
	cflag := halfCarry(data, rval, final)
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
		(c.registers.GetFlag(CpuFlagN) == 0 && c.registers.A > 99) {

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
	emu_cycle(2)

	if Instruction.Register1 == RegisterTypeAF {
		c.writeToRegister(Instruction.Register1, data&0xFFF0)
	} else {
		c.writeToRegister(Instruction.Register1, data)
	}
}

func (c *Cpu) execPUSH(Instruction CpuInstriction) {
	data := c.readFromRegister(Instruction.Register1)
	emu_cycle(2)

	c.stackPush(data)
	emu_cycle(1)
}

func (c *Cpu) execRET(instruction CpuInstriction) {
	if instruction.Condition != ConditionTypeNone {
		emu_cycle(1)

		if !c.registers.CheckFlag(instruction.Condition) {
			return
		}
	}

	data := c.stackPop()
	emu_cycle(3)

	c.registers.PC = data
}

func (c *Cpu) execRETI(instruction CpuInstriction) {
	c.masterInterupt = true
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

	c.stackPush(c.registers.PC)
	c.registers.PC = uint16(instruction.Param) & 0x0F
	emu_cycle(3)
}

func (c *Cpu) execSCF() {
	c.registers.SetFlags(0xFF, 0, 0, 1)
}

func (c *Cpu) execSTOP(_ uint16) {
	fmt.Println("execSTOP")
}

func (c *Cpu) execHALT() {
	c.halted = true
}

func (c *Cpu) execCB(Instruction CpuInstriction, cbyte uint16) {
	emu_cycle(1)

	if bitOp := uint8(cbyte >> 6 & 0x03); bitOp != 0 {
		c.execCB_BitOp(bitOp, Instruction, cbyte)
		return
	}

	// TODO: implement me
}

func (c *Cpu) execCB_BitOp(bitOp uint8, instruction CpuInstriction, cbyte uint16) {
	idx := uint8(cbyte >> 3 & 0x07)
	reg := c.cbRegLookup(int(cbyte & 0x07))

	if reg == RegisterTypeHL {
		emu_cycle(2)
	}

	switch bitOp {
	case 0x1: // BIT
		c.registers.SetFlags(uint8(c.readFromRegister(reg)>>idx&0x01), 0, 1, 0xFF)
	case 0x2: // RES
		mask := ^uint16(1 << idx)
		c.writeToRegister(reg, c.readFromRegister(reg)&mask)
	case 0x3: // SET
		c.writeToRegister(reg, c.readFromRegister(reg)|1<<idx)
	}
}
func (c *Cpu) cbRegLookup(i int) RegisterType {
	return []RegisterType{
		RegisterTypeB,
		RegisterTypeC,
		RegisterTypeD,
		RegisterTypeE,
		RegisterTypeH,
		RegisterTypeL,
		RegisterTypeHL,
		RegisterTypeA,
	}[0]
}
