package main

import (
	"errors"
	"log"
)

type Cpu struct {
	Registers *cpuRegisters
	Memory    *MemoryBus

	Halted         bool
	MasterInterupt bool
}

func NewCpu(bus *MemoryBus) *Cpu {
	return &Cpu{
		Registers: &cpuRegisters{
			PC: 0x100,
			A:  0x01,
		},
		Memory:         bus,
		MasterInterupt: true,
	}
}

func (c *Cpu) Step() error {
	if c.Halted {
		return nil
	}

	pc := c.Registers.PC
	opcode, instruction := c.fetchIsntruction()
	data := c.fetchData(instruction)

	log.Printf("[PC: %X] %X => %v -- %X", pc, opcode, instruction, data)
	return c.executeInstruction(instruction, data)
}

func (c *Cpu) readFromRegister(r RegisterType) uint16 {
	switch r {
	case RegisterTypeA:
		return uint16(c.Registers.A)
	case RegisterTypeB:
		return uint16(c.Registers.B)
	case RegisterTypeC:
		return uint16(c.Registers.C)
	case RegisterTypeD:
		return uint16(c.Registers.D)
	case RegisterTypeE:
		return uint16(c.Registers.E)
	case RegisterTypeH:
		return uint16(c.Registers.H)
	case RegisterTypeL:
		return uint16(c.Registers.L)
	case RegisterTypeAF:
		return uint16(c.Registers.F) | (uint16(c.Registers.A) << 8)
	case RegisterTypeBC:
		return uint16(c.Registers.C) | (uint16(c.Registers.B) << 8)
	case RegisterTypeDE:
		return uint16(c.Registers.E) | (uint16(c.Registers.D) << 8)
	case RegisterTypeHL:
		return uint16(c.Registers.L) | (uint16(c.Registers.H) << 8)
	case RegisterTypeSP:
		return c.Registers.SP
	case RegisterTypePC:
		return c.Registers.PC
	}

	// NB: not really possible but keeps the compiler happy
	return 0
}

func (c *Cpu) fetchIsntruction() (uint8, CpuInstriction) {
	opcode := c.Memory.Read(c.Registers.PC)
	c.Registers.PC++

	return opcode, CpuInstructions[opcode]
}

func (c *Cpu) fetchData(instruction CpuInstriction) uint16 {
	var fetchedData uint16

	switch instruction.AddressMode {
	case AddressModeR:
		fetchedData = c.readFromRegister(instruction.Register1)
	case AddressModeR_N16:
	case AddressModeR_R:
	case AddressModeR_N8:
		fetchedData = uint16(c.Memory.Read(c.Registers.PC))
		c.Registers.PC++
		emu_cycle(1)
	case AddressModeR_MR:
	case AddressModeR_HLI:
	case AddressModeR_HLD:
	case AddressModeR_A16:
	case AddressModeR_A8:
	case AddressModeN8:
	case AddressModeN16:
		fetchedData = uint16(c.Memory.Read(c.Registers.PC)) | (uint16(c.Memory.Read(c.Registers.PC+1)) << 8)
		emu_cycle(2)
		c.Registers.PC += 2
	case AddressModeN16_R:
	case AddressModeMR:
	case AddressModeMR_N8:
	case AddressModeMR_N16:
	case AddressModeMR_R:
	case AddressModeA8_R:
	case AddressModeA16_R:
	case AddressModeHLI_R:
	case AddressModeHLD_R:
	case AddressModeHL_SPR:
	}

	return fetchedData
}

func (c *Cpu) executeInstruction(instruction CpuInstriction, data uint16) error {
	switch instruction.Type {
	case InstructionTypeJP:
		if c.Registers.CheckFlag(instruction.Condition) {
			c.Registers.PC = data
			emu_cycle(1)
		}
	case InstructionTypeXOR:
		var isZero uint8
		c.Registers.A ^= uint8(data)
		if c.Registers.A == 0 {
			isZero = 1
		}
		c.Registers.SetFlags(isZero, 0, 0, 0)
	case InstructionTypeNOP:
		return nil
	case InstructionTypeDI:
		c.MasterInterupt = false
	case InstructionTypeEI:
		c.MasterInterupt = true

	case InstructionTypeNone:
		return errors.New("instruction not defined")
	case InstructionTypeLD:
		fallthrough
	case InstructionTypeINC:
		fallthrough
	case InstructionTypeDEC:
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
	case InstructionTypeADD:
		fallthrough
	case InstructionTypeSUB:
		fallthrough
	case InstructionTypeADC:
		fallthrough
	case InstructionTypeSBC:
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
