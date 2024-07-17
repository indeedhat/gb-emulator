package main

import (
	"errors"
	"fmt"
	"log"
)

type Cpu struct {
	registers *cpuRegisters
	membus    *MemoryBus

	halted bool

	// interrupts
	ime               bool
	enablingIME       bool
	interuptFlags     uint8
	interruptRegister uint8
}

func (c *Cpu) String() string {
	zf, hf, nf, cf := "-", "-", "-", "-"
	if c.registers.GetFlag(CpuFlagZ) == 1 {
		zf = "Z"
	}
	if c.registers.GetFlag(CpuFlagN) == 1 {
		nf = "N"
	}
	if c.registers.GetFlag(CpuFlagH) == 1 {
		hf = "H"
	}
	if c.registers.GetFlag(CpuFlagC) == 1 {
		cf = "C"
	}

	opcode := CpuInstructions[c.membus.Read(c.registers.PC)]

	return fmt.Sprintf("[PC: 0x%04X] %v\t (0x%02X%02X) %s%s%s%s | A 0x%02X | BC 0x%04X | DE 0x%04X | HL 0x%04X | SP 0x%04X |",
		c.registers.PC,
		opcode.Type,
		c.membus.Read(c.registers.PC+1),
		c.membus.Read(c.registers.PC+2),
		zf,
		nf,
		hf,
		cf,
		c.registers.A,
		c.readFromRegister(RegisterTypeBC),
		c.readFromRegister(RegisterTypeDE),
		c.readFromRegister(RegisterTypeHL),
		c.readFromRegister(RegisterTypeSP),
	)
}

func NewCpu(bus *MemoryBus) *Cpu {
	return &Cpu{
		registers: &cpuRegisters{
			PC: 0x100,
			A:  0x01,
			SP: 0xFFFE,
		},
		membus: bus,
		ime:    true,
	}
}

func (c *Cpu) Step() error {
	if c.halted {
		emu_cycle(1)
		if c.interuptFlags != 0 {
			c.halted = false
		}
	} else {
		log.Print(c)
		_, instruction := c.fetchIsntruction()
		data, destAddress := c.fetchData(instruction)
		c.executeInstruction(instruction, data, destAddress)
	}

	if c.ime {
		c.interruptHandler()
		c.enablingIME = false
	}

	if c.enablingIME {
		c.ime = true
	}

	return nil
}

func (c *Cpu) interruptHandler() {

}

func (c *Cpu) stackPop() uint16 {
	val := c.membus.Read16(c.registers.SP)
	c.registers.SP += 2

	return val
}

func (c *Cpu) stackPush(value uint16) {
	c.registers.SP--
	c.membus.Write(c.registers.SP, uint8(value>>8))

	c.registers.SP--
	c.membus.Write(c.registers.SP, uint8(value))
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

	case AddressModeR_R:
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_N16:
		fallthrough
	case AddressModeN16:
		data = c.membus.Read16(c.registers.PC)
		emu_cycle(2)
		c.registers.PC += 2

	case AddressModeHL_SPR:
		fallthrough
	case AddressModeR_A8:
		fallthrough
	case AddressModeR_N8:
		fallthrough
	case AddressModeN8:
		data = uint16(c.membus.Read(c.registers.PC))
		c.registers.PC++
		emu_cycle(1)

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

	case AddressModeA8_R:
		destAddr = &CpuDestAddress{uint16(c.membus.Read(c.registers.PC)) | 0xFF00}
		c.registers.PC++
		emu_cycle(1)

	case AddressModeMR:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.membus.Read(c.readFromRegister(instruction.Register1)))
		emu_cycle(1)

	case AddressModeMR_N8:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.membus.Read(c.registers.PC))
		c.registers.PC++
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
		destAddr = &CpuDestAddress{hl}
		c.writeToRegister(RegisterTypeHL, hl+1)
		data = c.readFromRegister(instruction.Register2)

	case AddressModeHLD_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{hl}
		c.writeToRegister(RegisterTypeHL, hl-1)
		data = c.readFromRegister(instruction.Register2)

	case AddressModeA16_R:
		destAddr = &CpuDestAddress{
			c.membus.Read16(c.registers.PC),
		}
		emu_cycle(2)
		c.registers.PC += 2
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_A16:
		addr := c.membus.Read16(c.registers.PC)
		c.registers.PC += 2
		data = uint16(c.membus.Read(addr))
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
		c.execHALT()
	case InstructionTypeINC:
		c.execINC(instruction, data, destAddress)
	case InstructionTypeJP:
		c.execJP(instruction, data)
	case InstructionTypeJR:
		c.execJR(instruction, data)
	case InstructionTypeLD:
		c.execLD(instruction, data, destAddress)
	case InstructionTypeLDH:
		c.execLDH(instruction, data, destAddress)
	case InstructionTypeNOP:
		return nil // NOOP
	case InstructionTypeOR:
		c.execOR(instruction, data)
	case InstructionTypePOP:
		c.execPOP(instruction)
	case InstructionTypePUSH:
		c.execPUSH(instruction)
	case InstructionTypeRET:
		c.execRET(instruction)
	case InstructionTypeRETI:
		c.execRETI(instruction)
	case InstructionTypeRLA:
		c.execRLA()
	case InstructionTypeRLCA:
		c.execRLCA()
	case InstructionTypeRRA:
		c.execRRA()
	case InstructionTypeRRCA:
		c.execRRCA()
	case InstructionTypeRST:
		c.execRST(instruction)
	case InstructionTypeSBC:
		c.execSBC(instruction, data)
	case InstructionTypeSCF:
		c.execSCF()
	case InstructionTypeSTOP:
		c.execSTOP(data)
	case InstructionTypeSUB:
		c.execSUB(instruction, data)
	case InstructionTypeXOR:
		c.execXOR(data)
	case InstructionTypeCB:
		c.execCB(instruction, data)

	default:
		return errors.New("not implemented")
	}
	return nil
}
