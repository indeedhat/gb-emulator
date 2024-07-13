package main

import (
	"errors"
	"log"
)

type Cpu struct {
	registers *cpuRegisters
	membus    *MemoryBus

	halted         bool
	masterInterupt bool
}

func NewCpu(bus *MemoryBus) *Cpu {
	return &Cpu{
		registers: &cpuRegisters{
			PC: 0x100,
			A:  0x01,
			SP: 0xFFFE,
		},
		membus:         bus,
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

func (c *Cpu) stackPop() uint16 {
	lsb := c.membus.Read(c.registers.SP)
	c.registers.SP++

	msb := c.membus.Read(c.registers.SP)
	c.registers.SP++
	return uint16(msb)<<8 | uint16(lsb)
}

func (c *Cpu) stackPush(value uint16) {
	c.registers.SP--
	c.membus.Write(c.registers.SP, uint8(value>>8))

	c.registers.SP--
	c.membus.Write(c.registers.SP, uint8(value))
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
	opcode := c.membus.Read(c.registers.PC)
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
		data = uint16(c.membus.Read(c.registers.PC)) | (uint16(c.membus.Read(c.registers.PC+1)) << 8)
		emu_cycle(2)
		c.registers.PC += 2
	case AddressModeR_N8:
		fallthrough
	case AddressModeN8:
		data = uint16(c.membus.Read(c.registers.PC))
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
		data = uint16(c.membus.Read(address))
		emu_cycle(1)
	case AddressModeR_A8:
		data = uint16(c.membus.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)
	case AddressModeA8_R:
		destAddr = &CpuDestAddress{uint16(c.membus.Read(c.registers.PC)) | 0xFF00}
		c.registers.PC++
		data = c.readFromRegister(instruction.Register2)
		emu_cycle(1)
	case AddressModeMR:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
	case AddressModeMR_N8:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.membus.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)
	case AddressModeMR_N16:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.membus.Read(c.registers.PC)) | (uint16(c.membus.Read(c.registers.PC+1)) << 8)
		c.registers.PC += 2
		emu_cycle(1)
	case AddressModeR_HLI:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.membus.Read(hl))
		emu_cycle(1)
		c.writeToRegister(RegisterTypeHL, hl+1)
	case AddressModeR_HLD:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.membus.Read(hl))
		emu_cycle(1)
		c.writeToRegister(RegisterTypeHL, hl-1)
	case AddressModeHLI_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{uint16(c.membus.Read(hl))}
		c.writeToRegister(RegisterTypeHL, hl+1)
		data = c.readFromRegister(instruction.Register2)
	case AddressModeHLD_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{uint16(c.membus.Read(hl))}
		c.writeToRegister(RegisterTypeHL, hl-1)
		data = c.readFromRegister(instruction.Register2)
	case AddressModeHL_SPR:
		data = uint16(c.membus.Read(c.registers.PC))
		emu_cycle(1)
		c.registers.PC++
	case AddressModeA16_R:
		destAddr = &CpuDestAddress{
			uint16(c.membus.Read(c.registers.PC)) | (uint16(c.membus.Read(c.registers.PC+1)) << 8),
		}
		emu_cycle(2)
		c.registers.PC += 2
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_A16:
		addr := uint16(c.membus.Read(c.registers.PC)) | (uint16(c.membus.Read(c.registers.PC+1)) << 8)
		data = uint16(c.membus.Read(addr))
		c.registers.PC += 2
		emu_cycle(3)
	}

	return data, destAddr
}

func (c *Cpu) executeInstruction(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) error {
	switch instruction.Type {
	case InstructionTypeNone:
		return errors.New("instruction not defined")
	case InstructionTypeADC:
		c.execADC(instruction, data)
	case InstructionTypeADD:
		c.execADD(instruction, data)
	case InstructionTypeAND:
		c.execAND(instruction, data)
	case InstructionTypeCALL:
		c.execCALL(instruction, data)
	case InstructionTypeCCF:
		c.execCCF()
	case InstructionTypeCP:
		c.execCP(instruction, data)
	case InstructionTypeCPL:
		c.execCPL()
	case InstructionTypeDAA:
		c.execDAA()
	case InstructionTypeDEC:
		c.execDEC(instruction, data, destAddress)
	case InstructionTypeDI:
		c.execDI()
	case InstructionTypeEI:
		c.execEI()
	case InstructionTypeHALT:
		// TODO
		fallthrough
	case InstructionTypeINC:
		c.execINC(instruction, data, destAddress)
	case InstructionTypeJP:
		c.execJP(instruction, data)
	case InstructionTypeLD:
		c.execLD(instruction, data, destAddress)
	case InstructionTypeLDH:
		c.execLDH(instruction, data)
	case InstructionTypeNOP:
		return nil // NOOP
	case InstructionTypeOR:
		c.execOR(instruction, data)
	case InstructionTypePOP:
		c.execPOP(instruction)
	case InstructionTypePUSH:
		c.execPUSH(instruction)

	case InstructionTypeXOR:
		c.execXOR(data)
	case InstructionTypeSUB:
		c.execSUB(instruction, data)
	case InstructionTypeSBC:
		c.execSBC(instruction, data)
	case InstructionTypeSTOP:
		fallthrough
	case InstructionTypeRLA:
		fallthrough
	case InstructionTypeJR:
		fallthrough
	case InstructionTypeRRA:
		fallthrough
	case InstructionTypeRLCA:
		fallthrough
	case InstructionTypeSCF:
		fallthrough
	case InstructionTypeRET:
		fallthrough
	case InstructionTypeCB:
		fallthrough
	case InstructionTypeRETI:
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
