package main

import (
	"errors"
	"fmt"
	"log"
)

type Cpu struct {
	registers *cpuRegisters

	halted bool

	// interrupts
	ime               bool
	enablingIME       bool
	interruptFlags    uint8
	interruptRegister uint8

	ctx *Context
}

func NewCpu(ctx *Context) {
	ctx.cpu = &Cpu{
		registers: &cpuRegisters{
			PC: 0x100,
			SP: 0xFFFE,
			A:  0x01,
			F:  0xB0,
			B:  0x00,
			C:  0x13,
			D:  0x00,
			E:  0xD8,
			H:  0x01,
			L:  0x4D,
		},
		ctx: ctx,
	}
}

func (c *Cpu) String(pc uint16) string {
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

	opcode := CpuInstructions[c.ctx.membus.Read(pc)]

	return fmt.Sprintf("%08X - %04X: %-7s (%02X %02X %02X) A: %02X F: %s%s%s%s BC: %02X%02X DE: %02X%02X HL: %02X%02X SP: %04X\n",
		c.ctx.ticks,
		pc,
		opcode.Type,
		c.ctx.membus.Read(pc),
		c.ctx.membus.Read(pc+1),
		c.ctx.membus.Read(pc+2),
		c.registers.A,
		zf,
		nf,
		hf,
		cf,
		c.registers.B,
		c.registers.C,
		c.registers.D,
		c.registers.E,
		c.registers.H,
		c.registers.L,
		c.registers.SP,
	)
}

func (c *Cpu) Step() error {
	if c.halted {
		c.ctx.EmuCycle(1)
		if c.interruptFlags != 0 {
			c.halted = false
		}
	} else {
		pc := c.registers.PC
		_, instruction := c.fetchIsntruction()
		c.ctx.EmuCycle(1)

		data, destAddress := c.fetchData(instruction)

		if c.ctx.debug.enbled {
			log.Print(c.String(pc))

			c.ctx.debug.Update()
			c.ctx.debug.Print()
		}

		c.executeInstruction(instruction, data, destAddress)
	}

	if c.ime {
		c.handleInterrupts()
		c.enablingIME = false
	}

	if c.enablingIME {
		c.ime = true
	}

	return nil
}

func (c *Cpu) stackPop() uint16 {
	val := c.ctx.membus.Read16(c.registers.SP)
	c.registers.SP += 2

	return val
}

func (c *Cpu) stackPush(value uint16) {
	c.registers.SP--
	c.ctx.membus.Write(c.registers.SP, uint8(value>>8))

	c.registers.SP--
	c.ctx.membus.Write(c.registers.SP, uint8(value))
}

func (c *Cpu) fetchIsntruction() (uint8, CpuInstriction) {
	opcode := c.ctx.membus.Read(c.registers.PC)
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
		data = c.ctx.membus.Read16(c.registers.PC)
		c.ctx.EmuCycle(2)
		c.registers.PC += 2

	case AddressModeHL_SPR:
		fallthrough
	case AddressModeR_A8:
		fallthrough
	case AddressModeR_N8:
		fallthrough
	case AddressModeN8:
		data = uint16(c.ctx.membus.Read(c.registers.PC))
		c.ctx.EmuCycle(1)
		c.registers.PC++

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
		data = uint16(c.ctx.membus.Read(address))
		c.ctx.EmuCycle(1)

	case AddressModeA8_R:
		destAddr = &CpuDestAddress{uint16(c.ctx.membus.Read(c.registers.PC)) | 0xFF00}
		c.ctx.EmuCycle(1)
		c.registers.PC++

	case AddressModeMR:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.ctx.membus.Read(c.readFromRegister(instruction.Register1)))
		c.ctx.EmuCycle(1)

	case AddressModeMR_N8:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.ctx.membus.Read(c.registers.PC))
		c.ctx.EmuCycle(1)
		c.registers.PC++

	case AddressModeR_HLI:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.ctx.membus.Read(hl))
		c.ctx.EmuCycle(1)
		c.writeToRegister(RegisterTypeHL, hl+1)

	case AddressModeR_HLD:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.ctx.membus.Read(hl))
		c.ctx.EmuCycle(1)
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
			c.ctx.membus.Read16(c.registers.PC),
		}
		c.ctx.EmuCycle(2)
		c.registers.PC += 2
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_A16:
		addr := c.ctx.membus.Read16(c.registers.PC)
		c.ctx.EmuCycle(2)
		c.registers.PC += 2
		data = uint16(c.ctx.membus.Read(addr))
		c.ctx.EmuCycle(1)
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
