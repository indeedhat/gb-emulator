package main

import (
	"errors"
	"log"
)

type Cpu struct {
	registers *cpuRegisters
	memory    *MemoryBus

	halted         bool
	masterInterupt bool
}

func NewCpu(bus *MemoryBus) *Cpu {
	return &Cpu{
		registers: &cpuRegisters{
			PC: 0x100,
			A:  0x01,
		},
		memory:         bus,
		masterInterupt: true,
	}
}

func (c *Cpu) Step() error {
	if c.halted {
		return nil
	}

	pc := c.registers.PC
	opcode, instruction := c.fetchIsntruction()
	data, destAddress := c.fetchData(instruction)

	log.Printf("[PC: 0x%X] 0x%X => %v -- 0x%X", pc, opcode, instruction, data)
	return c.executeInstruction(instruction, data, destAddress)
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
		return uint16(c.registers.F) | (uint16(c.registers.A) << 8)
	case RegisterTypeBC:
		return uint16(c.registers.C) | (uint16(c.registers.B) << 8)
	case RegisterTypeDE:
		return uint16(c.registers.E) | (uint16(c.registers.D) << 8)
	case RegisterTypeHL:
		return uint16(c.registers.L) | (uint16(c.registers.H) << 8)
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
		c.registers.F = uint8(val)
		c.registers.A = uint8(val >> 8)
	case RegisterTypeBC:
		c.registers.B = uint8(val)
		c.registers.C = uint8(val >> 8)
	case RegisterTypeDE:
		c.registers.D = uint8(val)
		c.registers.E = uint8(val >> 8)
	case RegisterTypeHL:
		c.registers.H = uint8(val)
		c.registers.L = uint8(val >> 8)
	case RegisterTypeSP:
		c.registers.SP = val
	case RegisterTypePC:
		c.registers.PC = val
	}
}

func (c *Cpu) fetchIsntruction() (uint8, CpuInstriction) {
	opcode := c.memory.Read(c.registers.PC)
	c.registers.PC++

	return opcode, CpuInstructions[opcode]
}

type CpuDestAddress struct {
	Address uint16
}

func (c *Cpu) fetchData(instruction CpuInstriction) (data uint16, destAddr *CpuDestAddress) {
	switch instruction.AddressMode {
	case AddressModeR:
		data = c.readFromRegister(instruction.Register1)

	case AddressModeR_N16:
		fallthrough
	case AddressModeN16:
		data = uint16(c.memory.Read(c.registers.PC)) | (uint16(c.memory.Read(c.registers.PC+1)) << 8)
		emu_cycle(2)
		c.registers.PC += 2
	case AddressModeR_N8:
		fallthrough
	case AddressModeN8:
		data = uint16(c.memory.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)
	case AddressModeR_R:
		data = c.readFromRegister(instruction.Register2)
	case AddressModeMR_R:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = c.readFromRegister(instruction.Register2)
		if instruction.Register1 == RegisterTypeC {
			destAddr.Address |= 0xFF00
		}
	case AddressModeR_MR:
		address := c.readFromRegister(instruction.Register2)
		if instruction.Register2 == RegisterTypeC {
			address |= 0xFF00
		}
		data = uint16(c.memory.Read(address))
		emu_cycle(1)
	case AddressModeR_A8:
		data = uint16(c.memory.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)
	case AddressModeA8_R:
		destAddr = &CpuDestAddress{uint16(c.memory.Read(c.registers.PC)) | 0xFF00}
		c.registers.PC++
		data = c.readFromRegister(instruction.Register2)
		emu_cycle(1)
	case AddressModeMR:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
	case AddressModeMR_N8:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.memory.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)
	case AddressModeMR_N16:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.memory.Read(c.registers.PC)) | (uint16(c.memory.Read(c.registers.PC+1)) << 8)
		c.registers.PC += 2
		emu_cycle(1)
	case AddressModeR_HLI:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.memory.Read(hl))
		emu_cycle(1)
		c.writeToRegister(RegisterTypeHL, hl+1)
	case AddressModeR_HLD:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.memory.Read(hl))
		emu_cycle(1)
		c.writeToRegister(RegisterTypeHL, hl-1)
	case AddressModeHLI_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{uint16(c.memory.Read(hl))}
		c.writeToRegister(RegisterTypeHL, hl+1)
		data = c.readFromRegister(instruction.Register2)
	case AddressModeHLD_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{uint16(c.memory.Read(hl))}
		c.writeToRegister(RegisterTypeHL, hl-1)
		data = c.readFromRegister(instruction.Register2)
	case AddressModeHL_SPR:
		data = uint16(c.memory.Read(c.registers.PC))
		emu_cycle(1)
		c.registers.PC++
	case AddressModeA16_R:
		destAddr = &CpuDestAddress{
			uint16(c.memory.Read(c.registers.PC)) | (uint16(c.memory.Read(c.registers.PC+1)) << 8),
		}
		emu_cycle(2)
		c.registers.PC += 2
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_A16:
		addr := uint16(c.memory.Read(c.registers.PC)) | (uint16(c.memory.Read(c.registers.PC+1)) << 8)
		data = uint16(c.memory.Read(addr))
		c.registers.PC += 2
		emu_cycle(3)
	}

	return data, destAddr
}

func (c *Cpu) executeInstruction(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) error {
	switch instruction.Type {
	case InstructionTypeJP:
		if c.registers.CheckFlag(instruction.Condition) {
			c.registers.PC = data
			emu_cycle(1)
		}
	case InstructionTypeXOR:
		var isZero uint8
		c.registers.A ^= uint8(data)
		if c.registers.A == 0 {
			isZero = 1
		}
		c.registers.SetFlags(isZero, 0, 0, 0)
	case InstructionTypeNOP:
		return nil
	case InstructionTypeDI:
		c.masterInterupt = false
	case InstructionTypeEI:
		c.masterInterupt = true
	case InstructionTypeNone:
		return errors.New("instruction not defined")
	case InstructionTypeLD:
		if nil != destAddress {
			if instruction.Register2.Is16bit() {
				c.memory.Write16(destAddress.Address, data)
				emu_cycle(1)
			} else {
				c.memory.Write(destAddress.Address, uint8(data&0xFF))
			}
			return nil
		}
		if instruction.AddressMode == AddressModeHL_SPR {
			final := c.readFromRegister(instruction.Register2) + data
			hflag := halfCarry(c.readFromRegister(instruction.Register2), data, final)
			cflag := carry(c.readFromRegister(instruction.Register2), data, final)
			c.registers.SetFlags(0, 0, hflag, cflag)
			c.writeToRegister(instruction.Register1, final)
			return nil
		}
		c.writeToRegister(instruction.Register1, data)
	case InstructionTypeINC:
		if nil != destAddress {
			c.memory.Write(destAddress.Address, uint8(data&0xFF)+1)
			emu_cycle(1)
			return nil
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
	case InstructionTypeDEC:
		if nil != destAddress {
			c.memory.Write(destAddress.Address, uint8(data&0xFF)-1)
			emu_cycle(1)
			return nil
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
	case InstructionTypeADD:
		var rval uint16
		var zflag, hflag, cflag uint8
		if instruction.Register1.Is16bit() {
			rval = c.readFromRegister(instruction.Register1)
			c.writeToRegister(instruction.Register1, rval+data)
			if instruction.Register1 != RegisterTypeSP {
				zflag = 0xFF
			}
		} else {
			rval = c.readFromRegister(instruction.Register1)
			final := uint8(rval + data)
			c.writeToRegister(instruction.Register1, uint16(data))
			if final == 0 {
				zflag = 1
			}
		}
		hflag = halfCarry(data, rval, rval+data)
		cflag = halfCarry(data, rval, rval+data)
		c.registers.SetFlags(zflag, 0, hflag, cflag)
	case InstructionTypeSUB:
		if instruction.Register2 == RegisterTypeA {
			c.writeToRegister(instruction.Register1, 0)
			c.registers.SetFlags(1, 1, 0, 0)
			return nil
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
		c.registers.SetFlags(zflag, 0, hflag, cflag)
	case InstructionTypeADC:
		fallthrough
	case InstructionTypeSBC:
		fallthrough
	case InstructionTypeRLCA:
		fallthrough
	case InstructionTypeSTOP:
		fallthrough
	case InstructionTypeRLA:
		fallthrough
	case InstructionTypeJR:
		fallthrough
	case InstructionTypeRRA:
		fallthrough
	case InstructionTypeDAA:
		fallthrough
	case InstructionTypeCPL:
		fallthrough
	case InstructionTypeSCF:
		fallthrough
	case InstructionTypeCCF:
		fallthrough
	case InstructionTypeHALT:
		fallthrough
	case InstructionTypeAND:
		fallthrough
	case InstructionTypeOR:
		fallthrough
	case InstructionTypeCP:
		fallthrough
	case InstructionTypePOP:
		fallthrough
	case InstructionTypePUSH:
		fallthrough
	case InstructionTypeRET:
		fallthrough
	case InstructionTypeCB:
		fallthrough
	case InstructionTypeCALL:
		fallthrough
	case InstructionTypeRETI:
		fallthrough
	case InstructionTypeLDH:
		fallthrough
	case InstructionTypeJPHL:
		fallthrough
	case InstructionTypeRST:
		fallthrough
	case InstructionTypeERR:
		fallthrough
	case InstructionTypePREFIX:
		fallthrough

	// CB instructions...
	case InstructionTypeRLC:
		fallthrough
	case InstructionTypeRRC:
		fallthrough
	case InstructionTypeRRCA:
		fallthrough
	case InstructionTypeRL:
		fallthrough
	case InstructionTypeRR:
		fallthrough
	case InstructionTypeSLA:
		fallthrough
	case InstructionTypeSRA:
		fallthrough
	case InstructionTypeSWAP:
		fallthrough
	case InstructionTypeSRL:
		fallthrough
	case InstructionTypeBIT:
		fallthrough
	case InstructionTypeRES:
		fallthrough
	case InstructionTypeSET:
		fallthrough
	default:
		return errors.New("not implemented")
	}
	return nil
}
